package store

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type ServiceProfileStore interface {
	DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteServiceProfile(ctx context.Context, id uuid.UUID) error

	// validator
	CheckCreateServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckUpdateDeleteServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckReadServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
}

func (h *Handler) CheckUpdateDeleteServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteServiceProfileAccess(ctx, username, id, userID)
}

func (h *Handler) CheckReadServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadServiceProfileAccess(ctx, username, id, userID)
}

func (h *Handler) CheckCreateServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckCreateServiceProfilesAccess(ctx, username, organizationID, userID)
}

func (h *Handler) CheckListServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckListServiceProfilesAccess(ctx, username, organizationID, userID)
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
