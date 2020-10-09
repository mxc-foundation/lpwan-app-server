package storage

import (
	"context"

	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"

	fds "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// FUOTA deployment states.
const (
	FUOTADeploymentMulticastCreate        = fds.FUOTADeploymentMulticastCreate
	FUOTADeploymentMulticastSetup         = fds.FUOTADeploymentMulticastSetup
	FUOTADeploymentFragmentationSessSetup = fds.FUOTADeploymentFragmentationSessSetup
	FUOTADeploymentMulticastSessCSetup    = fds.FUOTADeploymentMulticastSessCSetup
	FUOTADeploymentEnqueue                = fds.FUOTADeploymentEnqueue
	FUOTADeploymentStatusRequest          = fds.FUOTADeploymentStatusRequest
	FUOTADeploymentSetDeviceStatus        = fds.FUOTADeploymentSetDeviceStatus
	FUOTADeploymentCleanup                = fds.FUOTADeploymentCleanup
	FUOTADeploymentDone                   = fds.FUOTADeploymentDone
)

// FUOTA deployment device states.
const (
	FUOTADeploymentDevicePending = fds.FUOTADeploymentDevicePending
	FUOTADeploymentDeviceSuccess = fds.FUOTADeploymentDeviceSuccess
	FUOTADeploymentDeviceError   = fds.FUOTADeploymentDeviceError
)

// FUOTA deployment group types.
const (
	FUOTADeploymentGroupTypeB = fds.FUOTADeploymentGroupTypeB
	FUOTADeploymentGroupTypeC = fds.FUOTADeploymentGroupTypeC
)

// FUOTADeployment defiles a firmware update over the air deployment.
type FUOTADeployment fds.FUOTADeployment

// FUOTADeploymentListItem defines a FUOTA deployment item for listing.
type FUOTADeploymentListItem fds.FUOTADeploymentListItem

// FUOTADeploymentDevice defines the device record of a FUOTA deployment.
type FUOTADeploymentDevice fds.FUOTADeploymentDevice

// FUOTADeploymentDeviceListItem defines the Device as FUOTA deployment list item.
type FUOTADeploymentDeviceListItem fds.FUOTADeploymentDeviceListItem

// FUOTADeploymentFilters provides filters that can be used to filter on
// FUOTA deployments. Note that empty values are not used as filters.
type FUOTADeploymentFilters fds.FUOTADeploymentFilters

// SQL returns the SQL filter.
func (f FUOTADeploymentFilters) SQL() string {
	return fds.FUOTADeploymentFilters(f).SQL()
}

// Validate validates the FUOTADeployment data.
func (fd FUOTADeployment) Validate() error {
	return fds.FUOTADeployment(fd).Validate()
}

// CreateFUOTADeploymentForDevice creates and initializes a FUOTA deployment
// for the given device.
func CreateFUOTADeploymentForDevice(ctx context.Context, handler *store.Handler, fd *FUOTADeployment, devEUI lorawan.EUI64) error {
	return handler.CreateFUOTADeploymentForDevice(ctx, (*fds.FUOTADeployment)(fd), devEUI)
}

// GetFUOTADeployment returns the FUOTA deployment for the given ID.
func GetFUOTADeployment(ctx context.Context, handler *store.Handler, id uuid.UUID, forUpdate bool) (FUOTADeployment, error) {
	res, err := handler.GetFUOTADeployment(ctx, id, forUpdate)
	return FUOTADeployment(res), err
}

// GetPendingFUOTADeployments returns the pending FUOTA deployments.
func GetPendingFUOTADeployments(ctx context.Context, handler *store.Handler, batchSize int) ([]FUOTADeployment, error) {
	res, err := handler.GetPendingFUOTADeployments(ctx, batchSize)
	if err != nil {
		return nil, err
	}

	var fuotaList []FUOTADeployment
	for _, v := range res {
		fuotaItem := FUOTADeployment(v)
		fuotaList = append(fuotaList, fuotaItem)
	}

	return fuotaList, nil
}

// UpdateFUOTADeployment updates the given FUOTA deployment.
func UpdateFUOTADeployment(ctx context.Context, handler *store.Handler, fd *FUOTADeployment) error {
	return handler.UpdateFUOTADeployment(ctx, (*fds.FUOTADeployment)(fd))
}

// GetFUOTADeploymentCount returns the number of FUOTA deployments.
func GetFUOTADeploymentCount(ctx context.Context, handler *store.Handler, filters FUOTADeploymentFilters) (int, error) {
	return handler.GetFUOTADeploymentCount(ctx, fds.FUOTADeploymentFilters(filters))
}

// GetFUOTADeployments returns a slice of fuota deployments.
func GetFUOTADeployments(ctx context.Context, handler *store.Handler, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error) {
	res, err := handler.GetFUOTADeployments(ctx, fds.FUOTADeploymentFilters(filters))
	if err != nil {
		return nil, err
	}

	var fuotaList []FUOTADeploymentListItem
	for _, v := range res {
		fuotaItem := FUOTADeploymentListItem(v)
		fuotaList = append(fuotaList, fuotaItem)
	}

	return fuotaList, nil
}

// GetFUOTADeploymentDevice returns the FUOTA deployment record for the given
// device.
func GetFUOTADeploymentDevice(ctx context.Context, handler *store.Handler, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	res, err := handler.GetFUOTADeploymentDevice(ctx, fuotaDeploymentID, devEUI)
	return FUOTADeploymentDevice(res), err
}

// GetPendingFUOTADeploymentDevice returns the pending FUOTA deployment record
// for the given DevEUI.
func GetPendingFUOTADeploymentDevice(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	res, err := handler.GetPendingFUOTADeploymentDevice(ctx, devEUI)
	return FUOTADeploymentDevice(res), err
}

// UpdateFUOTADeploymentDevice updates the given fuota deployment device record.
func UpdateFUOTADeploymentDevice(ctx context.Context, handler *store.Handler, fdd *FUOTADeploymentDevice) error {
	return handler.UpdateFUOTADeploymentDevice(ctx, (*fds.FUOTADeploymentDevice)(fdd))
}

// GetFUOTADeploymentDeviceCount returns the device count for the given
// FUOTA deployment ID.
func GetFUOTADeploymentDeviceCount(ctx context.Context, handler *store.Handler, fuotaDeploymentID uuid.UUID) (int, error) {
	return handler.GetFUOTADeploymentDeviceCount(ctx, fuotaDeploymentID)
}

// GetFUOTADeploymentDevices returns a slice of devices for the given FUOTA
// deployment ID.
func GetFUOTADeploymentDevices(ctx context.Context, handler *store.Handler, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error) {
	res, err := handler.GetFUOTADeploymentDevices(ctx, fuotaDeploymentID, limit, offset)
	if err != nil {
		return nil, err
	}

	var fuotaList []FUOTADeploymentDeviceListItem
	for _, v := range res {
		fuotaItem := FUOTADeploymentDeviceListItem(v)
		fuotaList = append(fuotaList, fuotaItem)
	}

	return fuotaList, nil
}

// GetServiceProfileIDForFUOTADeployment returns the service-profile ID for the given FUOTA deployment.
func GetServiceProfileIDForFUOTADeployment(ctx context.Context, handler *store.Handler, fuotaDeploymentID uuid.UUID) (uuid.UUID, error) {
	return handler.GetServiceProfileIDForFUOTADeployment(ctx, fuotaDeploymentID)
}
