package store

import (
	"context"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/lib/pq/hstore"

	"github.com/gofrs/uuid"
)

type DeviceProfileStore interface {
	DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error
	GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (DeviceProfile, error)

	// validator
	CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error)
	CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error)

	CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
}

func (h *Handler) CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadDeviceProfileAccess(ctx, username, id, userID)
}

func (h *Handler) CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteDeviceProfileAccess(ctx, username, id, userID)
}

func (h *Handler) CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	return h.store.CheckCreateDeviceProfilesAccess(ctx, username, organizationID, applicationID, userID)
}

func (h *Handler) CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	return h.store.CheckListDeviceProfilesAccess(ctx, username, organizationID, applicationID, userID)
}

func (h *Handler) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.store.DeleteAllDeviceProfilesForOrganizationID(ctx, organizationID)
}

func (h *Handler) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	return h.store.DeleteDeviceProfile(ctx, id)
}

func (h *Handler) GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (DeviceProfile, error) {
	return h.store.GetDeviceProfile(ctx, id, forUpdate)
}

// DeviceProfileMeta defines the device-profile meta record.
type DeviceProfileMeta struct {
	DeviceProfileID   uuid.UUID `db:"device_profile_id"`
	NetworkServerID   int64     `db:"network_server_id"`
	OrganizationID    int64     `db:"organization_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Name              string    `db:"name"`
	NetworkServerName string    `db:"network_server_name"`
}

// DeviceProfile defines the device-profile.
type DeviceProfile struct {
	NetworkServerID      int64            `db:"network_server_id"`
	OrganizationID       int64            `db:"organization_id"`
	CreatedAt            time.Time        `db:"created_at"`
	UpdatedAt            time.Time        `db:"updated_at"`
	Name                 string           `db:"name"`
	PayloadCodec         string           `db:"payload_codec"`
	PayloadEncoderScript string           `db:"payload_encoder_script"`
	PayloadDecoderScript string           `db:"payload_decoder_script"`
	Tags                 hstore.Hstore    `db:"tags"`
	UplinkInterval       time.Duration    `db:"uplink_interval"`
	DeviceProfile        ns.DeviceProfile `db:"-"`
}
