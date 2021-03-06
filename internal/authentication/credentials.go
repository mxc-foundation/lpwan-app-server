package authentication

import (
	"context"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

// User contains information about the user
type User struct {
	ID            int64
	Email         string
	IsGlobalAdmin bool
}

// OrgUser contains information about the role of the user in organisation
type OrgUser struct {
	IsOrgUser      bool
	IsOrgAdmin     bool
	IsDeviceAdmin  bool
	IsGatewayAdmin bool
}

// Flag defines the authorization flag.
type Flag int

// Authorization flags.
const (
	Create Flag = iota
	Read
	Update
	Delete
	List
	UpdateProfile
	UpdatePassword
	FinishRegistration
)

func SetupCred(st auth.Store, jwtValidator *jwt.Validator, otpValidator *otp.Validator) {
	hl.st = st
	hl.jwtValidator = jwtValidator
	hl.otpValidator = otpValidator
}

var hl handler

type handler struct {
	st           auth.Store
	jwtValidator *jwt.Validator
	otpValidator *otp.Validator
}

type Credentials struct {
	// init when service starts
	h *handler
	// change based on ctx
	user    User
	orgUser OrgUser
}

func NewCredentials() *Credentials {
	return &Credentials{
		h: &hl,
	}
}

type options struct {
	audience    string
	requireOTP  bool
	limitedCred bool
	orgID       int64
}

// Option is used to configure validator checks
type Option func(opts *options)

// WithAudience requires that credentials presented included all the listed audiences
func WithAudience(audience string) Option {
	return func(opts *options) {
		opts.audience = audience
	}
}

// WithValidOTP requires that the request included valid OTP code
func WithValidOTP() Option {
	return func(opts *options) {
		opts.requireOTP = true
	}
}

// WithLimitedCredentials creates limited credentials
//
// Deprecated: do not use, this is only for the purposes of the registration
// process
func WithLimitedCredentials() Option {
	return func(opts *options) {
		opts.limitedCred = true
	}
}

func WithOrganizationID(orgID int64) Option {
	return func(opts *options) {
		opts.orgID = orgID
	}
}

// getCredentials returns a new Credentials object for the user, assuming that
// the user exists and active
func (c *Credentials) getCredentials(ctx context.Context, opts ...Option) (Credentials, error) {
	var cred Credentials

	cfg := options{audience: "lora-app-server"}
	for _, o := range opts {
		o(&cfg)
	}
	token, err := getTokenFromContext(ctx)
	if err != nil {
		return cred, err
	}
	jwtClaims, err := c.h.jwtValidator.GetClaims(token, cfg.audience)
	if err != nil {
		return cred, errors.Wrap(err, "getCredentials")
	}
	if jwtClaims.Username == "" {
		return cred, fmt.Errorf("username is required in claims")
	}

	otpClaims, err := c.h.otpValidator.GetClaims(ctx)
	if err != nil {
		return cred, errors.Wrap(err, "getCredentials")
	}

	if cfg.requireOTP {
		if otpClaims.OTP == "" {
			return cred, errors.Wrap(err, "getCredentials: OTP is required")
		}
		if enabled, err := c.h.otpValidator.IsEnabled(ctx, jwtClaims.Username); !enabled || err != nil {
			return cred, errors.Wrap(err, "getCredentials: two-factor authentication is not enabled")
		}
		if err := c.h.otpValidator.Validate(ctx, jwtClaims.Username, otpClaims.OTP); err != nil {
			return cred, errors.Wrap(err, "getCredentials: OTP is not valid")
		}
	}

	if cfg.limitedCred {
		// when cfg.limitedCred is true, func returns Credentials that only contain username. All
		// other checks will fail and possibly panic.
		//
		// Deprecated: this is only should be used for the user registration process,
		// and user registration process should be fixed to not require this hack
		cred.h = c.h
		cred.user.ID = -1
		cred.user.Email = jwtClaims.Username
		cred.user.IsGlobalAdmin = false
		return cred, nil
	}

	u, err := c.h.st.AuthGetUser(ctx, jwtClaims.Username)
	if err != nil {
		return cred, errors.Wrap(err, "getCredentials")
	}

	cred.h = c.h
	cred.user.ID = u.ID
	cred.user.Email = jwtClaims.Username
	cred.user.IsGlobalAdmin = u.IsGlobalAdmin

	if cred.user.IsGlobalAdmin {
		return cred, nil
	}

	if cfg.orgID != 0 {
		orgUser, err := c.h.st.AuthGetOrgUser(ctx, u.ID, cfg.orgID)
		if err != nil {
			return cred, errors.Wrap(err, "getCredentials")
		}
		if !orgUser.IsOrgUser {
			return cred, fmt.Errorf("not an org user")
		}

		cred.orgUser.IsGatewayAdmin = orgUser.IsGatewayAdmin
		cred.orgUser.IsDeviceAdmin = orgUser.IsDeviceAdmin
		cred.orgUser.IsOrgAdmin = orgUser.IsOrgAdmin
		cred.orgUser.IsOrgUser = orgUser.IsOrgUser
	}

	return cred, nil
}

var validAuthorizationRegexp = regexp.MustCompile(`(?i)^bearer (.*)$`)

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	token, ok := md["authorization"]
	if !ok || len(token) == 0 {
		return "", fmt.Errorf("missing authorization header")
	}

	match := validAuthorizationRegexp.FindStringSubmatch(token[0])

	// authorization header should respect RFC1945
	if len(match) == 0 {
		l := len(token)
		if l > 16 {
			l = 16
		}
		logrus.Warnf("Deprecated Authorization header: %s", token[0:l])
		return token[0], nil
	}

	return match[1], nil
}

