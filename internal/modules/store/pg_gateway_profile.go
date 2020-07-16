package store

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	gpmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
)

func (ps *pgstore) CheckCreateUpdateDeleteGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
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

func (ps *pgstore) CheckReadListGatewayProfileAccess(ctx context.Context, username string, userID int64) (bool, error) {
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
func (ps *pgstore) CreateGatewayProfile(ctx context.Context, gp *gpmod.GatewayProfile) error {
	gpID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "new uuid v4 error")
	}

	now := time.Now()

	gp.GatewayProfile.Id = gpID.Bytes()
	gp.CreatedAt = now
	gp.UpdatedAt = now

	_, err = ps.db.ExecContext(ctx, `
		insert into gateway_profile (
			gateway_profile_id,
			network_server_id,
			created_at,
			updated_at,
			name
		) values ($1, $2, $3, $4, $5)`,

		gpID,
		gp.NetworkServerID,
		gp.CreatedAt,
		gp.UpdatedAt,
		gp.Name,
	)
	if err != nil {
		return err
	}

	n, err := ps.GetNetworkServer(ctx, gp.NetworkServerID)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.CreateGatewayProfile(ctx, &ns.CreateGatewayProfileRequest{
		GatewayProfile: &gp.GatewayProfile,
	})
	if err != nil {
		return errors.Wrap(err, "create gateway-profile error")
	}

	log.WithFields(log.Fields{
		"id":     gpID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway-profile created")

	return nil
}

// GetGatewayProfile returns the gateway-profile matching the given id.
func (ps *pgstore) GetGatewayProfile(ctx context.Context, id uuid.UUID) (gpmod.GatewayProfile, error) {
	var gp gpmod.GatewayProfile
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
		return gp, err
	}

	n, err := ps.GetNetworkServer(ctx, gp.NetworkServerID)
	if err != nil {
		return gp, errors.Wrap(err, "get network-server error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return gp, errors.Wrap(err, "get network-server client error")
	}

	resp, err := nsClient.GetGatewayProfile(ctx, &ns.GetGatewayProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return gp, errors.Wrap(err, "get gateway-profile error")
	}

	if resp.GatewayProfile == nil {
		return gp, errors.New("gateway_profile must not be nil")
	}

	gp.GatewayProfile = *resp.GatewayProfile

	return gp, nil
}

// UpdateGatewayProfile updates the given gateway-profile.
func (ps *pgstore) UpdateGatewayProfile(ctx context.Context, gp *gpmod.GatewayProfile) error {
	gp.UpdatedAt = time.Now()
	gpID, err := uuid.FromBytes(gp.GatewayProfile.Id)
	if err != nil {
		return errors.Wrap(err, "uuid from bytes error")
	}

	res, err := ps.db.ExecContext(ctx, `
		update gateway_profile
		set
			updated_at = $2,
			network_server_id = $3,
			name = $4
		where
			gateway_profile_id = $1`,
		gpID,
		gp.UpdatedAt,
		gp.NetworkServerID,
		gp.Name,
	)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	n, err := ps.GetNetworkServer(ctx, gp.NetworkServerID)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.UpdateGatewayProfile(context.Background(), &ns.UpdateGatewayProfileRequest{
		GatewayProfile: &gp.GatewayProfile,
	})
	if err != nil {
		return errors.Wrap(err, "update gateway-profile error")
	}

	return nil
}

// DeleteGatewayProfile deletes the gateway-profile matching the given id.
func (ps *pgstore) DeleteGatewayProfile(ctx context.Context, id uuid.UUID) error {
	n, err := ps.GetNetworkServerForGatewayProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := ps.db.ExecContext(ctx, `
		delete from gateway_profile
		where
			gateway_profile_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.DeleteGatewayProfile(ctx, &ns.DeleteGatewayProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return errors.Wrap(err, "delete gateway-profile error")
	}

	return nil
}

// GetGatewayProfileCount returns the total number of gateway-profiles.
func (ps *pgstore) GetGatewayProfileCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(*)
		from gateway_profile`)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetGatewayProfileCountForNetworkServerID returns the total number of
// gateway-profiles given a network-server ID.
func (ps *pgstore) GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error) {
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
		return 0, err
	}

	return count, nil
}

// GetGatewayProfiles returns a slice of gateway-profiles.
func (ps *pgstore) GetGatewayProfiles(ctx context.Context, limit, offset int) ([]gpmod.GatewayProfileMeta, error) {
	var gps []gpmod.GatewayProfileMeta
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
		return nil, err
	}

	return gps, nil
}

// GetGatewayProfilesForNetworkServerID returns a slice of gateway-profiles
// for the given network-server ID.
func (ps *pgstore) GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]gpmod.GatewayProfileMeta, error) {
	var gps []gpmod.GatewayProfileMeta
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
		return nil, err
	}

	return gps, nil
}
