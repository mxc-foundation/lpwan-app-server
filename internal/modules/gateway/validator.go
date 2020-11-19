package gateway

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *auth.Credentials
}

type Validate interface {
	ValidateGlobalGatewaysAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateGatewayAccess(ctx context.Context, flag auth.Flag, mac lorawan.EUI64) (bool, error)
	ValidateOrganizationNetworkServerAccess(ctx context.Context, flag auth.Flag, organizationID, networkServerID int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
	IsGlobalAdmin(ctx context.Context, opts ...auth.Option) error
	IsOrgAdmin(ctx context.Context, organizationID int64, opts ...auth.Option) error
}

func NewValidator() Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

func (v *Validator) IsGlobalAdmin(ctx context.Context, opts ...auth.Option) error {
	return v.Credentials.IsGlobalAdmin(ctx, opts...)
}

func (v *Validator) IsOrgAdmin(ctx context.Context, organizationID int64, opts ...auth.Option) error {
	return v.Credentials.IsOrgAdmin(ctx, organizationID, opts...)
}

// ValidateGlobalGatewaysAccess validates if the client has access to the gateways.
func (v *Validator) ValidateGlobalGatewaysAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalGatewaysAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateGatewayAccess(ctx, u.Email, organizationID, u.ID)
	case auth.List:
		return ctrl.st.CheckListGatewayAccess(ctx, u.Email, organizationID, u.ID)
	default:
		panic("ValidateGlobalGatewaysAccess: unsupported flag")
	}
}

// ValidateGatewayAccess validates if the client has access to the given gateway.
func (v *Validator) ValidateGatewayAccess(ctx context.Context, flag auth.Flag, mac lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGatewayAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadGatewayAccess(ctx, u.Email, mac, u.ID)
	case auth.Update, auth.Delete:
		return ctrl.st.CheckUpdateDeleteGatewayAccess(ctx, u.Email, mac, u.ID)
	default:
		panic("ValidateGatewayAccess: unsupported flag")
	}
}

// ValidateOrganizationNetworkServerAccess validates if the given client has
// access to the given organization id / network server id combination.
func (v *Validator) ValidateOrganizationNetworkServerAccess(ctx context.Context, flag auth.Flag, organizationID, networkServerID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationNetworkServerAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadOrganizationNetworkServerAccess(ctx, u.Email, organizationID, networkServerID, u.ID)
	default:
		panic("ValidateOrganizationNetworkServerAccess: unsupported flag")
	}
}
