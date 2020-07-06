package organization

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
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

	// validator
	CheckReadOrganizationAccess(username string, userID int64, organizationID int64) (bool, error)
	CheckUpdateOrganizationAccess(username string, userID int64, organizationID int64) (bool, error)
	CheckDeleteOrganizationAccess(username string, userID int64, organizationID int64) (bool, error)

	CheckCreateOrganizationAccess(username string, userID int64) (bool, error)
	CheckListOrganizationAccess(username string, userID int64) (bool, error)

	CheckCreateOrganizationUserAccess(username string, userID int64, organizationID int64) (bool, error)
	CheckListOrganizationUserAccess(username string, userID int64, organizationID int64) (bool, error)

	CheckReadOrganizationUserAccess(username string, organizationID int64, userID, operatorUserID int64) (bool, error)
	CheckUpdateOrganizationUserAccess(username string, organizationID int64, userID, operatorUserID int64) (bool, error)
	CheckDeleteOrganizationUserAccess(username string, organizationID int64, userID, operatorUserID int64) (bool, error)
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
		return nil, status.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if valid, err := a.Validator.ValidateOrganizationsAccess(ctx, Create); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	org := Organization{
		Name:            req.Organization.Name,
		DisplayName:     req.Organization.DisplayName,
		CanHaveGateways: req.Organization.CanHaveGateways,
		MaxDeviceCount:  int(req.Organization.MaxDeviceCount),
		MaxGatewayCount: int(req.Organization.MaxGatewayCount),
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.CreateOrganization(ctx, &org)

		return errors.Wrap(err, "")
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrganizationResponse{
		Id: org.ID,
	}, nil
}

// Get returns the organization matching the given ID.
func (a *OrganizationAPI) Get(ctx context.Context, req *pb.GetOrganizationRequest) (*pb.GetOrganizationResponse, error) {
	if valid, err := a.Validator.ValidateOrganizationAccess(ctx, Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
	if valid, err := a.Validator.ValidateOrganizationsAccess(ctx, List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	filters := OrganizationFilters{
		Search: req.Search,
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	u, err := a.Validator.Credentials.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if !u.IsGlobalAdmin {
		filters.UserID = u.ID
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
		return nil, status.Errorf(codes.InvalidArgument, "organization must not be nil")
	}

	if valid, err := a.Validator.ValidateOrganizationAccess(ctx, Update, req.Organization.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := a.Validator.Credentials.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		org, err := a.Store.GetOrganization(ctx, req.Organization.Id, false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		org.Name = req.Organization.Name
		org.DisplayName = req.Organization.DisplayName

		if u.IsGlobalAdmin {
			org.CanHaveGateways = req.Organization.CanHaveGateways
			org.MaxGatewayCount = int(req.Organization.MaxGatewayCount)
			org.MaxDeviceCount = int(req.Organization.MaxDeviceCount)
		}

		err = a.Store.UpdateOrganization(ctx, &org)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the organization matching the given ID.
// Note: this should never happen, when there are still items in the organization, the organization should not be deleted
func (a *OrganizationAPI) Delete(ctx context.Context, req *pb.DeleteOrganizationRequest) (*empty.Empty, error) {
	/*	if err := a.Validator.otpValidator.jwtValidator.Validate(ctx,
			ValidateOrganizationAccess(Delete, req.Id)); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
	if valid, err := a.Validator.ValidateOrganizationUsersAccess(ctx, List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return nil, status.Errorf(codes.InvalidArgument, "organization_user must not be nil")
	}

	if valid, err := a.Validator.ValidateOrganizationUsersAccess(ctx, Create, req.OrganizationUser.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.CreateOrganizationUser(ctx,
			req.OrganizationUser.OrganizationId,
			req.OrganizationUser.Username,
			req.OrganizationUser.IsAdmin,
			req.OrganizationUser.IsDeviceAdmin,
			req.OrganizationUser.IsGatewayAdmin,
		)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
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

	if valid, err := a.Validator.ValidateOrganizationAccess(ctx, Update, req.OrganizationUser.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.UpdateOrganizationUser(ctx,
			req.OrganizationUser.OrganizationId,
			req.OrganizationUser.UserId,
			req.OrganizationUser.IsAdmin,
			req.OrganizationUser.IsDeviceAdmin,
			req.OrganizationUser.IsGatewayAdmin,
		)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// DeleteUser deletes the given user from the organization.
func (a *OrganizationAPI) DeleteUser(ctx context.Context, req *pb.DeleteOrganizationUserRequest) (*empty.Empty, error) {
	if valid, err := a.Validator.ValidateOrganizationAccess(ctx, Delete, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.DeleteOrganizationUser(ctx, req.OrganizationId, req.UserId)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// GetUser returns the user details for the given user ID.
func (a *OrganizationAPI) GetUser(ctx context.Context, req *pb.GetOrganizationUserRequest) (*pb.GetOrganizationUserResponse, error) {
	if valid, err := a.Validator.ValidateOrganizationAccess(ctx, Read, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
