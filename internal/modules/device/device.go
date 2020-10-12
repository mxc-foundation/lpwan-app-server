package device

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"

	"github.com/apex/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	appd "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "device"

type controller struct {
	st *store.Handler
	s  Config

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) (err error) {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}
	ctrl = &controller{
		s: Config{
			ApplicationServerID: conf.ApplicationServer.ID,
		},
	}
	return nil
}

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl.st = h

	return nil
}

func GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (Device, error) {
	return ctrl.st.GetDevice(ctx, devEUI, forUpdate)
}

func CreateDevice(ctx context.Context, d *Device, app *appd.Application, applicationServerID uuid.UUID) error {
	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max device count.
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {

		org, err := handler.GetOrganization(ctx, app.OrganizationID, true)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		// Validate max. device count when != 0.
		if org.MaxDeviceCount != 0 {
			count, err := handler.GetDeviceCount(ctx, DeviceFilters{ApplicationID: app.OrganizationID})
			if err != nil {
				return status.Errorf(codes.Unknown, "%v", err)
			}

			if count >= org.MaxDeviceCount {
				return status.Errorf(codes.Unknown, "%v", errHandler.ErrOrganizationMaxDeviceCount)
			}
		}

		err = handler.CreateDevice(ctx, d)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		timestampCreatedAt, _ := ptypes.TimestampProto(time.Now())

		// add this device to m2m server
		dvClient, err := m2mcli.GetM2MDeviceServiceClient()
		if err != nil {
			log.WithError(err).Error("Create device")
			return status.Errorf(codes.Unavailable, err.Error())
		}

		_, err = dvClient.AddDeviceInM2MServer(context.Background(), &pb.AddDeviceInM2MServerRequest{
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
		n, err := handler.GetNetworkServerForDevEUI(ctx, d.DevEUI)
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

	}); err != nil {
		return errors.Wrap(err, "faile to create new device")
	}

	return nil
}

func DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.DeleteDevice(ctx, devEUI); err != nil {
			return err
		}

		// delete device from m2m server, this procedure should not block delete device from appserver once it's deleted from
		// network server successfully
		dvClient, err := m2mcli.GetM2MDeviceServiceClient()
		if err != nil {
			log.WithError(err).Error("get m2m-server client error")
			return err
		}

		_, err = dvClient.DeleteDeviceInM2MServer(context.Background(), &pb.DeleteDeviceInM2MServerRequest{
			DevEui: devEUI.String(),
		})
		if err != nil && status.Code(err) != codes.NotFound {
			log.WithError(err).Error("m2m-server delete device api error")
			return err
		}

		n, err := handler.GetNetworkServerForDevEUI(ctx, devEUI)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
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

		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to delete device")
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

	log.WithFields(log.Fields{
		"f_cnt":     resp.FCnt,
		"dev_eui":   devEUI,
		"confirmed": confirmed,
	}).Info("downlink device-queue item handled")

	return resp.FCnt, nil
}
