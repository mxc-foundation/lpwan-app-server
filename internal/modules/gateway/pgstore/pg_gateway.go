package pgstore

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	gwmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
)

type GWHandler struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *GWHandler {
	return &GWHandler{
		db: db,
	}
}

func (h *GWHandler) AddGatewayFirmware(gwFw *gwmod.GatewayFirmware) (model string, err error) {
	err = h.db.QueryRowx(`
		insert into gateway_firmware (
			model, 
			resource_link, 
			md5_hash
		) values ($1, $2, $3)
		returning 
		    model;
		`,
		gwFw.Model,
		gwFw.ResourceLink,
		gwFw.FirmwareHash[:]).Scan(&model)

	if err != nil {
		return "", errors.Wrap(err, "AddGatewayFirmware")
	}
	return model, nil
}

func (h *GWHandler) GetGatewayFirmware(model string, forUpdate bool) (gwFw gwmod.GatewayFirmware, err error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	err = sqlx.Get(h.db, &gwFw, "select * from gateway_firmware where model = $1 "+fu, model)
	if err != nil {
		if err == sql.ErrNoRows {
			return gwFw, ErrDoesNotExist
		}
		return gwFw, err
	}
	return gwFw, nil
}

func (h *GWHandler) GetGatewayFirmwareList() (list []gwmod.GatewayFirmware, err error) {
	res, err := h.db.Query(`
		select 
			model, 
			resource_link, 
			md5_hash 
		from 
		     gateway_firmware ;
	`)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, ErrDoesNotExist
		}
		return nil, errors.Wrap(err, "GetGatewayFirmwareList")
	}

	defer res.Close()
	for res.Next() {
		var tmp []byte
		gatewayFirmware := gwmod.GatewayFirmware{}
		err := res.Scan(&gatewayFirmware.Model,
			&gatewayFirmware.ResourceLink,
			&tmp)
		if err != nil {
			return nil, errors.Wrap(err, "GetGatewayFirmwareList")
		}

		copy(gatewayFirmware.FirmwareHash[:], tmp)

		list = append(list, gatewayFirmware)
	}

	return list, nil
}

func (h *GWHandler) UpdateGatewayFirmware(gwFw *gwmod.GatewayFirmware) (model string, err error) {
	err = h.db.QueryRowx(`
		update 
		    gateway_firmware 
		set 
		    resource_link=$1, md5_hash=$2 
		where 
		      model =$3
		returning 
		    model;
		`,
		gwFw.ResourceLink,
		gwFw.FirmwareHash[:],
		gwFw.Model).Scan(&model)

	if err != nil {
		return "", errors.Wrap(err, "UpdateGatewayFirmware")
	}
	return model, nil
}

func (h *GWHandler) UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error {
	res, err := h.db.Exec(`
		update gateway
			set config = $1
		where
			mac = $2`,
		config,
		mac[:])
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	return nil
}

