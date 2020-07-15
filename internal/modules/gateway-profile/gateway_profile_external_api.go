package gatewayprofile

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// GatewayProfileAPI exports the GatewayProfile related functions.
type GatewayProfileAPI struct {
	st   GatewayProfileStore
	txSt store.Store
}

// NewGatewayProfileAPI creates a new GatewayProfileAPI.
func NewGatewayProfileAPI() *GatewayProfileAPI {
	st := store.New(storage.DB().DB)
	return &GatewayProfileAPI{
		st:   st,
		txSt: st,
	}
}

// Create creates the given gateway-profile.
func (a *GatewayProfileAPI) Create(ctx context.Context, req *pb.CreateGatewayProfileRequest) (*pb.CreateGatewayProfileResponse, error) {
	if req.GatewayProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "gateway_profile must not be nil")
	}

	if valid, err := NewValidator().ValidateGatewayProfileAccess(ctx, authcus.Create); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gp := GatewayProfile{
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

	tx, err := a.txSt.TxBegin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	defer tx.TxRollback(ctx)

	err = tx.CreateGatewayProfile(ctx, &gp)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if err := tx.TxCommit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	gpID, err := uuid.FromBytes(gp.GatewayProfile.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateGatewayProfileResponse{
		Id: gpID.String(),
	}, nil
}

// Get returns the gateway-profile matching the given id.
func (a *GatewayProfileAPI) Get(ctx context.Context, req *pb.GetGatewayProfileRequest) (*pb.GetGatewayProfileResponse, error) {
	if valid, err := NewValidator().ValidateGatewayProfileAccess(ctx, authcus.Read); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	gp, err := a.st.GetGatewayProfile(ctx, gpID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	out := pb.GetGatewayProfileResponse{
		GatewayProfile: &pb.GatewayProfile{
			Id:              gpID.String(),
			Name:            gp.Name,
			NetworkServerId: gp.NetworkServerID,
			Channels:        gp.GatewayProfile.Channels,
		},
	}

	out.CreatedAt, err = ptypes.TimestampProto(gp.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	out.UpdatedAt, err = ptypes.TimestampProto(gp.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

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

	if valid, err := NewValidator().ValidateGatewayProfileAccess(ctx, authcus.Update); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gpID, err := uuid.FromString(req.GatewayProfile.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	gp, err := a.st.GetGatewayProfile(ctx, gpID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
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

	tx, err := a.txSt.TxBegin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	defer tx.TxRollback(ctx)

	err = tx.UpdateGatewayProfile(ctx, &gp)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if err := tx.TxCommit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the gateway-profile matching the given id.
func (a *GatewayProfileAPI) Delete(ctx context.Context, req *pb.DeleteGatewayProfileRequest) (*empty.Empty, error) {
	if valid, err := NewValidator().ValidateGatewayProfileAccess(ctx, authcus.Delete); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	err = a.st.DeleteGatewayProfile(ctx, gpID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// List returns the existing gateway-profiles.
func (a *GatewayProfileAPI) List(ctx context.Context, req *pb.ListGatewayProfilesRequest) (*pb.ListGatewayProfilesResponse, error) {
	if valid, err := NewValidator().ValidateGatewayProfileAccess(ctx, authcus.List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var err error
	var count int
	var gps []GatewayProfileMeta

	if req.NetworkServerId == 0 {
		count, err = a.st.GetGatewayProfileCount(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		gps, err = a.st.GetGatewayProfiles(ctx, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	} else {
		count, err = a.st.GetGatewayProfileCountForNetworkServerID(ctx, req.NetworkServerId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		gps, err = a.st.GetGatewayProfilesForNetworkServerID(ctx, req.NetworkServerId, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
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

		row.CreatedAt, err = ptypes.TimestampProto(gp.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(gp.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		out.Result = append(out.Result, &row)
	}

	return &out, nil
}
