package store

import (
	"context"

	"github.com/gofrs/uuid"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pg pgstore.PgStore) *gwps {
	return &gwps{
		pg: pg,
	}
}

type gwps struct {
	pg pgstore.GatewayProfilePgStore
}

type GatewayProfileStore interface {
	CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) (uuid.UUID, error)
	GetGatewayProfile(ctx context.Context, id uuid.UUID) (GatewayProfile, error)
	UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error
	DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error
	GetGatewayProfileCount(ctx context.Context) (int, error)
	GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error)
	GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error)
	GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error)

	// validator
	CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error)
	CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error)
}

func (h *gwps) CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) (uuid.UUID, error) {
	return h.pg.CreateGatewayProfile(ctx, gp)
}
func (h *gwps) GetGatewayProfile(ctx context.Context, id uuid.UUID) (GatewayProfile, error) {
	return h.pg.GetGatewayProfile(ctx, id)
}
func (h *gwps) UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	return h.pg.UpdateGatewayProfile(ctx, gp)
}
func (h *gwps) DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error {
	return h.pg.DeleteGatewayProfile(ctx, id)
}
func (h *gwps) GetGatewayProfileCount(ctx context.Context) (int, error) {
	return h.pg.GetGatewayProfileCount(ctx)
}

// GetGatewayProfileCountForNetworkServerID returns the total number of
// gateway-profiles given a network-server ID.
func (h *gwps) GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error) {
	return h.pg.GetGatewayProfileCountForNetworkServerID(ctx, networkServerID)
}
func (h *gwps) GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error) {
	return h.pg.GetGatewayProfiles(ctx, limit, offset)
}
func (h *gwps) GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error) {
	return h.pg.GetGatewayProfilesForNetworkServerID(ctx, networkServerID, limit, offset)
}

// validator
func (h *gwps) CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.pg.CheckCreateUpdateDeleteGatewayProfileAccess(ctx, username, userID)
}
func (h *gwps) CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.pg.CheckReadListGatewayProfileAccess(ctx, username, userID)
}
