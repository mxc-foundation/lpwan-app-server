package pgstore

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
)

type NetworkServerHandler struct {
	db *storage.DBLogger
}

func New(db *storage.DBLogger) *NetworkServerHandler {
	networkServerHandler = NetworkServerHandler{
		db: db,
	}
	return &networkServerHandler
}

var networkServerHandler NetworkServerHandler

func Handler() *NetworkServerHandler {
	return &networkServerHandler
}

// GetDefaultNetworkServer returns the network-server matching the given name.
func (h *NetworkServerHandler) GetDefaultNetworkServer(ctx context.Context, db sqlx.Queryer) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(db, &n, "select * from network_server where name = 'default_network_server'")
	if err != nil {
		return n, errors.Wrap(err, "select error")
	}

	return n, nil
}

// CreateNetworkServer creates the given network-server.
func (h *NetworkServerHandler) CreateNetworkServer(ctx context.Context, n *nsmod.NetworkServer) error {
	if err := n.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now

	err := sqlx.Get(h.db, &n.ID, `
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
			gateway_discovery_dr
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
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
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	rpID, err := uuid.FromString(config.C.ApplicationServer.ID)
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
func (h *NetworkServerHandler) GetNetworkServer(ctx context.Context, id int64) (nsmod.NetworkServer, error) {
	var ns nsmod.NetworkServer
	err := sqlx.Get(h.db, &ns, "select * from network_server where id = $1", id)
	if err != nil {
		return ns, errors.Wrap(err, "select error")
	}

	return ns, nil
}

// UpdateNetworkServer updates the given network-server.
func (h *NetworkServerHandler) UpdateNetworkServer(ctx context.Context, n *nsmod.NetworkServer) error {
	if err := n.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	n.UpdatedAt = time.Now()

	res, err := h.db.Exec(`
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
		return errors.Wrap(err, "update error")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("ErrDoesNotExist")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	rpID, err := uuid.FromString(config.C.ApplicationServer.ID)
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
func (h *NetworkServerHandler) DeleteNetworkServer(ctx context.Context, id int64) error {
	n, err := h.GetNetworkServer(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := h.db.Exec("delete from network_server where id = $1", id)
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

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	rpID, err := uuid.FromString(config.C.ApplicationServer.ID)
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
func (h *NetworkServerHandler) GetNetworkServerCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.Get(h.db, &count, "select count(*) from network_server")
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return count, nil
}

// GetNetworkServerCountForOrganizationID returns the total number of
// network-servers accessible for the given organization id.
// A network-server is accessible for an organization when it is used by one
// of its service-profiles.
func (h *NetworkServerHandler) GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error) {
	var count int
	err := sqlx.Get(h.db, &count, `
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
func (h *NetworkServerHandler) GetNetworkServers(ctx context.Context, limit, offset int) ([]nsmod.NetworkServer, error) {
	var nss []nsmod.NetworkServer
	err := sqlx.Select(h.db, &nss, `
		select *
		from network_server
		order by name
		limit $1 offset $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return nss, nil
}

// GetNetworkServersForOrganizationID returns a slice of network-server
// accessible for the given organization id.
// A network-server is accessible for an organization when it is used by one
// of its service-profiles.
func (h *NetworkServerHandler) GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]nsmod.NetworkServer, error) {
	var nss []nsmod.NetworkServer
	err := sqlx.Select(h.db, &nss, `
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
func (h *NetworkServerHandler) GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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
		return n, errors.Wrap(err, "select error")
	}
	return n, nil
}

// GetNetworkServerForDeviceProfileID returns the network-server for the given
// device-profile id.
func (h *NetworkServerHandler) GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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

// GetNetworkServerForServiceProfileID returns the network-server for the given
// service-profile id.
func (h *NetworkServerHandler) GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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
		return n, errors.Wrap(err, "select error")
	}
	return n, nil
}

// GetNetworkServerForGatewayMAC returns the network-server for a given
// gateway mac.
func (h *NetworkServerHandler) GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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
		return n, errors.Wrap(err, "select error")
	}
	return n, nil
}

// GetNetworkServerForGatewayProfileID returns the network-server for the given
// gateway-profile id.
func (h *NetworkServerHandler) GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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
		return n, errors.Wrap(err, "select errror")
	}
	return n, nil
}

// GetNetworkServerForMulticastGroupID returns the network-server for the given
// multicast-group id.
func (h *NetworkServerHandler) GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (nsmod.NetworkServer, error) {
	var n nsmod.NetworkServer
	err := sqlx.Get(h.db, &n, `
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
		return n, errors.Wrap(err, "select error")
	}
	return n, nil
}
