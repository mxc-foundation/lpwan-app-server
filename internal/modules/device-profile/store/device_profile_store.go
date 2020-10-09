package store

import (
	"context"

	"github.com/gofrs/uuid"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pg pgstore.PgStore) *dps {
	return &dps{
		pg: pg,
	}
}

type dps struct {
	pg pgstore.DeviceProfilePgStore
}

type DeviceProfileStore interface {
	DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error
	DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error
	GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (DeviceProfile, error)
	CreateDeviceProfile(ctx context.Context, dp *DeviceProfile) error
	UpdateDeviceProfile(ctx context.Context, dp *DeviceProfile) error
	GetDeviceProfileCount(ctx context.Context, filters DeviceProfileFilters) (int, error)
	GetDeviceProfiles(ctx context.Context, filters DeviceProfileFilters) ([]DeviceProfileMeta, error)

	// validator
	CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error)
	CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error)

	CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
}

func (h *dps) GetDeviceProfiles(ctx context.Context, filters DeviceProfileFilters) ([]DeviceProfileMeta, error) {
	return h.pg.GetDeviceProfiles(ctx, filters)
}

func (h *dps) GetDeviceProfileCount(ctx context.Context, filters DeviceProfileFilters) (int, error) {
	return h.pg.GetDeviceProfileCount(ctx, filters)
}

func (h *dps) UpdateDeviceProfile(ctx context.Context, dp *DeviceProfile) error {
	return h.pg.UpdateDeviceProfile(ctx, dp)
}

func (h *dps) CreateDeviceProfile(ctx context.Context, dp *DeviceProfile) error {
	return h.pg.CreateDeviceProfile(ctx, dp)
}

func (h *dps) CheckReadDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckReadDeviceProfileAccess(ctx, username, id, userID)
}

func (h *dps) CheckUpdateDeleteDeviceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckUpdateDeleteDeviceProfileAccess(ctx, username, id, userID)
}

func (h *dps) CheckCreateDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	return h.pg.CheckCreateDeviceProfilesAccess(ctx, username, organizationID, applicationID, userID)
}

func (h *dps) CheckListDeviceProfilesAccess(ctx context.Context, username string, organizationID, applicationID, userID int64) (bool, error) {
	return h.pg.CheckListDeviceProfilesAccess(ctx, username, organizationID, applicationID, userID)
}

func (h *dps) DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, organizationID int64) error {
	return h.pg.DeleteAllDeviceProfilesForOrganizationID(ctx, organizationID)
}

func (h *dps) DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error {
	return h.pg.DeleteDeviceProfile(ctx, id)
}

func (h *dps) GetDeviceProfile(ctx context.Context, id uuid.UUID, forUpdate bool) (DeviceProfile, error) {
	return h.pg.GetDeviceProfile(ctx, id, forUpdate)
}
