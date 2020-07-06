package authentication

import (
	"context"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

// Store provides access to information about users and their roles
type Store interface {
	// GetUser returns user's information given that there is an active user
	// with the given username
	GetUser(ctx context.Context, username string) (User, error)
	// GetOrgUser returns user's role in the listed organization
	GetOrgUser(ctx context.Context, userID int64, orgID int64) (OrgUser, error)
}

/*// Credentials provides methods to assert the user's Credentials
type Credentials interface {
	// Username returns user's username
	Username() string
	// UserID returns id of the user
	UserID() int64
	// IsGlobalAdmin returns an error if user is not global admin
	IsGlobalAdmin(context.Context) error
	// IsOrgUser returns an error if user does not belong to the organization
	IsOrgUser(context.Context, int64) error
	// IsOrgAdmin returns an error if user is not admin of the organization
	IsOrgAdmin(context.Context, int64) error
	// IsDeviceAdmin returns an error if user is not device admin of the organization
	IsDeviceAdmin(context.Context, int64) error
	// IsGatewayAdmin returns an error if user is not gateway admin of the organization
	IsGatewayAdmin(context.Context, int64) error
}*/

func SetupCred(st Store, jwtValidator *jwt.JWTValidator, otpValidator *otp.Validator) {
	hl.st = st
	hl.jwtValidator = jwtValidator
	hl.otpValidator = otpValidator
}

var hl handler

type handler struct {
	st           Store
	jwtValidator *jwt.JWTValidator
	otpValidator *otp.Validator
}

type UserCredentials struct {
	id            int64
	username      string
	isGlobalAdmin bool
}

type OrgUserCredentials struct {
	id            int64
	username      string
	isGlobalAdmin bool
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
	jwtClaims, err := c.h.jwtValidator.GetClaims(ctx, cfg.audience)
	if err != nil {
		return cred, errors.Wrap(err, "getCredentials")
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

	u, err := c.h.st.GetUser(ctx, jwtClaims.Username)
	if err != nil {
		return cred, errors.Wrap(err, "getCredentials")
	}

	if cfg.orgID != 0 {
		orgUser, err := c.h.st.GetOrgUser(ctx, u.ID, cfg.orgID)
		if err != nil {
			return cred, errors.Wrap(err, "getCredentials")
		}

		cred.orgUser.IsGatewayAdmin = orgUser.IsGatewayAdmin
		cred.orgUser.IsDeviceAdmin = orgUser.IsDeviceAdmin
		cred.orgUser.IsOrgAdmin = orgUser.IsOrgAdmin
		cred.orgUser.IsOrgUser = true
	}

	if cfg.limitedCred {
		// when cfg.limitedCred is true, func returns Credentials that only contain username. All
		// other checks will fail and possibly panic.
		//
		// Deprecated: this is only should be used for the user registration process,
		// and user registration process should be fixed to not require this hack
		cred.h = c.h
		cred.user.ID = -1
		cred.user.Username = jwtClaims.Username
		cred.user.IsGlobalAdmin = false

	} else {
		cred.h = c.h
		cred.user.ID = u.ID
		cred.user.Username = jwtClaims.Username
		cred.user.IsGlobalAdmin = u.IsGlobalAdmin
	}

	return cred, nil
}

func (c *Credentials) GetUser(ctx context.Context) (User, error) {
	cred, err := c.getCredentials(ctx)
	if err != nil {
		return User{}, errors.Wrap(err, "failed to get credentials")
	}

	return cred.user, nil
}

func (c *Credentials) GetUserPermissionWithOrgID(ctx context.Context, orgID int64) (OrgUser, error) {
	cred, err := c.getCredentials(ctx, WithOrganizationID(orgID))
	if err != nil {
		return OrgUser{}, errors.Wrap(err, "failed to get credentials")
	}

	return cred.orgUser, nil
}

// Username returns the name of the user
func (c *Credentials) Username(ctx context.Context) (string, error) {
	cred, err := c.getCredentials(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to get credentials")
	}

	return cred.user.Username, nil
}

// UserID returns user id of the user
func (c *Credentials) UserID(ctx context.Context) (int64, error) {
	cred, err := c.getCredentials(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get credentials")
	}
	return cred.user.ID, nil
}

// IsGlobalAdmin checks that the user is a global admin and returns an error if
// he's not
func (c *Credentials) IsGlobalAdmin(ctx context.Context) (bool, error) {
	cred, err := c.getCredentials(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get credentials")
	}

	return cred.user.IsGlobalAdmin, nil
}

// IsOrgUser checks that the user belongs to the organisation, if not it
// returns an error
func (c *Credentials) IsOrgUser(ctx context.Context, orgID int64) (bool, error) {
	cred, err := c.getCredentials(ctx, WithOrganizationID(orgID))
	if err != nil {
		return false, errors.Wrap(err, "failed to get credentials")
	}

	if cred.user.IsGlobalAdmin || cred.orgUser.IsOrgUser {
		return true, nil
	}

	return false, nil
}

// IsOrgAdmin checks that the user is admin for the organisation, if not it
// returns an error
func (c *Credentials) IsOrgAdmin(ctx context.Context, orgID int64) (bool, error) {
	cred, err := c.getCredentials(ctx, WithOrganizationID(orgID))
	if err != nil {
		return false, errors.Wrap(err, "failed to get credentials")
	}

	if cred.user.IsGlobalAdmin || cred.orgUser.IsOrgAdmin {
		return true, nil
	}

	return false, nil
}

// IsDeviceAdmin checks that the user is device admin for the organisation, if
// not it returns an error
func (c *Credentials) IsDeviceAdmin(ctx context.Context, orgID int64) (bool, error) {
	cred, err := c.getCredentials(ctx, WithOrganizationID(orgID))
	if err != nil {
		return false, errors.Wrap(err, "failed to get credentials")
	}

	if cred.user.IsGlobalAdmin || cred.orgUser.IsDeviceAdmin {
		return true, nil
	}

	return false, nil
}

// IsGatewayAdmin checks that the user is gateway admin for the organisation,
// if not it returns an error
func (c *Credentials) IsGatewayAdmin(ctx context.Context, orgID int64) (bool, error) {
	cred, err := c.getCredentials(ctx, WithOrganizationID(orgID))
	if err != nil {
		return false, errors.Wrap(err, "failed to get credentials")
	}
	if cred.user.IsGlobalAdmin || cred.orgUser.IsGatewayAdmin {
		return true, nil
	}

	return false, nil
}
