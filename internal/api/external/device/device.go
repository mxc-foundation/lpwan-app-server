package device

import (
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
	dps "github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	orgd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/organization"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	appd "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	devd "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	nsportal "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Store defines db API used by device provision server
type Store interface {
	GetApplication(ctx context.Context, id int64) (appd.Application, error)
	Tx(ctx context.Context, f func(context.Context, *store.Handler) error) error
	UpdateDeviceWithDevProvisioingAttr(ctx context.Context, device *devd.Device) error
	CreateDeviceKeys(ctx context.Context, dc *devd.DeviceKeys) error
	GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (dps.DeviceProfile, error)
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (orgd.Organization, error)
	GetDeviceCount(ctx context.Context, filters devd.DeviceFilters) (int, error)
	CreateDevice(ctx context.Context, d *devd.Device) error
	GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (nsd.NetworkServer, error)
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (devd.Device, error)
	DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error
	GetApplicationWithIDAndOrganizationID(ctx context.Context, id, orgID int64) (appd.Application, error)
	GetDeviceProfileWithIDAndOrganizationID(ctx context.Context, id uuid.UUID, orgID int64, forUpdate bool) (dps.DeviceProfile, error)
	GetDefaultDeviceProfileID(ctx context.Context, orgID, nsID int64, forUpdate bool) (*uuid.UUID, error)
	GetDefaultNetworkServer(ctx context.Context) (nsd.NetworkServer, error)
}

// CreateDevice add new device and sync across all relevant servers
func CreateDevice(ctx context.Context, st Store, d *devd.Device, app *appd.Application,
	applicationServerID uuid.UUID, mxpCli pb.DSDeviceServiceClient, nsCli ns.NetworkServerServiceClient) error {
	org, err := st.GetOrganization(ctx, app.OrganizationID, true)
	if err != nil {
		return status.Errorf(codes.Unknown, "%v", err)
	}

	// Validate max. device count when != 0.
	if org.MaxDeviceCount != 0 {
		count, err := st.GetDeviceCount(ctx, devd.DeviceFilters{ApplicationID: app.OrganizationID})
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		if count >= org.MaxDeviceCount {
			return status.Errorf(codes.Unknown, "%v", errHandler.ErrOrganizationMaxDeviceCount)
		}
	}

	err = st.CreateDevice(ctx, d)
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
	_, err = nsCli.CreateDevice(ctx, &ns.CreateDeviceRequest{
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

// DeleteDevice deletes device and sync across all relevant servers. Must be called from within transaction
func DeleteDevice(ctx context.Context, st Store, devEUI lorawan.EUI64, mxpCli pb.DSDeviceServiceClient,
	psCli psPb.DeviceProvisionClient, nsCli *nscli.Client) error {
	n, err := st.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}
	nsClient, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
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
	_, err = nsClient.DeleteDevice(ctx, &ns.DeleteDeviceRequest{
		DevEui: devEUI[:],
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device error")
	}

	d, err := st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return err
	}
	// Set device to no server at PS
	if d.ProvisionID != "" {
		log.Debugf("DeleteDevice() Clear server addr for %v at PS", d.ProvisionID)
		_, err = psCli.SetDeviceServer(ctx, &psPb.SetDeviceServerRequest{ProvisionId: d.ProvisionID, Server: ""})
		if err != nil {
			return err
		}
	}

	if err := st.DeleteDevice(ctx, devEUI); err != nil {
		return err
	}
	return nil
}

// EnqueueDownlinkPayload adds the downlink payload to the network-server
// device-queue.
func EnqueueDownlinkPayload(ctx context.Context, st Store, devEUI lorawan.EUI64, confirmed bool, fPort uint8, data []byte) (uint32, error) {
	// get network-server and network-server api client
	n, err := st.GetNetworkServerForDevEUI(ctx, devEUI)
	if err != nil {
		return 0, errors.Wrap(err, "get network-server error")
	}

	nstruct := nsportal.NSStruct{
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
	d, err := st.GetDevice(ctx, devEUI, false)
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
