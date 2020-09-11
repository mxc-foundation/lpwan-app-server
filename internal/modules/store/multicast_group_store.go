package store

import (
	"context"
	"strings"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
)

type MulticastGroupStore interface {
	CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error
	GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate, localOnly bool) (MulticastGroup, error)
	UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error
	DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error
	GetMulticastGroupCount(ctx context.Context, filters MulticastGroupFilters) (int, error)
	GetMulticastGroups(ctx context.Context, filters MulticastGroupFilters) ([]MulticastGroupListItem, error)
	AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error
	RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error
	GetDeviceCountForMulticastGroup(ctx context.Context, multicastGroup uuid.UUID) (int, error)
	GetDevicesForMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, limit, offset int) ([]DeviceListItem, error)

	// validator
	CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error)
	CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
	CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error)
}

func (h *Handler) GetMulticastGroup(ctx context.Context, id uuid.UUID, forUpdate, localOnly bool) (MulticastGroup, error) {
	return h.store.GetMulticastGroup(ctx, id, forUpdate, localOnly)
}
func (h *Handler) UpdateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	return h.store.UpdateMulticastGroup(ctx, mg)
}
func (h *Handler) DeleteMulticastGroup(ctx context.Context, id uuid.UUID) error {
	return h.store.DeleteMulticastGroup(ctx, id)
}
func (h *Handler) GetMulticastGroupCount(ctx context.Context, filters MulticastGroupFilters) (int, error) {
	return h.store.GetMulticastGroupCount(ctx, filters)
}
func (h *Handler) GetMulticastGroups(ctx context.Context, filters MulticastGroupFilters) ([]MulticastGroupListItem, error) {
	return h.store.GetMulticastGroups(ctx, filters)
}
func (h *Handler) AddDeviceToMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return h.store.AddDeviceToMulticastGroup(ctx, multicastGroupID, devEUI)
}
func (h *Handler) RemoveDeviceFromMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return h.store.RemoveDeviceFromMulticastGroup(ctx, multicastGroupID, devEUI)
}
func (h *Handler) GetDeviceCountForMulticastGroup(ctx context.Context, multicastGroup uuid.UUID) (int, error) {
	return h.store.GetDeviceCountForMulticastGroup(ctx, multicastGroup)
}
func (h *Handler) GetDevicesForMulticastGroup(ctx context.Context, multicastGroupID uuid.UUID, limit, offset int) ([]DeviceListItem, error) {
	return h.store.GetDevicesForMulticastGroup(ctx, multicastGroupID, limit, offset)
}
func (h *Handler) CreateMulticastGroup(ctx context.Context, mg *MulticastGroup) error {
	return h.store.CreateMulticastGroup(ctx, mg)
}
func (h *Handler) CheckCreateMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckCreateMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckListMulticastGroupsAccess(ctx context.Context, username string, organizationID, userID int64) (bool, error) {
	return h.store.CheckListMulticastGroupsAccess(ctx, username, organizationID, userID)
}
func (h *Handler) CheckReadMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *Handler) CheckUpdateDeleteMulticastGroupAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckUpdateDeleteMulticastGroupAccess(ctx, username, multicastGroupID, userID)
}
func (h *Handler) CheckMulticastGroupQueueAccess(ctx context.Context, username string, multicastGroupID uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckMulticastGroupQueueAccess(ctx, username, multicastGroupID, userID)
}

// MulticastGroup defines the multicast-group.
type MulticastGroup struct {
	CreatedAt        time.Time         `db:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at"`
	Name             string            `db:"name"`
	MCAppSKey        lorawan.AES128Key `db:"mc_app_s_key"`
	MCKey            lorawan.AES128Key `db:"mc_key"`
	ServiceProfileID uuid.UUID         `db:"service_profile_id"`
	MulticastGroup   ns.MulticastGroup `db:"-"`
}

// MulticastGroupListItem defines the multicast-group for listing.
type MulticastGroupListItem struct {
	ID                 uuid.UUID `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Name               string    `db:"name"`
	ServiceProfileID   uuid.UUID `db:"service_profile_id"`
	ServiceProfileName string    `db:"service_profile_name"`
}

// Validate validates the service-profile data.
func (mg MulticastGroup) Validate() error {
	if strings.TrimSpace(mg.Name) == "" || len(mg.Name) > 100 {
		return ErrMulticastGroupInvalidName
	}
	return nil
}

// MulticastGroupFilters provide filters that can be used to filter on
// multicast-groups. Note that empty values are not used as filters.
type MulticastGroupFilters struct {
	OrganizationID   int64         `db:"organization_id"`
	ServiceProfileID uuid.UUID     `db:"service_profile_id"`
	DevEUI           lorawan.EUI64 `db:"dev_eui"`
	Search           string        `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filter.
func (f MulticastGroupFilters) SQL() string {
	var filters []string
	var nilEUI lorawan.EUI64

	if f.OrganizationID != 0 {
		filters = append(filters, "o.id = :organization_id")
	}
	if f.ServiceProfileID != uuid.Nil {
		filters = append(filters, "mg.service_profile_id = :service_profile_id")
	}
	if f.DevEUI != nilEUI {
		filters = append(filters, "dmg.dev_eui = :dev_eui")
	}
	if f.Search != "" {
		filters = append(filters, "mg.name ilike :search")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}
