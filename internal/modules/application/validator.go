package application

import (
	"context"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Store       ApplicationStore
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
	FinishRegistration
)

// ValidateGlobalApplicationsAccess validates if the client has access to the
// global applications resource.
func (v *Validator) ValidateGlobalApplicationsAccess(ctx context.Context, flag Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalApplicationsAccess")
	}

	switch flag {
	case Create:
		return v.Store.CheckCreateApplicationAccess(u.Username, u.ID, organizationID)
	case List:
		return v.Store.CheckListApplicationAccess(u.Username, u.ID, organizationID)
	default:
		panic("ValidateGlobalApplicationsAccess: unsupported flag")
	}
}

// ValidateApplicationAccess validates if the client has access to the given
// application.
func (v *Validator) ValidateApplicationAccess(ctx context.Context, flag Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateApplicationAccess: failed to get username ")
	}

	switch flag {
	case Read:
		return v.Store.CheckReadApplicationAccess(u.Username, u.ID, applicationID)
	case Update:
		return v.Store.CheckUpdateApplicationAccess(u.Username, u.ID, applicationID)
	case Delete:
		return v.Store.CheckDeleteApplicationAccess(u.Username, u.ID, applicationID)
	default:
		panic("ValidateApplicationAccess: unsupported flag")
	}
}
