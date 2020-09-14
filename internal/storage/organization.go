package storage

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Organization represents an organization.
type Organization store.Organization

// Validate validates the data of the Organization.
func (o Organization) Validate() error {
	return store.Organization(o).Validate()
}

// OrganizationUser represents an organization user.
type OrganizationUser store.OrganizationUser

// CreateOrganization creates the given Organization.
func CreateOrganization(ctx context.Context, handler *store.Handler, org *Organization) error {
	return handler.CreateOrganization(ctx, (*store.Organization)(org))
}

// GetOrganization returns the Organization for the given id.
// When forUpdate is set to true, then db must be a db transaction.
func GetOrganization(ctx context.Context, handler *store.Handler, id int64, forUpdate bool) (Organization, error) {
	res, err := handler.GetOrganization(ctx, id, forUpdate)
	return Organization(res), err
}

// OrganizationFilters provides filters for filtering organizations.
type OrganizationFilters store.OrganizationFilters

// SQL returns the SQL filters.
func (f OrganizationFilters) SQL() string {
	return store.OrganizationFilters(f).SQL()
}

// GetOrganizationCount returns the total number of organizations.
func GetOrganizationCount(ctx context.Context, handler *store.Handler, filters OrganizationFilters) (int, error) {
	return handler.GetOrganizationCount(ctx, store.OrganizationFilters(filters))
}

// GetOrganizations returns a slice of organizations, sorted by name.
func GetOrganizations(ctx context.Context, handler *store.Handler, filters OrganizationFilters) ([]Organization, error) {
	res, err := handler.GetOrganizations(ctx, store.OrganizationFilters(filters))
	if err != nil {
		return nil, err
	}

	var orgList []Organization
	for _, v := range res {
		orgItem := Organization(v)
		orgList = append(orgList, orgItem)
	}
	return orgList, nil
}

// UpdateOrganization updates the given organization.
func UpdateOrganization(ctx context.Context, handler *store.Handler, org *Organization) error {
	return handler.UpdateOrganization(ctx, (*store.Organization)(org))
}

// DeleteOrganization deletes the organization matching the given id.
func DeleteOrganization(ctx context.Context, handler *store.Handler, id int64) error {
	return handler.DeleteOrganization(ctx, id)
}

// CreateOrganizationUser adds the given user to the organization.
func CreateOrganizationUser(ctx context.Context, handler *store.Handler, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return handler.CreateOrganizationUser(ctx, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)

}

// UpdateOrganizationUser updates the given user of the organization.
func UpdateOrganizationUser(ctx context.Context, handler *store.Handler, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return handler.UpdateOrganizationUser(ctx, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)
}

// DeleteOrganizationUser deletes the given organization user.
func DeleteOrganizationUser(ctx context.Context, handler *store.Handler, organizationID, userID int64) error {
	return handler.DeleteOrganizationUser(ctx, organizationID, userID)
}

// GetOrganizationUser gets the information of the given organization user.
func GetOrganizationUser(ctx context.Context, handler *store.Handler, organizationID, userID int64) (OrganizationUser, error) {
	res, err := handler.GetOrganizationUser(ctx, organizationID, userID)
	return OrganizationUser(res), err

}

// GetOrganizationUserCount returns the number of users for the given organization.
func GetOrganizationUserCount(ctx context.Context, handler *store.Handler, organizationID int64) (int, error) {
	return handler.GetOrganizationUserCount(ctx, organizationID)
}

// GetOrganizationUsers returns the users for the given organization.
func GetOrganizationUsers(ctx context.Context, handler *store.Handler, organizationID int64, limit, offset int) ([]OrganizationUser, error) {
	res, err := handler.GetOrganizationUsers(ctx, organizationID, limit, offset)
	if err != nil {
		return nil, err
	}

	var orgUserList []OrganizationUser
	for _, v := range res {
		orgUserItem := OrganizationUser(v)
		orgUserList = append(orgUserList, orgUserItem)
	}

	return orgUserList, nil

}
