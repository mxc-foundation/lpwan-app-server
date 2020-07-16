package organization

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type OrganizationStore interface {
	CreateOrganization(ctx context.Context, org *Organization) error
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (Organization, error)
	GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error)
	GetOrganizations(ctx context.Context, filters OrganizationFilters) ([]Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id int64) error
	CreateOrganizationUser(ctx context.Context, organizationID int64, username string, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error
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

type Controller struct {
	St        OrganizationStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}
