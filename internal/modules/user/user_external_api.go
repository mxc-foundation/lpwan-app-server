package user

import (
	"context"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// UserAPI exports the User related functions.
type UserAPI struct {
	st   UserStore
	txSt store.Store
}

// NewUserAPI creates a new UserAPI.
func NewUserAPI() *UserAPI {
	st := store.New(storage.DB().DB)
	return &UserAPI{
		st:   st,
		txSt: st,
	}
}

// Create creates the given user.
func (a *UserAPI) Create(ctx context.Context, req *inpb.CreateUserRequest) (*inpb.CreateUserResponse, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}

	if valid, err := NewValidator().ValidateUsersGlobalAccess(ctx, authcus.Create); !valid || err != nil {
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

	tx, err := a.txSt.TxBegin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't start transaction: %v", err)
	}
	defer tx.TxRollback(ctx)

	err = tx.CreateUser(ctx, &user)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	for _, org := range req.Organizations {
		if err := tx.CreateOrganizationUser(ctx, org.OrganizationId, user.Username, org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
			return nil, status.Errorf(codes.Unknown, "%v", err)
		}
	}

	if err := tx.TxCommit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't commit transaction: %v", err)
	}

	return &inpb.CreateUserResponse{Id: user.ID}, nil
}

// Get returns the user matching the given ID.
func (a *UserAPI) Get(ctx context.Context, req *inpb.GetUserRequest) (*inpb.GetUserResponse, error) {
	if valid, err := NewValidator().ValidateUserAccess(ctx, authcus.Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.st.GetUser(ctx, req.Id)
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
	u, err := a.st.GetUserByEmail(ctx, req.UserEmail)
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
	if valid, err := NewValidator().ValidateUsersGlobalAccess(ctx, authcus.List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	users, err := a.st.GetUsers(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := a.st.GetUserCount(ctx)
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

	if valid, err := NewValidator().ValidateUserAccess(ctx, authcus.Update, req.User.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	user, err := a.st.GetUser(ctx, req.User.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	user.IsAdmin = req.User.IsAdmin
	user.IsActive = req.User.IsActive
	user.SessionTTL = req.User.SessionTtl
	user.Email = req.User.Email
	user.Note = req.User.Note

	err = a.st.UpdateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the user matching the given ID.
func (a *UserAPI) Delete(ctx context.Context, req *inpb.DeleteUserRequest) (*empty.Empty, error) {
	if valid, err := NewValidator().ValidateUserAccess(ctx, authcus.Delete, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// UpdatePassword updates the password for the user matching the given ID.
func (a *UserAPI) UpdatePassword(ctx context.Context, req *inpb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	if valid, err := NewValidator().ValidateUserAccess(ctx, authcus.UpdatePassword, req.UserId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.UpdatePassword(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

func (a *UserAPI) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	otp, err := a.st.GetTokenByUsername(ctx, req.UserEmail)
	if err != nil {
		return nil, err
	}

	return &inpb.GetOTPCodeResponse{OtpCode: otp}, nil
}
