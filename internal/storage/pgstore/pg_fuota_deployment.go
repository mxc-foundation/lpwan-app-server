package pgstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment/data"
)

func (ps *PgStore) CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join application a
			on a.organization_id = ou.organization_id
		left join device d
			on a.id = d.application_id
		left join fuota_deployment_device fdd
			on d.dev_eui = fdd.dev_eui
		left join fuota_deployment fd
			on fdd.fuota_deployment_id = fd.id
	`
	// global admin
	// organization user
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "fd.id = $2"},
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

func (ps *PgStore) CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join application a
			on a.organization_id = ou.organization_id
		left join device d
			on a.id = d.application_id
	`
	// global admin
	// organization admin
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 > 0", "a.id = $2"},
		{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 = 0", "d.dev_eui = $3"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, applicationID, devEUI, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]ds.DeviceKeys, error) {
	// query all device-keys that relate to this FUOTA deployment
	var deviceKeys []ds.DeviceKeys
	err := sqlx.SelectContext(ctx, ps.db, &deviceKeys, `
		select
			dk.*
		from
			fuota_deployment_device dd
		inner join
			device_keys dk
			on dd.dev_eui = dk.dev_eui
		where
			dd.fuota_deployment_id = $1`,
		id,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return deviceKeys, nil
}

// CreateFUOTADeploymentForDevice creates and initializes a FUOTA deployment
// for the given device.
func (ps *PgStore) CreateFUOTADeploymentForDevice(ctx context.Context, fd *FUOTADeployment, devEUI lorawan.EUI64) error {
	if err := fd.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()
	var err error
	fd.ID, err = uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "new uuid error")
	}

	fd.CreatedAt = now
	fd.UpdatedAt = now
	fd.NextStepAfter = now
	if fd.State == "" {
		fd.State = FUOTADeploymentMulticastCreate
	}

	_, err = ps.db.ExecContext(ctx, `
		insert into fuota_deployment (
			id,
			created_at,
			updated_at,
			name,
			multicast_group_id,

			fragmentation_matrix,
			descriptor,
			payload,
			state,
			next_step_after,
			unicast_timeout,
			frag_size,
			redundancy,
			block_ack_delay,
			multicast_timeout,
			group_type,
			dr,
			frequency,
			ping_slot_period
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		fd.ID,
		fd.CreatedAt,
		fd.UpdatedAt,
		fd.Name,
		fd.MulticastGroupID,
		[]byte{fd.FragmentationMatrix},
		fd.Descriptor[:],
		fd.Payload,
		fd.State,
		fd.NextStepAfter,
		fd.UnicastTimeout,
		fd.FragSize,
		fd.Redundancy,
		fd.BlockAckDelay,
		fd.MulticastTimeout,
		fd.GroupType,
		fd.DR,
		fd.Frequency,
		fd.PingSlotPeriod,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	_, err = ps.db.ExecContext(ctx, `
		insert into fuota_deployment_device (
			fuota_deployment_id,
			dev_eui,
			created_at,
			updated_at,
			state,
			error_message
		) values ($1, $2, $3, $4, $5, $6)`,
		fd.ID,
		devEUI,
		now,
		now,
		FUOTADeploymentDevicePending,
		"",
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"dev_eui": devEUI,
		"id":      fd.ID,
		"ctx_id":  ctx.Value(logging.ContextIDKey),
	}).Info("fuota deploymented created for device")

	return nil
}

// GetFUOTADeployment returns the FUOTA deployment for the given ID.
func (ps *PgStore) GetFUOTADeployment(ctx context.Context, id uuid.UUID, forUpdate bool) (FUOTADeployment, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	row := ps.db.QueryRowxContext(ctx, `
		select
			id,
			created_at,
			updated_at,
			name,
			multicast_group_id,
			fragmentation_matrix,
			descriptor,
			payload,
			state,
			next_step_after,
			unicast_timeout,
			frag_size,
			redundancy,
			block_ack_delay,
			multicast_timeout,
			group_type,
			dr,
			frequency,
			ping_slot_period
		from
			fuota_deployment
		where
			id = $1`+fu,
		id,
	)

	return ps.scanFUOTADeployment(row)
}

// GetPendingFUOTADeployments returns the pending FUOTA deployments.
func (ps *PgStore) GetPendingFUOTADeployments(ctx context.Context, batchSize int) ([]FUOTADeployment, error) {
	var out []FUOTADeployment

	rows, err := ps.db.QueryxContext(ctx, `
		select
			id,
			created_at,
			updated_at,
			name,
			multicast_group_id,
			fragmentation_matrix,
			descriptor,
			payload,
			state,
			next_step_after,
			unicast_timeout,
			frag_size,
			redundancy,
			block_ack_delay,
			multicast_timeout,
			group_type,
			dr,
			frequency,
			ping_slot_period
		from
			fuota_deployment
		where
			state != $1
			and next_step_after <= $2
		limit $3
		for update
		skip locked`,
		FUOTADeploymentDone,
		time.Now(),
		batchSize,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}
	defer rows.Close()

	for rows.Next() {
		item, err := ps.scanFUOTADeployment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}

	return out, nil
}

// UpdateFUOTADeployment updates the given FUOTA deployment.
func (ps *PgStore) UpdateFUOTADeployment(ctx context.Context, fd *FUOTADeployment) error {
	if err := fd.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	fd.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update fuota_deployment
		set
			updated_at = $2,
			name = $3,
			multicast_group_id = $4,
			fragmentation_matrix = $5,
			descriptor = $6,
			payload = $7,
			state = $8,
			next_step_after = $9,
			unicast_timeout = $10,
			frag_size = $11,
			redundancy = $12,
			block_ack_delay = $13,
			multicast_timeout = $14,
			group_type = $15,
			dr = $16,
			frequency = $17,
			ping_slot_period = $18
		where
			id = $1`,
		fd.ID,
		fd.UpdatedAt,
		fd.Name,
		fd.MulticastGroupID,
		[]byte{fd.FragmentationMatrix},
		fd.Descriptor[:],
		fd.Payload,
		fd.State,
		fd.NextStepAfter,
		fd.UnicastTimeout,
		fd.FragSize,
		fd.Redundancy,
		fd.BlockAckDelay,
		fd.MulticastTimeout,
		fd.GroupType,
		fd.DR,
		fd.Frequency,
		fd.PingSlotPeriod,
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
		"id":     fd.ID,
		"state":  fd.State,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("fuota deployment updated")

	return nil
}

