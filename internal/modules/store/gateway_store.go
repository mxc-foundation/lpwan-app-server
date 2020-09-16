package store

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq/hstore"

	"github.com/brocaar/lorawan"

	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type GatewayStore interface {
	AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error
	GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error

	AddGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw GatewayFirmware, err error)
	GetGatewayFirmwareList(ctx context.Context) (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error
	CreateGateway(ctx context.Context, gw *Gateway) error
	UpdateGateway(ctx context.Context, gw *Gateway) error
	UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error
	DeleteGateway(ctx context.Context, mac lorawan.EUI64) error
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error)
	GetGatewayCount(ctx context.Context, filters GatewayFilters) (int, error)
	GetGateways(ctx context.Context, filters GatewayFilters) ([]GatewayListItem, error)
	GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error)
	GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error
	GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error)
	GetGatewaysLoc(ctx context.Context, limit int) ([]GatewayLocation, error)
	GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error)
	CreateGatewayPing(ctx context.Context, ping *GatewayPing) error
	GetGatewayPing(ctx context.Context, id int64) (GatewayPing, error)
	CreateGatewayPingRX(ctx context.Context, rx *GatewayPingRX) error
	DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error
	GetAllGatewayMacList(ctx context.Context) ([]string, error)
	GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]GatewayPingRX, error)
	GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error)
	GetGatewaysActiveInactive(ctx context.Context, organizationID int64) (GatewaysActiveInactive, error)

	GetGatewayForPing(ctx context.Context) (*Gateway, error)

	// validator
	CheckCreateGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)

	CheckReadGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateDeleteGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error)

	CheckReadOrganizationNetworkServerAccess(ctx context.Context, username string, organizationID, networkserverID, userID int64) (bool, error)
}

func (h *Handler) GetGatewayForPing(ctx context.Context) (*Gateway, error) {
	return h.store.GetGatewayForPing(ctx)
}
func (h *Handler) GetGatewaysActiveInactive(ctx context.Context, organizationID int64) (GatewaysActiveInactive, error) {
	return h.store.GetGatewaysActiveInactive(ctx, organizationID)
}

func (h *Handler) AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.store.AddNewDefaultGatewayConfig(ctx, defaultConfig)
}
func (h *Handler) UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.store.UpdateDefaultGatewayConfig(ctx, defaultConfig)
}
func (h *Handler) GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	return h.store.GetDefaultGatewayConfig(ctx, defaultConfig)
}

