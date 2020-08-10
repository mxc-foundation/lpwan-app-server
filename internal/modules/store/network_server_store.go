package store

import (
	"context"
	"errors"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
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
	CheckCreateNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckReadNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error)
	CheckUpdateDeleteNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error)
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
	if ns.GatewayDiscoveryEnabled && ns.GatewayDiscoveryInterval <= 0 {
		return errors.New("ErrInvalidGatewayDiscoveryInterval")
	}
	return nil
}
