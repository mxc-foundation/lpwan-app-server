package user

import (
	"context"

	"github.com/pkg/errors"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type Validator struct {
	Credentials *cred.Credentials
}

type Validate interface {
	ValidateActiveUser(ctx context.Context) (bool, error)
	ValidateUsersGlobalAccess(ctx context.Context, flag auth.Flag) (bool, error)
	ValidateUserAccess(ctx context.Context, flag auth.Flag, userID int64) (bool, error)
	IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error
	Is2FAEnabled(ctx context.Context, userEmail string) (bool, error)
	SignJWToken(userEmail string, ttl int64, audience []string) (string, error)
	GetUser(ctx context.Context, opts ...cred.Option) (auth.User, error)
	NewConfiguration(ctx context.Context, userEmail string) (*otp.Configuration, error)
	Enable2FA(ctx context.Context) error
	Disable2FA(ctx context.Context) error
	OTPGetRecoveryCodes(ctx context.Context, regenerate bool) ([]string, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
	}
}

func (v *Validator) OTPGetRecoveryCodes(ctx context.Context, regenerate bool) ([]string, error) {
	u, err := v.Credentials.GetUser(ctx, cred.WithValidOTP())
	if err != nil {
		return nil, err
	}

	return v.Credentials.OTPGetRecoveryCodes(ctx, u.Email, regenerate)
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

	if err := v.Credentials.EnableOTP(ctx, u.Email, OTP); err != nil {
		return err
	}

	return nil
}

func (v *Validator) Disable2FA(ctx context.Context) error {
	u, err := v.Credentials.GetUser(ctx, cred.WithValidOTP())
	if err != nil {
		return err
	}

	if err := v.Credentials.DisableOTP(ctx, u.Email); err != nil {
		return err
	}
	return nil
}

func (v *Validator) NewConfiguration(ctx context.Context, userEmail string) (*otp.Configuration, error) {
	return v.Credentials.NewConfiguration(ctx, userEmail)
}

func (v *Validator) GetUser(ctx context.Context, opts ...cred.Option) (auth.User, error) {
	return v.Credentials.GetUser(ctx, opts...)
}

func (v *Validator) SignJWToken(userEmail string, ttl int64, audience []string) (string, error) {
	return v.Credentials.SignJWToken(userEmail, ttl, audience)
}

func (v *Validator) Is2FAEnabled(ctx context.Context, userEmail string) (bool, error) {
	return v.Credentials.Is2FAEnabled(ctx, userEmail)
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

// ValidateActiveUser validates if the user in the JWT claim is active.
func (v *Validator) ValidateActiveUser(ctx context.Context) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateActiveUser")
	}

	return ctrl.st.CheckActiveUser(ctx, u.Email, u.ID)
}

// ValidateUsersGlobalAccess validates if the client has access to the global users
// resource.
func (v *Validator) ValidateUsersGlobalAccess(ctx context.Context, flag auth.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUsersGlobalAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateUserAcess(ctx, u.Email, u.ID)
	case auth.List:
		return ctrl.st.CheckListUserAcess(ctx, u.Email, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateUserAccess validates if the client has access to the given user
// resource.
func (v *Validator) ValidateUserAccess(ctx context.Context, flag auth.Flag, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUserAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadUserAccess(ctx, u.Email, userID, u.ID)
	case auth.Update:
		return ctrl.st.CheckUpdateUserAccess(ctx, u.Email, userID, u.ID)
	case auth.Delete:
		return ctrl.st.CheckDeleteUserAccess(ctx, u.Email, userID, u.ID)
	case auth.UpdateProfile:
		return ctrl.st.CheckUpdateProfileUserAccess(ctx, u.Email, userID, u.ID)
	case auth.UpdatePassword:
		return ctrl.st.CheckUpdatePasswordUserAccess(ctx, u.Email, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