func (h *Handler) AddGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error) {
	return h.store.AddGatewayFirmware(ctx, gwFw)
}
func (h *Handler) GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw GatewayFirmware, err error) {
	return h.store.GetGatewayFirmware(ctx, model, forUpdate)
}
func (h *Handler) GetGatewayFirmwareList(ctx context.Context) (list []GatewayFirmware, err error) {
	return h.store.GetGatewayFirmwareList(ctx)
}
func (h *Handler) UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error) {
	return h.store.UpdateGatewayFirmware(ctx, gwFw)
}
func (h *Handler) UpdateGatewayConfigByGwId(ctx context.Context, config string, mac lorawan.EUI64) error {
	return h.store.UpdateGatewayConfigByGwId(ctx, config, mac)
}
func (h *Handler) CreateGateway(ctx context.Context, gw *Gateway) error {
	return h.store.CreateGateway(ctx, gw)
}
func (h *Handler) UpdateGateway(ctx context.Context, gw *Gateway) error {
	return h.store.UpdateGateway(ctx, gw)
}
func (h *Handler) UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	return h.store.UpdateFirstHeartbeat(ctx, mac, time)
}
func (h *Handler) UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error {
	return h.store.UpdateLastHeartbeat(ctx, mac, time)
}
func (h *Handler) SetAutoUpdateFirmware(ctx context.Context, mac lorawan.EUI64, autoUpdateFirmware bool) error {
	return h.store.SetAutoUpdateFirmware(ctx, mac, autoUpdateFirmware)
}
func (h *Handler) DeleteGateway(ctx context.Context, mac lorawan.EUI64) error {
	return h.store.DeleteGateway(ctx, mac)
}
func (h *Handler) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error) {
	return h.store.GetGateway(ctx, mac, forUpdate)
}
func (h *Handler) GetGatewayCount(ctx context.Context, filters GatewayFilters) (int, error) {
	return h.store.GetGatewayCount(ctx, filters)
}
func (h *Handler) GetGateways(ctx context.Context, filters GatewayFilters) ([]GatewayListItem, error) {
	return h.store.GetGateways(ctx, filters)
}
func (h *Handler) GetGatewayConfigByGwId(ctx context.Context, mac lorawan.EUI64) (string, error) {
	return h.store.GetGatewayConfigByGwId(ctx, mac)
}
func (h *Handler) GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	return h.store.GetFirstHeartbeat(ctx, mac)
}
func (h *Handler) UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error {
	return h.store.UpdateFirstHeartbeatToZero(ctx, mac)
}
func (h *Handler) GetLastHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error) {
	return h.store.GetLastHeartbeat(ctx, mac)
}
func (h *Handler) GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error) {
	return h.store.GetGatewayMiningList(ctx, time, limit)
}
func (h *Handler) GetGatewaysLoc(ctx context.Context, limit int) ([]GatewayLocation, error) {
	return h.store.GetGatewaysLoc(ctx, limit)
}
func (h *Handler) GetGatewaysForMACs(ctx context.Context, macs []lorawan.EUI64) (map[lorawan.EUI64]Gateway, error) {
	return h.store.GetGatewaysForMACs(ctx, macs)
}
func (h *Handler) CreateGatewayPing(ctx context.Context, ping *GatewayPing) error {
	return h.store.CreateGatewayPing(ctx, ping)
}
func (h *Handler) GetGatewayPing(ctx context.Context, id int64) (GatewayPing, error) {
	return h.store.GetGatewayPing(ctx, id)
}
func (h *Handler) CreateGatewayPingRX(ctx context.Context, rx *GatewayPingRX) error {
	return h.store.CreateGatewayPingRX(ctx, rx)
}
func (h *Handler) DeleteAllGatewaysForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.store.DeleteAllGatewaysForOrganizationID(ctx, organizationID)
}
func (h *Handler) GetAllGatewayMacList(ctx context.Context) ([]string, error) {
	return h.store.GetAllGatewayMacList(ctx)
}
func (h *Handler) GetGatewayPingRXForPingID(ctx context.Context, pingID int64) ([]GatewayPingRX, error) {
	return h.store.GetGatewayPingRXForPingID(ctx, pingID)
}
func (h *Handler) GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (GatewayPing, []GatewayPingRX, error) {
	return h.store.GetLastGatewayPingAndRX(ctx, mac)
}

// validator
func (h *Handler) CheckCreateGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckCreateGatewayAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckListGatewayAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckListGatewayAccess(ctx, username, organizationID, userID)
}

func (h *Handler) CheckReadGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckReadGatewayAccess(ctx, username, mac, userID)
}
func (h *Handler) CheckUpdateDeleteGatewayAccess(ctx context.Context, username string, mac lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteGatewayAccess(ctx, username, mac, userID)
}

func (h *Handler) CheckReadOrganizationNetworkServerAccess(ctx context.Context, username string, organizationID, networkserverID, userID int64) (bool, error) {
	return h.store.CheckReadOrganizationNetworkServerAccess(ctx, username, organizationID, networkserverID, userID)
}

var (
	gatewayNameRegexp          = regexp.MustCompile(`^[\w-]+$`)
	serialNumberOldGWValidator = regexp.MustCompile(`^MX([A-Z1-9]){7}$`)
	serialNumberNewGWValidator = regexp.MustCompile(`^M2X([A-Z1-9]){8}$`)
)

type DefaultGatewayConfig struct {
	ID            int64      `db:"id"`
	Model         string     `db:"model"`
	Region        string     `db:"region"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DefaultConfig string     `db:"default_config"`
}

// Gateway represents a gateway.
type Gateway struct {
	MAC                lorawan.EUI64 `db:"mac"`
	CreatedAt          time.Time     `db:"created_at"`
	UpdatedAt          time.Time     `db:"updated_at"`
	FirstSeenAt        *time.Time    `db:"first_seen_at"`
	LastSeenAt         *time.Time    `db:"last_seen_at"`
	Name               string        `db:"name"`
	Description        string        `db:"description"`
	OrganizationID     int64         `db:"organization_id"`
	Ping               bool          `db:"ping"`
	LastPingID         *int64        `db:"last_ping_id"`
	LastPingSentAt     *time.Time    `db:"last_ping_sent_at"`
	NetworkServerID    int64         `db:"network_server_id"`
	GatewayProfileID   *string       `db:"gateway_profile_id"`
	Latitude           float64       `db:"latitude"`
	Longitude          float64       `db:"longitude"`
	Altitude           float64       `db:"altitude"`
	Tags               hstore.Hstore `db:"tags"`
	Metadata           hstore.Hstore `db:"metadata"`
	Model              string        `db:"model"`
	FirstHeartbeat     int64         `db:"first_heartbeat"`
	LastHeartbeat      int64         `db:"last_heartbeat"`
	Config             string        `db:"config"`
	OsVersion          string        `db:"os_version"`
	Statistics         string        `db:"statistics"`
	SerialNumber       string        `db:"sn"`
	FirmwareHash       types.MD5SUM  `db:"firmware_hash"`
	AutoUpdateFirmware bool          `db:"auto_update_firmware"`
}

// GatewayListItem defines the gateway as list item.
type GatewayListItem struct {
	MAC               lorawan.EUI64 `db:"mac"`
	Name              string        `db:"name"`
	Description       string        `db:"description"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
	FirstSeenAt       *time.Time    `db:"first_seen_at"`
	LastSeenAt        *time.Time    `db:"last_seen_at"`
	OrganizationID    int64         `db:"organization_id"`
	NetworkServerID   int64         `db:"network_server_id"`
	Latitude          float64       `db:"latitude"`
	Longitude         float64       `db:"longitude"`
	Altitude          float64       `db:"altitude"`
	NetworkServerName string        `db:"network_server_name"`
	Model             string        `db:"model"`
	Config            string        `db:"config"`
}

