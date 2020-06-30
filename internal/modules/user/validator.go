package user

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *validator {
	return &validator{otpValidator: otpValidator}
}

// GetIsAdmin returns if the authenticated user is a global amin.
func (v *validator) GetIsAdmin(ctx context.Context) (bool, error) {
	claims, err := v.otpValidator.JwtValidator.GetClaims(ctx)
	if err != nil {
		return false, err
	}

	user, err := userAPI.Store.GetUserByUsername(ctx, claims.Username)
	if err != nil {
		return false, errors.Wrap(err, "get user by username error")
	}

	return user.IsAdmin, nil
}

// GetUser returns the user object.
func (v *validator) GetUser(ctx context.Context) (User, error) {
	claims, err := v.otpValidator.JwtValidator.GetClaims(ctx)
	if err != nil {
		return User{}, err
	}

	if claims.Subject != "user" {
		return User{}, errors.New("subject must be user")
	}

	if claims.UserID != 0 {
		return userAPI.Store.GetUser(ctx, claims.UserID)
	}

	if claims.Username != "" {
		return userAPI.Store.GetUserByEmail(ctx, claims.Username)
	}

	return User{}, errors.New("no username or user_id in claims")
}
