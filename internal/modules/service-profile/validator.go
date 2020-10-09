package serviceprofile

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
)

type Validator struct {
	Credentials *cred.Credentials
}

type Validate interface {
	GetUser(ctx context.Context) (auth.User, error)
	ValidateServiceProfilesAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateServiceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateServiceProfilesAccess validates if the client has access to the
// service-profiles.
func (v *Validator) ValidateServiceProfilesAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateServiceProfilesAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateServiceProfilesAccess(ctx, u.Email, organizationID, u.ID)
	case auth.List:
		return ctrl.st.CheckListServiceProfilesAccess(ctx, u.Email, organizationID, u.ID)
	default:
		panic("ValidateServiceProfilesAccess: not supported flag")
	}

}

// ValidateServiceProfileAccess validates if the client has access to the
// given service-profile.
func (v *Validator) ValidateServiceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateServiceProfileAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadServiceProfileAccess(ctx, u.Email, id, u.ID)
	case auth.Update, auth.Delete:
		return ctrl.st.CheckUpdateDeleteServiceProfileAccess(ctx, u.Email, id, u.ID)
	default:
		panic("ValidateServiceProfileAccess: not supported flag")
	}

}
