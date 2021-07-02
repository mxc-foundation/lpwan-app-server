package devprofile

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// CreateDeviceProfile creates the given device-profile.
// This will create the device-profile at the network-server side and will
// create a local reference record.
func CreateDeviceProfile(ctx context.Context, st *store.Handler, nsCli *nscli.Client, dp *DeviceProfile) error {
	dpID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "new uuid v4 error")
	}
	now := time.Now()
	dp.DeviceProfile.Id = dpID.Bytes()
	dp.CreatedAt = now
	dp.UpdatedAt = now

	nsClient, err := nsCli.GetNetworkServerServiceClient(dp.NetworkServerID)
	if err != nil {
		return err
	}
	_, err = nsClient.CreateDeviceProfile(ctx, &ns.CreateDeviceProfileRequest{
		DeviceProfile: &dp.DeviceProfile,
	})
	if err != nil {
		return errors.Wrap(err, "create device-profile errror")
	}
	if err := st.CreateDeviceProfile(ctx, dp); err != nil {
		return err
	}

	return nil
}

// UpdateDeviceProfile updates the given device-profile.
func UpdateDeviceProfile(ctx context.Context, st *store.Handler, nsCli *nscli.Client, dp *DeviceProfile) error {
	nsClient, err := nsCli.GetNetworkServerServiceClient(dp.NetworkServerID)
	if err != nil {
		return err
	}
	_, err = nsClient.UpdateDeviceProfile(ctx, &ns.UpdateDeviceProfileRequest{
		DeviceProfile: &dp.DeviceProfile,
	})
	if err != nil {
		return errors.Wrap(err, "update device-profile error")
	}
	if err := st.UpdateDeviceProfile(ctx, dp); err != nil {
		return err
	}
	return nil
}

// DeleteDeviceProfile deletes the device-profile matching the given id.
func DeleteDeviceProfile(ctx context.Context, st *store.Handler, nsCli *nscli.Client, id uuid.UUID) error {
	n, err := st.GetNetworkServerForDeviceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}
	nsClient, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return err
	}
	_, err = nsClient.DeleteDeviceProfile(ctx, &ns.DeleteDeviceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete device-profile error")
	}

	if err := st.DeleteDeviceProfile(ctx, id); err != nil {
		return err
	}
	return nil
}
