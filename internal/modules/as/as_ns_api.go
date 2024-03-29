package as

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq/hstore"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	pb "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/devprovision"
	"github.com/mxc-foundation/lpwan-app-server/internal/events/uplink"
	"github.com/mxc-foundation/lpwan-app-server/internal/gwping"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	dev "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	metricsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ApplicationServerAPI implements the as.ApplicationServerServer interface.
type ApplicationServerAPI struct {
	st             *store.Handler
	gIntegrations  []models.IntegrationHandler
	psCli          *pscli.Client
	nsCli          *nscli.Client
	devSessionList *devprovision.DeviceSessionList
}

// NewApplicationServerAPI returns a new ApplicationServerAPI.
func NewApplicationServerAPI(h *store.Handler, gIntegrations []models.IntegrationHandler,
	psCli *pscli.Client, nsCli *nscli.Client, devSessionList *devprovision.DeviceSessionList) *ApplicationServerAPI {
	return &ApplicationServerAPI{
		st:             h,
		gIntegrations:  gIntegrations,
		psCli:          psCli,
		nsCli:          nsCli,
		devSessionList: devSessionList,
	}
}

// HandleUplinkData handles incoming (uplink) data.
func (a *ApplicationServerAPI) HandleUplinkData(ctx context.Context, req *as.HandleUplinkDataRequest) (*empty.Empty, error) {
	if err := uplink.Handle(ctx, *req, a.st, a.gIntegrations, a.nsCli); err != nil {
		return nil, status.Errorf(codes.Internal, "handle uplink data error: %s", err)
	}

	return &empty.Empty{}, nil
}

