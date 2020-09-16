package pgstore

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

func (ps *pgstore) CheckCreateNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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

func (ps *pgstore) CheckListNetworkServersAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
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

func (ps *pgstore) CheckReadNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join service_profile sp
			on sp.organization_id = o.id
		left join network_server ns
			on ns.id = sp.network_server_id
	`
	// global admin
	// organization admin
	// organization gateway admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "ns.id = $2"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_gateway_admin = true", "ns.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, networkserverID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdateDeleteNetworkServerAccess(ctx context.Context, username string, networkserverID, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join service_profile sp
			on sp.organization_id = o.id
		left join network_server ns
			on ns.id = sp.network_server_id
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
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, username, networkserverID, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// GetDefaultNetworkServer returns the network-server matching the given name.
func (ps *pgstore) GetDefaultNetworkServer(ctx context.Context) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, "select * from network_server where name = 'default_network_server'")
	if err != nil {
		return n, errors.Wrap(err, "select error")
	}

	return n, nil
}

// CreateNetworkServer creates the given network-server.
func (ps *pgstore) CreateNetworkServer(ctx context.Context, n *store.NetworkServer) error {
	if err := n.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now

	err := sqlx.GetContext(ctx, ps.db, &n.ID, `
		insert into network_server (
			created_at,
			updated_at,
			name,
			server,
			ca_cert,
			tls_cert,
			tls_key,
			routing_profile_ca_cert,
			routing_profile_tls_cert,
			routing_profile_tls_key,
			gateway_discovery_enabled,
			gateway_discovery_interval,
			gateway_discovery_tx_frequency,
			gateway_discovery_dr,
			region,
			version
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		returning id`,
		n.CreatedAt,
		n.UpdatedAt,
		n.Name,
		n.Server,
		n.CACert,
		n.TLSCert,
		n.TLSKey,
		n.RoutingProfileCACert,
		n.RoutingProfileTLSCert,
		n.RoutingProfileTLSKey,
		n.GatewayDiscoveryEnabled,
		n.GatewayDiscoveryInterval,
		n.GatewayDiscoveryTXFrequency,
		n.GatewayDiscoveryDR,
		n.Region,
		n.Version,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
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

	rpID, err := uuid.FromString(external.GetApplicationServerID())
	if err != nil {
		return errors.Wrap(err, "uuid from string error")
	}

	_, err = nsClient.CreateRoutingProfile(ctx, &ns.CreateRoutingProfileRequest{
		RoutingProfile: &ns.RoutingProfile{
			Id:      rpID.Bytes(),
			AsId:    config.C.ApplicationServer.API.PublicHost,
			CaCert:  n.RoutingProfileCACert,
			TlsCert: n.RoutingProfileTLSCert,
			TlsKey:  n.RoutingProfileTLSKey,
		},
	})
	if err != nil {
		return errors.Wrap(err, "create routing-profile error")
	}

	log.WithFields(log.Fields{
		"id":     n.ID,
		"name":   n.Name,
		"server": n.Server,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("network-server created")
	return nil
}

// GetNetworkServer returns the network-server matching the given id.
func (ps *pgstore) GetNetworkServer(ctx context.Context, id int64) (store.NetworkServer, error) {
	var networkServer store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &networkServer, "select * from network_server where id = $1", id)
	if err != nil {
		return networkServer, handlePSQLError(Select, err, "select error")
	}

	return networkServer, nil
}

// UpdateNetworkServer updates the given network-server.
func (ps *pgstore) UpdateNetworkServer(ctx context.Context, n *store.NetworkServer) error {
	if err := n.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	n.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update network_server
		set
			updated_at = $2,
			name = $3,
			server = $4,
			ca_cert = $5,
			tls_cert = $6,
			tls_key = $7,
			routing_profile_ca_cert = $8,
			routing_profile_tls_cert = $9,
			routing_profile_tls_key = $10,
			gateway_discovery_enabled = $11,
			gateway_discovery_interval = $12,
			gateway_discovery_tx_frequency = $13,
			gateway_discovery_dr = $14
		where id = $1`,
		n.ID,
		n.UpdatedAt,
		n.Name,
		n.Server,
		n.CACert,
		n.TLSCert,
		n.TLSKey,
		n.RoutingProfileCACert,
		n.RoutingProfileTLSCert,
		n.RoutingProfileTLSKey,
		n.GatewayDiscoveryEnabled,
		n.GatewayDiscoveryInterval,
		n.GatewayDiscoveryTXFrequency,
		n.GatewayDiscoveryDR,
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

	rpID, err := uuid.FromString(external.GetApplicationServerID())
	if err != nil {
		return errors.Wrap(err, "uuid from string error")
	}

	_, err = nsClient.UpdateRoutingProfile(ctx, &ns.UpdateRoutingProfileRequest{
		RoutingProfile: &ns.RoutingProfile{
			Id:      rpID.Bytes(),
			AsId:    config.C.ApplicationServer.API.PublicHost,
			CaCert:  n.RoutingProfileCACert,
			TlsCert: n.RoutingProfileTLSCert,
			TlsKey:  n.RoutingProfileTLSKey,
		},
	})
	if err != nil {
		return errors.Wrap(err, "update routing-profile error")
	}

	log.WithFields(log.Fields{
		"id":     n.ID,
		"name":   n.Name,
		"server": n.Server,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("network-server updated")
	return nil
}

// DeleteNetworkServer deletes the network-server matching the given id.
func (ps *pgstore) DeleteNetworkServer(ctx context.Context, id int64) error {
	n, err := ps.GetNetworkServer(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := ps.db.ExecContext(ctx, "delete from network_server where id = $1", id)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
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

	rpID, err := uuid.FromString(external.GetApplicationServerID())
	if err != nil {
		return errors.Wrap(err, "uuid from string error")
	}

	_, err = nsClient.DeleteRoutingProfile(ctx, &ns.DeleteRoutingProfileRequest{
		Id: rpID.Bytes(),
	})
	if err != nil {
		return errors.Wrap(err, "delete routing-profile error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("network-server deleted")
	return nil
}

// GetNetworkServerCount returns the total number of network-servers.
func (ps *pgstore) GetNetworkServerCount(ctx context.Context, filters store.NetworkServerFilters) (int, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct ns.id)
		from
			network_server ns
		left join service_profile sp
			on ns.id = sp.network_server_id
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

// GetNetworkServerCountForOrganizationID returns the total number of
// network-servers accessible for the given organization id.
// A network-server is accessible for an organization when it is used by one
// of its service-profiles.
func (ps *pgstore) GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
		select
			count (distinct ns.id)
		from
			network_server ns
		inner join service_profile sp
			on sp.network_server_id = ns.id
		where
			sp.organization_id = $1`,
		organizationID,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetNetworkServers returns a slice of network-servers.
func (ps *pgstore) GetNetworkServers(ctx context.Context, filters store.NetworkServerFilters) ([]store.NetworkServer, error) {
	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			distinct ns.*
		from
			network_server ns
		left join service_profile sp
			on ns.id = sp.network_server_id
	`+filters.SQL()+`
		order by ns.name
		limit :limit
		offset :offset
	`, filters)
	if err != nil {
		return nil, errors.Wrap(err, "named query error")
	}

	var nss []store.NetworkServer
	err = sqlx.SelectContext(ctx, ps.db, &nss, query, args...)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return nss, nil
}

// GetNetworkServersForOrganizationID returns a slice of network-server
// accessible for the given organization id.
// A network-server is accessible for an organization when it is used by one
// of its service-profiles.
func (ps *pgstore) GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]store.NetworkServer, error) {
	var nss []store.NetworkServer
	err := sqlx.SelectContext(ctx, ps.db, &nss, `
		select
			ns.*
		from
			network_server ns
		inner join service_profile sp
			on sp.network_server_id = ns.id
		where
			sp.organization_id = $1
		group by ns.id
		order by name
		limit $2 offset $3`,
		organizationID,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return nss, nil
}

// GetNetworkServerForDevEUI returns the network-server for the given DevEUI.
func (ps *pgstore) GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from
			network_server ns
		inner join device_profile dp
			on dp.network_server_id = ns.id
		inner join device d
			on d.device_profile_id = dp.device_profile_id
		where
			d.dev_eui = $1`,
		devEUI,
	)
	if err != nil {
		return n, handlePSQLError(Select, err, "select error")
	}
	return n, nil
}

// GetNetworkServerForDeviceProfileID returns the network-server for the given
// device-profile id.
func (ps *pgstore) GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from
			network_server ns
		inner join device_profile dp
			on dp.network_server_id = ns.id
		where
			dp.device_profile_id = $1`,
		id,
	)
	if err != nil {
		return n, errors.Wrap(err, "select error")
	}
	return n, nil
}

