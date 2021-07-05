package pgstore

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	. "github.com/mxc-foundation/lpwan-app-server/internal/api/external/gp"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
)

func (ps *PgStore) CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	query := `
		select
			1
		from
			"user" u
	`
	// global admin
	var where = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
	}

	var ors []string
	for _, ands := range where {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	query = "select count(*) from (" + query + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := ps.db.QueryRowContext(ctx, query, username, userID).Scan(&count); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *PgStore) CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
	query := `
		select
			1
		from
			"user" u
	`
	// any active user
	var where = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true"},
	}

	var ors []string
	for _, ands := range where {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	query = "select count(*) from (" + query + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := ps.db.QueryRowContext(ctx, query, username, userID).Scan(&count); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CreateGatewayProfile creates the given gateway-profile.
// This will create the gateway-profile at the network-server side and will
// create a local reference record.
func (ps *PgStore) CreateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	var statsInterval time.Duration
	if gp.GatewayProfile.StatsInterval != nil {
		err := gp.GatewayProfile.StatsInterval.CheckValid()
		if err != nil {
			return errors.Wrap(err, "stats interval error")
		}
		statsInterval = gp.GatewayProfile.StatsInterval.AsDuration()
	}

	gpID, err := uuid.FromBytes(gp.GatewayProfile.Id)
	if err != nil {
		return err
	}

	_, err = ps.db.ExecContext(ctx, `
		insert into gateway_profile (
			gateway_profile_id,
			network_server_id,
			created_at,
			updated_at,
			name,
			stats_interval
		) values ($1, $2, $3, $4, $5, $6)`,

		gpID,
		gp.NetworkServerID,
		gp.CreatedAt,
		gp.UpdatedAt,
		gp.Name,
		statsInterval,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"id":     gpID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway-profile created")

	return nil
}

// GetGatewayProfile returns the gateway-profile matching the given id.
func (ps *PgStore) GetGatewayProfile(ctx context.Context, id uuid.UUID) (GatewayProfile, error) {
	var gp GatewayProfile
	err := sqlx.GetContext(ctx, ps.db, &gp, `
		select
			network_server_id,
			name,
			created_at,
			updated_at
		from gateway_profile
		where
			gateway_profile_id = $1`,
		id,
	)
	if err != nil {
		return gp, handlePSQLError(Select, err, "select error")
	}

	return gp, nil
}

// UpdateGatewayProfile updates the given gateway-profile.
func (ps *PgStore) UpdateGatewayProfile(ctx context.Context, gp *GatewayProfile) error {
	gp.UpdatedAt = time.Now()
	gpID, err := uuid.FromBytes(gp.GatewayProfile.Id)
	if err != nil {
		return errors.Wrap(err, "uuid from bytes error")
	}

	var statsInterval time.Duration
	if gp.GatewayProfile.StatsInterval != nil {
		err = gp.GatewayProfile.StatsInterval.CheckValid()
		if err != nil {
			return errors.Wrap(err, "stats interval error")
		}
		statsInterval = gp.GatewayProfile.StatsInterval.AsDuration()
	}

	res, err := ps.db.ExecContext(ctx, `
		update gateway_profile
		set
			updated_at = $2,
			network_server_id = $3,
			name = $4,
			stats_interval = $5
		where
			gateway_profile_id = $1`,
		gpID,
		gp.UpdatedAt,
		gp.NetworkServerID,
		gp.Name,
		statsInterval,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update gateway-profile error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	return nil
}

// DeleteGatewayProfile deletes the gateway-profile matching the given id.
func (ps *PgStore) DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from gateway_profile
		where
			gateway_profile_id = $1`,
		id,
	)
	if err != nil {
		return handlePSQLError(Delete, err, "delete gateway-profile error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}

	return nil
}

// GetGatewayProfileCount returns the total number of gateway-profiles.
func (ps *PgStore) GetGatewayProfileCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(*)
		from gateway_profile`)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errHandler.ErrDoesNotExist
		}
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetGatewayProfileCountForNetworkServerID returns the total number of
// gateway-profiles given a network-server ID.
func (ps *PgStore) GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(*)
		from gateway_profile
		where
			network_server_id = $1`,
		networkServerID,
	)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}

	return count, nil
}

// GetGatewayProfiles returns a slice of gateway-profiles.
func (ps *PgStore) GetGatewayProfiles(ctx context.Context, limit, offset int) ([]GatewayProfileMeta, error) {
	var gps []GatewayProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &gps, `
		select
			gp.*,
			n.name as network_server_name
		from
			gateway_profile gp
		inner join
			network_server n
		on
			n.id = gp.network_server_id
		order by
			name
		limit $1 offset $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return gps, nil
}

// GetGatewayProfilesForNetworkServerID returns a slice of gateway-profiles
// for the given network-server ID.
func (ps *PgStore) GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]GatewayProfileMeta, error) {
	var gps []GatewayProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &gps, `
		select
			gp.*,
			n.name as network_server_name
		from
			gateway_profile gp
		inner join
			network_server n
		on
			n.id = gp.network_server_id
		where
			network_server_id = $1
		order by
			name
		limit $2 offset $3`,
		networkServerID,
		limit,
		offset,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return gps, nil
}
