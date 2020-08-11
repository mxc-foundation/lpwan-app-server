package store

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type ServiceProfileStore interface {
	DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteServiceProfile(ctx context.Context, id uuid.UUID) error
}

func (h *Handler) DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.store.DeleteAllServiceProfilesForOrganizationID(ctx, organizationID)
}

func (h *Handler) DeleteServiceProfile(ctx context.Context, id uuid.UUID) error {
	return h.store.DeleteServiceProfile(ctx, id)
}

// ServiceProfileMeta defines the service-profile meta record.
type ServiceProfileMeta struct {
	ServiceProfileID  uuid.UUID `db:"service_profile_id"`
	NetworkServerID   int64     `db:"network_server_id"`
	OrganizationID    int64     `db:"organization_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Name              string    `db:"name"`
	NetworkServerName string    `db:"network_server_name"`
}
