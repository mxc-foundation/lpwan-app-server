package store

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
)

func NewStore(pg pgstore.PgStore) *sps {
	return &sps{
		pg: pg,
	}
}

type sps struct {
	pg pgstore.ServiceProfilePgStore
}

type ServiceProfileStore interface {
	CreateServiceProfile(ctx context.Context, sp *ServiceProfile) error
	GetServiceProfile(ctx context.Context, id uuid.UUID) (ServiceProfile, error)
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

func (h *sps) CheckUpdateDeleteServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckUpdateDeleteServiceProfileAccess(ctx, username, id, userID)
}

func (h *sps) CheckReadServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckReadServiceProfileAccess(ctx, username, id, userID)
}

func (h *sps) CheckCreateServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckCreateServiceProfilesAccess(ctx, username, organizationID, userID)
}

func (h *sps) CheckListServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckListServiceProfilesAccess(ctx, username, organizationID, userID)
}

func (h *sps) DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.pg.DeleteAllServiceProfilesForOrganizationID(ctx, organizationID)
}

func (h *sps) DeleteServiceProfile(ctx context.Context, id uuid.UUID) error {
	return h.pg.DeleteServiceProfile(ctx, id)
}
func (h *sps) CreateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	return h.pg.CreateServiceProfile(ctx, sp)
}
func (h *sps) GetServiceProfile(ctx context.Context, id uuid.UUID) (ServiceProfile, error) {
	return h.pg.GetServiceProfile(ctx, id)
}
func (h *sps) UpdateServiceProfile(ctx context.Context, sp *ServiceProfile) error {
	return h.pg.UpdateServiceProfile(ctx, sp)
}
func (h *sps) GetServiceProfileCount(ctx context.Context) (int, error) {
	return h.pg.GetServiceProfileCount(ctx)
}
func (h *sps) GetServiceProfileCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	return h.pg.GetServiceProfileCountForOrganizationID(ctx, organizationID)
}
func (h *sps) GetServiceProfileCountForUser(ctx context.Context, userID int64) (int, error) {
	return h.pg.GetServiceProfileCountForUser(ctx, userID)
}
func (h *sps) GetServiceProfiles(ctx context.Context, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.pg.GetServiceProfiles(ctx, limit, offset)
}
func (h *sps) GetServiceProfilesForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.pg.GetServiceProfilesForOrganizationID(ctx, organizationID, limit, offset)
}
func (h *sps) GetServiceProfilesForUser(ctx context.Context, userID int64, limit, offset int) ([]ServiceProfileMeta, error) {
	return h.pg.GetServiceProfilesForUser(ctx, userID, limit, offset)
}
