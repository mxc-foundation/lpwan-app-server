package store

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
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

var organizationNameRegexp = regexp.MustCompile(`^[\w-]+$`)

// Organization represents an organization.
type Organization struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Name            string    `db:"name"`
	DisplayName     string    `db:"display_name"`
	CanHaveGateways bool      `db:"can_have_gateways"`
	MaxDeviceCount  int       `db:"max_device_count"`
	MaxGatewayCount int       `db:"max_gateway_count"`
}

// Validate validates the data of the Organization.
func (o Organization) Validate() error {
	if !organizationNameRegexp.MatchString(o.Name) {
		return errors.New("ErrOrganizationInvalidName")
	}
	return nil
}

// OrganizationUser represents an organization user.
type OrganizationUser struct {
	UserID         int64     `db:"user_id"`
	Email          string    `db:"email"`
	IsAdmin        bool      `db:"is_admin"`
	IsDeviceAdmin  bool      `db:"is_device_admin"`
	IsGatewayAdmin bool      `db:"is_gateway_admin"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// OrganizationFilters provides filters for filtering organizations.
type OrganizationFilters struct {
	UserID int64  `db:"user_id"`
	Search string `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f OrganizationFilters) SQL() string {
	var filters []string

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if f.Search != "" {
		filters = append(filters, "o.display_name ilike :search")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}
