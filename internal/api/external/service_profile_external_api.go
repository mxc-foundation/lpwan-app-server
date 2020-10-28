package external

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"

	spmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ServiceProfileServiceAPI export the ServiceProfile related functions.
type ServiceProfileServiceAPI struct {
	st *store.Handler
}

// NewServiceProfileServiceAPI creates a new ServiceProfileServiceAPI.
func NewServiceProfileServiceAPI(h *store.Handler) *ServiceProfileServiceAPI {
	return &ServiceProfileServiceAPI{
		st: h,
	}
}

// Create creates the given service-profile.
func (a *ServiceProfileServiceAPI) Create(ctx context.Context, req *pb.CreateServiceProfileRequest) (*pb.CreateServiceProfileResponse, error) {
	if req.ServiceProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "service_profile must not be nil")
	}

	if valid, err := spmod.NewValidator().ValidateServiceProfilesAccess(ctx, auth.Create, req.ServiceProfile.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	sp := ServiceProfile{
		OrganizationID:  req.ServiceProfile.OrganizationId,
		NetworkServerID: req.ServiceProfile.NetworkServerId,
		Name:            req.ServiceProfile.Name,
		ServiceProfile: ns.ServiceProfile{
			UlRate:                 req.ServiceProfile.UlRate,
			UlBucketSize:           req.ServiceProfile.UlBucketSize,
			DlRate:                 req.ServiceProfile.DlRate,
			DlBucketSize:           req.ServiceProfile.DlBucketSize,
			AddGwMetadata:          req.ServiceProfile.AddGwMetadata,
			DevStatusReqFreq:       req.ServiceProfile.DevStatusReqFreq,
			ReportDevStatusBattery: req.ServiceProfile.ReportDevStatusBattery,
			ReportDevStatusMargin:  req.ServiceProfile.ReportDevStatusMargin,
			DrMin:                  req.ServiceProfile.DrMin,
			DrMax:                  req.ServiceProfile.DrMax,
			ChannelMask:            req.ServiceProfile.ChannelMask,
			PrAllowed:              req.ServiceProfile.PrAllowed,
			HrAllowed:              req.ServiceProfile.HrAllowed,
			RaAllowed:              req.ServiceProfile.RaAllowed,
			NwkGeoLoc:              req.ServiceProfile.NwkGeoLoc,
			TargetPer:              req.ServiceProfile.TargetPer,
			MinGwDiversity:         req.ServiceProfile.MinGwDiversity,
			UlRatePolicy:           ns.RatePolicy(req.ServiceProfile.UlRatePolicy),
			DlRatePolicy:           ns.RatePolicy(req.ServiceProfile.DlRatePolicy),
		},
	}

	// as this also performs a remote call to create the service-profile
	// on the network-server, wrap it in a transaction
	if err := spmod.CreateServiceProfile(ctx, &sp); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	spID, err := uuid.FromBytes(sp.ServiceProfile.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateServiceProfileResponse{
		Id: spID.String(),
	}, nil
}

// Get returns the service-profile matching the given id.
func (a *ServiceProfileServiceAPI) Get(ctx context.Context, req *pb.GetServiceProfileRequest) (*pb.GetServiceProfileResponse, error) {
	spID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := spmod.NewValidator().ValidateServiceProfileAccess(ctx, auth.Read, spID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	sp, err := spmod.GetServiceProfile(ctx, spID, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetServiceProfileResponse{
		ServiceProfile: &pb.ServiceProfile{
			Id:                     spID.String(),
			Name:                   sp.Name,
			OrganizationId:         sp.OrganizationID,
			NetworkServerId:        sp.NetworkServerID,
			UlRate:                 sp.ServiceProfile.UlRate,
			UlBucketSize:           sp.ServiceProfile.UlBucketSize,
			DlRate:                 sp.ServiceProfile.DlRate,
			DlBucketSize:           sp.ServiceProfile.DlBucketSize,
			AddGwMetadata:          sp.ServiceProfile.AddGwMetadata,
			DevStatusReqFreq:       sp.ServiceProfile.DevStatusReqFreq,
			ReportDevStatusBattery: sp.ServiceProfile.ReportDevStatusBattery,
			ReportDevStatusMargin:  sp.ServiceProfile.ReportDevStatusMargin,
			DrMin:                  sp.ServiceProfile.DrMin,
			DrMax:                  sp.ServiceProfile.DrMax,
			ChannelMask:            sp.ServiceProfile.ChannelMask,
			PrAllowed:              sp.ServiceProfile.PrAllowed,
			HrAllowed:              sp.ServiceProfile.HrAllowed,
			RaAllowed:              sp.ServiceProfile.RaAllowed,
			NwkGeoLoc:              sp.ServiceProfile.NwkGeoLoc,
			TargetPer:              sp.ServiceProfile.TargetPer,
			MinGwDiversity:         sp.ServiceProfile.MinGwDiversity,
			UlRatePolicy:           pb.RatePolicy(sp.ServiceProfile.UlRatePolicy),
			DlRatePolicy:           pb.RatePolicy(sp.ServiceProfile.DlRatePolicy),
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(sp.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(sp.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}

// Update updates the given serviceprofile.
func (a *ServiceProfileServiceAPI) Update(ctx context.Context, req *pb.UpdateServiceProfileRequest) (*empty.Empty, error) {
	if req.ServiceProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "service_profile must not be nil")
	}

	spID, err := uuid.FromString(req.ServiceProfile.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := spmod.NewValidator().ValidateServiceProfileAccess(ctx, auth.Update, spID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	sp, err := spmod.GetServiceProfile(ctx, spID, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	sp.Name = req.ServiceProfile.Name
	sp.ServiceProfile = ns.ServiceProfile{
		Id:                     spID.Bytes(),
		UlRate:                 req.ServiceProfile.UlRate,
		UlBucketSize:           req.ServiceProfile.UlBucketSize,
		DlRate:                 req.ServiceProfile.DlRate,
		DlBucketSize:           req.ServiceProfile.DlBucketSize,
		AddGwMetadata:          req.ServiceProfile.AddGwMetadata,
		DevStatusReqFreq:       req.ServiceProfile.DevStatusReqFreq,
		ReportDevStatusBattery: req.ServiceProfile.ReportDevStatusBattery,
		ReportDevStatusMargin:  req.ServiceProfile.ReportDevStatusMargin,
		DrMin:                  req.ServiceProfile.DrMin,
		DrMax:                  req.ServiceProfile.DrMax,
		ChannelMask:            req.ServiceProfile.ChannelMask,
		PrAllowed:              req.ServiceProfile.PrAllowed,
		HrAllowed:              req.ServiceProfile.HrAllowed,
		RaAllowed:              req.ServiceProfile.RaAllowed,
		NwkGeoLoc:              req.ServiceProfile.NwkGeoLoc,
		TargetPer:              req.ServiceProfile.TargetPer,
		MinGwDiversity:         req.ServiceProfile.MinGwDiversity,
		UlRatePolicy:           ns.RatePolicy(req.ServiceProfile.UlRatePolicy),
		DlRatePolicy:           ns.RatePolicy(req.ServiceProfile.DlRatePolicy),
	}

	// as this also performs a remote call to create the service-profile
	// on the network-server, wrap it in a transaction
	if err = spmod.UpdateServiceProfile(ctx, &sp); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the service-profile matching the given id.
func (a *ServiceProfileServiceAPI) Delete(ctx context.Context, req *pb.DeleteServiceProfileRequest) (*empty.Empty, error) {
	spID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := spmod.NewValidator().ValidateServiceProfileAccess(ctx, auth.Delete, spID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err = spmod.DeleteServiceProfile(ctx, spID); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// List lists the available service-profiles.
func (a *ServiceProfileServiceAPI) List(ctx context.Context, req *pb.ListServiceProfileRequest) (*pb.ListServiceProfileResponse, error) {
	if valid, err := spmod.NewValidator().ValidateServiceProfilesAccess(ctx, auth.List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := spmod.NewValidator().GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	var count int
	var sps []ServiceProfileMeta

	if req.OrganizationId == 0 {
		if user.IsGlobalAdmin {
			sps, err = a.st.GetServiceProfiles(ctx, int(req.Limit), int(req.Offset))
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}

			count, err = a.st.GetServiceProfileCount(ctx)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		} else {
			sps, err = a.st.GetServiceProfilesForUser(ctx, user.ID, int(req.Limit), int(req.Offset))
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}

			count, err = a.st.GetServiceProfileCountForUser(ctx, user.ID)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}
	} else {
		sps, err = a.st.GetServiceProfilesForOrganizationID(ctx, req.OrganizationId, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		count, err = a.st.GetServiceProfileCountForOrganizationID(ctx, req.OrganizationId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	resp := pb.ListServiceProfileResponse{
		TotalCount: int64(count),
	}
	for _, sp := range sps {
		row := pb.ServiceProfileListItem{
			Id:                sp.ServiceProfileID.String(),
			Name:              sp.Name,
			OrganizationId:    sp.OrganizationID,
			NetworkServerId:   sp.NetworkServerID,
			NetworkServerName: sp.NetworkServerName,
		}

		row.CreatedAt, err = ptypes.TimestampProto(sp.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(sp.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}
