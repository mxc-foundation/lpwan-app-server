package store

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type DeviceProfileStore interface {
	DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error
}

func (h *Handler) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.store.DeleteAllDeviceProfilesForOrganizationID(ctx, organizationID)
}

func (h *Handler) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	return h.store.DeleteDeviceProfile(ctx, id)
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