type GatewayFirmware struct {
	Model        string       `db:"model"`
	ResourceLink string       `db:"resource_link"`
	FirmwareHash types.MD5SUM `db:"md5_hash"`
}

// GatewayLocation represents a gateway location.
type GatewayLocation struct {
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
	Altitude  float64 `db:"altitude"`
}

// GatewayPing represents a gateway ping.
type GatewayPing struct {
	ID         int64         `db:"id"`
	CreatedAt  time.Time     `db:"created_at"`
	GatewayMAC lorawan.EUI64 `db:"gateway_mac"`
	Frequency  int           `db:"frequency"`
	DR         int           `db:"dr"`
}

// CreatePingLookup creates an automatically expiring MIC to ping id lookup.
func (gp GatewayPing) CreatePingLookup(mic lorawan.MIC) error {
	redisClient := rs.RedisClientType{}
	return redisClient.RSSet(fmt.Sprintf("%s", mic), gp.ID)
}

func (gp GatewayPing) GetPingLookup(mic lorawan.MIC) (int64, error) {
	redisClient := rs.RedisClientType{}
	return redisClient.RSGet(fmt.Sprintf("%s", mic)).Int64()
}

func (gp GatewayPing) DeletePingLookup(mic lorawan.MIC) error {
	redisClient := rs.RedisClientType{}
	return redisClient.RSDelete(fmt.Sprintf("%s", mic))
}

// GatewayPingRX represents a ping received by one of the gateways.
type GatewayPingRX struct {
	ID         int64         `db:"id"`
	PingID     int64         `db:"ping_id"`
	CreatedAt  time.Time     `db:"created_at"`
	GatewayMAC lorawan.EUI64 `db:"gateway_mac"`
	ReceivedAt *time.Time    `db:"received_at"`
	RSSI       int           `db:"rssi"`
	LoRaSNR    float64       `db:"lora_snr"`
	Location   GPSPoint      `db:"location"`
	Altitude   float64       `db:"altitude"`
}

// GPSPoint contains a GPS point.
type GPSPoint struct {
	Latitude  float64
	Longitude float64
}

// GatewaysActiveInactive holds the avtive and inactive counts.
type GatewaysActiveInactive struct {
	NeverSeenCount uint32 `db:"never_seen_count"`
	ActiveCount    uint32 `db:"active_count"`
	InactiveCount  uint32 `db:"inactive_count"`
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

// Value implements the driver.Valuer interface.
func (l GPSPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", strconv.FormatFloat(l.Latitude, 'f', -1, 64), strconv.FormatFloat(l.Longitude, 'f', -1, 64)), nil
}

// Scan implements the sql.Scanner interface.
func (l *GPSPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}

	_, err := fmt.Sscanf(string(b), "(%f,%f)", &l.Latitude, &l.Longitude)
	return err
}

// Validate validates the gateway data.
func (g Gateway) Validate() error {
	if !gatewayNameRegexp.MatchString(g.Name) {
		return ErrGatewayInvalidName
	}

	if strings.HasPrefix(g.Model, "MX19") {
		if !serialNumberNewGWValidator.MatchString(g.SerialNumber) {
			return errors.New("invalid gateway serial number")
		}
	} else if g.Model != "" {
		if !serialNumberOldGWValidator.MatchString(g.SerialNumber) {
			return errors.New("invalid gateway serial number")
		}
	}

	return nil
}