// GetFUOTADeploymentCount returns the number of FUOTA deployments.
func (ps *PgStore) GetFUOTADeploymentCount(ctx context.Context, filters FUOTADeploymentFilters) (int, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct fd.*)
		from
			fuota_deployment fd
		inner join
			fuota_deployment_device fdd
		on
			fd.id = fdd.fuota_deployment_id
		inner join
			device d
		on
			fdd.dev_eui = d.dev_eui
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

// GetFUOTADeployments returns a slice of fuota deployments.
func (ps *PgStore) GetFUOTADeployments(ctx context.Context, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			distinct fd.id,
			fd.created_at,
			fd.updated_at,
			fd.name,
			fd.state,
			fd.next_step_after
		from
			fuota_deployment fd
		inner join
			fuota_deployment_device fdd
		on
			fd.id = fdd.fuota_deployment_id
		inner join
			device d
		on
			fdd.dev_eui = d.dev_eui
	`+filters.SQL()+`
	order by
		fd.created_at desc
	limit :limit
	offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var items []FUOTADeploymentListItem
	if err = sqlx.SelectContext(ctx, ps.db, &items, query, args...); err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return items, nil
}

// GetFUOTADeploymentDevice returns the FUOTA deployment record for the given
// device.
func (ps *PgStore) GetFUOTADeploymentDevice(ctx context.Context, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	var out FUOTADeploymentDevice
	err := sqlx.GetContext(ctx, ps.db, &out, `
		select
			*
		from
			fuota_deployment_device
		where
			fuota_deployment_id = $1
			and dev_eui = $2`,
		fuotaDeploymentID,
		devEUI,
	)
	if err != nil {
		return out, handlePSQLError(Select, err, "select error")
	}
	return out, nil
}

// GetPendingFUOTADeploymentDevice returns the pending FUOTA deployment record
// for the given DevEUI.
func (ps *PgStore) GetPendingFUOTADeploymentDevice(ctx context.Context, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	var out FUOTADeploymentDevice

	err := sqlx.GetContext(ctx, ps.db, &out, `
		select
			*
		from
			fuota_deployment_device
		where
			dev_eui = $1
			and state = $2`,
		devEUI,
		FUOTADeploymentDevicePending,
	)
	if err != nil {
		return out, handlePSQLError(Select, err, "select error")
	}

	return out, nil
}

// UpdateFUOTADeploymentDevice updates the given fuota deployment device record.
func (ps *PgStore) UpdateFUOTADeploymentDevice(ctx context.Context, fdd *FUOTADeploymentDevice) error {
	fdd.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update
			fuota_deployment_device
		set
			updated_at = $3,
			state = $4,
			error_message = $5
		where
			dev_eui = $1
			and fuota_deployment_id = $2`,
		fdd.DevEUI,
		fdd.FUOTADeploymentID,
		fdd.UpdatedAt,
		fdd.State,
		fdd.ErrorMessage,
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
		"dev_eui":             fdd.DevEUI,
		"fuota_deployment_id": fdd.FUOTADeploymentID,
		"state":               fdd.State,
		"ctx_id":              ctx.Value(logging.ContextIDKey),
	}).Info("fuota deployment device updated")

	return nil
}

// GetFUOTADeploymentDeviceCount returns the device count for the given
// FUOTA deployment ID.
func (ps *PgStore) GetFUOTADeploymentDeviceCount(ctx context.Context, fuotaDeploymentID uuid.UUID) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(*)
		from
			fuota_deployment_device
		where
			fuota_deployment_id = $1`,
		fuotaDeploymentID,
	)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetFUOTADeploymentDevices returns a slice of devices for the given FUOTA
