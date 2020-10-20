package networkserver_portal

import (
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
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

	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {

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
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h
	ctrl.p = &pool{
		nsClients: make(map[string]nsClient),
	}

	return nil
}

// CreateNetworkServer :
func CreateNetworkServer(ctx context.Context, n *NetworkServer, h *store.Handler) error {
	if err := h.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
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

		_, err = nsClient.CreateRoutingProfile(ctx, &ns.CreateRoutingProfileRequest{
			RoutingProfile: &ns.RoutingProfile{
				Id:      ctrl.applicationServerID.Bytes(),
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

func UpdateNetworkServer(ctx context.Context, n *NetworkServer, h *store.Handler) error {
	if err := h.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
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

func DeleteNetworkServer(ctx context.Context, id int64, h *store.Handler) error {
	if err := h.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		n, err := handler.GetNetworkServer(ctx, id)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		if err := handler.DeleteNetworkServer(ctx, id); err != nil {
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
