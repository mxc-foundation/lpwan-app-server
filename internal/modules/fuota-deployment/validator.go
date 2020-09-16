package fuotamod

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateFUOTADeploymentAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error)
	ValidateFUOTADeploymentsAccess(ctx context.Context, flag authcus.Flag, applicationID int64, devEUI lorawan.EUI64) (bool, error)
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

// ValidateFUOTADeploymentAccess validates if the client has access to the
// given fuota deployment.
func (v *Validator) ValidateFUOTADeploymentAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateFUOTADeploymentAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadFUOTADeploymentAccess(ctx, u.UserEmail, id, u.ID)
	default:
		panic("ValidateFUOTADeploymentAccess: unsupported flag")
	}
}

// ValidateFUOTADeploymentsAccess validates if the client has access to the
// fuota deployments.
func (v *Validator) ValidateFUOTADeploymentsAccess(ctx context.Context, flag authcus.Flag, applicationID int64, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateFUOTADeploymentsAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateFUOTADeploymentsAccess(ctx, u.UserEmail, applicationID, devEUI, u.ID)
	default:
		panic("ValidateFUOTADeploymentsAccess: unsupported flag")
	}
}
