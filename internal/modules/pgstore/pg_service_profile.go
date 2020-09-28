package pgstore

import (
	"context"
	"strings"
	"time"

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

func (ps *pgstore) CheckReadServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
	`
	// global admin
	// organization users to which the service-profile is linked
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "sp.service_profile_id = $2"},
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

func (ps *pgstore) CheckUpdateDeleteServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
	`
	// global admin
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
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

func (ps *pgstore) CheckCreateServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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
	userWhere := [][]string{
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

func (ps *pgstore) CheckListServiceProfilesAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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
	// organization user (when organization id is given)
	// any active user (filtered by user)
	userWhere := [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 > 0", "o.id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 = 0"},
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

// CreateServiceProfile creates the given service-profile.
func (ps *pgstore) CreateServiceProfile(ctx context.Context, sp *store.ServiceProfile) error {
	if err := sp.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	spID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "new uuid v4 error")
	}

	now := time.Now()
	sp.CreatedAt = now
	sp.UpdatedAt = now
	sp.ServiceProfile.Id = spID.Bytes()

	_, err = ps.db.ExecContext(ctx, `
		insert into service_profile (
			service_profile_id,
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name
		) values ($1, $2, $3, $4, $5, $6)`,
		spID,
		sp.NetworkServerID,
		sp.OrganizationID,
		sp.CreatedAt,
		sp.UpdatedAt,
		sp.Name,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	n, err := ps.GetNetworkServer(ctx, sp.NetworkServerID)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// delete device from networkserver
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
	_, err = nsClient.CreateServiceProfile(ctx, &ns.CreateServiceProfileRequest{
		ServiceProfile: &sp.ServiceProfile,
	})
	if err != nil {
		return errors.Wrap(err, "create service-profile error")
	}

	log.WithFields(log.Fields{
		"id":     spID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("service-profile created")
	return nil
}

// GetServiceProfile returns the service-profile matching the given id.
func (ps *pgstore) GetServiceProfile(ctx context.Context, id uuid.UUID, localOnly bool) (store.ServiceProfile, error) {
	var sp store.ServiceProfile
	row := ps.db.QueryRowxContext(ctx, `
		select
			network_server_id,
			organization_id,
			created_at,
			updated_at,
			name
		from service_profile
		where
			service_profile_id = $1`,
		id,
	)
	if err := row.Err(); err != nil {
		return sp, handlePSQLError(Select, err, "select error")
	}

	err := row.Scan(&sp.NetworkServerID, &sp.OrganizationID, &sp.CreatedAt, &sp.UpdatedAt, &sp.Name)
	if err != nil {
		return sp, handlePSQLError(Scan, err, "scan error")
	}

	if localOnly {
		return sp, nil
	}

	n, err := ps.GetNetworkServer(ctx, sp.NetworkServerID)
	if err != nil {
		return sp, errors.Wrap(err, "get network-server errror")
	}

	// delete device from networkserver
	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	nsClient, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return sp, errors.Wrap(err, "get network-server client error")
	}

	resp, err := nsClient.GetServiceProfile(ctx, &ns.GetServiceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return sp, errors.Wrap(err, "get service-profile error")
	}

	if resp.ServiceProfile == nil {
		return sp, errors.New("service_profile must not be nil")
	}

	sp.ServiceProfile = *resp.ServiceProfile

	return sp, nil
}

// UpdateServiceProfile updates the given service-profile.
func (ps *pgstore) UpdateServiceProfile(ctx context.Context, sp *store.ServiceProfile) error {
	if err := sp.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	spID, err := uuid.FromBytes(sp.ServiceProfile.Id)
	if err != nil {
		return errors.Wrap(err, "uuid from bytes error")
	}

	sp.UpdatedAt = time.Now()
	res, err := ps.db.ExecContext(ctx, `
		update service_profile
		set
			updated_at = $2,
			name = $3
		where service_profile_id = $1`,
		spID,
		sp.UpdatedAt,
		sp.Name,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return store.ErrDoesNotExist
	}

	n, err := ps.GetNetworkServer(ctx, sp.NetworkServerID)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// delete device from networkserver
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
	_, err = nsClient.UpdateServiceProfile(ctx, &ns.UpdateServiceProfileRequest{
		ServiceProfile: &sp.ServiceProfile,
	})
	if err != nil {
		return errors.Wrap(err, "update service-profile error")
	}

	log.WithFields(log.Fields{
		"id":     spID,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("service-profile updated")

	return nil
}

// DeleteServiceProfile deletes the service-profile matching the given id.
func (ps *pgstore) DeleteServiceProfile(ctx context.Context, id uuid.UUID) error {
	n, err := ps.GetNetworkServerForServiceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
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

	res, err := ps.db.ExecContext(ctx, "delete from service_profile where service_profile_id = $1", id)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return store.ErrDoesNotExist
	}

	_, err = nsClient.DeleteServiceProfile(ctx, &ns.DeleteServiceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete service-profile error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("service-profile deleted")

	return nil
}

// GetServiceProfileCount returns the total number of service-profiles.
func (ps *pgstore) GetServiceProfileCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, "select count(*) from service_profile")
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}
	return count, nil
}

// GetServiceProfileCountForOrganizationID returns the total number of
// service-profiles for the given organization id.
func (ps *pgstore) GetServiceProfileCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, "select count(*) from service_profile where organization_id = $1", organizationID)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}
	return count, nil
}

