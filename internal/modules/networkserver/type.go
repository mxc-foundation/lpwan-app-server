package networkserver

import (
	"errors"
	"time"
)

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
