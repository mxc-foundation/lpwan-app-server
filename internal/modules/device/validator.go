package device

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	cred "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication/data"
)

type Validator struct {
	h           *store.Handler
	Credentials *cred.Credentials
}

type Validate interface {
	ValidateGlobalNodesAccess(ctx context.Context, flag auth.Flag, applicationID int64) (bool, error)
	ValidateNodeAccess(ctx context.Context, flag auth.Flag, devEUI lorawan.EUI64) (bool, error)
	ValidateDeviceQueueAccess(ctx context.Context, devEUI lorawan.EUI64, flag auth.Flag) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)

	ValidateMulticastGroupAccess(ctx context.Context, flag auth.Flag, multicastGroupID uuid.UUID) (bool, error)
}

func NewValidator(h *store.Handler) Validate {
	return &Validator{
		h:           h,
		Credentials: cred.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateMulticastGroupAccess validates if the client has access to the given
// multicast-group.
func (v *Validator) ValidateMulticastGroupAccess(ctx context.Context, flag auth.Flag, multicastGroupID uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupAccess")
	}

	switch flag {
	case auth.Read:
		return v.h.CheckReadDeviceProfileAccess(ctx, u.Email, multicastGroupID, u.ID)
	case auth.Update, auth.Delete:
		return v.h.CheckUpdateDeleteDeviceProfileAccess(ctx, u.Email, multicastGroupID, u.ID)
	default:
		panic("ValidateMulticastGroupAccess: not supported flag")
	}

}

// ValidateNodesAccess validates if the client has access to the global nodes
// resource.
func (v *Validator) ValidateGlobalNodesAccess(ctx context.Context, flag auth.Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalNodesAccess")
	}

	switch flag {
	case auth.Create:
		return v.h.CheckCreateNodeAccess(ctx, u.Email, applicationID, u.ID)
	case auth.List:
		return v.h.CheckListNodeAccess(ctx, u.Email, applicationID, u.ID)
	default:
		panic("ValidateGlobalNodesAccess: unsupported flag")
	}

}

// ValidateNodeAccess validates if the client has access to the given node.
func (v *Validator) ValidateNodeAccess(ctx context.Context, flag auth.Flag, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodeAccess")
	}

	switch flag {
	case auth.Read:
		return v.h.CheckReadNodeAccess(ctx, u.Email, devEUI, u.ID)
	case auth.Update:
		return v.h.CheckUpdateNodeAccess(ctx, u.Email, devEUI, u.ID)
	case auth.Delete:
		return v.h.CheckDeleteNodeAccess(ctx, u.Email, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}

}

// ValidateDeviceQueueAccess validates if the client has access to the queue
// of the given node.
func (v *Validator) ValidateDeviceQueueAccess(ctx context.Context, devEUI lorawan.EUI64, flag auth.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodeAccess")
	}

	switch flag {
	case auth.Create, auth.List, auth.Delete:
		return v.h.CheckCreateListDeleteDeviceQueueAccess(ctx, u.Email, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}
}