// GetServiceProfileCountForUser returns the total number of service-profiles
// for the given user ID.
func (ps *pgstore) GetServiceProfileCountForUser(ctx context.Context, userID int64) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count(sp.*)
		from service_profile sp
		inner join organization o
			on o.id = sp.organization_id
		inner join organization_user ou
			on ou.organization_id = o.id
		inner join "user" u
			on u.id = ou.user_id
		where
			u.id = $1`,
		userID,
	)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}
	return count, nil
}

// GetServiceProfiles returns a slice of service-profiles.
func (ps *pgstore) GetServiceProfiles(ctx context.Context, limit, offset int) ([]store.ServiceProfileMeta, error) {
	var sps []store.ServiceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &sps, `
		select *
		from service_profile
		order by name
		limit $1 offset $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return sps, nil
}

// GetServiceProfilesForOrganizationID returns a slice of service-profiles
// for the given organization id.
func (ps *pgstore) GetServiceProfilesForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]store.ServiceProfileMeta, error) {
	var sps []store.ServiceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &sps, `
		select
			sp.*,
			ns.name as network_server_name
		from
			service_profile sp
		inner join network_server ns
			on sp.network_server_id = ns.id
		where
			sp.organization_id = $1
		order by sp.name
		limit $2 offset $3`,
		organizationID,
		limit,
		offset,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return sps, nil
}

// GetServiceProfilesForUser returns a slice of service-profile for the given
// user ID.
func (ps *pgstore) GetServiceProfilesForUser(ctx context.Context, userID int64, limit, offset int) ([]store.ServiceProfileMeta, error) {
	var sps []store.ServiceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &sps, `
		select
			sp.*
		from service_profile sp
		inner join organization o
			on o.id = sp.organization_id
		inner join organization_user ou
			on ou.organization_id = o.id
		inner join "user" u
			on u.id = ou.user_id
		where
			u.id = $1
		order by sp.name
		limit $2 offset $3`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return sps, nil
}

// DeleteAllServiceProfilesForOrganizationID deletes all service-profiles
// given an organization id.
func (ps *pgstore) DeleteAllServiceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	var sps []store.ServiceProfileMeta
	err := sqlx.SelectContext(ctx, ps.db, &sps, "select * from service_profile where organization_id = $1", organizationID)
	if err != nil {
		return handlePSQLError(Select, err, "select error")
	}

	for _, sp := range sps {
		err = ps.DeleteServiceProfile(ctx, sp.ServiceProfileID)
		if err != nil {
			return errors.Wrap(err, "delete service-profile error")
		}
	}

	return nil
}
