package dp

import (
	"context"
	"strings"
	"time"

	"github.com/lib/pq/hstore"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
)

// DeviceProfileMeta defines the device-profile meta record.
type DeviceProfileMeta struct {
	DeviceProfileID   uuid.UUID `db:"device_profile_id"`
	NetworkServerID   int64     `db:"network_server_id"`
	OrganizationID    int64     `db:"organization_id"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	Name              string    `db:"name"`
	NetworkServerName string    `db:"network_server_name"`
}

// DeviceProfile defines the device-profile.
type DeviceProfile struct {
	NetworkServerID      int64            `db:"network_server_id"`
	OrganizationID       int64            `db:"organization_id"`
	CreatedAt            time.Time        `db:"created_at"`
	UpdatedAt            time.Time        `db:"updated_at"`
	Name                 string           `db:"name"`
	PayloadCodec         string           `db:"payload_codec"`
	PayloadEncoderScript string           `db:"payload_encoder_script"`
	PayloadDecoderScript string           `db:"payload_decoder_script"`
	Tags                 hstore.Hstore    `db:"tags"`
	UplinkInterval       time.Duration    `db:"uplink_interval"`
	DeviceProfile        ns.DeviceProfile `db:"-"`
}

// Validate validates the device-profile data.
func (dp DeviceProfile) Validate() error {
	if strings.TrimSpace(dp.Name) == "" || len(dp.Name) > 100 {
		return errHandler.ErrDeviceProfileInvalidName
	}
	return nil
}

// DeviceProfileFilters provide filders for filtering device-profiles.
type DeviceProfileFilters struct {
	ApplicationID  int64 `db:"application_id"`
	OrganizationID int64 `db:"organization_id"`
	UserID         int64 `db:"user_id"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f DeviceProfileFilters) SQL() string {
	var filters []string

	if f.ApplicationID != 0 {
		// Filter on organization_id too since dp > network-server > service-profile > application
		// join.
		filters = append(filters, "a.id = :application_id and dp.organization_id = a.organization_id")
	}

	if f.OrganizationID != 0 {
		filters = append(filters, "o.id = :organization_id")
	}

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// Store defines db APIs used by this package
type Store interface {
	CreateDeviceProfile(ctx context.Context, dp *DeviceProfile) error
	UpdateDeviceProfile(ctx context.Context, dp *DeviceProfile) error
	GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (nsd.NetworkServer, error)
	DeleteDeviceProfile(ctx context.Context, id uuid.UUID) error
}

// CreateDeviceProfile creates the given device-profile.
// This will create the device-profile at the network-server side and will
// create a local reference record.
func CreateDeviceProfile(ctx context.Context, st Store, nsCli *nscli.Client, dp *DeviceProfile) error {
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
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return errors.Wrap(err, "create device-profile errror")
	}
	if err := st.CreateDeviceProfile(ctx, dp); err != nil {
		return err
	}

	return nil
}

// UpdateDeviceProfile updates the given device-profile.
func UpdateDeviceProfile(ctx context.Context, st Store, nsCli *nscli.Client, dp *DeviceProfile) error {
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
func DeleteDeviceProfile(ctx context.Context, st Store, nsCli *nscli.Client, id uuid.UUID) error {
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
