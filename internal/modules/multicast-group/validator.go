package multicast

import (
	"context"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *auth.Credentials
}

type Validate interface {
	ValidateMulticastGroupsAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateMulticastGroupAccess(ctx context.Context, flag auth.Flag, multicastGroupID uuid.UUID) (bool, error)
	ValidateMulticastGroupQueueAccess(ctx context.Context, flag auth.Flag, multicastGroupID uuid.UUID) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)

	ValidateNodeAccess(ctx context.Context, flag auth.Flag, devEUI lorawan.EUI64) (bool, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateNodeAccess validates if the client has access to the given node.
func (v *Validator) ValidateNodeAccess(ctx context.Context, flag auth.Flag, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodeAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadNodeAccess(ctx, u.Email, devEUI, u.ID)
	case auth.Update:
		return ctrl.st.CheckUpdateNodeAccess(ctx, u.Email, devEUI, u.ID)
	case auth.Delete:
		return ctrl.st.CheckDeleteNodeAccess(ctx, u.Email, devEUI, u.ID)
	default:
		panic("ValidateNodeAccess: unsupported flag")
	}

}

// ValidateMulticastGroupsAccess validates if the client has access to the
// multicast-groups.
func (v *Validator) ValidateMulticastGroupsAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupsAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateServiceProfilesAccess(ctx, u.Email, organizationID, u.ID)
	case auth.List:
		return ctrl.st.CheckListServiceProfilesAccess(ctx, u.Email, organizationID, u.ID)
	default:
		panic("ValidateMulticastGroupsAccess: not supported flag")
	}

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
		return ctrl.st.CheckReadDeviceProfileAccess(ctx, u.Email, multicastGroupID, u.ID)
	case auth.Update, auth.Delete:
		return ctrl.st.CheckUpdateDeleteDeviceProfileAccess(ctx, u.Email, multicastGroupID, u.ID)
	default:
		panic("ValidateMulticastGroupAccess: not supported flag")
	}

}

// ValidateMulticastGroupQueueAccess validates if the client has access to
// the given multicast-group queue.
func (v *Validator) ValidateMulticastGroupQueueAccess(ctx context.Context, flag auth.Flag, multicastGroupID uuid.UUID) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupQueueAccess")
	}

	switch flag {
	case auth.Create, auth.Read, auth.List, auth.Delete:
		return ctrl.st.CheckMulticastGroupQueueAccess(ctx, u.Email, multicastGroupID, u.ID)
	default:
		panic("ValidateMulticastGroupQueueAccess: not supported flag")
	}

}
