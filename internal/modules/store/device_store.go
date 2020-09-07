package store

import (
	"context"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/lib/pq/hstore"
)

type DeviceStore interface {
	CreateDevice(ctx context.Context, d *Device, applicationServerID uuid.UUID) error
	GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (Device, error)
	GetDeviceCount(ctx context.Context, filters DeviceFilters) (int, error)
	GetAllDeviceEuis(ctx context.Context) ([]string, error)
	GetDevices(ctx context.Context, filters DeviceFilters) ([]DeviceListItem, error)
	UpdateDevice(ctx context.Context, d *Device) error
	DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error
	CreateDeviceKeys(ctx context.Context, dc *DeviceKeys) error
	GetDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) (DeviceKeys, error)
	UpdateDeviceKeys(ctx context.Context, dc *DeviceKeys) error
	DeleteDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) error
	CreateDeviceActivation(ctx context.Context, da *DeviceActivation) error
	GetLastDeviceActivationForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (DeviceActivation, error)
	DeleteAllDevicesForApplicationID(ctx context.Context, applicationID int64) error
	UpdateDeviceActivation(ctx context.Context, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error

	// validator
	CheckCreateNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error)
	CheckListNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error)

	CheckReadNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckDeleteNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckCreateListDeleteDeviceQueueAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
}

func (h *Handler) CreateDevice(ctx context.Context, d *Device, applicationServerID uuid.UUID) error {
	return h.store.CreateDevice(ctx, d, applicationServerID)
}
func (h *Handler) GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (Device, error) {
	return h.store.GetDevice(ctx, devEUI, forUpdate)
}
func (h *Handler) GetDeviceCount(ctx context.Context, filters DeviceFilters) (int, error) {
	return h.store.GetDeviceCount(ctx, filters)
}
func (h *Handler) GetAllDeviceEuis(ctx context.Context) ([]string, error) {
	return h.store.GetAllDeviceEuis(ctx)
}
func (h *Handler) GetDevices(ctx context.Context, filters DeviceFilters) ([]DeviceListItem, error) {
	return h.store.GetDevices(ctx, filters)
}
func (h *Handler) UpdateDevice(ctx context.Context, d *Device) error {
	return h.store.UpdateDevice(ctx, d)
}
func (h *Handler) DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error {
	return h.store.DeleteDevice(ctx, devEUI)
}
func (h *Handler) CreateDeviceKeys(ctx context.Context, dc *DeviceKeys) error {
	return h.store.CreateDeviceKeys(ctx, dc)
}
func (h *Handler) GetDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) (DeviceKeys, error) {
	return h.store.GetDeviceKeys(ctx, devEUI)
}
func (h *Handler) UpdateDeviceKeys(ctx context.Context, dc *DeviceKeys) error {
	return h.store.UpdateDeviceKeys(ctx, dc)
}
func (h *Handler) DeleteDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) error {
	return h.store.DeleteDeviceKeys(ctx, devEUI)
}
func (h *Handler) CreateDeviceActivation(ctx context.Context, da *DeviceActivation) error {
	return h.store.CreateDeviceActivation(ctx, da)
}
func (h *Handler) GetLastDeviceActivationForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (DeviceActivation, error) {
	return h.store.GetLastDeviceActivationForDevEUI(ctx, devEUI)
}
func (h *Handler) DeleteAllDevicesForApplicationID(ctx context.Context, applicationID int64) error {
	return h.store.DeleteAllDevicesForApplicationID(ctx, applicationID)
}
func (h *Handler) UpdateDeviceActivation(ctx context.Context, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error {
	return h.store.UpdateDeviceActivation(ctx, devEUI, devAddr, appSKey)
}

func (h *Handler) CheckCreateNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error) {
	return h.store.CheckCreateNodeAccess(ctx, username, applicationID, userID)
}
func (h *Handler) CheckListNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error) {
	return h.store.CheckListNodeAccess(ctx, username, applicationID, userID)
}

func (h *Handler) CheckReadNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckReadNodeAccess(ctx, username, devEUI, userID)
}
func (h *Handler) CheckUpdateNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckUpdateNodeAccess(ctx, username, devEUI, userID)
}
func (h *Handler) CheckDeleteNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckDeleteNodeAccess(ctx, username, devEUI, userID)
}
func (h *Handler) CheckCreateListDeleteDeviceQueueAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckCreateListDeleteDeviceQueueAccess(ctx, username, devEUI, userID)
}

// DeviceFilters provide filters that can be used to filter on devices.
// Note that empty values are not used as filter.
type DeviceFilters struct {
	ApplicationID    int64     `db:"application_id"`
	MulticastGroupID uuid.UUID `db:"multicast_group_id"`
	ServiceProfileID uuid.UUID `db:"service_profile_id"`
	Search           string    `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filter.
func (f DeviceFilters) SQL() string {
	var filters []string

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
