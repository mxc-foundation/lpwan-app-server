package fuotamod

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
)

type Validator struct {
	Credentials *cred.Credentials
}

type Validate interface {
	ValidateFUOTADeploymentAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error)
	ValidateFUOTADeploymentsAccess(ctx context.Context, flag auth.Flag, applicationID int64, devEUI lorawan.EUI64) (bool, error)
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

// ValidateFUOTADeploymentAccess validates if the client has access to the
// given fuota deployment.
func (v *Validator) ValidateFUOTADeploymentAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateFUOTADeploymentAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadFUOTADeploymentAccess(ctx, u.Email, id, u.ID)
	default:
		panic("ValidateFUOTADeploymentAccess: unsupported flag")
	}
}

// ValidateFUOTADeploymentsAccess validates if the client has access to the
// fuota deployments.
func (v *Validator) ValidateFUOTADeploymentsAccess(ctx context.Context, flag auth.Flag, applicationID int64, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateFUOTADeploymentsAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateFUOTADeploymentsAccess(ctx, u.Email, applicationID, devEUI, u.ID)
	default:
		panic("ValidateFUOTADeploymentsAccess: unsupported flag")
	}
}
