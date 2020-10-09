package topup

import (
	"context"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *cred.Credentials
}

type Validate interface {
	IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error
	IsOrgAdmin(ctx context.Context, ordID int64, opts ...cred.Option) error
}

func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
	}
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

func (v *Validator) IsOrgAdmin(ctx context.Context, ordID int64, opts ...cred.Option) error {
	return v.Credentials.IsOrgAdmin(ctx, ordID, opts...)
}
