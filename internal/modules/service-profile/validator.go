package serviceprofile

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

type Validator struct {
	Credentials *auth.Credentials
	st          *store.Handler
}

type Validate interface {
	ValidateServiceProfileAccess(ctx context.Context, flag auth.Flag, id uuid.UUID) (bool, error)
}

func NewValidator(st *store.Handler) Validate {
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
