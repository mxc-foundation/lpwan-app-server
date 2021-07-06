package serviceprofile

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

// Validator defines struct type for vadidating user access to APIs provided by this package
type Validator struct {
	Credentials *auth.Credentials
	st          Store
}

// Validate defines methods used on struct Validator
type Validate interface {
	ValidateServiceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error)
}

// NewValidator returns new Validate instance for this package
func NewValidator(st Store) Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
		st:          st,
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
		return v.st.CheckReadServiceProfileAccess(ctx, u.Email, id, u.ID)
	case auth.Update, auth.Delete:
		return v.st.CheckUpdateDeleteServiceProfileAccess(ctx, u.Email, id, u.ID)
	default:
		panic("ValidateServiceProfileAccess: not supported flag")
	}

}
