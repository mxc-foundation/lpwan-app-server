package store

import (
	"context"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
)

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

func (h *Handler) CreateNetworkServer(ctx context.Context, n *NetworkServer) error {
	return h.store.CreateNetworkServer(ctx, n)
}
func (h *Handler) GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error) {
	return h.store.GetNetworkServer(ctx, id)
}
func (h *Handler) UpdateNetworkServer(ctx context.Context, n *NetworkServer) error {
	return h.store.UpdateNetworkServer(ctx, n)
}
func (h *Handler) DeleteNetworkServer(ctx context.Context, id int64) error {
	return h.store.DeleteNetworkServer(ctx, id)
}
func (h *Handler) GetNetworkServerCount(ctx context.Context, filters NetworkServerFilters) (int, error) {
	return h.store.GetNetworkServerCount(ctx, filters)
}
func (h *Handler) GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	return h.store.GetNetworkServerCountForOrganizationID(ctx, organizationID)
}
func (h *Handler) GetNetworkServers(ctx context.Context, filters NetworkServerFilters) ([]NetworkServer, error) {
	return h.store.GetNetworkServers(ctx, filters)
}
func (h *Handler) GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error) {
	return h.store.GetNetworkServersForOrganizationID(ctx, organizationID, limit, offset)
}
func (h *Handler) GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error) {
	return h.store.GetNetworkServerForDevEUI(ctx, devEUI)
}
func (h *Handler) GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.store.GetNetworkServerForServiceProfileID(ctx, id)
}
func (h *Handler) GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error) {
	return h.store.GetNetworkServerForGatewayMAC(ctx, mac)
}
func (h *Handler) GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.store.GetNetworkServerForGatewayProfileID(ctx, id)
}
func (h *Handler) GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (NetworkServer, error) {
	return h.store.GetNetworkServerForMulticastGroupID(ctx, id)
}
func (h *Handler) GetDefaultNetworkServer(ctx context.Context) (NetworkServer, error) {
	return h.store.GetDefaultNetworkServer(ctx)
}

// validator
func (h *Handler) CheckCreateNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckCreateNetworkServersAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckListNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckListNetworkServersAccess(ctx, username, organizationID, userID)
}

func (h *Handler) CheckReadNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	return h.store.CheckReadNetworkServerAccess(ctx, username, networkserverID, userID)
}
func (h *Handler) CheckUpdateDeleteNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteNetworkServerAccess(ctx, username, networkserverID, userID)
}

// NetworkServer defines the information to connect to a network-server.
type NetworkServer struct {
	ID                          int64     `db:"id"`
	CreatedAt                   time.Time `db:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at"`
	Name                        string    `db:"name"`
	Server                      string    `db:"server"`
	CACert                      string    `db:"ca_cert"`
	TLSCert                     string    `db:"tls_cert"`
	TLSKey                      string    `db:"tls_key"`
	RoutingProfileCACert        string    `db:"routing_profile_ca_cert"`
	RoutingProfileTLSCert       string    `db:"routing_profile_tls_cert"`
	RoutingProfileTLSKey        string    `db:"routing_profile_tls_key"`
	GatewayDiscoveryEnabled     bool      `db:"gateway_discovery_enabled"`
	GatewayDiscoveryInterval    int       `db:"gateway_discovery_interval"`
	GatewayDiscoveryTXFrequency int       `db:"gateway_discovery_tx_frequency"`
	GatewayDiscoveryDR          int       `db:"gateway_discovery_dr"`
	Version                     string    `db:"version"`
	Region                      string    `db:"region"`
}

// Validate validates the network-server data.
func (ns NetworkServer) Validate() error {
	if strings.TrimSpace(ns.Name) == "" || len(ns.Name) > 100 {
		return ErrNetworkServerInvalidName
	}

	if ns.GatewayDiscoveryEnabled && ns.GatewayDiscoveryInterval <= 0 {
		return ErrInvalidGatewayDiscoveryInterval
	}
	return nil
}

// NetworkServerFilters provides filters for filtering network-servers.
type NetworkServerFilters struct {
	OrganizationID int64 `db:"organization_id"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f NetworkServerFilters) SQL() string {
	var filters []string

	if f.OrganizationID != 0 {
		filters = append(filters, "sp.organization_id = :organization_id")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}
