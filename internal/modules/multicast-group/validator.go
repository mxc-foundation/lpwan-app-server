package multicast

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
	ValidateMulticastGroupsAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error)
	ValidateMulticastGroupAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error)
	ValidateMulticastGroupQueueAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)

	ValidateNodeAccess(ctx context.Context, flag authcus.Flag, devEUI lorawan.EUI64) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
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

// ValidateMulticastGroupsAccess validates if the client has access to the
// multicast-groups.
func (v *Validator) ValidateMulticastGroupsAccess(ctx context.Context, flag authcus.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupsAccess")
	}

	switch flag {
	case authcus.Create:
		return ctrl.st.CheckCreateServiceProfilesAccess(ctx, u.UserEmail, organizationID, u.ID)
	case authcus.List:
		return ctrl.st.CheckListServiceProfilesAccess(ctx, u.UserEmail, organizationID, u.ID)
	default:
		panic("ValidateMulticastGroupsAccess: not supported flag")
	}

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

// ValidateMulticastGroupQueueAccess validates if the client has access to
// the given multicast-group queue.
func (v *Validator) ValidateMulticastGroupQueueAccess(ctx context.Context, flag authcus.Flag, multicastGroupID uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupQueueAccess")
	}

	switch flag {
	case authcus.Create, authcus.Read, authcus.List, authcus.Delete:
		return ctrl.st.CheckMulticastGroupQueueAccess(ctx, u.UserEmail, multicastGroupID, u.ID)
	default:
		panic("ValidateMulticastGroupQueueAccess: not supported flag")
	}

}
