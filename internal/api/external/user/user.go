// Package user implements APIs for user's registration and login
package user

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	udata "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Server implements Internal User Service
type Server struct {
	st     *store.Handler
	config udata.Config
	auth   auth.Authenticator
	jwtv   *jwt.Validator
	otpv   *otp.Validator
}

// NewServer creates a new server instance
func NewServer(h *store.Handler, auth auth.Authenticator, jwtv *jwt.Validator, otpv *otp.Validator, c udata.Config) *Server {
	return &Server{
		st:     h,
		config: c,
		auth:   auth,
		jwtv:   jwtv,
		otpv:   otpv,
	}
}

// Create creates the given user.
func (a *Server) Create(ctx context.Context, req *inpb.CreateUserRequest) (*inpb.CreateUserResponse, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	user := udata.User{
		SessionTTL: req.User.SessionTtl,
		IsAdmin:    req.User.IsAdmin,
		IsActive:   req.User.IsActive,
		Email:      req.User.Email,
		Note:       req.User.Note,
		Password:   req.Password,
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		err := handler.CreateUser(ctx, &user)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		for _, org := range req.Organizations {
			if err := handler.CreateOrganizationUser(ctx, org.OrganizationId, user.ID,
				org.IsAdmin, org.IsDeviceAdmin, org.IsGatewayAdmin); err != nil {
				return status.Errorf(codes.Unknown, "%v", err)
			}
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &inpb.CreateUserResponse{Id: user.ID}, nil
}

// Get returns the user matching the given ID.
func (a *Server) Get(ctx context.Context, req *inpb.GetUserRequest) (*inpb.GetUserResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !(cred.IsGlobalAdmin || cred.UserID == req.Id) {
		// only user themselves and the global admin can do that
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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
			Username:   user.Email,
			Email:      user.EmailOld,
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
func (a *Server) GetUserEmail(ctx context.Context, req *inpb.GetUserEmailRequest) (*inpb.GetUserEmailResponse, error) {
	username := normalizeUsername(req.UserEmail)
	u, err := a.st.GetUserByEmail(ctx, username)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
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
func (a *Server) List(ctx context.Context, req *inpb.ListUserRequest) (*inpb.ListUserResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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
			Username:   u.Email,
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
func (a *Server) Update(ctx context.Context, req *inpb.UpdateUserRequest) (*empty.Empty, error) {
	if req.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user must not be nil")
	}
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !(cred.UserID == req.User.Id) {
		// only user themselves can do that
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	user, err := a.st.GetUser(ctx, req.User.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}
	
	user.IsActive = req.User.IsActive
	user.SessionTTL = req.User.SessionTtl
	user.Email = req.User.Username
	user.Note = req.User.Note

	err = a.st.UpdateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes the user matching the given ID.
func (a *Server) Delete(ctx context.Context, req *inpb.DeleteUserRequest) (*empty.Empty, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := a.st.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// UpdatePassword updates the password for the user matching the given ID.
func (a *Server) UpdatePassword(ctx context.Context, req *inpb.UpdateUserPasswordRequest) (*empty.Empty, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !(cred.UserID == req.UserId) {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := a.st.UpdatePassword(ctx, req.UserId, req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

func (a *Server) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	userEmail := normalizeUsername(req.UserEmail)
	otp, err := a.st.GetTokenByUsername(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return &inpb.GetOTPCodeResponse{OtpCode: otp}, nil
}
