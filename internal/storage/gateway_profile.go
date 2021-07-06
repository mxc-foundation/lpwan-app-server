package storage

import (
	"context"

	gwps "github.com/mxc-foundation/lpwan-app-server/internal/api/external/gp"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Modulations
const (
	ModulationFSK  = gwps.ModulationFSK
	ModulationLoRa = gwps.ModulationLoRa
)

// ExtraChannel defines an extra channel for the gateway-profile.
type ExtraChannel gwps.ExtraChannel

// GatewayProfile defines a gateway-profile.
type GatewayProfile gwps.GatewayProfile

// GatewayProfileMeta defines the gateway-profile meta record.
type GatewayProfileMeta gwps.GatewayProfileMeta

// GetGatewayProfileCount returns the total number of gateway-profiles.
func GetGatewayProfileCount(ctx context.Context, handler *store.Handler) (int, error) {
	return handler.GetGatewayProfileCount(ctx)
}

// GetGatewayProfileCountForNetworkServerID returns the total number of
// gateway-profiles given a network-server ID.
func GetGatewayProfileCountForNetworkServerID(ctx context.Context, handler *store.Handler, networkServerID int64) (int, error) {
	return handler.GetGatewayProfileCountForNetworkServerID(ctx, networkServerID)
}

// GetGatewayProfiles returns a slice of gateway-profiles.
func GetGatewayProfiles(ctx context.Context, handler *store.Handler, limit, offset int) ([]GatewayProfileMeta, error) {
	res, err := handler.GetGatewayProfiles(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	var gpMetaList []GatewayProfileMeta
	for _, v := range res {
		gpMetaItem := GatewayProfileMeta(v)
		gpMetaList = append(gpMetaList, gpMetaItem)
	}

	return gpMetaList, nil
}

// GetGatewayProfilesForNetworkServerID returns a slice of gateway-profiles
// for the given network-server ID.
func GetGatewayProfilesForNetworkServerID(ctx context.Context, handler *store.Handler, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error) {
	res, err := handler.GetGatewayProfilesForNetworkServerID(ctx, networkServerID, limit, offset)
	if err != nil {
		return nil, err
	}

	var gpMetaList []GatewayProfileMeta
	for _, v := range res {
		gpMetaItem := GatewayProfileMeta(v)
		gpMetaList = append(gpMetaList, gpMetaItem)
	}

	return gpMetaList, nil
}
