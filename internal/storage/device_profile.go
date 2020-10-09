package storage

import (
	"context"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"

	dps "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// DeviceProfile defines the device-profile.
type DeviceProfile dps.DeviceProfile

// DeviceProfileMeta defines the device-profile meta record.
type DeviceProfileMeta dps.DeviceProfileMeta

// Validate validates the device-profile data.
func (dp DeviceProfile) Validate() error {
	return dps.DeviceProfile(dp).Validate()
}

// GetDeviceProfile returns the device-profile matching the given id.
// When forUpdate is set to true, then db must be a db transaction.
// When localOnly is set to true, no call to the network-server is made to
// retrieve additional device data.
func GetDeviceProfile(ctx context.Context, handler *store.Handler, id uuid.UUID, forUpdate, localOnly bool) (DeviceProfile, error) {
	dp, err := handler.GetDeviceProfile(ctx, id, forUpdate)
	if err != nil {
		return DeviceProfile(dp), err
	}

	if localOnly {
		return DeviceProfile(dp), nil
	}

	n, err := handler.GetNetworkServer(ctx, dp.NetworkServerID)
	if err != nil {
		return DeviceProfile(dp), errors.Wrap(err, "get network-server error")
	}

	nstruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}

	nsClient, err := nstruct.GetNetworkServiceClient()
	if err != nil {
		return DeviceProfile(dp), errors.Wrap(err, "get network-server client error")
	}

	resp, err := nsClient.GetDeviceProfile(ctx, &ns.GetDeviceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return DeviceProfile(dp), errors.Wrap(err, "get device-profile error")
	}
	if resp.DeviceProfile == nil {
		return DeviceProfile(dp), errors.New("device_profile must not be nil")
	}

	dp.DeviceProfile = *resp.DeviceProfile

	return DeviceProfile(dp), nil
}

// DeviceProfileFilters provide filders for filtering device-profiles.
type DeviceProfileFilters dps.DeviceProfileFilters

// SQL returns the SQL filters.
func (f DeviceProfileFilters) SQL() string {
	return dps.DeviceProfileFilters(f).SQL()
}

// GetDeviceProfileCount returns the total number of device-profiles.
func GetDeviceProfileCount(ctx context.Context, handler *store.Handler, filters DeviceProfileFilters) (int, error) {
	return handler.GetDeviceProfileCount(ctx, (dps.DeviceProfileFilters)(filters))
}

// GetDeviceProfiles returns a slice of device-profiles.
func GetDeviceProfiles(ctx context.Context, handler *store.Handler, filters DeviceProfileFilters) ([]DeviceProfileMeta, error) {
	res, err := handler.GetDeviceProfiles(ctx, (dps.DeviceProfileFilters)(filters))
	if err != nil {
		return nil, err
	}

	var dpList []DeviceProfileMeta
	for _, v := range res {
		dpItem := DeviceProfileMeta(v)
		dpList = append(dpList, dpItem)
	}
	return dpList, nil
}

// DeleteAllDeviceProfilesForOrganizationID deletes all device-profiles
// given an organization id.
func DeleteAllDeviceProfilesForOrganizationID(ctx context.Context, handler *store.Handler, organizationID int64) error {
	return handler.DeleteAllDeviceProfilesForOrganizationID(ctx, organizationID)
}
