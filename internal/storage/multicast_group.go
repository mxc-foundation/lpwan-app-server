package storage

import (
	"context"

	"github.com/gofrs/uuid"

	mcgs "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// MulticastGroup defines the multicast-group.
type MulticastGroup mcgs.MulticastGroup

// MulticastGroupListItem defines the multicast-group for listing.
type MulticastGroupListItem mcgs.MulticastGroupListItem

// Validate validates the service-profile data.
func (mg MulticastGroup) Validate() error {
	return mcgs.MulticastGroup(mg).Validate()
}

// MulticastGroupFilters provide filters that can be used to filter on
// multicast-groups. Note that empty values are not used as filters.
type MulticastGroupFilters mcgs.MulticastGroupFilters

// SQL returns the SQL filter.
func (f MulticastGroupFilters) SQL() string {
	return mcgs.MulticastGroupFilters(f).SQL()
}

// GetMulticastGroupCount returns the total number of multicast-groups given
// the provided filters. Note that empty values are not used as filters.
func GetMulticastGroupCount(ctx context.Context, handler *store.Handler, filters MulticastGroupFilters) (int, error) {
	return handler.GetMulticastGroupCount(ctx, mcgs.MulticastGroupFilters(filters))
}

// GetMulticastGroups returns a slice of multicast-groups, given the privded
// filters. Note that empty values are not used as filters.
func GetMulticastGroups(ctx context.Context, handler *store.Handler, filters MulticastGroupFilters) ([]MulticastGroupListItem, error) {
	res, err := handler.GetMulticastGroups(ctx, mcgs.MulticastGroupFilters(filters))
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
