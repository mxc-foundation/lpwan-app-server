package pgstore

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	orgmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"

	applicationPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/pgstore"
)

type OrgHandler struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func New(tx *sqlx.Tx, db *sqlx.DB) *OrgHandler {
	return &OrgHandler{
		tx: tx,
		db: db,
	}
}

// GetOrganizationIDList returns a slice of organizations id, sorted by name and
// respecting the given limit and offset.
func (h *OrgHandler) GetOrganizationIDList(limit, offset int, search string) ([]int, error) {
	var orgIDList []int

	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Select(h.db, &orgIDList, `
		select id
		from organization
		where
			($3 != '' and display_name ilike $3)
			or ($3 = '')
		order by display_name
		limit $1 offset $2`, limit, offset, search)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return orgIDList, nil
}

// CreateOrganization creates the given Organization.
func (h *OrgHandler) CreateOrganization(ctx context.Context, org *orgmod.Organization) error {
	if err := org.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()

	err := sqlx.Get(h.tx, &org.ID, `
		insert into organization (
			created_at,
			updated_at,
			name,
			display_name,
			can_have_gateways,
			max_gateway_count,
			max_device_count
		) values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		now,
		now,
		org.Name,
		org.DisplayName,
		org.CanHaveGateways,
		org.MaxGatewayCount,
		org.MaxDeviceCount,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}
	org.CreatedAt = now
	org.UpdatedAt = now
	log.WithFields(log.Fields{
		"id":     org.ID,
		"name":   org.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("organization created")
	return nil
}

// GetOrganization returns the orgmod.Organization for the given id.
// When forUpdate is set to true, then tx must be a tx transaction.
func (h *OrgHandler) GetOrganization(ctx context.Context, id int64, forUpdate bool) (orgmod.Organization, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var org orgmod.Organization
	err := sqlx.Get(h.db, &org, "select * from organization where id = $1"+fu, id)
	if err != nil {
		return org, errors.Wrap(err, "select error")
	}
	return org, nil
}

// GetOrganizationCount returns the total number of organizations.
func (h *OrgHandler) GetOrganizationCount(ctx context.Context, filters orgmod.OrganizationFilters) (int, error) {
	query, args, err := sqlx.Named(`
		select
			count(distinct o.*)
		from
			organization o
		left join organization_user ou
			on o.id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL(), filters)
	if err != nil {
		return 0, errors.Wrap(err, "named query error")
	}

	var count int
	err = sqlx.Get(h.db, &count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return count, nil
}

// GetOrganizations returns a slice of organizations, sorted by name.
func (h *OrgHandler) GetOrganizations(ctx context.Context, filters orgmod.OrganizationFilters) ([]orgmod.Organization, error) {
	query, args, err := sqlx.Named(`
		select
			o.*
		from
			organization o
		left join organization_user ou
			on o.id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL()+`
		group by
			o.id
		order by
			o.display_name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var orgs []orgmod.Organization
	err = sqlx.Select(h.db, &orgs, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return orgs, nil
}

// UpdateOrganization updates the given organization.
func (h *OrgHandler) UpdateOrganization(ctx context.Context, org *orgmod.Organization) error {
	if err := org.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	now := time.Now()
	res, err := h.tx.Exec(`
		update organization
		set
			name = $2,
			display_name = $3,
			can_have_gateways = $4,
			updated_at = $5,
			max_gateway_count = $6,
			max_device_count = $7
		where id = $1`,
		org.ID,
		org.Name,
		org.DisplayName,
		org.CanHaveGateways,
		now,
		org.MaxGatewayCount,
		org.MaxDeviceCount,
	)

	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	org.UpdatedAt = now
	log.WithFields(log.Fields{
		"name":   org.Name,
		"id":     org.ID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("organization updated")
	return nil
}

// DeleteOrganization deletes the organization matching the given id.
func (h *OrgHandler) DeleteOrganization(ctx context.Context, id int64) error {
	err := applicationPg.Handler().DeleteAllApplicationsForOrganizationID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete all applications error")
	}

	err = storage.DeleteAllServiceProfilesForOrganizationID(ctx, h.tx, id)
	if err != nil {
		return errors.Wrap(err, "delete all service-profiles error")
	}

	err = storage.DeleteAllDeviceProfilesForOrganizationID(ctx, h.tx, id)
	if err != nil {
		return errors.Wrap(err, "delete all device-profiles error")
	}

	res, err := h.tx.Exec("delete from organization where id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("organization deleted")
	return nil
}

// CreateOrganizationUser adds the given user to the organization.
func (h *OrgHandler) CreateOrganizationUser(ctx context.Context, organizationID int64, username string, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	var userID int64
	err := h.db.QueryRow(`
		select id from user where username = $1;
	`, username).Scan(&userID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	_, err = h.tx.Exec(`
		insert into organization_user (
			organization_id,
			user_id,
			is_admin,
			is_device_admin,
			is_gateway_admin,
			created_at,
			updated_at
		) values ($1, $2, $3, $4, $5, now(), now())`,
		organizationID,
		userID,
		isAdmin,
		isDeviceAdmin,
		isGatewayAdmin,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	log.WithFields(log.Fields{
		"username":         username,
		"user id":          userID,
		"organization_id":  organizationID,
		"is_admin":         isAdmin,
		"is_device_admin":  isDeviceAdmin,
		"is_gateway_admin": isGatewayAdmin,
		"ctx_id":           ctx.Value(logging.ContextIDKey),
	}).Info("user added to organization")
	return nil
}

// UpdateOrganizationUser updates the given user of the organization.
func (h *OrgHandler) UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	res, err := h.tx.Exec(`
		update organization_user
		set
			is_admin = $3,
			is_device_admin = $4,
			is_gateway_admin = $5,
			updated_at = now()
		where
			organization_id = $1
			and user_id = $2
	`, organizationID, userID, isAdmin, isDeviceAdmin, isGatewayAdmin)
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"user_id":          userID,
		"organization_id":  organizationID,
		"is_admin":         isAdmin,
		"is_device_admin":  isDeviceAdmin,
		"is_gateway_admin": isGatewayAdmin,
		"ctx_id":           ctx.Value(logging.ContextIDKey),
	}).Info("organization user updated")
	return nil
}

// DeleteOrganizationUser deletes the given organization user.
func (h *OrgHandler) DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error {
	res, err := h.tx.Exec(`delete from organization_user where organization_id = $1 and user_id = $2`, organizationID, userID)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	log.WithFields(log.Fields{
		"user_id":         userID,
		"organization_id": organizationID,
		"ctx_id":          ctx.Value(logging.ContextIDKey),
	}).Info("organization user deleted")
	return nil
}

// GetOrganizationUser gets the information of the given organization user.
func (h *OrgHandler) GetOrganizationUser(ctx context.Context, organizationID, userID int64) (orgmod.OrganizationUser, error) {
	var u orgmod.OrganizationUser
	err := sqlx.Get(h.db, &u, `
		select
			u.id as user_id,
			u.email as email,
			ou.created_at as created_at,
			ou.updated_at as updated_at,
			ou.is_admin as is_admin,
			ou.is_device_admin as is_device_admin,
			ou.is_gateway_admin as is_gateway_admin
		from organization_user ou
		inner join "user" u
			on u.id = ou.user_id
		where
			ou.organization_id = $1
			and ou.user_id = $2`,
		organizationID,
		userID,
	)
	if err != nil {
		return u, errors.Wrap(err, "select error")
	}
	return u, nil
}

// GetOrganizationUserCount returns the number of users for the given organization.
func (h *OrgHandler) GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error) {
	var count int
	err := sqlx.Get(h.db, &count, `
		select count(*)
		from organization_user
		where
			organization_id = $1`,
		organizationID,
	)
	if err != nil {
		return count, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetOrganizationUsers returns the users for the given organization.
func (h *OrgHandler) GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]orgmod.OrganizationUser, error) {
	var users []orgmod.OrganizationUser
	err := sqlx.Select(h.db, &users, `
		select
			u.id as user_id,
			u.email as email,
			ou.created_at as created_at,
			ou.updated_at as updated_at,
			ou.is_admin as is_admin,
			ou.is_device_admin as is_device_admin,
			ou.is_gateway_admin as is_gateway_admin
		from organization_user ou
		inner join "user" u
			on u.id = ou.user_id
		where
			ou.organization_id = $1
		order by u.email
		limit $2 offset $3`,
		organizationID,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return users, nil
}
