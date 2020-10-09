package store

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
)

func NewStore(pg pgstore.PgStore) *nss {
	return &nss{
		pg: pg,
	}
}

type nss struct {
	pg pgstore.NetworkServerPgStore
}

type NetworkServerStore interface {
	CreateNetworkServer(ctx context.Context, n *NetworkServer) error
	GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error)
	UpdateNetworkServer(ctx context.Context, n *NetworkServer) error
	DeleteNetworkServer(ctx context.Context, id int64) error
	GetNetworkServerCount(ctx context.Context, filters NetworkServerFilters) (int, error)
	GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error)
	GetNetworkServers(ctx context.Context, filters NetworkServerFilters) ([]NetworkServer, error)
	GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error)
	GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetDefaultNetworkServer(ctx context.Context) (NetworkServer, error)

	// validator
	CheckCreateNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckReadNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error)
	CheckUpdateDeleteNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error)
}

func (h *nss) GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.pg.GetNetworkServerForDeviceProfileID(ctx, id)
}

func (h *nss) CreateNetworkServer(ctx context.Context, n *NetworkServer) error {
	return h.pg.CreateNetworkServer(ctx, n)
}
func (h *nss) GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error) {
	return h.pg.GetNetworkServer(ctx, id)
}
func (h *nss) UpdateNetworkServer(ctx context.Context, n *NetworkServer) error {
	return h.pg.UpdateNetworkServer(ctx, n)
}
func (h *nss) DeleteNetworkServer(ctx context.Context, id int64) error {
	return h.pg.DeleteNetworkServer(ctx, id)
}
func (h *nss) GetNetworkServerCount(ctx context.Context, filters NetworkServerFilters) (int, error) {
	return h.pg.GetNetworkServerCount(ctx, filters)
}
func (h *nss) GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	return h.pg.GetNetworkServerCountForOrganizationID(ctx, organizationID)
}
func (h *nss) GetNetworkServers(ctx context.Context, filters NetworkServerFilters) ([]NetworkServer, error) {
	return h.pg.GetNetworkServers(ctx, filters)
}
func (h *nss) GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error) {
	return h.pg.GetNetworkServersForOrganizationID(ctx, organizationID, limit, offset)
}
func (h *nss) GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error) {
	return h.pg.GetNetworkServerForDevEUI(ctx, devEUI)
}
func (h *nss) GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.pg.GetNetworkServerForServiceProfileID(ctx, id)
}
func (h *nss) GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error) {
	return h.pg.GetNetworkServerForGatewayMAC(ctx, mac)
}
func (h *nss) GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.pg.GetNetworkServerForGatewayProfileID(ctx, id)
}
func (h *nss) GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.pg.GetNetworkServerForMulticastGroupID(ctx, id)
}
func (h *nss) GetDefaultNetworkServer(ctx context.Context) (NetworkServer, error) {
	return h.pg.GetDefaultNetworkServer(ctx)
}

// validator
func (h *nss) CheckCreateNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckCreateNetworkServersAccess(ctx, username, organizationID, userID)
}
func (h *nss) CheckListNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckListNetworkServersAccess(ctx, username, organizationID, userID)
}

func (h *nss) CheckReadNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	return h.pg.CheckReadNetworkServerAccess(ctx, username, networkserverID, userID)
}
func (h *nss) CheckUpdateDeleteNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	return h.pg.CheckUpdateDeleteNetworkServerAccess(ctx, username, networkserverID, userID)
}
