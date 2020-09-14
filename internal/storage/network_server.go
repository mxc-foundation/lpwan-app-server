package storage

import (
	"context"
	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// NetworkServer defines the information to connect to a network-server.
type NetworkServer store.NetworkServer

// Validate validates the network-server data.
func (ns NetworkServer) Validate() error {
	return store.NetworkServer(ns).Validate()
}

// CreateNetworkServer creates the given network-server.
func CreateNetworkServer(ctx context.Context, handler *store.Handler, n *NetworkServer) error {
	return handler.CreateNetworkServer(ctx, (*store.NetworkServer)(n))
}

// GetNetworkServer returns the network-server matching the given id.
func GetNetworkServer(ctx context.Context, handler *store.Handler, id int64) (NetworkServer, error) {
	res, err := handler.GetNetworkServer(ctx, id)
	return NetworkServer(res), err
}

// UpdateNetworkServer updates the given network-server.
func UpdateNetworkServer(ctx context.Context, handler *store.Handler, n *NetworkServer) error {
	return handler.UpdateNetworkServer(ctx, (*store.NetworkServer)(n))
}

// DeleteNetworkServer deletes the network-server matching the given id.
func DeleteNetworkServer(ctx context.Context, handler *store.Handler, id int64) error {
	return handler.DeleteNetworkServer(ctx, id)
}

// NetworkServerFilters provides filters for filtering network-servers.
type NetworkServerFilters store.NetworkServerFilters

// SQL returns the SQL filters.
func (f NetworkServerFilters) SQL() string {
	return store.NetworkServerFilters(f).SQL()
}

// GetNetworkServerCount returns the total number of network-servers.
func GetNetworkServerCount(ctx context.Context, handler *store.Handler, filters NetworkServerFilters) (int, error) {
	return handler.GetNetworkServerCount(ctx, store.NetworkServerFilters(filters))
}

// GetNetworkServers returns a slice of network-servers.
func GetNetworkServers(ctx context.Context, handler *store.Handler, filters NetworkServerFilters) ([]NetworkServer, error) {
	res, err := handler.GetNetworkServers(ctx, store.NetworkServerFilters(filters))
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
