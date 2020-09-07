package pgstore

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

func (ps *pgstore) CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
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

func (ps *pgstore) CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
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

func (ps *pgstore) CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
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

func (ps *pgstore) CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
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
func (ps *pgstore) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	var dps []store.DeviceProfileMeta
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

// DeleteDeviceProfile deletes the device-profile matching the given id.
func (ps *pgstore) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	n, err := ps.GetNetworkServerForDeviceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := ps.db.ExecContext(ctx, "delete from device_profile where device_profile_id = $1", id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	nsClient, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.DeleteDeviceProfile(ctx, &ns.DeleteDeviceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device-profile error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("device-profile deleted")

	return nil
}

// GetDeviceProfile returns the device-profile matching the given id.
// When forUpdate is set to true, then db must be a db transaction.
// When localOnly is set to true, no call to the network-server is made to
// retrieve additional device data.
func (ps *pgstore) GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (store.DeviceProfile, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var dp store.DeviceProfile

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
		return dp, errors.Wrap(err, "select error")
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
		return dp, errors.Wrap(err, "scan error")
	}
	/*
		n, err := ps.GetNetworkServer(ctx, dp.NetworkServerID)
		if err != nil {
			return dp, errors.Wrap(err, "get network-server error")
		}

		nstruct := nscli.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}

		nsClient, err := nstruct.GetNetworkServiceClient()
		if err != nil {
			return dp, errors.Wrap(err, "get network-server client error")
		}

		resp, err := nsClient.GetDeviceProfile(ctx, &ns.GetDeviceProfileRequest{
			Id: id.Bytes(),
		})
		if err != nil {
			return dp, errors.Wrap(err, "get device-profile error")
		}
		if resp.DeviceProfile == nil {
			return dp, errors.New("device_profile must not be nil")
		}

		dp.DeviceProfile = *resp.DeviceProfile
	*/
	return dp, nil
}
