package device

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	appd "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func CreateDevice(ctx context.Context, h *store.Handler, d *Device, app *appd.Application,
	applicationServerID uuid.UUID, mxpCli pb.DSDeviceServiceClient) error {
	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max device count.
	if !h.InTx() {
		return fmt.Errorf("CreateDevice must be called from within a transaction")
	}

	org, err := h.GetOrganization(ctx, app.OrganizationID, true)
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	// Validate max. device count when != 0.
	if org.MaxDeviceCount != 0 {
		count, err := h.GetDeviceCount(ctx, DeviceFilters{ApplicationID: app.OrganizationID})
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		if count >= org.MaxDeviceCount {
			return status.Errorf(codes.Unknown, "%v", errHandler.ErrOrganizationMaxDeviceCount)
		}
	}

	err = h.CreateDevice(ctx, d)
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	timestampCreatedAt := timestamppb.New(time.Now())

	// add this device to m2m server
	_, err = mxpCli.AddDeviceInM2MServer(context.Background(), &pb.AddDeviceInM2MServerRequest{
		OrgId: app.OrganizationID,
		DevProfile: &pb.AppServerDeviceProfile{
			DevEui:        d.DevEUI.String(),
			ApplicationId: d.ApplicationID,
			Name:          d.Name,
			CreatedAt:     timestampCreatedAt,
		},
	})
	if err != nil {
		return status.Errorf(codes.Unknown, "m2m server create device api error: %v", err)
	}

	// add this device to network server
	n, err := h.GetNetworkServerForDevEUI(ctx, d.DevEUI)
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	// add this device to network server
	nStruct := &nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nStruct.GetNetworkServiceClient()
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	_, err = client.CreateDevice(ctx, &ns.CreateDeviceRequest{
		Device: &ns.Device{
			DevEui:            d.DevEUI[:],
			DeviceProfileId:   d.DeviceProfileID.Bytes(),
			ServiceProfileId:  app.ServiceProfileID.Bytes(),
			RoutingProfileId:  applicationServerID.Bytes(),
			SkipFCntCheck:     d.SkipFCntCheck,
			ReferenceAltitude: d.ReferenceAltitude,
		},
	})
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	return nil
}

func DeleteDevice(ctx context.Context, h *store.Handler, devEUI lorawan.EUI64, mxpCli pb.DSDeviceServiceClient,
	psCli psPb.DeviceProvisionClient) error {
	if !h.InTx() {
		return fmt.Errorf("DeleteDevice must be called from within a transaction")
	}

	n, err := h.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// Get Device Provision ID from Description
	devicepid := ""
	d, err := h.GetDevice(ctx, devEUI, false)
	if err == nil {
		devdesc := d.Description
		if strings.Contains(devdesc, "PID:") {
			lines := strings.Split(devdesc, "\n")
			for _, line := range lines {
				if strings.Index(line, "PID:") == 0 {
					fields := strings.Split(line, ":")
					if len(fields) >= 2 {
						devicepid = strings.Trim(fields[1], " ")
					}
				}
			}
		}
	}

	if err := h.DeleteDevice(ctx, devEUI); err != nil {
		return err
	}

	// delete device from m2m server, this procedure should not block delete device from appserver once it's deleted from
	// network server successfully
	_, err = mxpCli.DeleteDeviceInM2MServer(context.Background(), &pb.DeleteDeviceInM2MServerRequest{
		DevEui: devEUI.String(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		logrus.WithError(err).Error("m2m-server delete device api error")
		return err
	}

	// delete device from networkserver
	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = client.DeleteDevice(ctx, &ns.DeleteDeviceRequest{
		DevEui: devEUI[:],
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device error")
	}

	// Set device to no server at PS
	if devicepid != "" {
		log.Debugf("DeleteDevice() Clear server addr for %v at PS", devicepid)
		_, err = psCli.SetDeviceServer(ctx, &psPb.SetDeviceServerRequest{ProvisionId: devicepid, Server: ""})
		if err != nil {
			return err
		}
	}

	return nil
}

// EnqueueDownlinkPayload adds the downlink payload to the network-server
// device-queue.
func EnqueueDownlinkPayload(ctx context.Context, h *store.Handler, devEUI lorawan.EUI64, confirmed bool, fPort uint8, data []byte) (uint32, error) {
	// get network-server and network-server api client
	n, err := h.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return 0, errors.Wrap(err, "get network-server error")
	}

	nstruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}

	nsClient, err := nstruct.GetNetworkServiceClient()
	if err != nil {
		return 0, errors.Wrap(err, "get network-server client error")
	}

	// get fCnt to use for encrypting and enqueueing
	resp, err := nsClient.GetNextDownlinkFCntForDevEUI(context.Background(), &ns.GetNextDownlinkFCntForDevEUIRequest{
		DevEui: devEUI[:],
	})
	if err != nil {
		return 0, errors.Wrap(err, "get next downlink fcnt for deveui error")
	}

	// get device
	d, err := h.GetDevice(ctx, devEUI, false)
	if err != nil {
		return 0, errors.Wrap(err, "get device error")
	}

	// encrypt payload
	b, err := lorawan.EncryptFRMPayload(d.AppSKey, false, d.DevAddr, resp.FCnt, data)
	if err != nil {
		return 0, errors.Wrap(err, "encrypt frmpayload error")
	}

	// enqueue device-queue item
	_, err = nsClient.CreateDeviceQueueItem(ctx, &ns.CreateDeviceQueueItemRequest{
		Item: &ns.DeviceQueueItem{
			DevAddr:    d.DevAddr[:],
			DevEui:     devEUI[:],
			FrmPayload: b,
			FCnt:       resp.FCnt,
			FPort:      uint32(fPort),
			Confirmed:  confirmed,
		},
	})
	if err != nil {
		return 0, errors.Wrap(err, "create device-queue item error")
	}

	logrus.WithFields(logrus.Fields{
		"f_cnt":     resp.FCnt,
		"dev_eui":   devEUI,
		"confirmed": confirmed,
	}).Info("downlink device-queue item handled")

	return resp.FCnt, nil
}
