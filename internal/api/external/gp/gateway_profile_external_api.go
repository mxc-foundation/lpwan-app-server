package gp

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	gwpd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// GatewayProfileAPI exports the GatewayProfile related functions.
type GatewayProfileAPI struct {
	st    *store.Handler
	nsCli *nscli.Client
	auth  auth.Authenticator
}

// NewGatewayProfileAPI creates a new GatewayProfileAPI.
func NewGatewayProfileAPI(st *store.Handler, nsCli *nscli.Client, auth auth.Authenticator) *GatewayProfileAPI {
	return &GatewayProfileAPI{
		st:    st,
		nsCli: nsCli,
		auth:  auth,
	}
}

// GetGatewayProfile returns the gateway-profile matching the given id.
func GetGatewayProfile(ctx context.Context, id uuid.UUID, st *store.Handler, nsCli *nscli.Client) (*gwpd.GatewayProfile, error) {
	var gp gwpd.GatewayProfile
	var err error
	// this function only returns gp.NetworkServerID, gp.CreatedAt, gp.UpdatedAt, gp.Name
	if gp, err = st.GetGatewayProfile(ctx, id); err != nil {
		return nil, err
	}
	// get gp.GatewayProfile
	nsClient, err := nsCli.GetNetworkServerServiceClient(gp.NetworkServerID)
	if err != nil {
		return nil, err
	}
	resp, err := nsClient.GetGatewayProfile(ctx, &ns.GetGatewayProfileRequest{
		Id: id.Bytes(),
	})
	if err != nil {
		return nil, fmt.Errorf("get gateway-profile from network-server error: %v", err)
	}
	if resp.GatewayProfile == nil {
		return nil, fmt.Errorf("gateway-profile obtained from network-server is empty")
	}

	gp.GatewayProfile = *resp.GatewayProfile
	return &gp, nil
}

// CreateGatewayProfile creates new gateway profile locally and in network server
func CreateGatewayProfile(ctx context.Context, st *store.Handler, nsCli *nscli.Client,
	gp *gwpd.GatewayProfile) (*uuid.UUID, error) {
	var err error
	gpID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("new uuid v4 error: %v", err)
	}
	gp.GatewayProfile.Id = gpID.Bytes()
	now := time.Now()
	gp.CreatedAt = now
	gp.UpdatedAt = now

	nsClient, err := nsCli.GetNetworkServerServiceClient(gp.NetworkServerID)
	if err != nil {
		return nil, err
	}
	_, err = nsClient.CreateGatewayProfile(ctx, &ns.CreateGatewayProfileRequest{
		GatewayProfile: &gp.GatewayProfile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway profile in network server:%v", err)
	}

	// create gateway profile in network server first, if creating fails locally, user can re-create later without
	// causing trouble. Otherwise, if creating succeeds locally but fails in network server, user could add gateway with
	// this gateway profile but getting no packets from network server
	if err = st.CreateGatewayProfile(ctx, gp); err != nil {
		return nil, err
	}

	return &gpID, nil
}

// UpdateGatewayProfile updates channels, extra channel and statsinternal in network server
//  updates name, statsinternal in appserver
func UpdateGatewayProfile(ctx context.Context, st *store.Handler, nsCli *nscli.Client,
	gp *gwpd.GatewayProfile) error {

	nsClient, err := nsCli.GetNetworkServerServiceClient(gp.NetworkServerID)
	if err != nil {
		return err
	}
	_, err = nsClient.UpdateGatewayProfile(ctx, &ns.UpdateGatewayProfileRequest{
		GatewayProfile: &gp.GatewayProfile,
	})
	if err != nil {
		return err
	}

	if err := st.UpdateGatewayProfile(ctx, gp); err != nil {
		return err
	}

	return nil
}

// Create creates the given gateway-profile.
func (a *GatewayProfileAPI) Create(ctx context.Context, req *pb.CreateGatewayProfileRequest) (*pb.CreateGatewayProfileResponse, error) {
	if req.GatewayProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "gateway_profile must not be nil")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gp := gwpd.GatewayProfile{
		NetworkServerID: req.GatewayProfile.NetworkServerId,
		Name:            req.GatewayProfile.Name,
		GatewayProfile: ns.GatewayProfile{
			Channels: req.GatewayProfile.Channels,
		},
	}
	for _, ec := range req.GatewayProfile.ExtraChannels {
		gp.GatewayProfile.ExtraChannels = append(gp.GatewayProfile.ExtraChannels, &ns.GatewayProfileExtraChannel{
			Frequency:        ec.Frequency,
			Bandwidth:        ec.Bandwidth,
			Bitrate:          ec.Bitrate,
			SpreadingFactors: ec.SpreadingFactors,
			Modulation:       ec.Modulation,
		})
	}

	gpID, err := CreateGatewayProfile(ctx, a.st, a.nsCli, &gp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateGatewayProfileResponse{
		Id: gpID.String(),
	}, nil
}

