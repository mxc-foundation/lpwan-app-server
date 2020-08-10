package withdraw

import (
	"context"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	IsGlobalAdmin(ctx context.Context, opts ...authcus.Option) error
	IsOrgAdmin(ctx context.Context, orgID int64, opts ...authcus.Option) error
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...authcus.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

func (v *Validator) IsOrgAdmin(ctx context.Context, orgID int64, opts ...authcus.Option) error {
	return v.Credentials.IsOrgAdmin(ctx, orgID, opts...)
}
