package user

import (
	"context"

	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Store       UserStore
	Credentials *authcus.Credentials
}

func NewValidator(v Validator) *Validator {
	return &Validator{
		Store:       v.Store,
		Credentials: v.Credentials,
	}
}

// API key subjects.
const (
	SubjectUser   = "user"
	SubjectAPIKey = "api_key"
)

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
func (v *Validator) ValidateUsersGlobalAccess(ctx context.Context, flag Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUsersGlobalAccess")
	}

	switch flag {
	case Create:
		return v.Store.CheckCreateUserAcess(u.Username, u.ID)
	case List:
		return v.Store.CheckListUserAcess(u.Username, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateUserAccess validates if the client has access to the given user
// resource.
func (v *Validator) ValidateUserAccess(ctx context.Context, flag Flag, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateUserAccess")
	}

	switch flag {
	case Read:
		return v.Store.CheckReadUserAccess(u.Username, userID, u.ID)
	case Update, Delete:
		return v.Store.CheckUpdateDeleteUserAccess(u.Username, userID, u.ID)
	case UpdateProfile:
		return v.Store.CheckUpdateProfileUserAccess(u.Username, userID, u.ID)
	case UpdatePassword:
		return v.Store.CheckUpdatePasswordUserAccess(u.Username, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
