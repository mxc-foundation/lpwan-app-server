package multicast

import (
	"context"

	"github.com/brocaar/lorawan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "multicast_grou["

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

// CreateMulticastGroup creates the given multicast-group.
func CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.CreateMulticastGroup(ctx, mg); err != nil {
			return err
		}

		n, err := handler.GetNetworkServerForServiceProfileID(ctx, mg.ServiceProfileID)
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
			return err
		}

		_, err = nsClient.CreateMulticastGroup(ctx, &ns.CreateMulticastGroupRequest{
			MulticastGroup: &mg.MulticastGroup,
		})
		if err != nil {
			return errors.Wrap(err, "create multicast-group error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// GetMulticastGroup returns the multicast-group given an id.
func GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate, localOnly bool) (mg MulticastGroup, err error) {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if mg, err = handler.GetMulticastGroup(ctx, id, forUpdate); err != nil {
			return err
		}

		if localOnly {
			return nil
		}

		n, err := handler.GetNetworkServerForServiceProfileID(ctx, id)
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
			return err
		}

		resp, err := nsClient.GetMulticastGroup(ctx, &ns.GetMulticastGroupRequest{
			Id: id.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "get multicast-group error")
		}

		if resp.MulticastGroup == nil {
			return errors.New("multicast_group must not be nil")
		}

		mg.MulticastGroup = *resp.MulticastGroup

		return nil
	}); err != nil {
		return mg, err
	}

	return mg, nil
}

// UpdateMulticastGroup updates the given multicast-group.
func UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.UpdateMulticastGroup(ctx, mg); err != nil {
			return err
		}

		n, err := handler.GetNetworkServerForServiceProfileID(ctx, mg.ServiceProfileID)
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
			return err
		}

		_, err = nsClient.UpdateMulticastGroup(ctx, &ns.UpdateMulticastGroupRequest{
			MulticastGroup: &mg.MulticastGroup,
		})
		if err != nil {
			return errors.Wrap(err, "update multicast-group error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// DeleteMulticastGroup deletes a multicast-group given an id.
func DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		n, err := handler.GetNetworkServerForMulticastGroupID(ctx, id)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		if err := handler.DeleteMulticastGroup(ctx, id); err != nil {
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
			return err
		}

		_, err = nsClient.DeleteMulticastGroup(ctx, &ns.DeleteMulticastGroupRequest{
			Id: id.Bytes(),
		})
		if err != nil && grpc.Code(err) != codes.NotFound {
			return errors.Wrap(err, "delete multicast-group error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// AddDeviceToMulticastGroup adds the given device to the given multicast-group.
// It is recommended that db is a transaction.
func AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.AddDeviceToMulticastGroup(ctx, multicastGroupID, devEUI); err != nil {
			return err
		}

		n, err := handler.GetNetworkServerForMulticastGroupID(ctx, multicastGroupID)
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
			return err
		}

		_, err = nsClient.AddDeviceToMulticastGroup(ctx, &ns.AddDeviceToMulticastGroupRequest{
			DevEui:           devEUI[:],
			MulticastGroupId: multicastGroupID.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "add device to multicast-group error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// RemoveDeviceFromMulticastGroup removes the given device from the given
// multicast-group.
func RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.RemoveDeviceFromMulticastGroup(ctx, multicastGroupID, devEUI); err != nil {
			return err
		}

		n, err := handler.GetNetworkServerForMulticastGroupID(ctx, multicastGroupID)
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
			return err
		}

		_, err = nsClient.RemoveDeviceFromMulticastGroup(ctx, &ns.RemoveDeviceFromMulticastGroupRequest{
			DevEui:           devEUI[:],
			MulticastGroupId: multicastGroupID.Bytes(),
		})
		if err != nil && status.Code(err) != codes.NotFound {
			return errors.Wrap(err, "remove device from multicast-group error")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