// Get returns the gateway-profile matching the given id.
func (a *GatewayProfileAPI) Get(ctx context.Context, req *pb.GetGatewayProfileRequest) (*pb.GetGatewayProfileResponse, error) {
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	gpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	gp, err := GetGatewayProfile(ctx, gpID, a.st, a.nsCli)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	out := pb.GetGatewayProfileResponse{
		GatewayProfile: &pb.GatewayProfile{
			Id:              gpID.String(),
			Name:            gp.Name,
			NetworkServerId: gp.NetworkServerID,
			Channels:        gp.GatewayProfile.Channels,
		},
	}

	out.CreatedAt = timestamppb.New(gp.CreatedAt)
	out.UpdatedAt = timestamppb.New(gp.UpdatedAt)

	for _, ec := range gp.GatewayProfile.ExtraChannels {
		out.GatewayProfile.ExtraChannels = append(out.GatewayProfile.ExtraChannels, &pb.GatewayProfileExtraChannel{
			Frequency:        ec.Frequency,
			Bandwidth:        ec.Bandwidth,
			Bitrate:          ec.Bitrate,
			SpreadingFactors: ec.SpreadingFactors,
			Modulation:       ec.Modulation,
		})
	}

	return &out, nil
}

// Update updates the given gateway-profile.
func (a *GatewayProfileAPI) Update(ctx context.Context, req *pb.UpdateGatewayProfileRequest) (*empty.Empty, error) {
	if req.GatewayProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "gateway_profile must not be nil")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gpID, err := uuid.FromString(req.GatewayProfile.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	gp, err := GetGatewayProfile(ctx, gpID, a.st, a.nsCli)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}
	if req.GatewayProfile.NetworkServerId != gp.NetworkServerID {
		// before we figure out how to swtich gateway profile among different network server, we don't allow modifying
		// network server id for exsiting gateway profile
		// changing network server id is also not supported on UI, but this should be checked agains direct API call
		return nil, status.Errorf(codes.InvalidArgument,
			"cannot modify network server id for existing gateway profile")
	}

	gp.Name = req.GatewayProfile.Name
	gp.GatewayProfile.Channels = req.GatewayProfile.Channels
	gp.GatewayProfile.ExtraChannels = []*ns.GatewayProfileExtraChannel{}

	for _, ec := range req.GatewayProfile.ExtraChannels {
		gp.GatewayProfile.ExtraChannels = append(gp.GatewayProfile.ExtraChannels, &ns.GatewayProfileExtraChannel{
			Frequency:        ec.Frequency,
			Bandwidth:        ec.Bandwidth,
			Bitrate:          ec.Bitrate,
			SpreadingFactors: ec.SpreadingFactors,
			Modulation:       ec.Modulation,
		})
	}

	err = UpdateGatewayProfile(ctx, a.st, a.nsCli, gp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

// DeleteGatewayProfile deletes gateway profile from local server then from network server
func DeleteGatewayProfile(ctx context.Context, st *store.Handler, gpID uuid.UUID, nsCli *nscli.Client) error {
	// get network server before gateway profile gets deleted
	n, err := st.GetNetworkServerForGatewayProfileID(ctx, gpID)
	if err != nil {
		return err
	}
	// delete gateway profile from local server first, it is still acceptable to keep record in network server when deleting
	// from network server failed, user won't be able to configure their gateways with discarded gateway profile record
	err = st.DeleteGatewayProfile(ctx, gpID)
	if err != nil {
		return err
	}

	nsClient, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return err
	}
	_, err = nsClient.DeleteGatewayProfile(ctx, &ns.DeleteGatewayProfileRequest{
		Id: gpID.Bytes(),
	})
	if err != nil {
		return errors.Wrap(err, "delete gateway-profile error")
	}

	return nil
}

// Delete deletes the gateway-profile matching the given id.
func (a *GatewayProfileAPI) Delete(ctx context.Context, req *pb.DeleteGatewayProfileRequest) (*empty.Empty, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	err = DeleteGatewayProfile(ctx, a.st, gpID, a.nsCli)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

// List returns the existing gateway-profiles.
func (a *GatewayProfileAPI) List(ctx context.Context, req *pb.ListGatewayProfilesRequest) (*pb.ListGatewayProfilesResponse, error) {
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	var count int
	var gps []gwpd.GatewayProfileMeta

	if req.NetworkServerId == 0 {
		count, err = a.st.GetGatewayProfileCount(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}

		gps, err = a.st.GetGatewayProfiles(ctx, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}
	} else {
		count, err = a.st.GetGatewayProfileCountForNetworkServerID(ctx, req.NetworkServerId)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}

		gps, err = a.st.GetGatewayProfilesForNetworkServerID(ctx, req.NetworkServerId, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}
	}

	out := pb.ListGatewayProfilesResponse{
		TotalCount: int64(count),
	}

	for _, gp := range gps {
		row := pb.GatewayProfileListItem{
			Id:                gp.GatewayProfileID.String(),
			Name:              gp.Name,
			NetworkServerName: gp.NetworkServerName,
			NetworkServerId:   gp.NetworkServerID,
		}

		row.CreatedAt = timestamppb.New(gp.CreatedAt)
		row.UpdatedAt = timestamppb.New(gp.UpdatedAt)
		out.Result = append(out.Result, &row)
	}

	return &out, nil
}
