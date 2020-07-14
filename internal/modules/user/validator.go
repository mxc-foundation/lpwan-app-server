package user

import (
	"context"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateActiveUser(ctx context.Context) (bool, error)
	ValidateUsersGlobalAccess(ctx context.Context, flag authcus.Flag) (bool, error)
	ValidateUserAccess(ctx context.Context, flag authcus.Flag, userID int64) (bool, error)
	GetIsAdmin(ctx context.Context) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetIsAdmin(ctx context.Context) (bool, error) {
	return v.Credentials.IsGlobalAdmin(ctx)
}

// ValidateActiveUser validates if the user in the JWT claim is active.
func (v *Validator) ValidateActiveUser(ctx context.Context) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateActiveUser")
	}

	return v.Store.CheckActiveUser(u.Username, u.ID)
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
		return v.Store.CheckCreateUserAcess(u.Username, u.ID)
	case authcus.List:
		return v.Store.CheckListUserAcess(u.Username, u.ID)
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
		return v.Store.CheckReadUserAccess(u.Username, userID, u.ID)
	case authcus.Update, authcus.Delete:
		return v.Store.CheckUpdateDeleteUserAccess(u.Username, userID, u.ID)
	case authcus.UpdateProfile:
		return v.Store.CheckUpdateProfileUserAccess(u.Username, userID, u.ID)
	case authcus.UpdatePassword:
		return v.Store.CheckUpdatePasswordUserAccess(u.Username, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
