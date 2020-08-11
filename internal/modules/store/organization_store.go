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

func (h *Handler) CreateOrganization(ctx context.Context, org *Organization) error {
	return h.store.CreateOrganization(ctx, org)
}
func (h *Handler) GetOrganization(ctx context.Context, id int64, forUpdate bool) (Organization, error) {
	return h.store.GetOrganization(ctx, id, forUpdate)
}
func (h *Handler) GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error) {
	return h.store.GetOrganizationCount(ctx, filters)
}
func (h *Handler) GetOrganizations(ctx context.Context, filters OrganizationFilters) ([]Organization, error) {
	return h.store.GetOrganizations(ctx, filters)
}
func (h *Handler) UpdateOrganization(ctx context.Context, org *Organization) error {
	return h.store.UpdateOrganization(ctx, org)
}
func (h *Handler) DeleteOrganization(ctx context.Context, id int64) error {
	return h.store.DeleteOrganization(ctx, id)
}
func (h *Handler) CreateOrganizationUser(ctx context.Context, organizationID int64, username string, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return h.store.CreateOrganizationUser(ctx, organizationID, username, isAdmin, isDeviceAdmin, isGatewayAdmin)
}
func (h *Handler) UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	return h.store.UpdateOrganizationUser(ctx, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)
}
func (h *Handler) DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error {
	return h.store.DeleteOrganizationUser(ctx, organizationID, userID)
}
func (h *Handler) GetOrganizationUser(ctx context.Context, organizationID, userID int64) (OrganizationUser, error) {
	return h.store.GetOrganizationUser(ctx, organizationID, userID)
}
func (h *Handler) GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error) {
	return h.store.GetOrganizationUserCount(ctx, organizationID)
}
func (h *Handler) GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]OrganizationUser, error) {
	return h.store.GetOrganizationUsers(ctx, organizationID, limit, offset)
}
func (h *Handler) GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error) {
	return h.store.GetOrganizationIDList(ctx, limit, offset, search)
}

// validator
func (h *Handler) CheckReadOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.store.CheckReadOrganizationAccess(ctx, username, userID, organizationID)
}
func (h *Handler) CheckUpdateOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.store.CheckUpdateOrganizationAccess(ctx, username, userID, organizationID)
}
func (h *Handler) CheckDeleteOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.store.CheckDeleteOrganizationAccess(ctx, username, userID, organizationID)
}

func (h *Handler) CheckCreateOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.store.CheckCreateOrganizationAccess(ctx, username, userID)
}
func (h *Handler) CheckListOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	return h.store.CheckListOrganizationAccess(ctx, username, userID)
}

func (h *Handler) CheckCreateOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.store.CheckCreateOrganizationUserAccess(ctx, username, userID, organizationID)
}
func (h *Handler) CheckListOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	return h.store.CheckListOrganizationUserAccess(ctx, username, userID, organizationID)
}

func (h *Handler) CheckReadOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckReadOrganizationUserAccess(ctx, username, organizationID, userID, operatorUserID)
}
func (h *Handler) CheckUpdateOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckUpdateOrganizationUserAccess(ctx, username, organizationID, userID, operatorUserID)
}
func (h *Handler) CheckDeleteOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	return h.store.CheckDeleteOrganizationUserAccess(ctx, username, organizationID, userID, organizationID)
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
