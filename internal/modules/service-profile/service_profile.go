package serviceprofile

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "service_profile"

type controller struct {
	st *store.Handler

	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{}
	return nil
}
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h

	return nil
}

// CreateServiceProfile creates the given service-profile.
func CreateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.CreateServiceProfile(ctx, sp); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, sp.NetworkServerID)
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
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}
		_, err = nsClient.CreateServiceProfile(ctx, &ns.CreateServiceProfileRequest{
			ServiceProfile: &sp.ServiceProfile,
		})
		if err != nil {
			return errors.Wrap(err, "create service-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetServiceProfile returns the service-profile matching the given id.
func GetServiceProfile(ctx context.Context, id uuid.UUID, localOnly bool) (sp ServiceProfile, err error) {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if sp, err = handler.GetServiceProfile(ctx, id); err != nil {
			return err
		}

		if localOnly {
			return nil
		}

		n, err := handler.GetNetworkServer(ctx, sp.NetworkServerID)
		if err != nil {
			return errors.Wrap(err, "get network-server errror")
		}

		// delete device from networkserver
		nsStruct := nscli.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		resp, err := nsClient.GetServiceProfile(ctx, &ns.GetServiceProfileRequest{
			Id: id.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "get service-profile error")
		}

		if resp.ServiceProfile == nil {
			return errors.New("service_profile must not be nil")
		}

		sp.ServiceProfile = *resp.ServiceProfile

		return nil
	}); err != nil {
		return sp, err
	}

	return sp, nil
}

// UpdateServiceProfile updates the given service-profile.
func UpdateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.UpdateServiceProfile(ctx, sp); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, sp.NetworkServerID)
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
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}
		_, err = nsClient.UpdateServiceProfile(ctx, &ns.UpdateServiceProfileRequest{
			ServiceProfile: &sp.ServiceProfile,
		})
		if err != nil {
			return errors.Wrap(err, "update service-profile error")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// DeleteServiceProfile deletes the service-profile matching the given id.
func DeleteServiceProfile(ctx context.Context, id uuid.UUID) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		n, err := handler.GetNetworkServerForServiceProfileID(ctx, id)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		if err := handler.DeleteServiceProfile(ctx, id); err != nil {
			return err
		}

		nsStruct := nscli.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nsStruct.GetNetworkServiceClient()
		if err != nil {
			return errors.Wrap(err, "get network-server client error")
		}

		_, err = nsClient.DeleteServiceProfile(ctx, &ns.DeleteServiceProfileRequest{
			Id: id.Bytes(),
		})
		if err != nil && status.Code(err) != codes.NotFound {
			return errors.Wrap(err, "delete service-profile error")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
