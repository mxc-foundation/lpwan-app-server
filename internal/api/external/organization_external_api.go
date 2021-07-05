package external

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// OrganizationAPI exports the organization related functions.
type OrganizationAPI struct {
	st    *store.Handler
	nsCli *nscli.Client
}

// NewOrganizationAPI creates a new OrganizationAPI.
func NewOrganizationAPI(h *store.Handler, nsCli *nscli.Client) *OrganizationAPI {
	return &OrganizationAPI{
		st:    h,
		nsCli: nsCli,
	}
}

// Create creates the given organization.
func (a *OrganizationAPI) Create(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	if req.Organization == nil {
		return nil, status.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if valid, err := organization.NewValidator(a.st).ValidateOrganizationsAccess(ctx, auth.Create); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org := organization.Organization{
		Name:            req.Organization.Name,
		DisplayName:     req.Organization.DisplayName,
		CanHaveGateways: req.Organization.CanHaveGateways,
		MaxDeviceCount:  int(req.Organization.MaxDeviceCount),
		MaxGatewayCount: int(req.Organization.MaxGatewayCount),
	}

	err := a.st.CreateOrganization(ctx, &org)
	if err != nil {
		return nil, err
	}
	// set all default settings for new organization, this step should not interrupt creating organization
	organization.ActivateOrganization(ctx, a.st, a.st, a.st, org.ID, a.nsCli)
	return &pb.CreateOrganizationResponse{
		Id: org.ID,
	}, nil
}

// Get returns the organization matching the given ID.
func (a *OrganizationAPI) Get(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org, err := a.st.GetOrganization(ctx, req.Id, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetOrganizationResponse{
		Organization: &pb.Organization{
			Id:              org.ID,
			Name:            org.Name,
			DisplayName:     org.DisplayName,
			CanHaveGateways: org.CanHaveGateways,
			MaxDeviceCount:  uint32(org.MaxDeviceCount),
			MaxGatewayCount: uint32(org.MaxGatewayCount),
		},
	}

	resp.CreatedAt = timestamppb.New(org.CreatedAt)
	resp.UpdatedAt = timestamppb.New(org.UpdatedAt)

	return &resp, nil
}

// List lists the organizations to which the user has access.
func (a *OrganizationAPI) List(ctx context.Context, req *pb.ListOrganizationRequest) (*pb.ListOrganizationResponse, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationsAccess(ctx, auth.List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	filters := organization.OrganizationFilters{
		Search: req.Search,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	u, err := organization.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if !u.IsGlobalAdmin {
		filters.UserID = u.ID
	}

	count, err := a.st.GetOrganizationCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	orgs, err := a.st.GetOrganizations(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListOrganizationResponse{
		TotalCount: int64(count),
	}

	for _, org := range orgs {
		row := pb.OrganizationListItem{
			Id:              org.ID,
			Name:            org.Name,
			DisplayName:     org.DisplayName,
			CanHaveGateways: org.CanHaveGateways,
		}

		row.CreatedAt = timestamppb.New(org.CreatedAt)
		row.UpdatedAt = timestamppb.New(org.UpdatedAt)

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// Update updates the given organization.
func (a *OrganizationAPI) Update(ctx context.Context, req *pb.UpdateOrganizationRequest) (*empty.Empty, error) {
	if req.Organization == nil {
		return nil, status.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Update, req.Organization.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := organization.NewValidator(a.st).GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	org, err := a.st.GetOrganization(ctx, req.Organization.Id, false)
	if err != nil {
		return nil, err
	}

	org.Name = req.Organization.Name
	org.DisplayName = req.Organization.DisplayName

	if u.IsGlobalAdmin {
		org.CanHaveGateways = req.Organization.CanHaveGateways
		org.MaxGatewayCount = int(req.Organization.MaxGatewayCount)
		org.MaxDeviceCount = int(req.Organization.MaxDeviceCount)
	}

	err = a.st.UpdateOrganization(ctx, &org)

	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the organization matching the given ID.
// Note: this should never happen, when there are still items in the organization, the organization should not be deleted
func (a *OrganizationAPI) Delete(ctx context.Context, req *pb.DeleteOrganizationRequest) (*empty.Empty, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Delete, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		if err := handler.DeleteAllGatewaysForOrganizationID(ctx, req.Id); err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		err := handler.DeleteAllApplicationsForOrganizationID(ctx, req.Id)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		err = handler.DeleteAllServiceProfilesForOrganizationID(ctx, req.Id)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		err = handler.DeleteAllDeviceProfilesForOrganizationID(ctx, req.Id)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		if err := handler.DeleteOrganization(ctx, req.Id); err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &empty.Empty{}, nil
}

// ListUsers lists the users assigned to the given organization.
func (a *OrganizationAPI) ListUsers(ctx context.Context, req *pb.ListOrganizationUsersRequest) (*pb.ListOrganizationUsersResponse, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationUsersAccess(ctx, auth.List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := a.st.GetOrganizationUsers(ctx, req.OrganizationId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	userCount, err := a.st.GetOrganizationUserCount(ctx, req.OrganizationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListOrganizationUsersResponse{
		TotalCount: int64(userCount),
	}

	for _, u := range users {
		row := pb.OrganizationUserListItem{
			Username:       u.Email,
			UserId:         u.UserID,
			IsAdmin:        u.IsAdmin,
			IsDeviceAdmin:  u.IsDeviceAdmin,
			IsGatewayAdmin: u.IsGatewayAdmin,
		}

		row.CreatedAt = timestamppb.New(u.CreatedAt)
		row.UpdatedAt = timestamppb.New(u.UpdatedAt)

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// AddUser creates the given organization-user link.
func (a *OrganizationAPI) AddUser(ctx context.Context, req *pb.AddOrganizationUserRequest) (*empty.Empty, error) {
	if req.OrganizationUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "organization_user must not be nil")
	}

	if valid, err := organization.NewValidator(a.st).ValidateOrganizationUsersAccess(ctx, auth.Create, req.OrganizationUser.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.CreateOrganizationUser(ctx,
		req.OrganizationUser.OrganizationId,
		req.OrganizationUser.UserId,
		req.OrganizationUser.IsAdmin,
		req.OrganizationUser.IsDeviceAdmin,
		req.OrganizationUser.IsGatewayAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// UpdateUser updates the given user.
func (a *OrganizationAPI) UpdateUser(ctx context.Context, req *pb.UpdateOrganizationUserRequest) (*empty.Empty, error) {
	if req.OrganizationUser == nil {
		return nil, status.Errorf(codes.InvalidArgument, "organization_user must not be nil")
	}

	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Update, req.OrganizationUser.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.UpdateOrganizationUser(ctx,
		req.OrganizationUser.OrganizationId,
		req.OrganizationUser.UserId,
		req.OrganizationUser.IsAdmin,
		req.OrganizationUser.IsDeviceAdmin,
		req.OrganizationUser.IsGatewayAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// DeleteUser deletes the given user from the organization.
func (a *OrganizationAPI) DeleteUser(ctx context.Context, req *pb.DeleteOrganizationUserRequest) (*empty.Empty, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Delete, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.DeleteOrganizationUser(ctx, req.OrganizationId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// GetUser returns the user details for the given user ID.
func (a *OrganizationAPI) GetUser(ctx context.Context, req *pb.GetOrganizationUserRequest) (*pb.GetOrganizationUserResponse, error) {
	if valid, err := organization.NewValidator(a.st).ValidateOrganizationAccess(ctx, auth.Read, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := a.st.GetOrganizationUser(ctx, req.OrganizationId, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetOrganizationUserResponse{
		OrganizationUser: &pb.OrganizationUser{
			Username:       u.Email,
			OrganizationId: req.OrganizationId,
			UserId:         req.UserId,
			IsAdmin:        u.IsAdmin,
			IsDeviceAdmin:  u.IsDeviceAdmin,
			IsGatewayAdmin: u.IsGatewayAdmin,
		},
	}

	resp.CreatedAt = timestamppb.New(u.CreatedAt)
	resp.UpdatedAt = timestamppb.New(u.UpdatedAt)

	return &resp, nil
}
