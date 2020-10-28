package organization

import (
	"context"

	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *auth.Credentials
}

type Validate interface {
	ValidateOrganizationAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateOrganizationsAccess(ctx context.Context, flag auth.Flag) (bool, error)
	ValidateOrganizationUsersAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateOrganizationUserAccess(ctx context.Context, flag auth.Flag, organizationID, userID int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateOrganizationAccess validates if the client has access to the
// given organization.
func (v *Validator) ValidateOrganizationAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.Update:
		return ctrl.st.CheckUpdateOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.Delete:
		return ctrl.st.CheckDeleteOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationsAccess validates if the client has access to the
// organizations.
func (v *Validator) ValidateOrganizationsAccess(ctx context.Context, flag auth.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationsAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateOrganizationAccess(ctx, u.Email, u.ID)
	case auth.List:
		return ctrl.st.CheckListOrganizationAccess(ctx, u.Email, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUsersAccess validates if the client has access to
// the organization users.
func (v *Validator) ValidateOrganizationUsersAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateOrganizationUserAccess(ctx, u.Email, u.ID, organizationID)
	case auth.List:
		return ctrl.st.CheckListOrganizationUserAccess(ctx, u.Email, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUserAccess validates if the client has access to the
// given user of the given organization.
func (v *Validator) ValidateOrganizationUserAccess(ctx context.Context, flag auth.Flag, organizationID, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadOrganizationUserAccess(ctx, u.Email, organizationID, userID, u.ID)
	case auth.Update:
		return ctrl.st.CheckUpdateOrganizationUserAccess(ctx, u.Email, organizationID, userID, u.ID)
	case auth.Delete:
		return ctrl.st.CheckDeleteOrganizationUserAccess(ctx, u.Email, organizationID, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
