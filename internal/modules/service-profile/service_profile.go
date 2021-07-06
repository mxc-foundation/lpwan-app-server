package serviceprofile

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
)

// Store defines the db APIs used by this package
type Store interface {
	CreateServiceProfile(ctx context.Context, sp *ServiceProfile) (*uuid.UUID, error)
	GetServiceProfile(ctx context.Context, id uuid.UUID) (ServiceProfile, error)
	UpdateServiceProfile(ctx context.Context, sp *ServiceProfile) error
	GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (nsd.NetworkServer, error)
	DeleteServiceProfile(ctx context.Context, id uuid.UUID) error

	CheckReadServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckUpdateDeleteServiceProfileAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
}

// CreateServiceProfile creates the given service-profile.
func CreateServiceProfile(ctx context.Context, st Store, sp *ServiceProfile, nsCli *nscli.Client) (*uuid.UUID, error) {
	spID, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "new uuid v4 error")
	}
	now := time.Now()
	sp.CreatedAt = now
	sp.UpdatedAt = now
	sp.ServiceProfile.Id = spID.Bytes()

	nsClient, err := nsCli.GetNetworkServerServiceClient(sp.NetworkServerID)
	if err != nil {
		return nil, err
	}
	_, err = nsClient.CreateServiceProfile(ctx, &ns.CreateServiceProfileRequest{
		ServiceProfile: &sp.ServiceProfile,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create service-profile error")
	}

	return st.CreateServiceProfile(ctx, sp)
}

// GetServiceProfile returns the service-profile matching the given id.
func GetServiceProfile(ctx context.Context, st Store, id uuid.UUID, nsCli *nscli.Client, localOnly bool) (*ServiceProfile, error) {
	var sp ServiceProfile
	var err error

	if sp, err = st.GetServiceProfile(ctx, id); err != nil {
		return nil, err
	}

	if localOnly {
		return &sp, nil
	}

	nsClient, err := nsCli.GetNetworkServerServiceClient(sp.NetworkServerID)
	if err != nil {
		return nil, err
	}

	resp, err := nsClient.GetServiceProfile(ctx, &ns.GetServiceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "get service-profile error")
	}

	if resp.ServiceProfile == nil {
		return nil, errors.New("service_profile must not be nil")
	}

	sp.ServiceProfile = *resp.ServiceProfile

	return &sp, nil
}

// UpdateServiceProfile updates the given service-profile.
func UpdateServiceProfile(ctx context.Context, st Store, nsCli *nscli.Client, sp *ServiceProfile) error {
	nsClient, err := nsCli.GetNetworkServerServiceClient(sp.NetworkServerID)
	if err != nil {
		return err
	}
	_, err = nsClient.UpdateServiceProfile(ctx, &ns.UpdateServiceProfileRequest{
		ServiceProfile: &sp.ServiceProfile,
	})
	if err != nil {
		return errors.Wrap(err, "update service-profile error")
	}

	return st.UpdateServiceProfile(ctx, sp)
}

// DeleteServiceProfile deletes the service-profile matching the given id.
func DeleteServiceProfile(ctx context.Context, st Store, nsCli *nscli.Client, id uuid.UUID) error {
	n, err := st.GetNetworkServerForServiceProfileID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}
	nsClient, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return err
	}
	_, err = nsClient.DeleteServiceProfile(ctx, &ns.DeleteServiceProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete service-profile error")
	}
	return st.DeleteServiceProfile(ctx, id)
}
