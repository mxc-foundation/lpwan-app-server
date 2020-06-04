package external

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// OrganizationAPI exports the organization related functions.
type OrganizationAPI struct {
	validator auth.Validator
}

// NewOrganizationAPI creates a new OrganizationAPI.
func NewOrganizationAPI(validator auth.Validator) *OrganizationAPI {
	return &OrganizationAPI{
		validator: validator,
	}
}

// Create creates the given organization.
func (a *OrganizationAPI) Create(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	if req.Organization == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsGlobalAdmin(ctx); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be a global admin")
	}

	org := storage.Organization{
		Name:            req.Organization.Name,
		DisplayName:     req.Organization.DisplayName,
		CanHaveGateways: req.Organization.CanHaveGateways,
	}

	err = storage.CreateOrganization(ctx, storage.DB(), &org)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateOrganizationResponse{
		Id: org.ID,
	}, nil
}

// Get returns the organization matching the given ID.
func (a *OrganizationAPI) Get(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgUser(ctx, req.Id); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization user")
	}

	org, err := storage.GetOrganization(ctx, storage.DB(), req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetOrganizationResponse{
		Organization: &pb.Organization{
			Id:              org.ID,
			Name:            org.Name,
			DisplayName:     org.DisplayName,
			CanHaveGateways: org.CanHaveGateways,
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}

// List lists the organizations to which the user has access.
func (a *OrganizationAPI) List(ctx context.Context, req *pb.ListOrganizationRequest) (*pb.ListOrganizationResponse, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	var count int
	var orgs []storage.Organization

	if err := cred.IsGlobalAdmin(ctx); err == nil {
		count, err = storage.GetOrganizationCount(ctx, storage.DB(), req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		orgs, err = storage.GetOrganizations(ctx, storage.DB(), int(req.Limit), int(req.Offset), req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	} else {
		count, err = storage.GetOrganizationCountForUser(ctx, storage.DB(), cred.Username(), req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		orgs, err = storage.GetOrganizationsForUser(ctx, storage.DB(), cred.Username(), int(req.Limit), int(req.Offset), req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
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

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// Update updates the given organization.
func (a *OrganizationAPI) Update(ctx context.Context, req *pb.UpdateOrganizationRequest) (*empty.Empty, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	if req.Organization == nil {
		return nil, status.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if err := cred.IsOrgAdmin(ctx, req.Organization.Id); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization admin")
	}

	org, err := storage.GetOrganization(ctx, storage.DB(), req.Organization.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	org.Name = req.Organization.Name
	org.DisplayName = req.Organization.DisplayName
	if cred.IsGlobalAdmin(ctx) == nil {
		org.CanHaveGateways = req.Organization.CanHaveGateways
	}

	err = storage.UpdateOrganization(ctx, storage.DB(), &org)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the organization matching the given ID.
func (a *OrganizationAPI) Delete(ctx context.Context, req *pb.DeleteOrganizationRequest) (*empty.Empty, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsGlobalAdmin(ctx); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be a global admin")
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		if err := storage.DeleteAllGatewaysForOrganizationID(ctx, tx, req.Id); err != nil {
			return helpers.ErrToRPCError(err)
		}

		if err := storage.DeleteOrganization(ctx, tx, req.Id); err != nil {
			return helpers.ErrToRPCError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ListUsers lists the users assigned to the given organization.
func (a *OrganizationAPI) ListUsers(ctx context.Context, req *pb.ListOrganizationUsersRequest) (*pb.ListOrganizationUsersResponse, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgUser(ctx, req.OrganizationId); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization user")
	}

	users, err := storage.GetOrganizationUsers(ctx, storage.DB(), req.OrganizationId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	userCount, err := storage.GetOrganizationUserCount(ctx, storage.DB(), req.OrganizationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListOrganizationUsersResponse{
		TotalCount: int64(userCount),
	}

	for _, u := range users {
		row := pb.OrganizationUserListItem{
			UserId:         u.UserID,
			Username:       u.Username,
			IsAdmin:        u.IsAdmin,
			IsDeviceAdmin:  u.IsDeviceAdmin,
			IsGatewayAdmin: u.IsGatewayAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(u.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(u.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// AddUser creates the given organization-user link.
func (a *OrganizationAPI) AddUser(ctx context.Context, req *pb.AddOrganizationUserRequest) (*empty.Empty, error) {
	if req.OrganizationUser == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization_user must not be nil")
	}

	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgAdmin(ctx, req.OrganizationUser.OrganizationId); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization admin")
	}

	err = storage.CreateOrganizationUser(ctx,
		storage.DB(),
		req.OrganizationUser.OrganizationId,
		req.OrganizationUser.UserId,
		req.OrganizationUser.IsAdmin,
		req.OrganizationUser.IsDeviceAdmin,
		req.OrganizationUser.IsGatewayAdmin,
	)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// UpdateUser updates the given user.
func (a *OrganizationAPI) UpdateUser(ctx context.Context, req *pb.UpdateOrganizationUserRequest) (*empty.Empty, error) {
	if req.OrganizationUser == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization_user must not be nil")
	}

	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgAdmin(ctx, req.OrganizationUser.OrganizationId); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization admin")
	}

	err = storage.UpdateOrganizationUser(ctx,
		storage.DB(),
		req.OrganizationUser.OrganizationId,
		req.OrganizationUser.UserId,
		req.OrganizationUser.IsAdmin,
		req.OrganizationUser.IsDeviceAdmin,
		req.OrganizationUser.IsGatewayAdmin,
	)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// DeleteUser deletes the given user from the organization.
func (a *OrganizationAPI) DeleteUser(ctx context.Context, req *pb.DeleteOrganizationUserRequest) (*empty.Empty, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgAdmin(ctx, req.OrganizationId); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization admin")
	}

	err = storage.DeleteOrganizationUser(ctx, storage.DB(), req.OrganizationId, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetUser returns the user details for the given user ID.
func (a *OrganizationAPI) GetUser(ctx context.Context, req *pb.GetOrganizationUserRequest) (*pb.GetOrganizationUserResponse, error) {
	cred, err := a.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if cred.UserID() != req.UserId {
		if err := cred.IsOrgAdmin(ctx, req.OrganizationId); err != nil {
			return nil, status.Error(codes.PermissionDenied, "must be user themselves or an organization admin")
		}
	}

	user, err := storage.GetOrganizationUser(ctx, storage.DB(), req.OrganizationId, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetOrganizationUserResponse{
		OrganizationUser: &pb.OrganizationUser{
			OrganizationId: req.OrganizationId,
			UserId:         req.UserId,
			IsAdmin:        user.IsAdmin,
			IsDeviceAdmin:  user.IsDeviceAdmin,
			IsGatewayAdmin: user.IsGatewayAdmin,
			Username:       user.Username,
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}
