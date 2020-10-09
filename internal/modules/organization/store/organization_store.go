package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/data"
)

func NewStore(pg pgstore.PgStore) *orgs {
	return &orgs{
		pg: pg,
	}
}

type orgs struct {
	pg pgstore.OrganizationPgStore
}

type OrganizationStore interface {
	CreateOrganization(ctx context.Context, org *Organization) error
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (Organization, error)
	GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error)
	GetOrganizations(ctx context.Context, filters OrganizationFilters) ([]Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id int64) error
	CreateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error
	UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error
	DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error
	GetOrganizationUser(ctx context.Context, organizationID, userID int64) (OrganizationUser, error)
	GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error)
	GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]OrganizationUser, error)
	GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error)

	// validator
	CheckReadOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckUpdateOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckDeleteOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)

	CheckCreateOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error)
	CheckListOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error)

	CheckCreateOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckListOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)

	CheckReadOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error)
	CheckUpdateOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error)
	CheckDeleteOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error)
}

func (h *orgs) CreateOrganization(ctx context.Context, org *Organization) error {
	return h.pg.CreateOrganization(ctx, org)
}
func (h *orgs) GetOrganization(ctx context.Context, id int64, forUpdate bool) (Organization, error) {
	return h.pg.GetOrganization(ctx, id, forUpdate)
}
func (h *orgs) GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error) {
	return h.pg.GetOrganizationCount(ctx, filters)
}
func (h *orgs) GetOrganizations(ctx context.Context, filters OrganizationFilters) ([]Organization, error) {
	return h.pg.GetOrganizations(ctx, filters)
}
func (h *orgs) UpdateOrganization(ctx context.Context, org *Organization) error {
	return h.pg.UpdateOrganization(ctx, org)
}
func (h *orgs) DeleteOrganization(ctx context.Context, id int64) error {
	return h.pg.DeleteOrganization(ctx, id)
}
func (h *orgs) CreateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return h.pg.CreateOrganizationUser(ctx, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)
}
func (h *orgs) UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return h.pg.UpdateOrganizationUser(ctx, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)
}
func (h *orgs) DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error {
	return h.pg.DeleteOrganizationUser(ctx, organizationID, userID)
}
func (h *orgs) GetOrganizationUser(ctx context.Context, organizationID, userID int64) (OrganizationUser, error) {
	return h.pg.GetOrganizationUser(ctx, organizationID, userID)
}
func (h *orgs) GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error) {
	return h.pg.GetOrganizationUserCount(ctx, organizationID)
}
func (h *orgs) GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]OrganizationUser, error) {
	return h.pg.GetOrganizationUsers(ctx, organizationID, limit, offset)
}
func (h *orgs) GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error) {
	return h.pg.GetOrganizationIDList(ctx, limit, offset, search)
}

// validator
func (h *orgs) CheckReadOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.pg.CheckReadOrganizationAccess(ctx, username, userID, organizationID)
}
func (h *orgs) CheckUpdateOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.pg.CheckUpdateOrganizationAccess(ctx, username, userID, organizationID)
}
func (h *orgs) CheckDeleteOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.pg.CheckDeleteOrganizationAccess(ctx, username, userID, organizationID)
}

func (h *orgs) CheckCreateOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.pg.CheckCreateOrganizationAccess(ctx, username, userID)
}
func (h *orgs) CheckListOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.pg.CheckListOrganizationAccess(ctx, username, userID)
}

func (h *orgs) CheckCreateOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.pg.CheckCreateOrganizationUserAccess(ctx, username, userID, organizationID)
}
func (h *orgs) CheckListOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.pg.CheckListOrganizationUserAccess(ctx, username, userID, organizationID)
}

func (h *orgs) CheckReadOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckReadOrganizationUserAccess(ctx, username, organizationID, userID, operatorUserID)
}
func (h *orgs) CheckUpdateOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckUpdateOrganizationUserAccess(ctx, username, organizationID, userID, operatorUserID)
}
func (h *orgs) CheckDeleteOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.pg.CheckDeleteOrganizationUserAccess(ctx, username, organizationID, userID, organizationID)
}
