package storage

import (
	"context"
	"database/sql/driver"

	"github.com/brocaar/lorawan"

	gws "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Gateway represents a gateway.
type Gateway gws.Gateway

// GatewayListItem defines the gateway as list item.
type GatewayListItem gws.GatewayListItem

// GatewayPing represents a gateway ping.
type GatewayPing gws.GatewayPing

// GatewayPingRX represents a ping received by one of the gateways.
type GatewayPingRX gws.GatewayPingRX

// GPSPoint contains a GPS point.
type GPSPoint gws.GPSPoint

// GatewaysActiveInactive holds the avtive and inactive counts.
type GatewaysActiveInactive gws.GatewaysActiveInactive

// Value implements the driver.Valuer interface.
func (l GPSPoint) Value() (driver.Value, error) {
	return gws.GPSPoint(l).Value()
}

// Scan implements the sql.Scanner interface.
func (l *GPSPoint) Scan(src interface{}) error {
	return (*gws.GPSPoint)(l).Scan(src)
}

// Validate validates the gateway data.
func (g Gateway) Validate() error {
	return gws.Gateway(g).Validate()
}

// GatewayFilters provides filters for filtering gateways.
type GatewayFilters gws.GatewayFilters

// SQL returns the SQL filters.
func (f GatewayFilters) SQL() string {
	return gws.GatewayFilters(f).SQL()
}

// GetGatewaysForMACs returns a map of gateways given a slice of MACs.
func GetGatewaysForMACs(ctx context.Context, handler *store.Handler, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error) {
	res, err := handler.GetGatewaysForMACs(ctx, macs)
	if err != nil {
		return nil, err
	}

	out := make(map[lorawan.EUI64]Gateway)
	for k, v := range res {
		out[k] = Gateway(v)
	}

	return out, nil
}
