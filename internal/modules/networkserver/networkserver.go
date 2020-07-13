package networkserver

import (
	"context"
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type NetworkServerStore interface {
	CreateNetworkServer(ctx context.Context, n *NetworkServer) error
	GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error)
	UpdateNetworkServer(ctx context.Context, n *NetworkServer) error
	DeleteNetworkServer(ctx context.Context, id int64) error
	GetNetworkServerCount(ctx context.Context) (int, error)
	GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error)
	GetNetworkServers(ctx context.Context, limit, offset int) ([]NetworkServer, error)
	GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error)
	GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetDefaultNetworkServer(ctx context.Context) (NetworkServer, error)

	// validator
	CheckCreateNetworkServersAccess(username string, organizationID, userID int64) (bool, error)
	CheckListNetworkServersAccess(username string, organizationID, userID int64) (bool, error)

	CheckReadNetworkServerAccess(username string, networkserverID, userID int64) (bool, error)
	CheckUpdateDeleteNetworkServerAccess(username string, networkserverID, userID int64) (bool, error)
}

type Controller struct {
	St        NetworkServerStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}
