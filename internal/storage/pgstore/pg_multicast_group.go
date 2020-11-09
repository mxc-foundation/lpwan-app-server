package pgstore

import (
	"context"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
)

func (ps *PgStore) CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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
	userWhere := [][]string{
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

func (ps *PgStore) CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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
	userWhere := [][]string{
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

func (ps *PgStore) CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
		left join multicast_group .
			on sp.service_profile_id = ..service_profile_id
	`
	// global admin
	// organization users
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "..id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, multicastGroupID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
		left join multicast_group .
			on sp.service_profile_id = ..service_profile_id
	`
	// global admin
	// organization admin users
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "..id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, multicastGroupID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
		left join multicast_group .
			on sp.service_profile_id = ..service_profile_id
	`
	// global admin
	// organization user
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "..id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, multicastGroupID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil

}

// CreateMulticastGroup creates the given multicast-group.
func (ps *PgStore) CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	if err := mg.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	mgID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "new uuid v4 error")
	}

	now := time.Now()
	mg.MulticastGroup.Id = mgID.Bytes()
	mg.CreatedAt = now
	mg.UpdatedAt = now

	_, err = ps.db.ExecContext(ctx, `
		insert into multicast_group (
			id,
			created_at,
			updated_at,
			name,
			service_profile_id,
			mc_app_s_key,
			mc_key
		) values ($1, $2, $3, $4, $5, $6, $7)
	`,
		mgID,
		mg.CreatedAt,
		mg.UpdatedAt,
		mg.Name,
		mg.ServiceProfileID,
		mg.MCAppSKey,
		mg.MCKey,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"id":     mgID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("multicast-group created")

	return nil
}

// GetMulticastGroup returns the multicast-group given an id.
func (ps *PgStore) GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate bool) (MulticastGroup, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var mg MulticastGroup

	err := sqlx.GetContext(ctx, ps.db, &mg, `
		select
			created_at,
			updated_at,
			name,
			service_profile_id,
			mc_app_s_key,
			mc_key
		from
			multicast_group
		where
			id = $1
	`+fu, id)
	if err != nil {
		return mg, handlePSQLError(Select, err, "select error")
	}

	return mg, nil
}

// UpdateMulticastGroup updates the given multicast-group.
func (ps *PgStore) UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	if err := mg.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	mgID, err := uuid.FromBytes(mg.MulticastGroup.Id)
	if err != nil {
		return errors.Wrap(err, "uuid from bytes error")
	}

	mg.UpdatedAt = time.Now()
	res, err := ps.db.ExecContext(ctx, `
		update
			multicast_group
		set
			updated_at = $2,
			name = $3,
			mc_app_s_key = $4,
			mc_key = $5
		where
			id = $1
	`,
		mgID,
		mg.UpdatedAt,
		mg.Name,
		mg.MCAppSKey,
		mg.MCKey,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"id":     mgID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("multicast-group updated")

	return nil
}

// DeleteMulticastGroup deletes a multicast-group given an id.
func (ps *PgStore) DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error {
	res, err := ps.db.ExecContext(ctx, `
		delete
		from
			multicast_group
		where
			id = $1
	`, id)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("multicast-group deleted")

	return nil
}

// GetMulticastGroupCount returns the total number of multicast-groups given
// the provided filters. Note that empty values are not used as filters.
func (ps *PgStore) GetMulticastGroupCount(ctx context.Context, filters MulticastGroupFilters) (int, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct ..*)
		from
			multicast_group .
		inner join service_profile sp
			on sp.service_profile_id = ..service_profile_id
		inner join organization o
			on o.id = sp.organization_id
		left join device_multicast_group dmg
			on ..id = dmg.multicast_group_id
	`+filters.SQL(), filters)
	if err != nil {
		return 0, errors.Wrap(err, "named query error")
	}

	var count int
	err = sqlx.GetContext(ctx, ps.db, &count, query, args...)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetMulticastGroups returns a slice of multicast-groups, given the privded
// filters. Note that empty values are not used as filters.
func (ps *PgStore) GetMulticastGroups(ctx context.Context, filters MulticastGroupFilters) ([]MulticastGroupListItem, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			distinct ..id,
			..created_at,
			..updated_at,
			..name,
			..service_profile_id,
			sp.name as service_profile_name
		from
			multicast_group .
		inner join service_profile sp
			on sp.service_profile_id = ..service_profile_id
		inner join organization o
			on o.id = sp.organization_id
		left join device_multicast_group dmg
			on ..id = dmg.multicast_group_id
	`+filters.SQL()+`
		order by
			..name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var mgs []MulticastGroupListItem
	err = sqlx.SelectContext(ctx, ps.db, &mgs, query, args...)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return mgs, nil
}

// AddDeviceToMulticastGroup adds the given device to the given multicast-group.
// It is recommended that db is a transaction.
func (ps *PgStore) AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	_, err := ps.db.ExecContext(ctx, `
		insert into device_multicast_group (
			dev_eui,
			multicast_group_id,
			created_at
		) values ($1, $2, $3)
	`, devEUI, multicastGroupID, time.Now())
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"multicast_group_id": multicastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("device added to multicast-group")

	return nil
}

// RemoveDeviceFromMulticastGroup removes the given device from the given
// multicast-group.
func (ps *PgStore) RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from
			device_multicast_group
		where
			dev_eui = $1
			and multicast_group_id = $2
	`, devEUI, multicastGroupID)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"multicast_group_id": multicastGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("Device removed from multicast-group")

	return nil
}

// GetDeviceCountForMulticastGroup returns the number of devices for the given
// multicast-group.
func (ps *PgStore) GetDeviceCountForMulticastGroup(ctx context.Context, multicastGroup uuid.UUID) (int, error) {
	var count int

	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(*)
		from
			device_multicast_group
		where
			multicast_group_id = $1
	`, multicastGroup)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetDevicesForMulticastGroup returns a slice of devices for the given
// multicast-group.
func (ps *PgStore) GetDevicesForMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, limit, offset int) ([]ds.DeviceListItem, error) {
	var devices []ds.DeviceListItem

	err := sqlx.SelectContext(ctx, ps.db, &devices, `
		select
			d.*,
			dp.name as device_profile_name
		from
			device d
		inner join device_profile dp
			on dp.device_profile_id = d.device_profile_id
		inner join device_multicast_group dmg
			on dmg.dev_eui = d.dev_eui
		where
			dmg.multicast_group_id = $1
		order by
			d.name
		limit $2
		offset $3
	`, multicastGroupID, limit, offset)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return devices, nil
}
