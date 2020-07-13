package user

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByExternalID(ctx context.Context, externalID string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserCount(ctx context.Context) (int, error)
	GetUsers(ctx context.Context, limit, offset int) ([]User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id int64) error
	LoginUserByPassword(ctx context.Context, email string, password string) error
	GetProfile(ctx context.Context, id int64) (UserProfile, error)
	GetUserToken(u User) (string, error)
	RegisterUser(user *User, token string) error
	GetUserByToken(token string) (User, error)
	GetTokenByUsername(ctx context.Context, username string) (string, error)
	FinishRegistration(userID int64, password string) error
	UpdatePassword(ctx context.Context, id int64, newpassword string) error

	// validator
	CheckActiveUser(username string, userID int64) (bool, error)

	CheckCreateUserAcess(username string, userID int64) (bool, error)
	CheckListUserAcess(username string, userID int64) (bool, error)

	CheckReadUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdateDeleteUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdateProfileUserAccess(username string, userID, operatorUserID int64) (bool, error)
	CheckUpdatePasswordUserAccess(username string, userID, operatorUserID int64) (bool, error)
}

// UserAPI exports the User related functions.
type UserAPI struct {
	Validator *Validator
	Store     UserStore
}

// NewUserAPI creates a new UserAPI.
func NewUserAPI(api UserAPI) *UserAPI {
	userAPI = UserAPI{
		Validator: api.Validator,
		Store:     api.Store,
	}

	return &userAPI
}

var (
	userAPI UserAPI
)

func GetUserAPI() *UserAPI {
	return &userAPI
}

// Create creates the given user.
func (a *UserAPI) Create(ctx context.Context, req *inpb.CreateUserRequest) (*inpb.CreateUserResponse, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if valid, err := a.Validator.ValidateUsersGlobalAccess(ctx, Create); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user := User{
		Username:   req.User.Username,
		SessionTTL: req.User.SessionTtl,
		IsAdmin:    req.User.IsAdmin,
		IsActive:   req.User.IsActive,
		Email:      req.User.Email,
		Note:       req.User.Note,
		Password:   req.Password,
	}

	err := a.Store.Transaction(ctx, nil, func() error {
		err := a.Store.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		for _, org := range req.Organizations {
			if err := organization.GetOrganizationAPI().Store.CreateOrganizationUser(ctx, org.OrganizationId, user.Username, org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &inpb.CreateUserResponse{Id: user.ID}, nil
}

// Get returns the user matching the given ID.
func (a *UserAPI) Get(ctx context.Context, req *inpb.GetUserRequest) (*inpb.GetUserResponse, error) {
	if valid, err := a.Validator.ValidateUserAccess(ctx, Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.Store.GetUser(ctx, req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.GetUserResponse{
		User: &inpb.User{
			Id:         user.ID,
			SessionTtl: user.SessionTTL,
			IsAdmin:    user.IsAdmin,
			IsActive:   user.IsActive,
			Email:      user.Email,
			Note:       user.Note,
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

// GetUserEmail returns true if user does not exist
func (a *UserAPI) GetUserEmail(ctx context.Context, req *inpb.GetUserEmailRequest) (*inpb.GetUserEmailResponse, error) {
	u, err := a.Store.GetUserByEmail(ctx, req.UserEmail)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			return &inpb.GetUserEmailResponse{Status: true}, nil
		}
		return nil, helpers.ErrToRPCError(err)
	}
	if u.SecurityToken != nil {
		// user exists but has not finished registration
		return &inpb.GetUserEmailResponse{Status: true}, nil
	}

	return &inpb.GetUserEmailResponse{Status: false}, nil
}

// List lists the users.
func (a *UserAPI) List(ctx context.Context, req *inpb.ListUserRequest) (*inpb.ListUserResponse, error) {
	if valid, err := a.Validator.ValidateUsersGlobalAccess(ctx, List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := a.Store.GetUsers(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := a.Store.GetUserCount(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ListUserResponse{
		TotalCount: int64(totalUserCount),
	}

	for _, u := range users {
		row := inpb.UserListItem{
			Id:         u.ID,
			SessionTtl: u.SessionTTL,
			IsAdmin:    u.IsAdmin,
			IsActive:   u.IsActive,
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

// Update updates the given user.
func (a *UserAPI) Update(ctx context.Context, req *inpb.UpdateUserRequest) (*empty.Empty, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if valid, err := a.Validator.ValidateUserAccess(ctx, Update, req.User.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		user, err := a.Store.GetUser(ctx, req.User.Id)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		user.IsAdmin = req.User.IsAdmin
		user.IsActive = req.User.IsActive
		user.SessionTTL = req.User.SessionTtl
		user.Email = req.User.Email
		user.Note = req.User.Note

		err = a.Store.UpdateUser(ctx, &user)
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

// Delete deletes the user matching the given ID.
func (a *UserAPI) Delete(ctx context.Context, req *inpb.DeleteUserRequest) (*empty.Empty, error) {
	if valid, err := a.Validator.ValidateUserAccess(ctx, Delete, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		err := a.Store.DeleteUser(ctx, req.Id)
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

// UpdatePassword updates the password for the user matching the given ID.
func (a *UserAPI) UpdatePassword(ctx context.Context, req *inpb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	if valid, err := a.Validator.ValidateUserAccess(ctx, UpdatePassword, req.UserId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.Store.UpdatePassword(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

func (a *UserAPI) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	otp, err := a.Store.GetTokenByUsername(ctx, req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &inpb.GetOTPCodeResponse{OtpCode: otp}, nil
}