// GetNetworkServerForGatewayMAC returns the network-server for a given
// gateway mac.
func (ps *pgstore) GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from network_server ns
		inner join gateway gw
			on gw.network_server_id = ns.id
		where
			gw.mac = $1`,
		mac[:],
	)
	if err != nil {
		return n, handlePSQLError(Select, err, "select error")
	}
	return n, nil
}

// GetNetworkServerForGatewayProfileID returns the network-server for the given
// gateway-profile id.
func (ps *pgstore) GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from
			network_server ns
		inner join gateway_profile gp
			on gp.network_server_id = ns.id
		where
			gp.gateway_profile_id = $1`,
		id,
	)
	if err != nil {
		return n, handlePSQLError(Select, err, "select errror")
	}
	return n, nil
}

// GetNetworkServerForMulticastGroupID returns the network-server for the given
// multicast-group id.
func (ps *pgstore) GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from
			network_server ns
		inner join service_profile sp
			on sp.network_server_id = ns.id
		inner join multicast_group mg
			on mg.service_profile_id = sp.service_profile_id
		where
			mg.id = $1
	`, id)
	if err != nil {
		return n, handlePSQLError(Select, err, "select error")
	}
	return n, nil
}

// GetNetworkServerForServiceProfileID returns the network-server for the given
// service-profile id.
func (ps *pgstore) GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (store.NetworkServer, error) {
	var n store.NetworkServer
	err := sqlx.GetContext(ctx, ps.db, &n, `
		select
			ns.*
		from
			network_server ns
		inner join service_profile sp
			on sp.network_server_id = ns.id
		where
			sp.service_profile_id = $1`,
		id,
	)
	if err != nil {
		return n, handlePSQLError(Select, err, "select error")
	}
	return n, nil
}