// CreateGateway creates the given Gateway.
func (h *GWHandler) CreateGateway(ctx context.Context, gw *gwmod.Gateway) error {
	if err := gw.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()
	gw.CreatedAt = now
	timestampCreatedAt, _ := ptypes.TimestampProto(gw.CreatedAt)

	gw.UpdatedAt = now

	_, err := h.db.Exec(`
		insert into gateway (
			mac,
			created_at,
			updated_at,
			name,
			description,
			organization_id,
			ping,
			last_ping_id,
			last_ping_sent_at,
			network_server_id,
			gateway_profile_id,
			first_seen_at,
			last_seen_at,
			latitude,
			longitude,
			altitude,
		    model,
		    first_heartbeat,
		    last_heartbeat,
		    config,
		    os_version,
			sn
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
		          $20, $21, $22)`,
		gw.MAC[:],
		gw.CreatedAt,
		gw.UpdatedAt,
		gw.Name,
		gw.Description,
		gw.OrganizationID,
		gw.Ping,
		gw.LastPingID,
		gw.LastPingSentAt,
		gw.NetworkServerID,
		gw.GatewayProfileID,
		gw.FirstSeenAt,
		gw.LastSeenAt,
		gw.Latitude,
		gw.Longitude,
		gw.Altitude,
		gw.Model,
		gw.FirstHeartbeat,
		gw.LastHeartbeat,
		gw.Config,
		gw.OsVersion,
		gw.SerialNumber)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	// add this gateway to m2m server
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	gwClient := m2m_api.NewM2MServerServiceClient(m2mClient)
	if err == nil {
		_, err = gwClient.AddGatewayInM2MServer(context.Background(), &m2m_api.AddGatewayInM2MServerRequest{
			OrgId: gw.OrganizationID,
			GwProfile: &m2m_api.AppServerGatewayProfile{
				Mac:         gw.MAC.String(),
				OrgId:       gw.OrganizationID,
				Description: gw.Description,
				Name:        gw.Name,
				CreatedAt:   timestampCreatedAt,
			},
		})
		if err != nil {
			log.WithError(err).Error("m2m server create gateway api error")
		}
	} else {
		log.WithError(err).Error("get m2m-server client error")
	}

	log.WithFields(log.Fields{
		"id":     gw.MAC,
		"name":   gw.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway created")
	return nil
}

// UpdateGateway updates the given Gateway.
func (h *GWHandler) UpdateGateway(ctx context.Context, gw *gwmod.Gateway) error {
	if err := gw.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	now := time.Now()

	res, err := h.db.Exec(`
		update gateway
			set updated_at = $2,
			name = $3,
			description = $4,
			organization_id = $5,
			ping = $6,
			last_ping_id = $7,
			last_ping_sent_at = $8,
			network_server_id = $9,
			gateway_profile_id = $10,
			first_seen_at = $11,
			last_seen_at = $12,
			latitude = $13,
			longitude = $14,
			altitude = $15,
		    model = $16,
		    config = $17,
		    os_version = $18,
		    statistics = $19,
			firmware_hash = $20
		where
			mac = $1`,
		gw.MAC[:],
		time.Now().UTC(),
		gw.Name,
		gw.Description,
		gw.OrganizationID,
		gw.Ping,
		gw.LastPingID,
		gw.LastPingSentAt,
		gw.NetworkServerID,
		gw.GatewayProfileID,
		gw.FirstSeenAt,
		gw.LastSeenAt,
		gw.Latitude,
		gw.Longitude,
		gw.Altitude,
		gw.Model,
		gw.Config,
		gw.OsVersion,
		gw.Statistics,
		gw.FirmwareHash[:])
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	gw.UpdatedAt = now
	log.WithFields(log.Fields{
		"id":     gw.MAC,
		"name":   gw.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway updated")
	return nil
}

// UpdateFirstHeartbeat updates the first heartbeat by mac
func (h *GWHandler) UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	res, err := h.db.Exec(`
		update gateway
			set first_heartbeat = $1
		where
			mac = $2`,
		time,
		mac,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update first heartbeat error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	return nil
}

// UpdateLastHeartbeat updates the last heartbeat by mac
func (h *GWHandler) UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	res, err := h.db.Exec(`
		update gateway
			set last_heartbeat = $1
		where
			mac = $2`,
		time,
		mac,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update last heartbeat error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	return nil
}

func (h *GWHandler) SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error {
	res, err := h.db.Exec(`
		update gateway
			set auto_update_firmware = $1
		where
			mac = $2`,
		autoUpdateFirmware,
		mac[:],
	)
	if err != nil {
		return handlePSQLError(Update, err, "update auto_update_firmware error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	return nil
}

// DeleteGateway deletes the gateway matching the given MAC.
func (h *GWHandler) DeleteGateway(ctx context.Context, mac lorawan.EUI64) error {
	n, err := storage.GetNetworkServerForGatewayMAC(ctx, h.db, mac)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// if the gateway is MatchX gateway, unregister it from provisioning server
	obj, err := h.GetGateway(ctx, mac, false)
	if err != nil {
		return errors.Wrap(err, "get gateway error")
	}
	if strings.HasPrefix(obj.Model, "MX") {
		provConf := config.C.ProvisionServer
		provClient, err := provisionserver.CreateClientWithCert(provConf.ProvisionServer, provConf.CACert,
			provConf.TLSCert, provConf.TLSKey)
		if err != nil {
			return errors.Wrap(err, "failed to connect to provisioning server")
		}

		_, err = provClient.UnregisterGw(context.Background(), &psPb.UnregisterGwRequest{
			Sn:  obj.SerialNumber,
			Mac: obj.MAC.String(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to unregister from provisioning server")
		}
	}

	res, err := h.db.Exec("delete from gateway where mac = $1", mac[:])
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = nsClient.DeleteGateway(ctx, &ns.DeleteGatewayRequest{
		Id: mac[:],
	})
	if err != nil && grpc.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete gateway error")
	}

	// delete this gateway from m2m-server
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	gwClient := m2m_api.NewM2MServerServiceClient(m2mClient)
	if err == nil {
		_, err = gwClient.DeleteGatewayInM2MServer(context.Background(), &m2m_api.DeleteGatewayInM2MServerRequest{
			MacAddress: mac.String(),
		})
		if err != nil && grpc.Code(err) != codes.NotFound {
			log.WithError(err).Error("delete gateway from m2m-server error")
		}
	} else {
		log.WithError(err).Error("get m2m-server client error")
	}

	log.WithFields(log.Fields{
		"id":     mac,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway deleted")
	return nil
}

// GetGateway returns the gateway for the given mac.
func (h *GWHandler) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gwmod.Gateway, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var gw gwmod.Gateway
	err := sqlx.Get(h.db, &gw, "select * from gateway where mac = $1"+fu, mac[:])
	if err != nil {
		if err == sql.ErrNoRows {
			return gw, ErrDoesNotExist
		}
		return gw, err
	}
	return gw, nil
}

// GetGatewayCount returns the total number of gateways.
func (h *GWHandler) GetGatewayCount(ctx context.Context, search string) (int, error) {
	var count int
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Get(h.db, &count, `
		select
			count(*)
		from gateway
		where
			$1 = ''
			or (
				$1 != ''
				and (
					name ilike $1
					or encode(mac, 'hex') ilike $1
				)
			)
		`,
		search,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetGateways returns a slice of gateways sorted by name.
func (h *GWHandler) GetGateways(ctx context.Context, limit, offset int32, search string) ([]gwmod.Gateway, error) {
	var gws []gwmod.Gateway
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Select(h.db, &gws, `
		select
			*
		from gateway
		where
			$3 = ''
			or (
				$3 != ''
				and (
					name ilike $3
					or encode(mac, 'hex') ilike $3
				)
			)
		order by
			name
		limit $1 offset $2`,
		limit,
		offset,
		search,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return gws, nil
}

func (h *GWHandler) GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error) {
	var gwConfig string
	err := sqlx.Get(h.db, &gwConfig, `
		select
			config
		from gateway
		where mac = $1`,
		mac[:],
	)
	if err != nil {
		return "", errors.Wrap(err, "select error")
	}

	return gwConfig, nil
}

// GetFirstHeartbeat returns the first heartbeat
func (h *GWHandler) GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	var firstHeartbeat int64
	err := sqlx.Get(h.db, &firstHeartbeat, `
		select 
			first_heartbeat
		from gateway
		where mac = $1
        limit 1`,
		mac,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return firstHeartbeat, nil
}

func (h *GWHandler) UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error {
	res, err := h.db.Exec(`
		update gateway
			set first_heartbeat = 0
		where
			mac = $1`,
		mac,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update first heartbeat to zero error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return ErrDoesNotExist
	}

	return nil
}

// GetLastHeartbeat returns the last heartbeat
func (h *GWHandler) GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	var lastHeartbeat int64

	err := sqlx.Get(h.db, &lastHeartbeat, `
		select 
			last_heartbeat
		from gateway
		where mac = $1
		limit 1`,
		mac,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}

	return lastHeartbeat, nil
}

func (h *GWHandler) GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error) {
	var macs []lorawan.EUI64

	err := sqlx.Select(h.db, &macs, `
		select 
			mac
		from gateway
		where first_heartbeat not in (0)
        and $1 - first_heartbeat > $2`,
		time, limit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	return macs, nil
}

// GetGatewaysLoc returns a slice of gateways locations.
func (h *GWHandler) GetGatewaysLoc(ctx context.Context, limit int) ([]gwmod.GatewayLocation, error) {
	var gwsLoc []gwmod.GatewayLocation

	err := sqlx.Select(h.db, &gwsLoc, `
		select
			latitude,
			longitude,
			altitude
		from gateway
		where latitude > 0 and longitude > 0
		limit $1`,
		limit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return gwsLoc, nil
}

// GetGatewaysForMACs returns a map of gateways given a slice of MACs.
func (h *GWHandler) GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]gwmod.Gateway, error) {
	out := make(map[lorawan.EUI64]gwmod.Gateway)
	var macsB [][]byte
	for i := range macs {
		macsB = append(macsB, macs[i][:])
	}

	var gws []gwmod.Gateway
	err := sqlx.Select(h.db, &gws, "select * from gateway where mac = any($1)", pq.ByteaArray(macsB))
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	if len(gws) != len(macs) {
		log.WithFields(log.Fields{
			"expected": len(macs),
			"returned": len(gws),
			"ctx_id":   ctx.Value(logging.ContextIDKey),
		}).Warning("requested number of gateways does not match returned")
	}

	for i := range gws {
		out[gws[i].MAC] = gws[i]
	}

	return out, nil
}

// GetGatewayCountForOrganizationID returns the total number of gateways
// given an organization ID.
func (h *GWHandler) GetGatewayCountForOrganizationID(ctx context.Context, organizationID int64, search string) (int, error) {
	var count int
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Get(h.db, &count, `
		select
			count(*)
		from gateway
		where
			organization_id = $1
			and (
				$2 = ''
				or (
					$2 != ''
					and (
						name ilike $2
						or encode(mac, 'hex') ilike $2
					)
				)
			)`,
		organizationID,
		search,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetGatewaysForOrganizationID returns a slice of gateways sorted by name
// for the given organization ID.
func (h *GWHandler) GetGatewaysForOrganizationID(ctx context.Context, organizationID int64, limit, offset int, search string) ([]gwmod.Gateway, error) {
	var gws []gwmod.Gateway
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Select(h.db, &gws, `
		select
			*
		from gateway
		where
			organization_id = $1
			and (
				$4 = ''
				or (
					$4 != ''
					and (
						name ilike $4
						or encode(mac, 'hex') ilike $4
					)
				)
			)
		order by
			name
		limit $2 offset $3`,
		organizationID,
		limit,
		offset,
		search,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return gws, nil
}

// GetGatewayCountForUser returns the total number of gateways to which the
// given user has access.
func (h *GWHandler) GetGatewayCountForUser(ctx context.Context, username string, search string) (int, error) {
	var count int
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Get(h.db, &count, `
		select
			count(g.*)
		from gateway g
		inner join organization o
			on o.id = g.organization_id
		inner join organization_user ou
			on ou.organization_id = o.id
		inner join "user" u
			on u.id = ou.user_id
		where
			u.username = $1
			and (
				$2 = ''
				or (
					$2 != ''
					and (
						g.name ilike $2
						or encode(g.mac, 'hex') ilike $2
					)
				)
			)`,
		username,
		search,
	)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetGatewaysForUser returns a slice of gateways sorted by name to which the
// given user has access.
func (h *GWHandler) GetGatewaysForUser(ctx context.Context, username string, limit, offset int, search string) ([]gwmod.Gateway, error) {
	var gws []gwmod.Gateway
	if search != "" {
		search = "%" + search + "%"
	}

	err := sqlx.Select(h.db, &gws, `
		select
			g.*
		from gateway g
		inner join organization o
			on o.id = g.organization_id
		inner join organization_user ou
			on ou.organization_id = o.id
		inner join "user" u
			on u.id = ou.user_id
		where
			u.username = $1
			and (
				$4 = ''
				or (
					$4 != ''
					and (
						g.name ilike $4
						or encode(g.mac, 'hex') ilike $4
					)
				)
			)
		order by
			g.name
		limit $2 offset $3`,
		username,
		limit,
		offset,
		search,
	)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return gws, nil
}

// CreateGatewayPing creates the given gateway ping.
func (h *GWHandler) CreateGatewayPing(ctx context.Context, ping *gwmod.GatewayPing) error {
	ping.CreatedAt = time.Now()

	err := sqlx.Get(h.db, &ping.ID, `
		insert into gateway_ping (
			created_at,
			gateway_mac,
			frequency,
			dr
		) values ($1, $2, $3, $4)
		returning id`,
		ping.CreatedAt,
		ping.GatewayMAC[:],
		ping.Frequency,
		ping.DR,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	log.WithFields(log.Fields{
		"gateway_mac": ping.GatewayMAC,
		"frequency":   ping.Frequency,
		"dr":          ping.DR,
		"id":          ping.ID,
		"ctx_id":      ctx.Value(logging.ContextIDKey),
	}).Info("gateway ping created")

	return nil
}

// GetGatewayPing returns the ping matching the given id.
func (h *GWHandler) GetGatewayPing(ctx context.Context, id int64) (gwmod.GatewayPing, error) {
	var ping gwmod.GatewayPing
	err := sqlx.Get(h.db, &ping, "select * from gateway_ping where id = $1", id)
	if err != nil {
		return ping, handlePSQLError(Select, err, "select error")
	}

	return ping, nil
}

// CreateGatewayPingRX creates the received ping.
func (h *GWHandler) CreateGatewayPingRX(ctx context.Context, rx *gwmod.GatewayPingRX) error {
	rx.CreatedAt = time.Now()

	err := sqlx.Get(h.db, &rx.ID, `
		insert into gateway_ping_rx (
			ping_id,
			created_at,
			gateway_mac,
			received_at,
			rssi,
			lora_snr,
			location,
			altitude
		) values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning id`,
		rx.PingID,
		rx.CreatedAt,
		rx.GatewayMAC[:],
		rx.ReceivedAt,
		rx.RSSI,
		rx.LoRaSNR,
		rx.Location,
		rx.Altitude,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	return nil
}

// DeleteAllGatewaysForOrganizationID deletes all gateways for a given
// organization id.
func (h *GWHandler) DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error {
	var gws []gwmod.Gateway
	err := sqlx.Select(h.db, &gws, "select * from gateway where organization_id = $1", organizationID)
	if err != nil {
		return handlePSQLError(Select, err, "select error")
	}

	for _, gw := range gws {
		err = h.DeleteGateway(ctx, gw.MAC)
		if err != nil {
			return errors.Wrap(err, "delete gateway error")
		}
	}

	return nil
}

// GetAllGatewayMacList get a list of all gateway mac
func (h *GWHandler) GetAllGatewayMacList(ctx context.Context) ([]string, error) {
	var gwMacList []string
	var list []lorawan.EUI64
	err := sqlx.Select(h.db, &list, `select mac from gateway order by created_at desc`)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}

	for _, gwMac := range list {
		gwMacList = append(gwMacList, gwMac.String())
	}
	return gwMacList, nil
}

// GetGatewayPingRXForPingID returns the received gateway pings for the given
// ping ID.
func (h *GWHandler) GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]gwmod.GatewayPingRX, error) {
	var rx []gwmod.GatewayPingRX

	err := sqlx.Select(h.db, &rx, "select * from gateway_ping_rx where ping_id = $1", pingID)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return rx, nil
}

// GetLastGatewayPingAndRX returns the last gateway ping and RX for the given
// gateway MAC.
func (h *GWHandler) GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (gwmod.GatewayPing, []gwmod.GatewayPingRX, error) {
	gw, err := h.GetGateway(ctx, mac, false)
	if err != nil {
		return gwmod.GatewayPing{}, nil, errors.Wrap(err, "get gateway error")
	}

	if gw.LastPingID == nil {
		return gwmod.GatewayPing{}, nil, ErrDoesNotExist
	}

	ping, err := h.GetGatewayPing(ctx, *gw.LastPingID)
	if err != nil {
		return gwmod.GatewayPing{}, nil, errors.Wrap(err, "get gateway ping error")
	}

	rx, err := h.GetGatewayPingRXForPingID(ctx, ping.ID)
	if err != nil {
		return gwmod.GatewayPing{}, nil, errors.Wrap(err, "get gateway ping rx for ping id error")
	}

	return ping, rx, nil
}
