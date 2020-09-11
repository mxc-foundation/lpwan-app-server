package storage

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/gofrs/uuid"
)

// Modulations
const (
	ModulationFSK  = store.ModulationFSK
	ModulationLoRa = store.ModulationLoRa
)

// ExtraChannel defines an extra channel for the gateway-profile.
type ExtraChannel store.ExtraChannel

// GatewayProfile defines a gateway-profile.
type GatewayProfile store.GatewayProfile

// GatewayProfileMeta defines the gateway-profile meta record.
type GatewayProfileMeta store.GatewayProfileMeta

// CreateGatewayProfile creates the given gateway-profile.
// This will create the gateway-profile at the network-server side and will
// create a local reference record.
func CreateGatewayProfile(ctx context.Context, handler *store.Handler, gp *GatewayProfile) error {
	return handler.CreateGatewayProfile(ctx, (*store.GatewayProfile)(gp))
}

// GetGatewayProfile returns the gateway-profile matching the given id.
func GetGatewayProfile(ctx context.Context, handler *store.Handler, id uuid.UUID) (GatewayProfile, error) {
	res, err := handler.GetGatewayProfile(ctx, id)
	return GatewayProfile(res), err
}

// UpdateGatewayProfile updates the given gateway-profile.
func UpdateGatewayProfile(ctx context.Context, handler *store.Handler, gp *GatewayProfile) error {
	return handler.UpdateGatewayProfile(ctx, (*store.GatewayProfile)(gp))
}

// DeleteGatewayProfile deletes the gateway-profile matching the given id.
func DeleteGatewayProfile(ctx context.Context, handler *store.Handler, id uuid.UUID) error {
	return handler.DeleteGatewayProfile(ctx, id)
}

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
