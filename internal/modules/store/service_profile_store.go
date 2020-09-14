package store

import (
	"context"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type ServiceProfileStore interface {
	CreateServiceProfile(ctx context.Context, sp *ServiceProfile) error
	GetServiceProfile(ctx context.Context, id uuid.UUID, localOnly bool) (ServiceProfile, error)
	UpdateServiceProfile(ctx context.Context, sp *ServiceProfile) error
	DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteServiceProfile(ctx context.Context, id uuid.UUID) error
	GetServiceProfileCount(ctx context.Context) (int, error)
	GetServiceProfileCountForOrganizationID(ctx context.Context, organizationID int64) (int, error)
	GetServiceProfileCountForUser(ctx context.Context, userID int64) (int, error)
	GetServiceProfiles(ctx context.Context, limit, offset int) ([]ServiceProfileMeta, error)
	GetServiceProfilesForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]ServiceProfileMeta, error)
	GetServiceProfilesForUser(ctx context.Context, userID int64, limit, offset int) ([]ServiceProfileMeta, error)

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
func (h *Handler) CreateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	return h.store.CreateServiceProfile(ctx, sp)
}
func (h *Handler) GetServiceProfile(ctx context.Context, id uuid.UUID, localOnly bool) (ServiceProfile, error) {
	return h.store.GetServiceProfile(ctx, id, localOnly)
}
func (h *Handler) UpdateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	return h.store.UpdateServiceProfile(ctx, sp)
}
func (h *Handler) GetServiceProfileCount(ctx context.Context) (int, error) {
	return h.store.GetServiceProfileCount(ctx)
}
func (h *Handler) GetServiceProfileCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	return h.store.GetServiceProfileCountForOrganizationID(ctx, organizationID)
}
func (h *Handler) GetServiceProfileCountForUser(ctx context.Context, userID int64) (int, error) {
	return h.store.GetServiceProfileCountForUser(ctx, userID)
}
func (h *Handler) GetServiceProfiles(ctx context.Context, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.store.GetServiceProfiles(ctx, limit, offset)
}
func (h *Handler) GetServiceProfilesForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.store.GetServiceProfilesForOrganizationID(ctx, organizationID, limit, offset)
}
func (h *Handler) GetServiceProfilesForUser(ctx context.Context, userID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.store.GetServiceProfilesForUser(ctx, userID, limit, offset)
}

// ServiceProfile defines the service-profile.
type ServiceProfile struct {
	NetworkServerID int64             `db:"network_server_id"`
	OrganizationID  int64             `db:"organization_id"`
	CreatedAt       time.Time         `db:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at"`
	Name            string            `db:"name"`
	ServiceProfile  ns.ServiceProfile `db:"-"`
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

// Validate validates the service-profile data.
func (sp ServiceProfile) Validate() error {
	if strings.TrimSpace(sp.Name) == "" || len(sp.Name) > 100 {
		return ErrServiceProfileInvalidName
	}
	return nil
}
