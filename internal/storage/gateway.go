package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"
	"strings"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
)

var gatewayNameRegexp = regexp.MustCompile(`^[\w-]+$`)

// Gateway represents a gateway.
type Gateway store.Gateway

// GatewayListItem defines the gateway as list item.
type GatewayListItem store.GatewayListItem

// GatewayPing represents a gateway ping.
type GatewayPing store.GatewayPing

// GatewayPingRX represents a ping received by one of the gateways.
type GatewayPingRX store.GatewayPingRX

// GPSPoint contains a GPS point.
type GPSPoint store.GPSPoint

// GatewaysActiveInactive holds the avtive and inactive counts.
type GatewaysActiveInactive store.GatewaysActiveInactive

// Value implements the driver.Valuer interface.
func (l GPSPoint) Value() (driver.Value, error) {
	return store.GPSPoint(l).Value()
}

// Scan implements the sql.Scanner interface.
func (l *GPSPoint) Scan(src interface{}) error {
	return (*store.GPSPoint)(l).Scan(src)
}

// Validate validates the gateway data.
func (g Gateway) Validate() error {
	return store.Gateway(g).Validate()
}

// CreateGateway creates the given Gateway.
func CreateGateway(ctx context.Context, handler *store.Handler, gw *Gateway) error {
	return handler.CreateGateway(ctx, (*store.Gateway)(gw))
}

