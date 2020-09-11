package store

import (
	"context"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
)

type GatewayProfileStore interface {
	CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) error
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

func (h *Handler) CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	return h.store.CreateGatewayProfile(ctx, gp)
}
func (h *Handler) GetGatewayProfile(ctx context.Context, id uuid.UUID) (GatewayProfile, error) {
	return h.store.GetGatewayProfile(ctx, id)
}
func (h *Handler) UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	return h.store.UpdateGatewayProfile(ctx, gp)
}
func (h *Handler) DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error {
	return h.store.DeleteGatewayProfile(ctx, id)
}
func (h *Handler) GetGatewayProfileCount(ctx context.Context) (int, error) {
	return h.store.GetGatewayProfileCount(ctx)
}
func (h *Handler) GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error) {
	return h.store.GetGatewayProfileCountForNetworkServerID(ctx, networkServerID)
}
func (h *Handler) GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error) {
	return h.store.GetGatewayProfiles(ctx, limit, offset)
}
func (h *Handler) GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error) {
	return h.store.GetGatewayProfilesForNetworkServerID(ctx, networkServerID, limit, offset)
}

// validator
func (h *Handler) CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.store.CheckCreateUpdateDeleteGatewayProfileAccess(ctx, username, userID)
}
func (h *Handler) CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.store.CheckReadListGatewayProfileAccess(ctx, username, userID)
}

// Modulations
const (
	ModulationFSK  = "FSK"
	ModulationLoRa = "LORA"
)

// ExtraChannel defines an extra channel for the gateway-profile.
type ExtraChannel struct {
	Modulation       string
	Frequency        int
	Bandwidth        int
	Bitrate          int
	SpreadingFactors []int
}

// GatewayProfile defines a gateway-profile.
type GatewayProfile struct {
	NetworkServerID int64             `db:"network_server_id"`
	CreatedAt       time.Time         `db:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"`
	Name            string            `db:"name"`
	GatewayProfile  ns.GatewayProfile `db:"-"`
}

// GatewayProfileMeta defines the gateway-profile meta record.
type GatewayProfileMeta struct {
	GatewayProfileID  uuid.UUID     `db:"gateway_profile_id"`
	NetworkServerID   int64         `db:"network_server_id"`
	NetworkServerName string        `db:"network_server_name"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
	Name              string        `db:"name"`
	StatsInterval     time.Duration `db:"stats_interval"`
}
