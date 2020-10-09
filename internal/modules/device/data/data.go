package data

import (
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/lib/pq/hstore"
)

type Config struct {
	ApplicationServerID string
}

// DeviceFilters provide filters that can be used to filter on devices.
// Note that empty values are not used as filter.
type DeviceFilters struct {
	OrganizationID   int64         `db:"organization_id"`
	ApplicationID    int64         `db:"application_id"`
	MulticastGroupID uuid.UUID     `db:"multicast_group_id"`
	ServiceProfileID uuid.UUID     `db:"service_profile_id"`
	Search           string        `db:"search"`
	Tags             hstore.Hstore `db:"tags"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filter.
func (f DeviceFilters) SQL() string {
	var filters []string

	if f.OrganizationID != 0 {
		filters = append(filters, "a.organization_id = :organization_id")
	}

	if f.ApplicationID != 0 {
		filters = append(filters, "d.application_id = :application_id")
	}

	if f.MulticastGroupID != uuid.Nil {
		filters = append(filters, "dmg.multicast_group_id = :multicast_group_id")
	}

	if f.ServiceProfileID != uuid.Nil {
		filters = append(filters, "a.service_profile_id = :service_profile_id")
	}

	if f.Search != "" {
		filters = append(filters, "(d.name ilike :search or encode(d.dev_eui, 'hex') ilike :search)")
	}

	if len(f.Tags.Map) != 0 {
		filters = append(filters, "d.tags @> :tags")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// Device defines a LoRaWAN device.
type Device struct {
	DevEUI                    lorawan.EUI64     `db:"dev_eui"`
	CreatedAt                 time.Time         `db:"created_at"`
	UpdatedAt                 time.Time         `db:"updated_at"`
	LastSeenAt                *time.Time        `db:"last_seen_at"`
	ApplicationID             int64             `db:"application_id"`
	DeviceProfileID           uuid.UUID         `db:"device_profile_id"`
	Name                      string            `db:"name"`
	Description               string            `db:"description"`
	SkipFCntCheck             bool              `db:"-"`
	ReferenceAltitude         float64           `db:"-"`
	DeviceStatusBattery       *float32          `db:"device_status_battery"`
	DeviceStatusMargin        *int              `db:"device_status_margin"`
	DeviceStatusExternalPower bool              `db:"device_status_external_power_source"`
	DR                        *int              `db:"dr"`
	Latitude                  *float64          `db:"latitude"`
	Longitude                 *float64          `db:"longitude"`
	Altitude                  *float64          `db:"altitude"`
	Variables                 hstore.Hstore     `db:"variables"`
	Tags                      hstore.Hstore     `db:"tags"`
	DevAddr                   lorawan.DevAddr   `db:"dev_addr"`
	AppSKey                   lorawan.AES128Key `db:"app_s_key"`
	IsDisabled                bool              `db:"-"`
}

// DeviceListItem defines the Device as list item.
type DeviceListItem struct {
	Device
	DeviceProfileName string `db:"device_profile_name"`
}

// Validate validates the device data.
func (d Device) Validate() error {
	return nil
}

// DeviceKeys defines the keys for a LoRaWAN device.
type DeviceKeys struct {
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
	DevEUI    lorawan.EUI64     `db:"dev_eui"`
	NwkKey    lorawan.AES128Key `db:"nwk_key"`
	AppKey    lorawan.AES128Key `db:"app_key"`
	GenAppKey lorawan.AES128Key `db:"gen_app_key"`
	JoinNonce int               `db:"join_nonce"`
}

// DeviceActivation defines the device-activation for a LoRaWAN device.
type DeviceActivation struct {
	ID        int64             `db:"id"`
	CreatedAt time.Time         `db:"created_at"`
	DevEUI    lorawan.EUI64     `db:"dev_eui"`
	DevAddr   lorawan.DevAddr   `db:"dev_addr"`
	AppSKey   lorawan.AES128Key `db:"app_s_key"`
}

// DevicesActiveInactive holds the active and inactive counts.
type DevicesActiveInactive struct {
	NeverSeenCount uint32 `db:"never_seen_count"`
	ActiveCount    uint32 `db:"active_count"`
	InactiveCount  uint32 `db:"inactive_count"`
}

type DevicesDataRates map[uint32]uint32
