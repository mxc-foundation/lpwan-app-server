package gateway

import (
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	"github.com/brocaar/lorawan"
)

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
