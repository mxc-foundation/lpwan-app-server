package pgstore

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	orgmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

func (ps *pgstore) CheckReadOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization user
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdateOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckDeleteOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckCreateOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// global admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckListOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// any active user (results are filtered by the api)
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckCreateOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckListOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization user
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckReadOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization admin
	// user itself
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.user_id = $3", "ou.user_id = u.id"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdateOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true", "$3 = $3"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckDeleteOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`
	// global admin
	// organization admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true", "$3 = $3"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// GetOrganizationIDList returns a slice of organizations id, sorted by name and
// respecting the given limit and offset.
func (ps *pgstore) GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error) {
	var orgIDList []int

	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.SelectContext(ctx, ps.db, &orgIDList, `
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
func (ps *pgstore) CreateOrganization(ctx context.Context, org *orgmod.Organization) error {
	if err := org.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()

	err := sqlx.GetContext(ctx, ps.db, &org.ID, `
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
func (ps *pgstore) GetOrganization(ctx context.Context, id int64, forUpdate bool) (orgmod.Organization, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var org orgmod.Organization
	err := sqlx.GetContext(ctx, ps.db, &org, "select * from organization where id = $1"+fu, id)
	if err != nil {
		return org, errors.Wrap(err, "select error")
	}
	return org, nil
}

// GetOrganizationCount returns the total number of organizations.
func (ps *pgstore) GetOrganizationCount(ctx context.Context, filters orgmod.OrganizationFilters) (int, error) {
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
	err = sqlx.GetContext(ctx, ps.db, &count, query, args...)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return count, nil
}

// GetOrganizations returns a slice of organizations, sorted by name.
func (ps *pgstore) GetOrganizations(ctx context.Context, filters orgmod.OrganizationFilters) ([]orgmod.Organization, error) {
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
	err = sqlx.SelectContext(ctx, ps.db, &orgs, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return orgs, nil
}

// UpdateOrganization updates the given organization.
func (ps *pgstore) UpdateOrganization(ctx context.Context, org *orgmod.Organization) error {
	if err := org.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	now := time.Now()
	res, err := ps.db.ExecContext(ctx, `
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
func (ps *pgstore) DeleteOrganization(ctx context.Context, id int64) error {
	res, err := ps.db.ExecContext(ctx, "delete from organization where id = $1", id)
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
func (ps *pgstore) CreateOrganizationUser(ctx context.Context, organizationID int64, username string, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	var userID int64
	err := ps.db.QueryRowContext(ctx, `
		select id from user where username = $1;
	`, username).Scan(&userID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	_, err = ps.db.ExecContext(ctx, `
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
func (ps *pgstore) UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error {
	res, err := ps.db.ExecContext(ctx, `
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
func (ps *pgstore) DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error {
	res, err := ps.db.ExecContext(ctx, `delete from organization_user where organization_id = $1 and user_id = $2`, organizationID, userID)
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
func (ps *pgstore) GetOrganizationUser(ctx context.Context, organizationID, userID int64) (orgmod.OrganizationUser, error) {
	var u orgmod.OrganizationUser
	err := sqlx.GetContext(ctx, ps.db, &u, `
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
func (ps *pgstore) GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
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
func (ps *pgstore) GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]orgmod.OrganizationUser, error) {
	var users []orgmod.OrganizationUser
	err := sqlx.SelectContext(ctx, ps.db, &users, `
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