func (c *Credentials) GetUser(ctx context.Context, opts ...Option) (User, error) {
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return User{}, errors.Wrap(err, "failed to get credentials")
	}

	return cred.user, nil
}

func (c *Credentials) GetUserPermissionWithOrgID(ctx context.Context, orgID int64, opts ...Option) (OrgUser, error) {
	opts = append(opts, WithOrganizationID(orgID))
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return OrgUser{}, errors.Wrap(err, "failed to get credentials")
	}

	return cred.orgUser, nil
}

// Username returns the name of the user
func (c *Credentials) Username(ctx context.Context, opts ...Option) (string, error) {
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return "", errors.Wrap(err, "failed to get credentials")
	}

	return cred.user.Email, nil
}

// UserID returns user id of the user
func (c *Credentials) UserID(ctx context.Context, opts ...Option) (int64, error) {
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get credentials")
	}
	return cred.user.ID, nil
}

// IsGlobalAdmin checks that the user is a global admin and returns an error if
// he's not
func (c *Credentials) IsGlobalAdmin(ctx context.Context, opts ...Option) error {
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}

	if !cred.user.IsGlobalAdmin {
		return errors.New("user is not global admin")
	}

	return nil
}

// IsOrgUser checks that the user belongs to the organisation, if not it
// returns an error
func (c *Credentials) IsOrgUser(ctx context.Context, orgID int64, opts ...Option) error {
	opts = append(opts, WithOrganizationID(orgID))
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}

	if !cred.user.IsGlobalAdmin && !cred.orgUser.IsOrgUser {
		return errors.New("user is neither org user nor global admin")
	}

	return nil
}

// IsOrgAdmin checks that the user is admin for the organisation, if not it
// returns an error
func (c *Credentials) IsOrgAdmin(ctx context.Context, orgID int64, opts ...Option) error {
	opts = append(opts, WithOrganizationID(orgID))
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}

	if !cred.user.IsGlobalAdmin && !cred.orgUser.IsOrgAdmin {
		return errors.New("user is neither org admin nor global admin")
	}

	return nil
}

// IsDeviceAdmin checks that the user is device admin for the organisation, if
// not it returns an error
func (c *Credentials) IsDeviceAdmin(ctx context.Context, orgID int64, opts ...Option) error {
	opts = append(opts, WithOrganizationID(orgID))
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}

	if !cred.user.IsGlobalAdmin && !cred.orgUser.IsDeviceAdmin {
		return errors.New("user is neither device admin nor global admin")
	}

	return nil
}

// IsGatewayAdmin checks that the user is gateway admin for the organisation,
// if not it returns an error
func (c *Credentials) IsGatewayAdmin(ctx context.Context, orgID int64, opts ...Option) error {
	opts = append(opts, WithOrganizationID(orgID))
	cred, err := c.getCredentials(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials")
	}
	if !cred.user.IsGlobalAdmin && !cred.orgUser.IsGatewayAdmin {
		return errors.New("user is neither gateway admin nor global admin")
	}

	return nil
}

// Is2FAEnabled requires username, since ctx does not contain user info at this point
func (c *Credentials) Is2FAEnabled(ctx context.Context, username string) (bool, error) {
	return c.h.otpValidator.IsEnabled(ctx, username)
}

// SignJWToken requires username, since ctx does not contain user info at this point
func (c *Credentials) SignJWToken(userID int64, username string, ttl int64, audience []string) (string, error) {
	return c.h.jwtValidator.SignToken(jwt.Claims{UserID: userID, Username: username}, ttl, audience)
}

// New2FAConfiguration generates a new TOTP configuration for the user
func (c *Credentials) New2FAConfiguration(ctx context.Context, username string) (*otp.Configuration, error) {
	return c.h.otpValidator.NewConfiguration(ctx, username)
}

func (c *Credentials) GetOTP(ctx context.Context) (string, error) {
	otpClaims, err := c.h.otpValidator.GetClaims(ctx)
	if err != nil {
		return "", errors.Wrap(err, "fail to get otp from ctx")
	}

	return otpClaims.OTP, nil
}

func (c *Credentials) EnableOTP(ctx context.Context, username, otp string) error {
	return c.h.otpValidator.Enable(ctx, username, otp)
}

func (c *Credentials) DisableOTP(ctx context.Context, username string) error {
	return c.h.otpValidator.Disable(ctx, username)
}

func (c *Credentials) OTPGetRecoveryCodes(ctx context.Context, username string, regenerate bool) ([]string, error) {
	return c.h.otpValidator.GetRecoveryCodes(ctx, username, regenerate)
}
