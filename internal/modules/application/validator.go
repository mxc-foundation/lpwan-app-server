package application

import (
	"context"

	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateGlobalApplicationsAccess(ctx context.Context, flag Flag, organizationID int64) (bool, error)
	ValidateApplicationAccess(ctx context.Context, flag Flag, applicationID int64) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
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
)

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateGlobalApplicationsAccess validates if the client has access to the
// global applications resource.
func (v *Validator) ValidateGlobalApplicationsAccess(ctx context.Context, flag Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalApplicationsAccess")
	}

	switch flag {
	case Create:
		return Service.St.CheckCreateApplicationAccess(ctx, u.UserEmail, u.ID, organizationID)
	case List:
		return Service.St.CheckListApplicationAccess(ctx, u.UserEmail, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateApplicationAccess validates if the client has access to the given
// application.
func (v *Validator) ValidateApplicationAccess(ctx context.Context, flag Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateApplicationAccess")
	}

	switch flag {
	case Read:
		return Service.St.CheckReadApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	case Update:
		return Service.St.CheckUpdateApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	case Delete:
		return Service.St.CheckDeleteApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	default:
		panic("unsupported flag")
	}
}
