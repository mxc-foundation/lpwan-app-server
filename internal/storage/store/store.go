package store

import (
	"context"
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	fss "github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation/store"
	mcss "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/store"
	mcs "github.com/mxc-foundation/lpwan-app-server/internal/migrations/code/store"
	apps "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/store"
	dps "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile/store"
	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/store"
	fds "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment/store"
	gwps "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile/store"
	gws "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/store"
	mgs "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/store"
	orgs "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/store"
	sps "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/store"
	usrs "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/store"
	nss "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/store"
)

func NewStore() *Handler {
	return &Handler{
		basicStore:                        newStore(pgstore.New()),
		DeviceStore:                       ds.NewStore(pgstore.New()),
		DeviceProfileStore:                dps.NewStore(pgstore.New()),
		ApplicationStore:                  apps.NewApplicationStore(pgstore.New()),
		IntegrationsStore:                 apps.NewIntegrationStore(pgstore.New()),
		GatewayStore:                      gws.NewStore(pgstore.New()),
		OrganizationStore:                 orgs.NewStore(pgstore.New()),
		GatewayProfileStore:               gwps.NewStore(pgstore.New()),
		FUOTADeploymentStore:              fds.NewStore(pgstore.New()),
		MulticastGroupStore:               mgs.NewStore(pgstore.New()),
		NetworkServerStore:                nss.NewStore(pgstore.New()),
		MigrateCodeStore:                  mcs.NewStore(pgstore.New()),
		ServiceProfileStore:               sps.NewStore(pgstore.New()),
		RemoteFragmentaionSessionStore:    fss.NewStore(pgstore.New()),
		RemoteMulticastSetupStore:         mcss.NewSetupStore(pgstore.New()),
		RemoteMulticastClassCSessionStore: mcss.NewClassCSessionStore(pgstore.New()),
		UserStore:                         usrs.NewStore(pgstore.New()),
	}
}

type Handler struct {
	basicStore
	ds.DeviceStore
	dps.DeviceProfileStore
	apps.ApplicationStore
	apps.IntegrationsStore
	gws.GatewayStore
	orgs.OrganizationStore
	gwps.GatewayProfileStore
	fds.FUOTADeploymentStore
	mgs.MulticastGroupStore
	nss.NetworkServerStore
	mcs.MigrateCodeStore
	sps.ServiceProfileStore
	fss.RemoteFragmentaionSessionStore
	mcss.RemoteMulticastSetupStore
	mcss.RemoteMulticastClassCSessionStore
	usrs.UserStore
	inTX bool
}

// txBegin creates a transaction and returns a new instance of Handler that
// will either commit or rollback all the changes that done using this
// instance.
func (s *Handler) txBegin(ctx context.Context) (*Handler, error) {
	if s.inTX {
		return nil, fmt.Errorf("already in transaction")
	}
	store, err := s.basicStore.TxBegin(ctx)
	if err != nil {
		return nil, err
	}
	btx := *s
	btx.basicStore = newStore(store)
	btx.DeviceStore = ds.NewStore(store)
	btx.DeviceProfileStore = dps.NewStore(store)
	btx.ApplicationStore = apps.NewApplicationStore(store)
	btx.IntegrationsStore = apps.NewIntegrationStore(store)
	btx.GatewayStore = gws.NewStore(store)
	btx.OrganizationStore = orgs.NewStore(store)
	btx.GatewayProfileStore = gwps.NewStore(store)
	btx.FUOTADeploymentStore = fds.NewStore(store)
	btx.MulticastGroupStore = mgs.NewStore(store)
	btx.NetworkServerStore = nss.NewStore(store)
	btx.MigrateCodeStore = mcs.NewStore(store)
	btx.ServiceProfileStore = sps.NewStore(store)
	btx.RemoteFragmentaionSessionStore = fss.NewStore(store)
	btx.RemoteMulticastSetupStore = mcss.NewSetupStore(store)
	btx.RemoteMulticastClassCSessionStore = mcss.NewClassCSessionStore(store)
	btx.UserStore = usrs.NewStore(store)
	btx.inTX = true
	return &btx, nil
}

// Tx starts transaction and executes the function passing to it Handler
// using this transaction. It automatically rolls the transaction back if
// function returns an error. If the error has been caused by serialization
// error, it calls the function again. In order for serialization errors
// handling to work, the function should return Handler errors
// unchanged, or wrap them using %w.
func (h *Handler) Tx(ctx context.Context, f func(context.Context, *Handler) error) error {
	for {
		tx, err := h.txBegin(ctx)
		if err != nil {
			return err
		}
		err = f(ctx, tx)
		if err == nil {
			if err = tx.basicStore.TxCommit(ctx); err == nil {
				return nil
			}
		}
		_ = tx.basicStore.TxRollback(ctx)
		if h.IsErrorRepeat(err) {
			// failed due to the serialization error, try again
			continue
		}
		return err
	}
}

// IsErrorRepeat returns true if the error indicates that the action has failed
// because of the conflict with another transaction and that the application
// should try to repeat the action
func (h *Handler) IsErrorRepeat(err error) bool {
	return h.basicStore.IsErrorRepeat(err)
}

// TxBegin creates a transaction and returns a new instance of Handler that
// will either commit or rollback all the changes that done using this
// instance.
// This is only used in test files
func (s *Handler) TxBegin(ctx context.Context) (*Handler, error) {
	if s.inTX {
		return nil, fmt.Errorf("already in transaction")
	}
	store, err := s.basicStore.TxBegin(ctx)
	if err != nil {
		return nil, err
	}
	btx := *s
	btx.basicStore = newStore(store)
	btx.DeviceStore = ds.NewStore(store)
	btx.DeviceProfileStore = dps.NewStore(store)
	btx.ApplicationStore = apps.NewApplicationStore(store)
	btx.IntegrationsStore = apps.NewIntegrationStore(store)
	btx.GatewayStore = gws.NewStore(store)
	btx.OrganizationStore = orgs.NewStore(store)
	btx.GatewayProfileStore = gwps.NewStore(store)
	btx.FUOTADeploymentStore = fds.NewStore(store)
	btx.MulticastGroupStore = mgs.NewStore(store)
	btx.NetworkServerStore = nss.NewStore(store)
	btx.MigrateCodeStore = mcs.NewStore(store)
	btx.ServiceProfileStore = sps.NewStore(store)
	btx.RemoteFragmentaionSessionStore = fss.NewStore(store)
	btx.RemoteMulticastSetupStore = mcss.NewSetupStore(store)
	btx.RemoteMulticastClassCSessionStore = mcss.NewClassCSessionStore(store)
	btx.UserStore = usrs.NewStore(store)
	btx.inTX = true
	return &btx, nil
}

// TxRollback
// This is only used in test files
func (s *Handler) TxRollback(ctx context.Context) error {
	return s.basicStore.TxRollback(ctx)
}
