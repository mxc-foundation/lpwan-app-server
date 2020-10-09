package application

import (
	"context"

	"github.com/pkg/errors"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
)

type Validator struct {
	Credentials *cred.Credentials
}

type Validate interface {
	ValidateGlobalApplicationsAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateApplicationAccess(ctx context.Context, flag auth.Flag, applicationID int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
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
		return ctrl.st.CheckCreateApplicationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.List:
		return ctrl.st.CheckListApplicationAccess(ctx, u.Email, u.ID, organizationID)
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
		return ctrl.st.CheckReadApplicationAccess(ctx, u.Email, u.ID, applicationID)
	case auth.Update:
		return ctrl.st.CheckUpdateApplicationAccess(ctx, u.Email, u.ID, applicationID)
	case auth.Delete:
		return ctrl.st.CheckDeleteApplicationAccess(ctx, u.Email, u.ID, applicationID)
	default:
		panic("unsupported flag")
	}
}
