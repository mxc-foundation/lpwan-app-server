package gatewayprofile

import (
	"context"

	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateGatewayProfileAccess(ctx context.Context, flag authcus.Flag) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

// ValidateGatewayProfileAccess validates if the client has access
// to the gateway-profiles.
func (v *Validator) ValidateGatewayProfileAccess(ctx context.Context, flag authcus.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGatewayProfileAccess")
	}

	switch flag {
	case authcus.Create, authcus.Update, authcus.Delete:
		return Service.St.CheckCreateUpdateDeleteGatewayProfileAccess(ctx, u.UserEmail, u.ID)
	case authcus.Read, authcus.List:
		return Service.St.CheckReadListGatewayProfileAccess(ctx, u.UserEmail, u.ID)
	default:
		panic("ValidateGatewayProfileAccess: unsupported flag")
	}

}
