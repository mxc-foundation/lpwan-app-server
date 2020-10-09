package networkserver_portal

import (
	"fmt"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	nscliLegacy "github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "network_server"

type controller struct {
	st                          *store.Handler
	applicationServerID         uuid.UUID
	p                           Pool
	applicationServerPublicHost string
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	appServerID, err := uuid.FromString(conf.ApplicationServer.ID)
	if err != nil {
		return errors.Wrap(err, "failed to convert applicationserver id from string to uuid")
	}

	ctrl = &controller{
		applicationServerID:         appServerID,
		applicationServerPublicHost: conf.ApplicationServer.API.PublicHost,
	}

	return nil
}

func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	if err := nscliLegacy.Setup(); err != nil {
		return err
	}

	ctrl.st = h
	ctrl.p = &pool{
		nsClients: make(map[string]nsClient),
	}

	return nil
}

// GetNetworkServerForDevEUI :
func GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error) {
	return ctrl.st.GetNetworkServerForDevEUI(ctx, devEUI)
}

// GetNetworkServer :
func GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error) {
	return ctrl.st.GetNetworkServer(ctx, id)
}

// GetNetworkServerForGatewayProfileID :
func GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return ctrl.st.GetNetworkServerForGatewayProfileID(ctx, id)
}

// GetNetworkServerForGatewayMAC :
func GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error) {
	return ctrl.st.GetNetworkServerForGatewayMAC(ctx, mac)
}

func CreateNetworkServer(ctx context.Context, n *NetworkServer) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.CreateNetworkServer(ctx, n); err != nil {
			return err
		}

		nsStruct := NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		rpID := ctrl.applicationServerID

		_, err = nsClient.CreateRoutingProfile(ctx, &ns.CreateRoutingProfileRequest{
			RoutingProfile: &ns.RoutingProfile{
				Id:      rpID.Bytes(),
				AsId:    ctrl.applicationServerPublicHost,
				CaCert:  n.RoutingProfileCACert,
				TlsCert: n.RoutingProfileTLSCert,
				TlsKey:  n.RoutingProfileTLSKey,
			},
		})
		if err != nil {
			return errors.Wrap(err, "create routing-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func UpdateNetworkServer(ctx context.Context, n *NetworkServer) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.UpdateNetworkServer(ctx, n); err != nil {
			return err
		}

		nsStruct := NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		rpID := ctrl.applicationServerID

		_, err = nsClient.UpdateRoutingProfile(ctx, &ns.UpdateRoutingProfileRequest{
			RoutingProfile: &ns.RoutingProfile{
				Id:      rpID.Bytes(),
				AsId:    ctrl.applicationServerPublicHost,
				CaCert:  n.RoutingProfileCACert,
				TlsCert: n.RoutingProfileTLSCert,
				TlsKey:  n.RoutingProfileTLSKey,
			},
		})
		if err != nil {
			return errors.Wrap(err, "update routing-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func DeleteNetworkServer(ctx context.Context, id int64) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.DeleteNetworkServer(ctx, id); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, id)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		nsStruct := NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		rpID := ctrl.applicationServerID

		_, err = nsClient.DeleteRoutingProfile(ctx, &ns.DeleteRoutingProfileRequest{
			Id: rpID.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "delete routing-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