// HandleDownlinkACK handles an ack on a downlink transmission.
func (a *ApplicationServerAPI) HandleDownlinkACK(ctx context.Context, req *as.HandleDownlinkACKRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.DevEui)

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		errStr := fmt.Sprintf("get device error: %s", err)
		logrus.WithField("dev_eui", devEUI).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}
	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		errStr := fmt.Sprintf("get application error: %s", err)
		logrus.WithField("id", d.ApplicationID).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	logrus.WithFields(logrus.Fields{
		"dev_eui": devEUI,
	}).Info("downlink device-queue item acknowledged")

	pl := pb.AckEvent{
		ApplicationId:   uint64(app.ID),
		ApplicationName: app.Name,
		DeviceName:      d.Name,
		DevEui:          devEUI[:],
		Acknowledged:    req.Acknowledged,
		FCnt:            req.FCnt,
		Tags:            make(map[string]string),
	}

	// set tags
	for k, v := range d.Tags.Map {
		if v.Valid {
			pl.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	err = integration.ForApplicationID(ctx, app.ID, a.gIntegrations, a.st).HandleAckEvent(ctx, vars, pl)
	if err != nil {
		logrus.WithError(err).Error("send ack event error")
	}

	return &empty.Empty{}, nil
}

// HandleTxAck handles a tx acknowledgement.
func (a *ApplicationServerAPI) HandleTxAck(ctx context.Context, req *as.HandleTxAckRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.DevEui)

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		errStr := fmt.Sprintf("get device error: %s", err)
		logrus.WithField("dev_eui", devEUI).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}
	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		errStr := fmt.Sprintf("get application error: %s", err)
		logrus.WithField("id", d.ApplicationID).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	logrus.WithFields(logrus.Fields{
		"dev_eui": devEUI,
	}).Info("downlink tx acknowledged by gateway")

	pl := pb.TxAckEvent{
		ApplicationId:   uint64(app.ID),
		ApplicationName: app.Name,
		DeviceName:      d.Name,
		DevEui:          devEUI[:],
		FCnt:            req.FCnt,
		Tags:            make(map[string]string),
	}

	// set tags
	for k, v := range d.Tags.Map {
		if v.Valid {
			pl.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	err = integration.ForApplicationID(ctx, app.ID, a.gIntegrations, a.st).HandleTxAckEvent(ctx, vars, pl)
	if err != nil {
		logrus.WithError(err).Error("send tx ack event error")
	}

	return &empty.Empty{}, nil
}

// HandleError handles an incoming error.
func (a *ApplicationServerAPI) HandleError(ctx context.Context, req *as.HandleErrorRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.DevEui)

	d, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		errStr := fmt.Sprintf("get device error: %s", err)
		logrus.WithField("dev_eui", devEUI).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		errStr := fmt.Sprintf("get application error: %s", err)
		logrus.WithField("id", d.ApplicationID).Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	logrus.WithFields(logrus.Fields{
		"type":    req.Type,
		"dev_eui": devEUI,
	}).Error(req.Error)

	var errType pb.ErrorType
	switch req.Type {
	case as.ErrorType_OTAA:
		errType = pb.ErrorType_OTAA
	case as.ErrorType_DATA_UP_FCNT_RESET:
		errType = pb.ErrorType_UPLINK_FCNT_RESET
	case as.ErrorType_DATA_UP_FCNT_RETRANSMISSION:
		errType = pb.ErrorType_UPLINK_FCNT_RETRANSMISSION
	case as.ErrorType_DATA_UP_MIC:
		errType = pb.ErrorType_UPLINK_MIC
	case as.ErrorType_DEVICE_QUEUE_ITEM_SIZE:
		errType = pb.ErrorType_DOWNLINK_PAYLOAD_SIZE
	case as.ErrorType_DEVICE_QUEUE_ITEM_FCNT:
		errType = pb.ErrorType_DOWNLINK_FCNT
	case as.ErrorType_DATA_DOWN_GATEWAY:
		errType = pb.ErrorType_DOWNLINK_GATEWAY
	}

	pl := pb.ErrorEvent{
		ApplicationId:   uint64(app.ID),
		ApplicationName: app.Name,
		DeviceName:      d.Name,
		DevEui:          devEUI[:],
		Type:            errType,
		Error:           req.Error,
		FCnt:            req.FCnt,
		Tags:            make(map[string]string),
	}

	// set tags
	for k, v := range d.Tags.Map {
		if v.Valid {
			pl.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	err = integration.ForApplicationID(ctx, app.ID, a.gIntegrations, a.st).HandleErrorEvent(ctx, vars, pl)
	if err != nil {
		errStr := fmt.Sprintf("send error notification to integration error: %s", err)
		logrus.Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	return &empty.Empty{}, nil
}

// HandleProprietaryUplink handles proprietary uplink payloads.
func (a *ApplicationServerAPI) HandleProprietaryUplink(ctx context.Context, req *as.HandleProprietaryUplinkRequest) (*empty.Empty, error) {
	if req.TxInfo == nil {
		return nil, status.Errorf(codes.InvalidArgument, "tx_info must not be nil")
	}

	processed, errForDev := devprovision.HandleReceivedFrame(ctx, req, a.st, a.psCli.GetDeviceProvisionServiceClient(),
		a.nsCli, a.devSessionList)
	if errForDev != nil {
		errStr := fmt.Sprintf("handle received proprietary error: %s", errForDev)
		logrus.Error(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}
	if !processed {
		err := gwping.HandleReceivedPing(ctx, req)
		if err != nil {
			errStr := fmt.Sprintf("handle received ping error: %s", err)
			logrus.Error(errStr)
			return nil, status.Errorf(codes.Internal, errStr)
		}
	}

	return &empty.Empty{}, nil
}

// SetDeviceStatus updates the device-status for the given device.
func (a *ApplicationServerAPI) SetDeviceStatus(ctx context.Context, req *as.SetDeviceStatusRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	copy(devEUI[:], req.DevEui)

	var d dev.Device
	var err error

	err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		d, err = handler.GetDevice(ctx, devEUI, true)
		if err != nil {
			return helpers.ErrToRPCError(errors.Wrap(err, "get device error"))
		}

		marg := int(req.Margin)
		d.DeviceStatusMargin = &marg

		if req.BatteryLevelUnavailable {
			d.DeviceStatusBattery = nil
			d.DeviceStatusExternalPower = false
		} else if req.ExternalPowerSource {
			d.DeviceStatusExternalPower = true
			d.DeviceStatusBattery = nil
		} else {
			d.DeviceStatusExternalPower = false
			d.DeviceStatusBattery = &req.BatteryLevel
		}

		if err = handler.UpdateDevice(ctx, &d); err != nil {
			return helpers.ErrToRPCError(errors.Wrap(err, "update device error"))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "get application error"))
	}

	pl := pb.StatusEvent{
		ApplicationId:           uint64(app.ID),
		ApplicationName:         app.Name,
		DeviceName:              d.Name,
		DevEui:                  d.DevEUI[:],
		Margin:                  req.Margin,
		ExternalPowerSource:     req.ExternalPowerSource,
		BatteryLevel:            float32(math.Round(float64(req.BatteryLevel*100))) / 100,
		BatteryLevelUnavailable: req.BatteryLevelUnavailable,
		Tags:                    make(map[string]string),
	}

	// set tags
	for k, v := range d.Tags.Map {
		if v.Valid {
			pl.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	err = integration.ForApplicationID(ctx, app.ID, a.gIntegrations, a.st).HandleStatusEvent(ctx, vars, pl)
	if err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "send status notification to handler error"))
	}

	return &empty.Empty{}, nil
}

// SetDeviceLocation updates the device-location.
func (a *ApplicationServerAPI) SetDeviceLocation(ctx context.Context, req *as.SetDeviceLocationRequest) (*empty.Empty, error) {
	if req.Location == nil {
		return nil, status.Errorf(codes.InvalidArgument, "location must not be nil")
	}

	var devEUI lorawan.EUI64
	copy(devEUI[:], req.DevEui)

	var d dev.Device
	var err error

	err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		d, err = handler.GetDevice(ctx, devEUI, true)
		if err != nil {
			return helpers.ErrToRPCError(errors.Wrap(err, "get device error"))
		}

		d.Latitude = &req.Location.Latitude
		d.Longitude = &req.Location.Longitude
		d.Altitude = &req.Location.Altitude

		if err = handler.UpdateDevice(ctx, &d); err != nil {
			return helpers.ErrToRPCError(errors.Wrap(err, "update device error"))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	app, err := a.st.GetApplication(ctx, d.ApplicationID)
	if err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "get application error"))
	}

	pl := pb.LocationEvent{
		ApplicationId:   uint64(app.ID),
		ApplicationName: app.Name,
		DeviceName:      d.Name,
		DevEui:          d.DevEUI[:],
		Location:        req.Location,
		UplinkIds:       req.UplinkIds,
		Tags:            make(map[string]string),
	}

	// set tags
	for k, v := range d.Tags.Map {
		if v.Valid {
			pl.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	err = integration.ForApplicationID(ctx, app.ID, a.gIntegrations, a.st).HandleLocationEvent(ctx, vars, pl)
	if err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "send location notification to handler error"))
	}

	return &empty.Empty{}, nil
}

