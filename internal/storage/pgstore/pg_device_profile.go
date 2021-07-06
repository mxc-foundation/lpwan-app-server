package pgstore

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	dpapi "github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
)

func (ps *PgStore) CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
		left join device_profile dp
			on dp.organization_id = o.id
	`
	// gloabal admin
	// organization users
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "dp.device_profile_id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, id, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
		left join device_profile dp
			on dp.organization_id = o.id
	`
	// global admin
	// organization admin users
	// organization device admin users
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin=true", "dp.device_profile_id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin=true", "dp.device_profile_id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, id, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`

	// global admin
	// organization admin
	// organization device admin
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.is_admin = true", "$3 = 0"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ou.is_device_admin = true", "$3 = 0"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, applicationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`
	// global admin
	// organization user (when organization id is given)
	// user linked to a given application (when application id is given)
	// any active user (filtered by user)
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "$3 = 0", "$2 > 0", "o.id = $2"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "$2 = 0", "$3 > 0", "a.id = $3"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "$2 = 0", "$3 = 0"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, organizationID, applicationID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// DeleteAllDeviceProfilesForOrganizationID deletes all device-profiles
// given an organization id.
func (ps *PgStore) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	var dps []dpapi.DeviceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &dps, `
		select
			device_profile_id,
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name
		from
			device_profile
		where
			organization_id = $1`,
		organizationID)
	if err != nil {
		return errors.Wrap(err, "select error")
	}

	for _, dp := range dps {
		err = ps.DeleteDeviceProfile(ctx, dp.DeviceProfileID)
		if err != nil {
			return errors.Wrap(err, "delete device-profile error")
		}
	}

	return nil
}

