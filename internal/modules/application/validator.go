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
	ValidateGlobalApplicationsAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateApplicationAccess(ctx context.Context, flag authcus.Flag, applicationID int64) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateGlobalApplicationsAccess validates if the client has access to the
// global applications resource.
func (v *Validator) ValidateGlobalApplicationsAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalApplicationsAccess")
	}

	switch flag {
	case authcus.Create:
		return ctrl.st.CheckCreateApplicationAccess(ctx, u.UserEmail, u.ID, organizationID)
	case authcus.List:
		return ctrl.st.CheckListApplicationAccess(ctx, u.UserEmail, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateApplicationAccess validates if the client has access to the given
// application.
func (v *Validator) ValidateApplicationAccess(ctx context.Context, flag authcus.Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateApplicationAccess")
	}

	switch flag {
	case authcus.Read:
		return ctrl.st.CheckReadApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	case authcus.Update:
		return ctrl.st.CheckUpdateApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	case authcus.Delete:
		return ctrl.st.CheckDeleteApplicationAccess(ctx, u.UserEmail, u.ID, applicationID)
	default:
		panic("unsupported flag")
	}
}
