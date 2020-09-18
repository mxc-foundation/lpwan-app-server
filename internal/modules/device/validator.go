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
	ValidateDeviceQueueAccess(ctx context.Context, devEUI lorawan.EUI64, flag authcus.Flag) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)

	ValidateMulticastGroupAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateMulticastGroupAccess validates if the client has access to the given
// multicast-group.
func (v *Validator) ValidateMulticastGroupAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupAccess")
	}

	switch flag {
	case authcus.Read:
		return ctrl.st.CheckReadDeviceProfileAccess(ctx, u.UserEmail, multicastGroupID, u.ID)
	case authcus.Update, authcus.Delete:
		return ctrl.st.CheckUpdateDeleteDeviceProfileAccess(ctx, u.UserEmail, multicastGroupID, u.ID)
	default:
		panic("ValidateMulticastGroupAccess: not supported flag")
	}

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
		return ctrl.st.CheckCreateNodeAccess(ctx, u.UserEmail, applicationID, u.ID)
	case authcus.List:
		return ctrl.st.CheckListNodeAccess(ctx, u.UserEmail, applicationID, u.ID)
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
		return ctrl.st.CheckReadNodeAccess(ctx, u.UserEmail, devEUI, u.ID)
	case authcus.Update:
		return ctrl.st.CheckUpdateNodeAccess(ctx, u.UserEmail, devEUI, u.ID)
	case authcus.Delete:
		return ctrl.st.CheckDeleteNodeAccess(ctx, u.UserEmail, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}

}

// ValidateDeviceQueueAccess validates if the client has access to the queue
// of the given node.
func (v *Validator) ValidateDeviceQueueAccess(ctx context.Context, devEUI lorawan.EUI64, flag authcus.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodeAccess")
	}

	switch flag {
	case authcus.Create, authcus.List, authcus.Delete:
		return ctrl.st.CheckCreateListDeleteDeviceQueueAccess(ctx, u.UserEmail, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}
}
