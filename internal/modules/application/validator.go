package application

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

// Validator defines struct type for vadidating user access to APIs provided by this package
type Validator struct {
	Credentials *auth.Credentials
	st          *store.Handler
}

// Validate defines methods used on struct Validator
type Validate interface {
	ValidateGlobalApplicationsAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateApplicationAccess(ctx context.Context, flag auth.Flag, applicationID int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
}

// NewValidator returns new Validate instance for this package
func NewValidator(st *store.Handler) Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
		st:          st,
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateGlobalApplicationsAccess validates if the client has access to the
// global applications resource.
func (v *Validator) ValidateGlobalApplicationsAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalApplicationsAccess")
	}

	switch flag {
	case auth.Create:
		return v.st.CheckCreateApplicationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.List:
		return v.st.CheckListApplicationAccess(ctx, u.Email, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateApplicationAccess validates if the client has access to the given
// application.
func (v *Validator) ValidateApplicationAccess(ctx context.Context, flag auth.Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateApplicationAccess")
	}

	switch flag {
	case auth.Read:
		return v.st.CheckReadApplicationAccess(ctx, u.Email, u.ID, applicationID)
	case auth.Update:
		return v.st.CheckUpdateApplicationAccess(ctx, u.Email, u.ID, applicationID)
	case auth.Delete:
		return v.st.CheckDeleteApplicationAccess(ctx, u.Email, u.ID, applicationID)
	default:
		panic("unsupported flag")
	}
}
