package organization

import (
	"context"

	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateOrganizationAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateOrganizationsAccess(ctx context.Context, flag authcus.Flag) (bool, error)
	ValidateOrganizationUsersAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateOrganizationUserAccess(ctx context.Context, flag authcus.Flag, organizationID, userID int64) (bool, error)
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

// ValidateOrganizationAccess validates if the client has access to the
// given organization.
func (v *Validator) ValidateOrganizationAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadOrganizationAccess(u.Username, u.ID, organizationID)
	case authcus.Update:
		return Service.St.CheckUpdateOrganizationAccess(u.Username, u.ID, organizationID)
	case authcus.Delete:
		return Service.St.CheckDeleteOrganizationAccess(u.Username, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationsAccess validates if the client has access to the
// organizations.
func (v *Validator) ValidateOrganizationsAccess(ctx context.Context, flag authcus.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationsAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateOrganizationAccess(u.Username, u.ID)
	case authcus.List:
		return Service.St.CheckListOrganizationAccess(u.Username, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUsersAccess validates if the client has access to
// the organization users.
func (v *Validator) ValidateOrganizationUsersAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateOrganizationUserAccess(u.Username, u.ID, organizationID)
	case authcus.List:
		return Service.St.CheckListOrganizationUserAccess(u.Username, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUserAccess validates if the client has access to the
// given user of the given organization.
func (v *Validator) ValidateOrganizationUserAccess(ctx context.Context, flag authcus.Flag, organizationID, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	case authcus.Update:
		return Service.St.CheckUpdateOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	case authcus.Delete:
		return Service.St.CheckDeleteOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
