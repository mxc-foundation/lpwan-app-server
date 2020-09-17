package storage

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
)

// MulticastGroup defines the multicast-group.
type MulticastGroup store.MulticastGroup

// MulticastGroupListItem defines the multicast-group for listing.
type MulticastGroupListItem store.MulticastGroupListItem

// Validate validates the service-profile data.
func (mg MulticastGroup) Validate() error {
	return store.MulticastGroup(mg).Validate()
}

// CreateMulticastGroup creates the given multicast-group.
func CreateMulticastGroup(ctx context.Context, handler *store.Handler, mg *MulticastGroup) error {
	return handler.CreateMulticastGroup(ctx, (*store.MulticastGroup)(mg))
}

// GetMulticastGroup returns the multicast-group given an id.
func GetMulticastGroup(ctx context.Context, handler *store.Handler, id uuid.UUID, forUpdate, localOnly bool) (MulticastGroup, error) {
	res, err := handler.GetMulticastGroup(ctx, id, forUpdate, localOnly)
	return MulticastGroup(res), err
}

// UpdateMulticastGroup updates the given multicast-group.
func UpdateMulticastGroup(ctx context.Context, handler *store.Handler, mg *MulticastGroup) error {
	return handler.UpdateMulticastGroup(ctx, (*store.MulticastGroup)(mg))
}

// DeleteMulticastGroup deletes a multicast-group given an id.
func DeleteMulticastGroup(ctx context.Context, handler *store.Handler, id uuid.UUID) error {
	return handler.DeleteMulticastGroup(ctx, id)
}

// MulticastGroupFilters provide filters that can be used to filter on
// multicast-groups. Note that empty values are not used as filters.
type MulticastGroupFilters store.MulticastGroupFilters

// SQL returns the SQL filter.
func (f MulticastGroupFilters) SQL() string {
	return store.MulticastGroupFilters(f).SQL()
}

// GetMulticastGroupCount returns the total number of multicast-groups given
// the provided filters. Note that empty values are not used as filters.
func GetMulticastGroupCount(ctx context.Context, handler *store.Handler, filters MulticastGroupFilters) (int, error) {
	return handler.GetMulticastGroupCount(ctx, store.MulticastGroupFilters(filters))
}

// GetMulticastGroups returns a slice of multicast-groups, given the privded
// filters. Note that empty values are not used as filters.
func GetMulticastGroups(ctx context.Context, handler *store.Handler, filters MulticastGroupFilters) ([]MulticastGroupListItem, error) {
	res, err := handler.GetMulticastGroups(ctx, store.MulticastGroupFilters(filters))
	if err != nil {
		return nil, err
	}

	var mgList []MulticastGroupListItem
	for _, v := range res {
		mgItem := MulticastGroupListItem(v)
		mgList = append(mgList, mgItem)
	}
	return mgList, nil
}

// AddDeviceToMulticastGroup adds the given device to the given multicast-group.
// It is recommended that db is a transaction.
func AddDeviceToMulticastGroup(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return handler.AddDeviceToMulticastGroup(ctx, multicastGroupID, devEUI)
}

// RemoveDeviceFromMulticastGroup removes the given device from the given
// multicast-group.
func RemoveDeviceFromMulticastGroup(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, devEUI lorawan.EUI64) error {
	return handler.RemoveDeviceFromMulticastGroup(ctx, multicastGroupID, devEUI)
}

// GetDeviceCountForMulticastGroup returns the number of devices for the given
// multicast-group.
func GetDeviceCountForMulticastGroup(ctx context.Context, handler *store.Handler, multicastGroup uuid.UUID) (int, error) {
	return handler.GetDeviceCountForMulticastGroup(ctx, multicastGroup)
}

// GetDevicesForMulticastGroup returns a slice of devices for the given
// multicast-group.
func GetDevicesForMulticastGroup(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, limit, offset int) ([]DeviceListItem, error) {
	res, err := handler.GetDevicesForMulticastGroup(ctx, multicastGroupID, limit, offset)
	if err != nil {
		return nil, err
	}

	var devList []DeviceListItem
	for _, v := range res {
		devItem := DeviceListItem(v)
		devList = append(devList, devItem)
	}

	return devList, nil
}
