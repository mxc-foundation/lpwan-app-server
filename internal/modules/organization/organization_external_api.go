package organization

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
)

type OrganizationStore interface {
	CreateOrganization(ctx context.Context, org *Organization) error
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (Organization, error)
	GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error)
	GetOrganizations(ctx context.Context, filters OrganizationFilters) ([]Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) error
	DeleteOrganization(ctx context.Context, id int64) error
	CreateOrganizationUser(ctx context.Context, organizationID int64, username string, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error
	UpdateOrganizationUser(ctx context.Context, organizationID, userID int64, isAdmin, isDeviceAdmin, isGatewayAdmin bool) error
	DeleteOrganizationUser(ctx context.Context, organizationID, userID int64) error
	GetOrganizationUser(ctx context.Context, organizationID, userID int64) (OrganizationUser, error)
	GetOrganizationUserCount(ctx context.Context, organizationID int64) (int, error)
	GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]OrganizationUser, error)
	GetOrganizationIDList(limit, offset int, search string) ([]int, error)
}

// OrganizationAPI exports the organization related functions.
type OrganizationAPI struct {
	Validator *Validator
	Store     OrganizationStore
}

// NewOrganizationAPI creates a new OrganizationAPI.
func NewOrganizationAPI(api OrganizationAPI) *OrganizationAPI {
	organizationAPI = OrganizationAPI{
		Validator: api.Validator,
		Store:     api.Store,
	}
	return &organizationAPI
}

var (
	organizationAPI OrganizationAPI
)

func GetOrganizationAPI() *OrganizationAPI {
	return &organizationAPI
}

// Create creates the given organization.
func (a *OrganizationAPI) Create(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	if req.Organization == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationsAccess(Create)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org := Organization{
		Name:            req.Organization.Name,
		DisplayName:     req.Organization.DisplayName,
		CanHaveGateways: req.Organization.CanHaveGateways,
		MaxDeviceCount:  int(req.Organization.MaxDeviceCount),
		MaxGatewayCount: int(req.Organization.MaxGatewayCount),
	}

	err := a.Store.CreateOrganization(ctx, &org)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateOrganizationResponse{
		Id: org.ID,
	}, nil
}

// Get returns the organization matching the given ID.
func (a *OrganizationAPI) Get(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationAccess(Read, req.Id)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org, err := a.Store.GetOrganization(ctx, req.Id, false)
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
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationsAccess(List)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	filters := OrganizationFilters{
		Search: req.Search,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	sub, err := a.Validator.otpValidator.JwtValidator.GetSubject(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	switch sub {
	case SubjectUser:
		username, err := a.Validator.otpValidator.JwtValidator.GetUsername(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		filters.UserName = username
		/*		u, err := user.GetUserAPI().Validator.GetUser(ctx)
				if err != nil {
					return nil, helpers.ErrToRPCError(err)
				}

				if !u.IsAdmin {
					filters.UserID = u.ID
				}*/
	case SubjectAPIKey:
		// Nothing to do as the Validator function already validated that the
		// API key must be a global admin key.
	default:
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token subject: %s", err)
	}

	count, err := a.Store.GetOrganizationCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	orgs, err := a.Store.GetOrganizations(ctx, filters)
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
	if req.Organization == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationAccess(Update, req.Organization.Id)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	/*	u, err := user.GetUserAPI().Validator.GetUser(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}*/

	org, err := a.Store.GetOrganization(ctx, req.Organization.Id, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	org.Name = req.Organization.Name
	org.DisplayName = req.Organization.DisplayName
	/*	if u.IsAdmin {
		org.CanHaveGateways = req.Organization.CanHaveGateways
		org.MaxGatewayCount = int(req.Organization.MaxGatewayCount)
		org.MaxDeviceCount = int(req.Organization.MaxDeviceCount)
	}*/

	err = a.Store.UpdateOrganization(ctx, &org)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the organization matching the given ID.
// Note: this should never happen, when there are still items in the organization, the organization should not be deleted
func (a *OrganizationAPI) Delete(ctx context.Context, req *pb.DeleteOrganizationRequest) (*empty.Empty, error) {
	/*	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
			ValidateOrganizationAccess(Delete, req.Id)); err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}

		err := storage.Transaction(func(tx sqlx.Ext) error {
			if err := gateway.GetGatewayAPI().Store.DeleteAllGatewaysForOrganizationID(ctx, req.Id); err != nil {
				return helpers.ErrToRPCError(err)
			}

			if err := a.Store.DeleteOrganization(ctx, req.Id); err != nil {
				return helpers.ErrToRPCError(err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}*/

	return &empty.Empty{}, nil
}

// ListUsers lists the users assigned to the given organization.
func (a *OrganizationAPI) ListUsers(ctx context.Context, req *pb.ListOrganizationUsersRequest) (*pb.ListOrganizationUsersResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationUsersAccess(List, req.OrganizationId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := a.Store.GetOrganizationUsers(ctx, req.OrganizationId, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	userCount, err := a.Store.GetOrganizationUserCount(ctx, req.OrganizationId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListOrganizationUsersResponse{
		TotalCount: int64(userCount),
	}

	for _, u := range users {
		row := pb.OrganizationUserListItem{
			UserId:         u.UserID,
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

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationUsersAccess(Create, req.OrganizationUser.OrganizationId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.Store.CreateOrganizationUser(ctx,
		req.OrganizationUser.OrganizationId,
		req.OrganizationUser.Username,
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

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationUserAccess(Update, req.OrganizationUser.OrganizationId, req.OrganizationUser.UserId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.Store.UpdateOrganizationUser(ctx,
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
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationUserAccess(Delete, req.OrganizationId, req.UserId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.Store.DeleteOrganizationUser(ctx, req.OrganizationId, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// GetUser returns the user details for the given user ID.
func (a *OrganizationAPI) GetUser(ctx context.Context, req *pb.GetOrganizationUserRequest) (*pb.GetOrganizationUserResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx,
		ValidateOrganizationUserAccess(Read, req.OrganizationId, req.UserId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := a.Store.GetOrganizationUser(ctx, req.OrganizationId, req.UserId)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetOrganizationUserResponse{
		OrganizationUser: &pb.OrganizationUser{
			OrganizationId: req.OrganizationId,
			UserId:         req.UserId,
			IsAdmin:        u.IsAdmin,
			IsDeviceAdmin:  u.IsDeviceAdmin,
			IsGatewayAdmin: u.IsGatewayAdmin,
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(u.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(u.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}
