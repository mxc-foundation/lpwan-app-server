package topup

import (
	"context"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	GetIsAdmin(ctx context.Context) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetIsAdmin(ctx context.Context) (bool, error) {
	return v.Credentials.IsGlobalAdmin(ctx)
}
