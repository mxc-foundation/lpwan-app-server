package user

import (
	"context"

	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateActiveUser(ctx context.Context) (bool, error)
	ValidateUsersGlobalAccess(ctx context.Context, flag authcus.Flag) (bool, error)
	ValidateUserAccess(ctx context.Context, flag authcus.Flag, userID int64) (bool, error)
	IsGlobalAdmin(ctx context.Context, opts ...authcus.Option) error
	Is2FAEnabled(ctx context.Context, username string) (bool, error)
	SignJWToken(username string, ttl int64, audience []string) (string, error)
	GetUser(ctx context.Context, opts ...authcus.Option) (authcus.User, error)
	NewConfiguration(ctx context.Context, username string) (*otp.Configuration, error)
	Enable2FA(ctx context.Context) error
	Disable2FA(ctx context.Context) error
	OTPGetRecoveryCodes(ctx context.Context, username string, regenerate bool) ([]string, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) OTPGetRecoveryCodes(ctx context.Context, username string, regenerate bool) ([]string, error) {
	return v.Credentials.OTPGetRecoveryCodes(ctx, username, regenerate)
}

func (v *Validator) Enable2FA(ctx context.Context) error {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return err
	}

	OTP, err := v.Credentials.GetOTP(ctx)
	if err != nil {
		return err
	}

	if err := v.Credentials.EnableOTP(ctx, u.Username, OTP); err != nil {
		return err
	}

	return nil
}

func (v *Validator) Disable2FA(ctx context.Context) error {
	u, err := v.Credentials.GetUser(ctx, authcus.WithValidOTP())
	if err != nil {
		return err
	}

	if err := v.Credentials.DisableOTP(ctx, u.Username); err != nil {
		return err
	}
	return nil
}

func (v *Validator) NewConfiguration(ctx context.Context, username string) (*otp.Configuration, error) {
	return v.Credentials.NewConfiguration(ctx, username)
}

func (v *Validator) GetUser(ctx context.Context, opts ...authcus.Option) (authcus.User, error) {
	return v.Credentials.GetUser(ctx, opts...)
}

func (v *Validator) SignJWToken(username string, ttl int64, audience []string) (string, error) {
	return v.Credentials.SignJWToken(username, ttl, audience)
}

func (v *Validator) Is2FAEnabled(ctx context.Context, username string) (bool, error) {
	return v.Credentials.Is2FAEnabled(ctx, username)
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...authcus.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

// ValidateActiveUser validates if the user in the JWT claim is active.
func (v *Validator) ValidateActiveUser(ctx context.Context) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateActiveUser")
	}

	return Service.St.CheckActiveUser(ctx, u.Username, u.ID)
}

// ValidateUsersAccess validates if the client has access to the global users
// resource.
func (v *Validator) ValidateUsersGlobalAccess(ctx context.Context, flag authcus.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUsersGlobalAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateUserAcess(ctx, u.Username, u.ID)
	case authcus.List:
		return Service.St.CheckListUserAcess(ctx, u.Username, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateUserAccess validates if the client has access to the given user
// resource.
func (v *Validator) ValidateUserAccess(ctx context.Context, flag authcus.Flag, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUserAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadUserAccess(ctx, u.Username, userID, u.ID)
	case authcus.Update, authcus.Delete:
		return Service.St.CheckUpdateDeleteUserAccess(ctx, u.Username, userID, u.ID)
	case authcus.UpdateProfile:
		return Service.St.CheckUpdateProfileUserAccess(ctx, u.Username, userID, u.ID)
	case authcus.UpdatePassword:
		return Service.St.CheckUpdatePasswordUserAccess(ctx, u.Username, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
