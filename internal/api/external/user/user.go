// Package user implements APIs for user's registration and login
package user

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

// User defines the user structure.
type User struct {
	ID            int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Email         string
	PasswordHash  string
	IsAdmin       bool
	IsActive      bool
	EmailVerified bool
	SecurityToken string
}

type OrganizationUser struct {
	UserID           int64
	OrganizationID   int64
	OrganizationName string
	IsOrgAdmin       bool
	IsDeviceAdmin    bool
	IsGatewayAdmin   bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SearchResult defines a search result.
type SearchResult struct {
	Kind             string         `db:"kind"`
	Score            float64        `db:"score"`
	OrganizationID   *int64         `db:"organization_id"`
	OrganizationName *string        `db:"organization_name"`
	ApplicationID    *int64         `db:"application_id"`
	ApplicationName  *string        `db:"application_name"`
	DeviceDevEUI     *lorawan.EUI64 `db:"device_dev_eui"`
	DeviceName       *string        `db:"device_name"`
	GatewayMAC       *lorawan.EUI64 `db:"gateway_mac"`
	GatewayName      *string        `db:"gateway_name"`
}

type Store interface {
	// ActivateUser creates the organization for the new user, adds the user to
	// the org and activates the user
	ActivateUser(ctx context.Context, userID int64, passwordHash, orgName, orgDisplayName string) error
	// CreateUser creates a new user and adds it to all organization listed
	CreateUser(ctx context.Context, user User, orgUser []OrganizationUser) (User, error)
	// GetUserByID returns the user with the given ID
	GetUserByID(ctx context.Context, userID int64) (User, error)
	// GetUserByEmail returns the user with the given email
	GetUserByEmail(ctx context.Context, email string) (User, error)
	// GetUserByToken returns the user with the given security token
	GetUserByToken(ctx context.Context, token string) (User, error)
	// GetUserOrganizations returns the list of organizations to which the user
	// belongs and the roles of the user in these organizations
	GetUserOrganizations(ctx context.Context, userID int64) ([]OrganizationUser, error)
	// GetUserCount returns the total number of users
	GetUserCount(ctx context.Context) (int64, error)
	// GetUsers returns list of users
	GetUsers(ctx context.Context, offset, limit int) ([]User, error)
	// If password reset OTP has been generated already then returns it,
	// otherwise sets the new OTP and returns it
	GetOrSetPasswordResetOTP(ctx context.Context, userID int64, otp string) (string, error)
	// SetUserActiveStatus disables or enables the user
	SetUserActiveStatus(ctx context.Context, userID int64, isActive bool) error
	// SetUserEmail changes the email address of the user
	SetUserEmail(ctx context.Context, userID int64, email string) error
	// SetUserPasswordHash sets the password hash for the user
	SetUserPasswordHash(ctx context.Context, userID int64, passwordHash string) error
	// SetUserPasswordIfOTPMatch sets the user's password if the OTP provided is correct
	SetUserPasswordIfOTPMatch(ctx context.Context, userID int64, otp, passwordHash string) error
	// DeleteUser deletes the user
	DeleteUser(ctx context.Context, userID int64) error

	// GlobalSearch performs a search on organizations, applications, gateways
	// and devices
	GlobalSearch(ctx context.Context, userID int64, globalAdmin bool, search string, limit, offset int) ([]SearchResult, error)
}

// Server configuration
type Config struct {
	RecaptchaHost   string
	RecaptchaSecret string
	// If true, then users who have 2FA configured must enter OTP to login
	Enable2FALogin bool
}

// Server implements Internal User Service
type Server struct {
	store    Store
	config   Config
	auth     auth.Authenticator
	jwtv     *jwt.Validator
	otpv     *otp.Validator
	pwhasher *pwhash.PasswordHasher
}

// NewServer creates a new server instance
func NewServer(store Store, auth auth.Authenticator, jwtv *jwt.Validator, otpv *otp.Validator, config Config) *Server {
	return &Server{
		store:  store,
		config: config,
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

	if err := validatePass(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	passHash, err := a.pwhasher.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't hash password: %v", err)
	}

	user := User{
		IsAdmin:      req.User.IsAdmin,
		IsActive:     req.User.IsActive,
		Email:        req.User.Email,
		PasswordHash: passHash,
	}

	var orgUsers []OrganizationUser
	for _, org := range req.Organizations {
		orgUsers = append(orgUsers, OrganizationUser{
			OrganizationID: org.OrganizationId,
			IsOrgAdmin:     org.IsAdmin,
			IsDeviceAdmin:  org.IsDeviceAdmin,
			IsGatewayAdmin: org.IsGatewayAdmin,
		})
	}

	user, err = a.store.CreateUser(ctx, user, orgUsers)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create user: %v", err)
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

	user, err := a.store.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.GetUserResponse{
		User: &inpb.User{
			Id:       user.ID,
			IsAdmin:  user.IsAdmin,
			IsActive: user.IsActive,
			Username: user.Email,
		},
		CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
		UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
	}

	return &resp, nil
}

// GetUserEmail returns true if user does not exist
func (a *Server) GetUserEmail(ctx context.Context, req *inpb.GetUserEmailRequest) (*inpb.GetUserEmailResponse, error) {
	email := normalizeUsername(req.UserEmail)
	u, err := a.store.GetUserByEmail(ctx, email)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			return &inpb.GetUserEmailResponse{Status: true}, nil
		}
		return nil, helpers.ErrToRPCError(err)
	}
	if u.SecurityToken != "" {
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

	users, err := a.store.GetUsers(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	totalUserCount, err := a.store.GetUserCount(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := inpb.ListUserResponse{
		TotalCount: int64(totalUserCount),
	}

	for _, user := range users {
		row := inpb.UserListItem{
			Username:  user.Email,
			Id:        user.ID,
			IsAdmin:   user.IsAdmin,
			IsActive:  user.IsActive,
			CreatedAt: &timestamp.Timestamp{Seconds: user.CreatedAt.Unix()},
			UpdatedAt: &timestamp.Timestamp{Seconds: user.UpdatedAt.Unix()},
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

	user, err := a.store.GetUserByID(ctx, req.User.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	if req.User.IsActive != user.IsActive {
		if err := a.store.SetUserActiveStatus(ctx, user.ID, req.User.IsActive); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't set user status: %v", err)
		}
	}
	if req.User.Email != "" && req.User.Email != user.Email {
		if err := a.store.SetUserEmail(ctx, user.ID, req.User.Email); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update user's email: %v", err)
		}
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

	if err := a.store.DeleteUser(ctx, req.Id); err != nil {
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

	if err := validatePass(req.Password); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	h, err := a.pwhasher.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "problem updating the password: %v", err)
	}
	if err := a.store.SetUserPasswordHash(ctx, req.UserId, h); err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

func (a *Server) GetOTPCode(ctx context.Context, req *inpb.GetOTPCodeRequest) (*inpb.GetOTPCodeResponse, error) {
	userEmail := normalizeUsername(req.UserEmail)
	u, err := a.store.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if u.SecurityToken == "" {
		return nil, status.Errorf(codes.NotFound, "no token for the user")
	}

	return &inpb.GetOTPCodeResponse{OtpCode: u.SecurityToken}, nil
}
