package gatewayprofile

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
	ValidateGatewayProfileAccess(ctx context.Context, flag auth.Flag) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
	}
}

// ValidateGatewayProfileAccess validates if the client has access
// to the gateway-profiles.
func (v *Validator) ValidateGatewayProfileAccess(ctx context.Context, flag auth.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGatewayProfileAccess")
	}

	switch flag {
	case auth.Create, auth.Update, auth.Delete:
		return ctrl.st.CheckCreateUpdateDeleteGatewayProfileAccess(ctx, u.Email, u.ID)
	case auth.Read, auth.List:
		return ctrl.st.CheckReadListGatewayProfileAccess(ctx, u.Email, u.ID)
	default:
		panic("ValidateGatewayProfileAccess: unsupported flag")
	}

}
