package store

import (
	"context"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"

	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pg pgstore.PgStore) *fds {
	return &fds{
		pg: pg,
	}
}

type fds struct {
	pg pgstore.FuotaDeploymentPgStore
}

type FUOTADeploymentStore interface {
	GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]ds.DeviceKeys, error)
	CreateFUOTADeploymentForDevice(ctx context.Context, fd *FUOTADeployment, devEUI lorawan.EUI64) error
	GetFUOTADeployment(ctx context.Context, id uuid.UUID, forUpdate bool) (FUOTADeployment, error)
	GetPendingFUOTADeployments(ctx context.Context, batchSize int) ([]FUOTADeployment, error)
	UpdateFUOTADeployment(ctx context.Context, fd *FUOTADeployment) error
	GetFUOTADeploymentCount(ctx context.Context, filters FUOTADeploymentFilters) (int, error)
	GetFUOTADeployments(ctx context.Context, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error)
	GetFUOTADeploymentDevice(ctx context.Context, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error)
	GetPendingFUOTADeploymentDevice(ctx context.Context, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error)
	UpdateFUOTADeploymentDevice(ctx context.Context, fdd *FUOTADeploymentDevice) error
	GetFUOTADeploymentDeviceCount(ctx context.Context, fuotaDeploymentID uuid.UUID) (int, error)
	GetFUOTADeploymentDevices(ctx context.Context, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error)
	GetServiceProfileIDForFUOTADeployment(ctx context.Context, fuotaDeploymentID uuid.UUID) (uuid.UUID, error)
	SetFromRemoteMulticastSetup(ctx context.Context, fuotaDevelopmentID, multicastGroupID uuid.UUID) error
	SetFromRemoteFragmentationSession(ctx context.Context, fuotaDevelopmentID uuid.UUID, fragIdx int) error
	SetIncompleteFuotaDevelopment(ctx context.Context, fuotaDevelopmentID uuid.UUID) error

	// validator
	CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error)
}

func (h *fds) SetIncompleteFuotaDevelopment(ctx context.Context, fuotaDevelopmentID uuid.UUID) error {
	return h.pg.SetIncompleteFuotaDevelopment(ctx, fuotaDevelopmentID)
}

func (h *fds) SetFromRemoteFragmentationSession(ctx context.Context, fuotaDevelopmentID uuid.UUID, fragIdx int) error {
	return h.pg.SetFromRemoteFragmentationSession(ctx, fuotaDevelopmentID, fragIdx)
}
func (h *fds) SetFromRemoteMulticastSetup(ctx context.Context, fuotaDevelopmentID, multicastGroupID uuid.UUID) error {
	return h.pg.SetFromRemoteMulticastSetup(ctx, fuotaDevelopmentID, multicastGroupID)
}
func (h *fds) CreateFUOTADeploymentForDevice(ctx context.Context, fd *FUOTADeployment, devEUI lorawan.EUI64) error {
	return h.pg.CreateFUOTADeploymentForDevice(ctx, fd, devEUI)
}
func (h *fds) GetFUOTADeployment(ctx context.Context, id uuid.UUID, forUpdate bool) (FUOTADeployment, error) {
	return h.pg.GetFUOTADeployment(ctx, id, forUpdate)
}
func (h *fds) GetPendingFUOTADeployments(ctx context.Context, batchSize int) ([]FUOTADeployment, error) {
	return h.pg.GetPendingFUOTADeployments(ctx, batchSize)
}
func (h *fds) UpdateFUOTADeployment(ctx context.Context, fd *FUOTADeployment) error {
	return h.pg.UpdateFUOTADeployment(ctx, fd)
}
func (h *fds) GetFUOTADeploymentCount(ctx context.Context, filters FUOTADeploymentFilters) (int, error) {
	return h.pg.GetFUOTADeploymentCount(ctx, filters)
}
func (h *fds) GetFUOTADeployments(ctx context.Context, filters FUOTADeploymentFilters) ([]FUOTADeploymentListItem, error) {
	return h.pg.GetFUOTADeployments(ctx, filters)
}
func (h *fds) GetFUOTADeploymentDevice(ctx context.Context, fuotaDeploymentID uuid.UUID, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	return h.pg.GetFUOTADeploymentDevice(ctx, fuotaDeploymentID, devEUI)
}
func (h *fds) GetPendingFUOTADeploymentDevice(ctx context.Context, devEUI lorawan.EUI64) (FUOTADeploymentDevice, error) {
	return h.pg.GetPendingFUOTADeploymentDevice(ctx, devEUI)
}
func (h *fds) UpdateFUOTADeploymentDevice(ctx context.Context, fdd *FUOTADeploymentDevice) error {
	return h.pg.UpdateFUOTADeploymentDevice(ctx, fdd)
}
func (h *fds) GetFUOTADeploymentDeviceCount(ctx context.Context, fuotaDeploymentID uuid.UUID) (int, error) {
	return h.pg.GetFUOTADeploymentDeviceCount(ctx, fuotaDeploymentID)
}
func (h *fds) GetFUOTADeploymentDevices(ctx context.Context, fuotaDeploymentID uuid.UUID, limit, offset int) ([]FUOTADeploymentDeviceListItem, error) {
	return h.pg.GetFUOTADeploymentDevices(ctx, fuotaDeploymentID, limit, offset)
}
func (h *fds) GetServiceProfileIDForFUOTADeployment(ctx context.Context, fuotaDeploymentID uuid.UUID) (uuid.UUID, error) {
	return h.pg.GetServiceProfileIDForFUOTADeployment(ctx, fuotaDeploymentID)
}
func (h *fds) GetDeviceKeysFromFuotaDevelopmentDevice(ctx context.Context, id uuid.UUID) ([]ds.DeviceKeys, error) {
	return h.pg.GetDeviceKeysFromFuotaDevelopmentDevice(ctx, id)
}
func (h *fds) CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckReadFUOTADeploymentAccess(ctx, username, id, userID)
}

func (h *fds) CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.pg.CheckCreateFUOTADeploymentsAccess(ctx, username, applicationID, devEUI, userID)
}
