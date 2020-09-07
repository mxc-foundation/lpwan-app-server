package devprofile

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
	ValidateDeviceProfilesAccess(ctx context.Context, flag authcus.Flag, organizationID, applicationID int64) (bool, error)
	ValidateDeviceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error)
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

// ValidateDeviceProfilesAccess validates if the client has access to the
// device-profiles.
func (v *Validator) ValidateDeviceProfilesAccess(ctx context.Context, flag authcus.Flag, organizationID, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateDeviceProfilesAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateDeviceProfilesAccess(ctx, u.UserEmail, organizationID, applicationID, u.ID)
	case authcus.List:
		return Service.St.CheckListDeviceProfilesAccess(ctx, u.UserEmail, organizationID, applicationID, u.ID)
	default:
		panic("ValidateDeviceProfilesAccess: unsupported flag")
	}

}

// ValidateDeviceProfileAccess validates if the client has access to the
// given device-profile.
func (v *Validator) ValidateDeviceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateDeviceProfileAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadDeviceProfileAccess(ctx, u.UserEmail, id, u.ID)
	case authcus.Update, authcus.Delete:
		return Service.St.CheckUpdateDeleteDeviceProfileAccess(ctx, u.UserEmail, id, u.ID)
	default:
		panic("ValidateDeviceProfileAccess: unsupported flag")
	}

}