// CreateDeviceProfile creates the given device-profile.
// This will create the device-profile at the network-server side and will
// create a local reference record.
func (ps *PgStore) CreateDeviceProfile(ctx context.Context, dp *dpapi.DeviceProfile) error {
	if err := dp.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	dpID, err := uuid.FromBytes(dp.DeviceProfile.Id)
	if err != nil {
		return err
	}

	_, err = ps.db.ExecContext(ctx, `
        insert into device_profile (
            device_profile_id,
            network_server_id,
            organization_id,
            created_at,
            updated_at,
            name,
			payload_codec,
			payload_encoder_script,
			payload_decoder_script,
			tags,
			uplink_interval
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		dpID,
		dp.NetworkServerID,
		dp.OrganizationID,
		dp.CreatedAt,
		dp.UpdatedAt,
		dp.Name,
		dp.PayloadCodec,
		dp.PayloadEncoderScript,
		dp.PayloadDecoderScript,
		dp.Tags,
		dp.UplinkInterval,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"id":     dpID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("device-profile created")

	return nil
}

// DeleteDeviceProfile deletes the device-profile matching the given id.
func (ps *PgStore) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	res, err := ps.db.ExecContext(ctx, "delete from device_profile where device_profile_id = $1", id)
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
	}).Info("device-profile deleted")

	return nil
}

// GetDeviceProfileWithIDAndOrganizationID returns the device-profile matching the given id and organization id
func (ps *PgStore) GetDeviceProfileWithIDAndOrganizationID(ctx context.Context, id uuid.UUID,
	orgID int64, forUpdate bool) (dpapi.DeviceProfile, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dp dpapi.DeviceProfile

	row := ps.db.QueryRowxContext(ctx, `
		select
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name,
			payload_codec,
			payload_encoder_script,
			payload_decoder_script,
			tags,
			uplink_interval
		from device_profile
		where
			device_profile_id = $1 and organization_id = $2`+fu,
		id,
		orgID,
	)
	if err := row.Err(); err != nil {
		return dp, handlePSQLError(Select, err, "select error")
	}

	err := row.Scan(
		&dp.NetworkServerID,
		&dp.OrganizationID,
		&dp.CreatedAt,
		&dp.UpdatedAt,
		&dp.Name,
		&dp.PayloadCodec,
		&dp.PayloadEncoderScript,
		&dp.PayloadDecoderScript,
		&dp.Tags,
		&dp.UplinkInterval,
	)
	if err != nil {
		return dp, handlePSQLError(Scan, err, "scan error")
	}

	return dp, nil
}

// GetDeviceProfile returns the device-profile matching the given id.
// When forUpdate is set to true, then db must be a db transaction.
// When localOnly is set to true, no call to the network-server is made to
// retrieve additional device data.
func (ps *PgStore) GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (dpapi.DeviceProfile, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dp dpapi.DeviceProfile

	row := ps.db.QueryRowxContext(ctx, `
		select
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name,
			payload_codec,
			payload_encoder_script,
			payload_decoder_script,
			tags,
			uplink_interval
		from device_profile
		where
			device_profile_id = $1`+fu,
		id,
	)
	if err := row.Err(); err != nil {
		return dp, handlePSQLError(Select, err, "select error")
	}

	err := row.Scan(
		&dp.NetworkServerID,
		&dp.OrganizationID,
		&dp.CreatedAt,
		&dp.UpdatedAt,
		&dp.Name,
		&dp.PayloadCodec,
		&dp.PayloadEncoderScript,
		&dp.PayloadDecoderScript,
		&dp.Tags,
		&dp.UplinkInterval,
	)
	if err != nil {
		return dp, handlePSQLError(Scan, err, "scan error")
	}

	return dp, nil
}

// GetDefaultDeviceProfileID returns the default device profile id with given organization id and network server id
func (ps *PgStore) GetDefaultDeviceProfileID(ctx context.Context, orgID, nsID int64, forUpdate bool) (*uuid.UUID, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dpID uuid.UUID

	row := ps.db.QueryRowxContext(ctx, `
		select
			device_profile_id
		from device_profile
		where
			network_server_id = $1 and organization_id = $2`+fu,
		nsID, orgID,
	)
	if err := row.Err(); err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	err := row.Scan(&dpID)
	if err != nil {
		return nil, handlePSQLError(Scan, err, "scan error")
	}

	return &dpID, nil
}

// UpdateDeviceProfile updates the given device-profile.
func (ps *PgStore) UpdateDeviceProfile(ctx context.Context, dp *dpapi.DeviceProfile) error {
	if err := dp.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	dpID, err := uuid.FromBytes(dp.DeviceProfile.Id)
	if err != nil {
		return errors.Wrap(err, "uuid from bytes error")
	}

	dp.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
        update device_profile
        set
            updated_at = $2,
            name = $3,
			payload_codec = $4,
			payload_encoder_script = $5,
			payload_decoder_script = $6,
			tags = $7,
			uplink_interval = $8
		where device_profile_id = $1`,
		dpID,
		dp.UpdatedAt,
		dp.Name,
		dp.PayloadCodec,
		dp.PayloadEncoderScript,
		dp.PayloadDecoderScript,
		dp.Tags,
		dp.UplinkInterval,
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
		"id":     dpID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("device-profile updated")

	return nil
}

// GetDeviceProfileCount returns the total number of device-profiles.
func (ps *PgStore) GetDeviceProfileCount(ctx context.Context, filters dpapi.DeviceProfileFilters) (int, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct dp.*)
		from
			device_profile dp
		inner join network_server ns
			on dp.network_server_id = ns.id
		inner join organization o
			on dp.organization_id = o.id
		left join service_profile sp
			on ns.id = sp.network_server_id
		left join application a
			on sp.service_profile_id = a.service_profile_id
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
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetDeviceProfiles returns a slice of device-profiles.
func (ps *PgStore) GetDeviceProfiles(ctx context.Context, filters dpapi.DeviceProfileFilters) ([]dpapi.DeviceProfileMeta, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			dp.device_profile_id,
			dp.network_server_id,
			dp.organization_id,
			dp.created_at,
			dp.updated_at,
			dp.name,
			ns.name as network_server_name
		from
			device_profile dp
		inner join network_server ns
			on dp.network_server_id = ns.id
		inner join organization o
			on dp.organization_id = o.id
		left join service_profile sp
			on ns.id = sp.network_server_id
		left join application a
			on sp.service_profile_id = a.service_profile_id
		left join organization_user ou
			on o.id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL()+`
		group by
			dp.device_profile_id,
			ns.name
		order by
			dp.name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var dps []dpapi.DeviceProfileMeta
	err = sqlx.SelectContext(ctx, ps.db, &dps, query, args...)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return dps, nil
}
