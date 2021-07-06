package withdraw

import (
	"context"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

// Validator defines struct type for vadidating user access to APIs provided by this package
type Validator struct {
	Credentials *cred.Credentials
}

// Validate defines methods used on struct Validator
type Validate interface {
	IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error
	IsOrgAdmin(ctx context.Context, orgID int64, opts ...cred.Option) error
}

// NewValidator returns new Validate instance for this package
func NewValidator() Validate {
	return &Validator{
		Credentials: cred.NewCredentials(),
	}
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...cred.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

func (v *Validator) IsOrgAdmin(ctx context.Context, orgID int64, opts ...cred.Option) error {
	return v.Credentials.IsOrgAdmin(ctx, orgID, opts...)
}