// HandleGatewayStats handles the given gateway stats.
func (a *ApplicationServerAPI) HandleGatewayStats(ctx context.Context, req *as.HandleGatewayStatsRequest) (*empty.Empty, error) {
	var gatewayID lorawan.EUI64
	copy(gatewayID[:], req.GatewayId)

	ts := time.Now()
	gw, err := a.st.GetGateway(ctx, gatewayID, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "get gateway error"))
	}

	if gw.FirstSeenAt == nil {
		gw.FirstSeenAt = &ts
	}
	gw.LastSeenAt = &ts

	if loc := req.Location; loc != nil {
		gw.Latitude = loc.Latitude
		gw.Longitude = loc.Longitude
		gw.Altitude = loc.Altitude
	}

	gw.Metadata = hstore.Hstore{
		Map: make(map[string]sql.NullString),
	}
	for k, v := range req.Metadata {
		gw.Metadata.Map[k] = sql.NullString{Valid: true, String: v}
	}

	if err := a.st.UpdateGateway(ctx, &gw); err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "update gateway error"))
	}

	metrics := metricsmod.MetricsRecord{
		Time: ts,
		Metrics: map[string]float64{
			"rx_count":    float64(req.RxPacketsReceived),
			"rx_ok_count": float64(req.RxPacketsReceivedOk),
			"tx_count":    float64(req.TxPacketsReceived),
			"tx_ok_count": float64(req.TxPacketsEmitted),
		},
	}
	if err := metricsmod.SaveMetrics(ctx, fmt.Sprintf("gw:%s", gatewayID), metrics); err != nil {
		return nil, helpers.ErrToRPCError(errors.Wrap(err, "save metrics error"))
	}

	return &empty.Empty{}, nil
}
