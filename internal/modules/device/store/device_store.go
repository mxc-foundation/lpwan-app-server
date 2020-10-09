package store

import (
	"context"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	"github.com/brocaar/lorawan"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
)

func NewStore(pg pgstore.PgStore) *ds {
	return &ds{
		pg: pg,
	}
}

type ds struct {
	pg pgstore.DevicePgstore
}

type DeviceStore interface {
	CreateDevice(ctx context.Context, d *Device) error
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
	UpdateDeviceLastSeenAndDR(ctx context.Context, devEUI lorawan.EUI64, ts time.Time, dr int) error
	GetDevicesActiveInactive(ctx context.Context, organizationID int64) (DevicesActiveInactive, error)
	GetDevicesDataRates(ctx context.Context, organizationID int64) (DevicesDataRates, error)

	// validator
	CheckCreateNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error)
	CheckListNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error)

	CheckReadNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckUpdateNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckDeleteNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
	CheckCreateListDeleteDeviceQueueAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error)
}

func (h *ds) GetDevicesDataRates(ctx context.Context, organizationID int64) (DevicesDataRates, error) {
	return h.pg.GetDevicesDataRates(ctx, organizationID)
}

func (h *ds) GetDevicesActiveInactive(ctx context.Context, organizationID int64) (DevicesActiveInactive, error) {
	return h.pg.GetDevicesActiveInactive(ctx, organizationID)
}

func (h *ds) UpdateDeviceLastSeenAndDR(ctx context.Context, devEUI lorawan.EUI64, ts time.Time, dr int) error {
	return h.pg.UpdateDeviceLastSeenAndDR(ctx, devEUI, ts, dr)
}

func (h *ds) CreateDevice(ctx context.Context, d *Device) error {
	return h.pg.CreateDevice(ctx, d)
}
func (h *ds) GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (Device, error) {
	return h.pg.GetDevice(ctx, devEUI, forUpdate)
}
func (h *ds) GetDeviceCount(ctx context.Context, filters DeviceFilters) (int, error) {
	return h.pg.GetDeviceCount(ctx, filters)
}
func (h *ds) GetAllDeviceEuis(ctx context.Context) ([]string, error) {
	return h.pg.GetAllDeviceEuis(ctx)
}
func (h *ds) GetDevices(ctx context.Context, filters DeviceFilters) ([]DeviceListItem, error) {
	return h.pg.GetDevices(ctx, filters)
}
func (h *ds) UpdateDevice(ctx context.Context, d *Device) error {
	return h.pg.UpdateDevice(ctx, d)
}
func (h *ds) DeleteDevice(ctx context.Context, devEUI lorawan.EUI64) error {
	return h.pg.DeleteDevice(ctx, devEUI)
}
func (h *ds) CreateDeviceKeys(ctx context.Context, dc *DeviceKeys) error {
	return h.pg.CreateDeviceKeys(ctx, dc)
}
func (h *ds) GetDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) (DeviceKeys, error) {
	return h.pg.GetDeviceKeys(ctx, devEUI)
}
func (h *ds) UpdateDeviceKeys(ctx context.Context, dc *DeviceKeys) error {
	return h.pg.UpdateDeviceKeys(ctx, dc)
}
func (h *ds) DeleteDeviceKeys(ctx context.Context, devEUI lorawan.EUI64) error {
	return h.pg.DeleteDeviceKeys(ctx, devEUI)
}
func (h *ds) CreateDeviceActivation(ctx context.Context, da *DeviceActivation) error {
	return h.pg.CreateDeviceActivation(ctx, da)
}
func (h *ds) GetLastDeviceActivationForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (DeviceActivation, error) {
	return h.pg.GetLastDeviceActivationForDevEUI(ctx, devEUI)
}
func (h *ds) DeleteAllDevicesForApplicationID(ctx context.Context, applicationID int64) error {
	return h.pg.DeleteAllDevicesForApplicationID(ctx, applicationID)
}
func (h *ds) UpdateDeviceActivation(ctx context.Context, devEUI lorawan.EUI64, devAddr lorawan.DevAddr, appSKey lorawan.AES128Key) error {
	return h.pg.UpdateDeviceActivation(ctx, devEUI, devAddr, appSKey)
}

func (h *ds) CheckCreateNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error) {
	return h.pg.CheckCreateNodeAccess(ctx, username, applicationID, userID)
}
func (h *ds) CheckListNodeAccess(ctx context.Context, username string, applicationID int64, userID int64) (bool, error) {
	return h.pg.CheckListNodeAccess(ctx, username, applicationID, userID)
}

func (h *ds) CheckReadNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckReadNodeAccess(ctx, username, devEUI, userID)
}
func (h *ds) CheckUpdateNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckUpdateNodeAccess(ctx, username, devEUI, userID)
}
func (h *ds) CheckDeleteNodeAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckDeleteNodeAccess(ctx, username, devEUI, userID)
}
func (h *ds) CheckCreateListDeleteDeviceQueueAccess(ctx context.Context, username string, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckCreateListDeleteDeviceQueueAccess(ctx, username, devEUI, userID)
}
