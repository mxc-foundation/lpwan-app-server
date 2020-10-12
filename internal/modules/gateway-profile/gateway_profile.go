package gatewayprofile

import (
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/pkg/errors"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"

	"github.com/gofrs/uuid"
	"golang.org/x/net/context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "gateway_profile"

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

func GetGatewayProfileCount(ctx context.Context) (int, error) {
	return ctrl.st.GetGatewayProfileCount(ctx)
}

func GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error) {
	return ctrl.st.GetGatewayProfiles(ctx, limit, offset)
}

// CreateGatewayProfile creates the given gateway-profile.
// This will create the gateway-profile at the network-server side and will
// create a local reference record.
func CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.CreateGatewayProfile(ctx, gp); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, gp.NetworkServerID)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
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

		_, err = nsClient.CreateGatewayProfile(ctx, &ns.CreateGatewayProfileRequest{
			GatewayProfile: &gp.GatewayProfile,
		})
		if err != nil {
			return errors.Wrap(err, "create gateway-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetGatewayProfile returns the gateway-profile matching the given id.
func GetGatewayProfile(ctx context.Context, id uuid.UUID) (gp GatewayProfile, err error) {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if gp, err = handler.GetGatewayProfile(ctx, id); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, gp.NetworkServerID)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
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

		resp, err := nsClient.GetGatewayProfile(ctx, &ns.GetGatewayProfileRequest{
			Id: id.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "get gateway-profile error")
		}

		if resp.GatewayProfile == nil {
			return errors.New("gateway_profile must not be nil")
		}

		gp.GatewayProfile = *resp.GatewayProfile

		return nil
	}); err != nil {
		return gp, err
	}

	return gp, nil
}

// UpdateGatewayProfile updates the given gateway-profile.
func UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.UpdateGatewayProfile(ctx, gp); err != nil {
			return err
		}

		n, err := handler.GetNetworkServer(ctx, gp.NetworkServerID)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
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

		_, err = nsClient.UpdateGatewayProfile(context.Background(), &ns.UpdateGatewayProfileRequest{
			GatewayProfile: &gp.GatewayProfile,
		})
		if err != nil {
			return errors.Wrap(err, "update gateway-profile error")
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// DeleteGatewayProfile deletes the gateway-profile matching the given id.
func DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.DeleteGatewayProfile(ctx, id); err != nil {
			return err
		}

		n, err := handler.GetNetworkServerForGatewayProfileID(ctx, id)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
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
		_, err = nsClient.DeleteGatewayProfile(ctx, &ns.DeleteGatewayProfileRequest{
			Id: id.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "delete gateway-profile error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
