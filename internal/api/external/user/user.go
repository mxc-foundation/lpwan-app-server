// Package user implements APIs for user's registration and login
package user

import (
	"context"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/organization"

	spmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"

	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inpb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

// User defines the user structure.
type User struct {
	ID               int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Email            string
	DisplayName      string
	PasswordHash     string
	IsAdmin          bool
	IsActive         bool
	EmailVerified    bool
	SecurityToken    string
	LastLoginService string
}

type OrganizationUser struct {
	UserID                  int64
	OrganizationID          int64
	OrganizationName        string
	OrganizationDisplayName string
	IsOrgAdmin              bool
	IsDeviceAdmin           bool
	IsGatewayAdmin          bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// Organization represents an organization.
type Organization struct {
	ID              int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	DisplayName     string
	CanHaveGateways bool
	MaxDeviceCount  int
	MaxGatewayCount int
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
	ActivateUser(ctx context.Context, userID int64, passwordHash, orgName, orgDisplayName string) (OrganizationUser, error)
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
	GetUsers(ctx context.Context, limit, offset int64) ([]User, error)
	// GetOrSetPasswordResetOTP if the password reset OTP has been generated
	// already then returns it, otherwise sets the new OTP and returns it
	GetOrSetPasswordResetOTP(ctx context.Context, userID int64, otp string) (string, error)
	// SetUserActiveStatus disables or enables the user
	SetUserActiveStatus(ctx context.Context, userID int64, isActive bool) error
	// SetUserEmail changes the email address of the user
	SetUserEmail(ctx context.Context, userID int64, email string) error
	// SetUserPasswordHash sets the password hash for the user
	SetUserPasswordHash(ctx context.Context, userID int64, passwordHash string) error
	// SetUserPasswordIfOTPMatch sets the user's password if the OTP provided is correct
	SetUserPasswordIfOTPMatch(ctx context.Context, userID int64, otp, passwordHash string) error
	// SetUserDisplayName updates display name of the user
	SetUserDisplayName(ctx context.Context, displayName string, userID int64) error
	// DeleteUser deletes the user
	DeleteUser(ctx context.Context, userID int64) error
	// SetUserLastLogin updates display_name and last_login_service
	SetUserLastLogin(ctx context.Context, userID int64, displayName, service string) error

	// GetUserIDByExternalUserID gets user id from service name and external user id
	GetUserIDByExternalUserID(ctx context.Context, service string, externalUserID string) (int64, error)
	// GetExternalUserByUserIDAndService gets external user id from service name and user id
	GetExternalUserByUserIDAndService(ctx context.Context, service string, userID int64) (ExternalUser, error)
	// GetExternalUsersByUserID gets all external users bound with userID
	GetExternalUsersByUserID(ctx context.Context, userID int64) ([]ExternalUser, error)
	// AddExternalUserLogin inserts new external id and user id relation
	AddExternalUserLogin(ctx context.Context, extUser ExternalUser) error
	// DeleteExternalUserLogin removes binding relation between external account and supernode account
	DeleteExternalUserLogin(ctx context.Context, userID int64, service string) error
	// SetExternalUsername updates external user's username
	SetExternalUsername(ctx context.Context, service, externalUserID, externalUsername string) error
	// GetExternalUserByToken returns external user with given security token
	GetExternalUserByToken(ctx context.Context, service, token string) (ExternalUser, error)
	// GetExternalUserByUsername returns external user with given external username
	GetExternalUserByUsername(ctx context.Context, service, externalUsername string) (ExternalUser, error)
	// ConfirmExternalUserID updates external id of an external user and set verification to empty string
	ConfirmExternalUserID(ctx context.Context, extUser ExternalUser) error

	// GlobalSearch performs a search on organizations, applications, gateways
	// and devices
	GlobalSearch(ctx context.Context, userID int64, globalAdmin bool, search string, limit, offset int) ([]SearchResult, error)
	GetDefaultNetworkServer(ctx context.Context) (nsd.NetworkServer, error)
	// ShopifyStore defines db apis for shopify service
	ShopifyStore
}

// Mailer is an interface responsible for sending emails to the user
type Mailer interface {
	// SendRegistrationConfirmation sends email to the user confirming registration
	SendRegistrationConfirmation(email, lang, securityToken string) error
	// SendVerifyEmailConfirmation sends email with confirmation message to given address
	SendVerifyEmailConfirmation(email, lang, securityToken string) error
	// SendPasswordResetUnknown sends an email that password reset was requested,
	// but the user is unknown
	SendPasswordResetUnknown(email, lang string) error
	// SendPasswordReset sends password reset email
	SendPasswordReset(email, lang, otp string) error
}

// ExternalAuthentication defines configuration for external_auth section
type ExternalAuthentication struct {
	WechatAuth      auth.WeChatAuthentication `mapstructure:"wechat_auth"`
	DebugWechatAuth auth.WeChatAuthentication `mapstructure:"debug_wechat_auth"`
}

// Config defines configuration
type Config struct {
	Recaptcha RecaptchaConfig
	// If true, then users who have 2FA configured must enter OTP to login
	Enable2FALogin bool
	// path to logo
	OperatorLogoPath string
	// external user wechat login config
	WeChatLogin auth.WeChatAuthentication
	// external user wechat login config, debug mode
	DebugWeChatLogin auth.WeChatAuthentication
	// shopify private app configuration
	ShopifyConfig Shopify
}

// Server implements Internal User Service
type Server struct {
	store    Store
	orgStore organization.Store
	spStore  spmod.Store
	dpStore  dp.Store
	mailer   Mailer
	config   Config
	auth     auth.Authenticator
	jwtv     *jwt.Validator
	otpv     *otp.Validator
	pwhasher *pwhash.PasswordHasher
	nsCli    *nscli.Client
}

// NewServer creates a new server instance
func NewServer(store Store, orgStore organization.Store, spStore spmod.Store, dpStore dp.Store, mailer Mailer, auth auth.Authenticator,
	jwtv *jwt.Validator, otpv *otp.Validator, pwhasher *pwhash.PasswordHasher, config Config, nsCli *nscli.Client) *Server {
	return &Server{
		store:    store,
		spStore:  spStore,
		orgStore: orgStore,
		dpStore:  dpStore,
		mailer:   mailer,
		config:   config,
		auth:     auth,
		jwtv:     jwtv,
		otpv:     otpv,
		pwhasher: pwhasher,
		nsCli:    nsCli,
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

	if err := validateEmail(req.User.Email); err != nil {
		return nil, helpers.ErrToRPCError(err)
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
			Username: user.DisplayName,
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

	users, err := a.store.GetUsers(ctx, req.Limit, req.Offset)
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
			Username:  user.DisplayName,
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
func (a *Server) Update(ctx context.Context, req *inpb.UpdateUserRequest) (*inpb.UpdateUserResponse, error) {
	var jwToken string

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

	if req.User.Username != "" && req.User.Username != user.DisplayName {
		if err := a.store.SetUserDisplayName(ctx, req.User.Username, user.ID); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update user's display name : %v", err)
		}
	}

	var userEmail string
	newEmail := normalizeUsername(req.User.Email)
	if newEmail != "" && req.User.Email != user.Email {
		if err := validateEmail(newEmail); err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		if err := a.store.SetUserEmail(ctx, user.ID, newEmail); err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't update user's email: %v", err)
		}
		userEmail = newEmail
	} else {
		userEmail = user.Email
	}

	// username or email changed or not, always return jwt in response
	jwToken, err = a.jwtv.SignToken(jwt.Claims{Username: userEmail, UserID: user.ID, Service: cred.Service}, 0, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create a token: %v", err)
	}

	return &inpb.UpdateUserResponse{Jwt: jwToken}, nil
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
