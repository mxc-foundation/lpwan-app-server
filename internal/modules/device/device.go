package device

import (
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/brocaar/lorawan"
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
}

type Controller struct {
	St        DeviceStore
	Validator Validator
}

var Service *Controller

func Setup() error {
	Service.St = store.New(storage.DB().DB)
	return nil
}
