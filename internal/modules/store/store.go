package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

type Store interface {
	// TxBegin starts a new transaction and returns a new ApplicationStore instance
	TxBegin(ctx context.Context) (Store, error)
	// TxCommit commits the transaction, store is not usable after this call
	TxCommit(ctx context.Context) error
	// TxRollback rolls the transaction back, store is not usable after this call
	TxRollback(ctx context.Context) error

	application.ApplicationStore
	device.DeviceStore
	gateway.GatewayStore
	gatewayprofile.GatewayProfileStore
	networkserver.NetworkServerStore
	organization.OrganizationStore
	user.UserStore
}