// deployment ID.
func (ps *PgStore) GetFUOTADeploymentDevices(ctx context.Context, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error) {
	var out []FUOTADeploymentDeviceListItem

	err := sqlx.SelectContext(ctx, ps.db, &out, `
		select
			dd.created_at,
			dd.updated_at,
			dd.fuota_deployment_id,
			dd.dev_eui,
			d.name as device_name,
			dd.state,
			dd.error_message
		from
			fuota_deployment_device dd
		inner join
			device d
			on dd.dev_eui = d.dev_eui
		where
			dd.fuota_deployment_id = $3
		order by
			d.Name
		limit $1
		offset $2`,
		limit,
		offset,
		fuotaDeploymentID,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return out, nil
}

// GetServiceProfileIDForFUOTADeployment returns the service-profile ID for the given FUOTA deployment.
func (ps *PgStore) GetServiceProfileIDForFUOTADeployment(ctx context.Context, fuotaDeploymentID uuid.UUID) (uuid.UUID, error) {
	var out uuid.UUID

	err := sqlx.GetContext(ctx, ps.db, &out, `
		select
			a.service_profile_id
		from
			fuota_deployment_device fdd
		inner join
			device d
		on
			d.dev_eui = fdd.dev_eui
		inner join
			application a
		on
			a.id = d.application_id
		where
			fdd.fuota_deployment_id = $1
		limit 1`,
		fuotaDeploymentID,
	)
	if err != nil {
		return out, handlePSQLError(Select, err, "select error")
	}

	return out, nil
}

func (ps *PgStore) scanFUOTADeployment(row sqlx.ColScanner) (FUOTADeployment, error) {
	var fd FUOTADeployment

	var fragmentationMatrix []byte
	var descriptor []byte

	err := row.Scan(
		&fd.ID,
		&fd.CreatedAt,
		&fd.UpdatedAt,
		&fd.Name,
		&fd.MulticastGroupID,
		&fragmentationMatrix,
		&descriptor,
		&fd.Payload,
		&fd.State,
		&fd.NextStepAfter,
		&fd.UnicastTimeout,
		&fd.FragSize,
		&fd.Redundancy,
		&fd.BlockAckDelay,
		&fd.MulticastTimeout,
		&fd.GroupType,
		&fd.DR,
		&fd.Frequency,
		&fd.PingSlotPeriod,
	)
	if err != nil {
		return fd, handlePSQLError(Select, err, "select error")
	}

	if len(fragmentationMatrix) != 1 {
		return fd, fmt.Errorf("FragmentationMatrix must have length 1, got: %d", len(fragmentationMatrix))
	}
	fd.FragmentationMatrix = fragmentationMatrix[0]

	if len(descriptor) != len(fd.Descriptor) {
		return fd, fmt.Errorf("Descriptor must have length: %d, got: %d", len(fd.Descriptor), len(descriptor))
	}
	copy(fd.Descriptor[:], descriptor)

	return fd, nil
}

// SetFromRemoteMulticastSetup set remote multicast session error
func (ps *PgStore) SetFromRemoteMulticastSetup(ctx context.Context, fuotaDevelopmentID, multicastGroupID uuid.UUID) error {
	_, err := ps.db.ExecContext(ctx, `
		update
			fuota_deployment_device fdd
		set
			state = $5,
			error_message = $6
		from
			remote_multicast_setup rms
		where
			fdd.fuota_deployment_id = $1
			and rms.multicast_group_id = $2

			and fdd.state = $3
			and rms.state_provisioned = $4

			-- join the two tables
			and fdd.dev_eui = rms.dev_eui`,

		fuotaDevelopmentID,
		multicastGroupID,
		FUOTADeploymentDevicePending,
		false,
		FUOTADeploymentDeviceError,
		"The device failed to provision the remote multicast setup.",
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	return nil
}

// SetFromRemoteFragmentationSession set remote multicast session error
func (ps *PgStore) SetFromRemoteFragmentationSession(ctx context.Context, fuotaDevelopmentID uuid.UUID, fragIdx int) error {
	_, err := ps.db.ExecContext(ctx, `
		update
			fuota_deployment_device fdd
		set
			state = $5,
			error_message = $6
		from
			remote_fragmentation_session rfs
		where
			fdd.fuota_deployment_id = $1
			and rfs.frag_index = $2

			and fdd.state = $3
			and rfs.state_provisioned = $4

			-- join the two tables
			and fdd.dev_eui = rfs.dev_eui`,
		fuotaDevelopmentID,
		fragIdx,
		FUOTADeploymentDevicePending,
		false,
		FUOTADeploymentDeviceError,
		"The device failed to provision the fragmentation session setup.",
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	return nil
}

func (ps *PgStore) SetIncompleteFuotaDevelopment(ctx context.Context, fuotaDevelopmentID uuid.UUID) error {
	_, err := ps.db.ExecContext(ctx, `
		update
			fuota_deployment_device
		set
			state = $3,
			error_message = $4
		where
			fuota_deployment_id = $1
			and state = $2`,
		fuotaDevelopmentID,
		FUOTADeploymentDevicePending,
		FUOTADeploymentDeviceError,
		"Device did not complete the FUOTA deployment or did not confirm that it completed the FUOTA deployment.",
	)

	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	return nil
}