// UpdateGateway updates the given Gateway.
func UpdateGateway(ctx context.Context, handler *store.Handler, gw *Gateway) error {
	if err := gw.Validate(); err != nil {
		return errors.Wrap(err, "validate error")
	}

	gw.UpdatedAt = time.Now()

	res, err := db.Exec(`
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
			tags = $16,
			metadata = $17
		where
			mac = $1`,
		gw.MAC[:],
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
		gw.Tags,
		gw.Metadata,
	)
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

	log.WithFields(log.Fields{
		"id":     gw.MAC,
		"name":   gw.Name,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway updated")
	return nil
}

// DeleteGateway deletes the gateway matching the given MAC.
func DeleteGateway(ctx context.Context, handler *store.Handler, mac lorawan.EUI64) error {
	n, err := GetNetworkServerForGatewayMAC(ctx, db, mac)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	res, err := db.Exec("delete from gateway where mac = $1", mac[:])
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

	log.WithFields(log.Fields{
		"id":     mac,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway deleted")
	return nil
}

// GetGateway returns the gateway for the given mac.
func GetGateway(ctx context.Context, handler *store.Handler, mac lorawan.EUI64, forUpdate bool) (Gateway, error) {
	var fu string
	if forUpdate {
		fu = " for update"
	}

	var gw Gateway
	err := sqlx.Get(db, &gw, "select * from gateway where mac = $1"+fu, mac[:])
	if err != nil {
		if err == sql.ErrNoRows {
			return gw, ErrDoesNotExist
		}
	}
	return gw, nil
}

// GatewayFilters provides filters for filtering gateways.
type GatewayFilters struct {
	OrganizationID int64  `db:"organization_id"`
	UserID         int64  `db:"user_id"`
	Search         string `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f GatewayFilters) SQL() string {
	var filters []string

	if f.OrganizationID != 0 {
		filters = append(filters, "g.organization_id = :organization_id")
	}

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if f.Search != "" {
		filters = append(filters, "(g.name ilike :search or encode(g.mac, 'hex') ilike :search)")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// GetGatewayCount returns the total number of gateways.
func GetGatewayCount(ctx context.Context, handler *store.Handler, filters GatewayFilters) (int, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			count(distinct g.*)
		from
			gateway g
		inner join organization o
			on o.id = g.organization_id
		left join organization_user ou
			on o.id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL(), filters)
	if err != nil {
		return 0, errors.Wrap(err, "named query error")
	}

	var count int
	err = sqlx.Get(db, &count, query, args...)
	if err != nil {

		return 0, errors.Wrap(err, "named query error")
	}

	return count, nil
}

// GetGateways returns a slice of gateways sorted by name.
func GetGateways(ctx context.Context, handler *store.Handler, filters GatewayFilters) ([]GatewayListItem, error) {
	if filters.Search != "" {
		filters.Search = "%" + filters.Search + "%"
	}

	query, args, err := sqlx.BindNamed(sqlx.DOLLAR, `
		select
			distinct g.mac,
			g.name,
			g.description,
			g.created_at,
			g.updated_at,
			g.first_seen_at,
			g.last_seen_at,
			g.organization_id,
			g.network_server_id,
			g.latitude,
			g.longitude,
			g.altitude,
			n.name as network_server_name
		from
			gateway g
		inner join organization o
			on o.id = g.organization_id
		inner join network_server n
			on n.id = g.network_server_id
		left join organization_user ou
			on o.id = ou.organization_id
		left join "user" u
			on ou.user_id = u.id
	`+filters.SQL()+`
		order by
			g.name
		limit :limit
		offset :offset
	`, filters)

	var gws []GatewayListItem
	err = sqlx.Select(db, &gws, query, args...)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return gws, nil
}

// GetGatewaysForMACs returns a map of gateways given a slice of MACs.
func GetGatewaysForMACs(ctx context.Context, handler *store.Handler, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error) {
	out := make(map[lorawan.EUI64]Gateway)
	var macsB [][]byte
	for i := range macs {
		macsB = append(macsB, macs[i][:])
	}

	var gws []Gateway
	err := sqlx.Select(db, &gws, "select * from gateway where mac = any($1)", pq.ByteaArray(macsB))
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

// CreateGatewayPing creates the given gateway ping.
func CreateGatewayPing(ctx context.Context, handler *store.Handler, ping *GatewayPing) error {
	ping.CreatedAt = time.Now()

	err := sqlx.Get(db, &ping.ID, `
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
func GetGatewayPing(ctx context.Context, handler *store.Handler, id int64) (GatewayPing, error) {
	var ping GatewayPing
	err := sqlx.Get(db, &ping, "select * from gateway_ping where id = $1", id)
	if err != nil {
		return ping, handlePSQLError(Select, err, "select error")
	}

	return ping, nil
}

// CreateGatewayPingRX creates the received ping.
func CreateGatewayPingRX(ctx context.Context, handler *store.Handler, rx *GatewayPingRX) error {
	rx.CreatedAt = time.Now()

	err := sqlx.Get(db, &rx.ID, `
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
func DeleteAllGatewaysForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) error {
	var gws []Gateway
	err := sqlx.Select(db, &gws, "select * from gateway where organization_id = $1", organizationID)
	if err != nil {
		return handlePSQLError(Select, err, "select error")
	}

	for _, gw := range gws {
		err = DeleteGateway(ctx, db, gw.MAC)
		if err != nil {
			return errors.Wrap(err, "delete gateway error")
		}
	}

	return nil
}

// GetGatewayPingRXForPingID returns the received gateway pings for the given
// ping ID.
func GetGatewayPingRXForPingID(ctx context.Context, handler *store.Handler, pingID int64) ([]GatewayPingRX, error) {
	var rx []GatewayPingRX

	err := sqlx.Select(db, &rx, "select * from gateway_ping_rx where ping_id = $1", pingID)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return rx, nil
}

// GetLastGatewayPingAndRX returns the last gateway ping and RX for the given
// gateway MAC.
func GetLastGatewayPingAndRX(ctx context.Context, handler *store.Handler, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error) {
	gw, err := GetGateway(ctx, db, mac, false)
	if err != nil {
		return GatewayPing{}, nil, errors.Wrap(err, "get gateway error")
	}

	if gw.LastPingID == nil {
		return GatewayPing{}, nil, ErrDoesNotExist
	}

	ping, err := GetGatewayPing(ctx, db, *gw.LastPingID)
	if err != nil {
		return GatewayPing{}, nil, errors.Wrap(err, "get gateway ping error")
	}

	rx, err := GetGatewayPingRXForPingID(ctx, db, ping.ID)
	if err != nil {
		return GatewayPing{}, nil, errors.Wrap(err, "get gateway ping rx for ping id error")
	}

	return ping, rx, nil
}

// GetGatewaysActiveInactive returns the active / inactive gateways.
func GetGatewaysActiveInactive(ctx context.Context, handler *store.Handler, organizationID int64) (GatewaysActiveInactive, error) {
	var out GatewaysActiveInactive
	err := sqlx.Get(db, &out, `
		with gateway_active_inactive as (
			select
				g.last_seen_at as last_seen_at,
				make_interval(secs => coalesce(gp.stats_interval / 1000000000, 30)) * 1.5 as stats_interval
			from
				gateway g
			left join gateway_profile gp
				on g.gateway_profile_id = gp.gateway_profile_id
			where
				$1 = 0 or g.organization_id = $1
		)
		select
			coalesce(sum(case when last_seen_at is null then 1 end), 0) as never_seen_count,
			coalesce(sum(case when (now() - stats_interval) > last_seen_at then 1 end), 0) as inactive_count,
			coalesce(sum(case when (now() - stats_interval) <= last_seen_at then 1 end), 0) as active_count
		from
			gateway_active_inactive
	`, organizationID)
	if err != nil {
		return out, errors.Wrap(err, "get gateway active/inactive count error")
	}

	return out, nil
}
