package devprofile

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *auth.Credentials
}

type Validate interface {
	ValidateDeviceProfilesAccess(ctx context.Context, flag auth.Flag, organizationID, applicationID int64) (bool, error)
	ValidateDeviceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error)
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

// ValidateDeviceProfilesAccess validates if the client has access to the
// device-profiles.
func (v *Validator) ValidateDeviceProfilesAccess(ctx context.Context, flag auth.Flag, organizationID, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateDeviceProfilesAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateDeviceProfilesAccess(ctx, u.Email, organizationID, applicationID, u.ID)
	case auth.List:
		return ctrl.st.CheckListDeviceProfilesAccess(ctx, u.Email, organizationID, applicationID, u.ID)
	default:
		panic("ValidateDeviceProfilesAccess: unsupported flag")
	}

}

// ValidateDeviceProfileAccess validates if the client has access to the
// given device-profile.
func (v *Validator) ValidateDeviceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateDeviceProfileAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadDeviceProfileAccess(ctx, u.Email, id, u.ID)
	case auth.Update, auth.Delete:
		return ctrl.st.CheckUpdateDeleteDeviceProfileAccess(ctx, u.Email, id, u.ID)
	default:
		panic("ValidateDeviceProfileAccess: unsupported flag")
	}

}
