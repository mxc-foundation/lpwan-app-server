package gateway

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateGlobalGatewaysAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateGatewayAccess(ctx context.Context, flag authcus.Flag, mac lorawan.EUI64) (bool, error)
	ValidateOrganizationNetworkServerAccess(ctx context.Context, flag authcus.Flag, organizationID, networkServerID int64) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)
	IsGlobalAdmin(ctx context.Context) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

func (v *Validator) IsGlobalAdmin(ctx context.Context) (bool, error) {
	return v.Credentials.IsGlobalAdmin(ctx)
}

// ValidateGatewaysAccess validates if the client has access to the gateways.
func (v *Validator) ValidateGlobalGatewaysAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalGatewaysAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateGatewayAccess(u.Username, organizationID, u.ID)
	case authcus.List:
		return Service.St.CheckListGatewayAccess(u.Username, organizationID, u.ID)
	default:
		panic("ValidateGlobalGatewaysAccess: unsupported flag")
	}
}

// ValidateGatewayAccess validates if the client has access to the given gateway.
func (v *Validator) ValidateGatewayAccess(ctx context.Context, flag authcus.Flag, mac lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGatewayAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadGatewayAccess(u.Username, mac, u.ID)
	case authcus.Update, authcus.Delete:
		return Service.St.CheckUpdateDeleteGatewayAccess(u.Username, mac, u.ID)
	default:
		panic("ValidateGatewayAccess: unsupported flag")
	}
}

// ValidateOrganizationNetworkServerAccess validates if the given client has
// access to the given organization id / network server id combination.
func (v *Validator) ValidateOrganizationNetworkServerAccess(ctx context.Context, flag authcus.Flag, organizationID, networkServerID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationNetworkServerAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadOrganizationNetworkServerAccess(u.Username, organizationID, networkServerID, u.ID)
	default:
		panic("ValidateOrganizationNetworkServerAccess: unsupported flag")
	}
}
