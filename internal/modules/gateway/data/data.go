package data

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/lib/pq/hstore"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type GatewayBindStruct struct {
	NewGateway struct {
		Bind    string `mapstructure:"new_gateway_bind"`
		CACert  string `mapstructure:"ecc_ca_cert"`
		TLSCert string `mapstructure:"ecc_tls_cert"`
		TLSKey  string `mapstructure:"ecc_tls_key"`
	} `mapstructure:"new_gateway"`

	OldGateway struct {
		Bind    string `mapstructure:"old_gateway_bind"`
		CACert  string `mapstructure:"rsa_ca_cert"`
		TLSCert string `mapstructure:"rsa_tls_cert"`
		TLSKey  string `mapstructure:"rsa_tls_key"`
	} `mapstructure:"old_gateway"`
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
	STCOrgID           *int64        `db:"stc_org_id"`
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
	STCOrgID          *int64        `db:"stc_org_id"`
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
		return errHandler.ErrGatewayInvalidName
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
