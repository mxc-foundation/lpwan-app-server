package device

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateGlobalNodesAccess(ctx context.Context, flag authcus.Flag, applicationID int64) (bool, error)
	ValidateNodeAccess(ctx context.Context, flag authcus.Flag, devEUI lorawan.EUI64) (bool, error)
	ValidateMulticastGroupAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error)
	ValidateServiceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) ValidateServiceProfileAccess(ctx context.Context, flag authcus.Flag, id uuid.UUID) (bool, error) {
	return v.Credentials.ValidateServiceProfileAccess(ctx, flag, id)
}

func (v *Validator) ValidateMulticastGroupAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error) {
	return v.Credentials.ValidateMulticastGroupAccess(ctx, flag, multicastGroupID)
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateNodesAccess validates if the client has access to the global nodes
// resource.
func (v *Validator) ValidateGlobalNodesAccess(ctx context.Context, flag authcus.Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalNodesAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateNodeAccess(u.Username, applicationID, u.ID)
	case authcus.List:
		return Service.St.CheckListNodeAccess(u.Username, applicationID, u.ID)
	default:
		panic("ValidateGlobalNodesAccess: unsupported flag")
	}

}

// ValidateNodeAccess validates if the client has access to the given node.
func (v *Validator) ValidateNodeAccess(ctx context.Context, flag authcus.Flag, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodeAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadNodeAccess(u.Username, devEUI, u.ID)
	case authcus.Update:
		return Service.St.CheckUpdateNodeAccess(u.Username, devEUI, u.ID)
	case authcus.Delete:
		return Service.St.CheckDeleteNodeAccess(u.Username, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}

}
