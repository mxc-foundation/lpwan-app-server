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

// CreateGateway creates the given Gateway.
func CreateGateway(ctx context.Context, handler *store.Handler, gw *Gateway) error {
	return handler.CreateGateway(ctx, (*gws.Gateway)(gw))
}

// UpdateGateway updates the given Gateway.
func UpdateGateway(ctx context.Context, handler *store.Handler, gw *Gateway) error {
	return handler.UpdateGateway(ctx, (*gws.Gateway)(gw))
}

// DeleteGateway deletes the gateway matching the given MAC.
func DeleteGateway(ctx context.Context, handler *store.Handler, mac lorawan.EUI64) error {
	return handler.DeleteGateway(ctx, mac)
}

// GetGateway returns the gateway for the given mac.
func GetGateway(ctx context.Context, handler *store.Handler, mac lorawan.EUI64, forUpdate bool) (Gateway, error) {
	gw, err := handler.GetGateway(ctx, mac, forUpdate)
	return Gateway(gw), err
}

// GatewayFilters provides filters for filtering gateways.
type GatewayFilters gws.GatewayFilters

// SQL returns the SQL filters.
func (f GatewayFilters) SQL() string {
	return gws.GatewayFilters(f).SQL()
}

// GetGatewayCount returns the total number of gateways.
func GetGatewayCount(ctx context.Context, handler *store.Handler, filters GatewayFilters) (int, error) {
	return handler.GetGatewayCount(ctx, (gws.GatewayFilters)(filters))
}

// GetGateways returns a slice of gateways sorted by name.
func GetGateways(ctx context.Context, handler *store.Handler, filters GatewayFilters) ([]GatewayListItem, error) {
	res, err := handler.GetGateways(ctx, gws.GatewayFilters(filters))
	if err != nil {
		return nil, err
	}

	var gwList []GatewayListItem
	for _, v := range res {
		gwItem := GatewayListItem(v)
		gwList = append(gwList, gwItem)
	}
	return gwList, nil
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

// CreateGatewayPing creates the given gateway ping.
func CreateGatewayPing(ctx context.Context, handler *store.Handler, ping *GatewayPing) error {
	return handler.CreateGatewayPing(ctx, (*gws.GatewayPing)(ping))
}

// GetGatewayPing returns the ping matching the given id.
func GetGatewayPing(ctx context.Context, handler *store.Handler, id int64) (GatewayPing, error) {
	gwPing, err := handler.GetGatewayPing(ctx, id)
	return GatewayPing(gwPing), err
}

// CreateGatewayPingRX creates the received ping.
func CreateGatewayPingRX(ctx context.Context, handler *store.Handler, rx *GatewayPingRX) error {
	return handler.CreateGatewayPingRX(ctx, (*gws.GatewayPingRX)(rx))
}

// DeleteAllGatewaysForOrganizationID deletes all gateways for a given
// organization id.
func DeleteAllGatewaysForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) error {
	return handler.DeleteAllGatewaysForOrganizationID(ctx, organizationID)
}

// GetGatewayPingRXForPingID returns the received gateway pings for the given
// ping ID.
func GetGatewayPingRXForPingID(ctx context.Context, handler *store.Handler, pingID int64) ([]GatewayPingRX, error) {
	res, err := handler.GetGatewayPingRXForPingID(ctx, pingID)
	if err != nil {
		return nil, err
	}

	var gwPingRXList []GatewayPingRX
	for _, v := range res {
		gwPingRXItem := GatewayPingRX(v)
		gwPingRXList = append(gwPingRXList, gwPingRXItem)
	}
	return gwPingRXList, nil
}

// GetLastGatewayPingAndRX returns the last gateway ping and RX for the given
// gateway MAC.
func GetLastGatewayPingAndRX(ctx context.Context, handler *store.Handler, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error) {
	gwPing, gwPingRX, err := handler.GetLastGatewayPingAndRX(ctx, mac)
	if err != nil {
		return GatewayPing{}, nil, err
	}

	var gwPingRXList []GatewayPingRX
	for _, v := range gwPingRX {
		gwPingRXItem := GatewayPingRX(v)
		gwPingRXList = append(gwPingRXList, gwPingRXItem)
	}
	return GatewayPing(gwPing), gwPingRXList, err
}

// GetGatewaysActiveInactive returns the active / inactive gateways.
func GetGatewaysActiveInactive(ctx context.Context, handler *store.Handler, organizationID int64) (GatewaysActiveInactive, error) {
	res, err := handler.GetGatewaysActiveInactive(ctx, organizationID)
	return GatewaysActiveInactive(res), err

}
