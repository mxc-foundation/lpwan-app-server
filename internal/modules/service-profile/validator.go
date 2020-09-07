package serviceprofile

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	GetUser(ctx context.Context) (authcus.User, error)
	ValidateServiceProfilesAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateServiceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateServiceProfilesAccess validates if the client has access to the
// service-profiles.
func (v *Validator) ValidateServiceProfilesAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateServiceProfilesAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateServiceProfilesAccess(ctx, u.UserEmail, organizationID, u.ID)
	case authcus.List:
		return Service.St.CheckListServiceProfilesAccess(ctx, u.UserEmail, organizationID, u.ID)
	default:
		panic("ValidateServiceProfilesAccess: not supported flag")
	}

}

// ValidateServiceProfileAccess validates if the client has access to the
// given service-profile.
func (v *Validator) ValidateServiceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateServiceProfileAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadServiceProfileAccess(ctx, u.UserEmail, id, u.ID)
	case authcus.Update, authcus.Delete:
		return Service.St.CheckUpdateDeleteServiceProfileAccess(ctx, u.UserEmail, id, u.ID)
	default:
		panic("ValidateServiceProfileAccess: not supported flag")
	}

}
