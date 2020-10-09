package store

import (
	"context"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"

	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
)

func NewStore(pg pgstore.PgStore) *mgs {
	return &mgs{
		pg: pg,
	}
}

type mgs struct {
	pg pgstore.MulticastGroupPgStore
}

type MulticastGroupStore interface {
	CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error
	GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate bool) (MulticastGroup, error)
	UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error
	DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error
	GetMulticastGroupCount(ctx context.Context, filters MulticastGroupFilters) (int, error)
	GetMulticastGroups(ctx context.Context, filters MulticastGroupFilters) ([]MulticastGroupListItem, error)
	AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error
	RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error
	GetDeviceCountForMulticastGroup(ctx context.Context, multicastGroup uuid.UUID) (int, error)
	GetDevicesForMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, limit, offset int) ([]ds.DeviceListItem, error)

	// validator
	CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
}

func (h *mgs) GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate bool) (MulticastGroup, error) {
	return h.pg.GetMulticastGroup(ctx, id, forUpdate)
}
func (h *mgs) UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	return h.pg.UpdateMulticastGroup(ctx, mg)
}
func (h *mgs) DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error {
	return h.pg.DeleteMulticastGroup(ctx, id)
}
func (h *mgs) GetMulticastGroupCount(ctx context.Context, filters MulticastGroupFilters) (int, error) {
	return h.pg.GetMulticastGroupCount(ctx, filters)
}
func (h *mgs) GetMulticastGroups(ctx context.Context, filters MulticastGroupFilters) ([]MulticastGroupListItem, error) {
	return h.pg.GetMulticastGroups(ctx, filters)
}
func (h *mgs) AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return h.pg.AddDeviceToMulticastGroup(ctx, multicastGroupID, devEUI)
}
func (h *mgs) RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return h.pg.RemoveDeviceFromMulticastGroup(ctx, multicastGroupID, devEUI)
}
func (h *mgs) GetDeviceCountForMulticastGroup(ctx context.Context, multicastGroup uuid.UUID) (int, error) {
	return h.pg.GetDeviceCountForMulticastGroup(ctx, multicastGroup)
}
func (h *mgs) GetDevicesForMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, limit, offset int) ([]ds.DeviceListItem, error) {
	return h.pg.GetDevicesForMulticastGroup(ctx, multicastGroupID, limit, offset)
}
func (h *mgs) CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	return h.pg.CreateMulticastGroup(ctx, mg)
}
func (h *mgs) CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckCreateMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *mgs) CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.pg.CheckListMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *mgs) CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckReadMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *mgs) CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckUpdateDeleteMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *mgs) CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.pg.CheckMulticastGroupQueueAccess(ctx, username, multicastGroupID, userID)
}
