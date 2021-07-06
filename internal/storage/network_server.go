package storage

import (
	"context"

	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"

	nss "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// NetworkServer defines the information to connect to a network-server.
type NetworkServer nss.NetworkServer

// Validate validates the network-server data.
func (ns NetworkServer) Validate() error {
	return nss.NetworkServer(ns).Validate()
}

// GetNetworkServer returns the network-server matching the given id.
func GetNetworkServer(ctx context.Context, handler *store.Handler, id int64) (NetworkServer, error) {
	res, err := handler.GetNetworkServer(ctx, id)
	return NetworkServer(res), err
}

// NetworkServerFilters provides filters for filtering network-servers.
type NetworkServerFilters nss.NetworkServerFilters

// SQL returns the SQL filters.
func (f NetworkServerFilters) SQL() string {
	return nss.NetworkServerFilters(f).SQL()
}

// GetNetworkServerCount returns the total number of network-servers.
func GetNetworkServerCount(ctx context.Context, handler *store.Handler, filters NetworkServerFilters) (int, error) {
	return handler.GetNetworkServerCount(ctx, nss.NetworkServerFilters(filters))
}

// GetNetworkServers returns a slice of network-servers.
func GetNetworkServers(ctx context.Context, handler *store.Handler, filters NetworkServerFilters) ([]NetworkServer, error) {
	res, err := handler.GetNetworkServers(ctx, nss.NetworkServerFilters(filters))
	if err != nil {
		return nil, err
	}

	var nss []NetworkServer
	for _, v := range res {
		nssItem := NetworkServer(v)
		nss = append(nss, nssItem)
	}

	return nss, nil
}

// GetNetworkServerForDevEUI returns the network-server for the given DevEUI.
func GetNetworkServerForDevEUI(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64) (NetworkServer, error) {
	res, err := handler.GetNetworkServerForDevEUI(ctx, devEUI)
	return NetworkServer(res), err
}

// GetNetworkServerForGatewayMAC returns the network-server for a given
// gateway mac.
func GetNetworkServerForGatewayMAC(ctx context.Context, handler *store.Handler, mac lorawan.EUI64) (NetworkServer, error) {
	res, err := handler.GetNetworkServerForGatewayMAC(ctx, mac)
	return NetworkServer(res), err
}

// GetNetworkServerForGatewayProfileID returns the network-server for the given
// gateway-profile id.
func GetNetworkServerForGatewayProfileID(ctx context.Context, handler *store.Handler, id uuid.UUID) (NetworkServer, error) {
	res, err := handler.GetNetworkServerForGatewayProfileID(ctx, id)
	return NetworkServer(res), err
}

// GetNetworkServerForMulticastGroupID returns the network-server for the given
// multicast-group id.
func GetNetworkServerForMulticastGroupID(ctx context.Context, handler *store.Handler, id uuid.UUID) (NetworkServer, error) {
	res, err := handler.GetNetworkServerForMulticastGroupID(ctx, id)
	return NetworkServer(res), err
}
