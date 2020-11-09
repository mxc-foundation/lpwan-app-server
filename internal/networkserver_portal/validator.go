package networkserver_portal

import (
	"context"

	"github.com/pkg/errors"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *auth.Credentials
}

type Validate interface {
	ValidateGlobalNetworkServersAccess(ctx context.Context, flag auth.Flag, orginizationID int64) (bool, error)
	ValidateNetworkServerAccess(ctx context.Context, flag auth.Flag, id int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateGlobalNetworkServersAccess validates if the client has access to the network-servers.
func (v *Validator) ValidateGlobalNetworkServersAccess(ctx context.Context, flag auth.Flag, orginizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalNetworkServersAccess")
	}

	switch flag {
	case auth.Create:
		return ctrl.st.CheckCreateNetworkServersAccess(ctx, u.Email, orginizationID, u.ID)
	case auth.List:
		return ctrl.st.CheckListNetworkServersAccess(ctx, u.Email, orginizationID, u.ID)
	default:
		panic("ValidateGlobalNetworkServersAccess: unsupported flag")
	}
}

// ValidateNetworkServerAccess validates if the client has access to the
// given network-server.
func (v *Validator) ValidateNetworkServerAccess(ctx context.Context, flag auth.Flag, networkserverID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNetworkServerAccess")
	}

	switch flag {
	case auth.Read:
		return ctrl.st.CheckReadNetworkServerAccess(ctx, u.Email, networkserverID, u.ID)
	case auth.Update, auth.Delete:
		return ctrl.st.CheckUpdateDeleteNetworkServerAccess(ctx, u.Email, networkserverID, u.ID)
	default:
		panic("ValidateNetworkServerAccess: unsupported flag")
	}
}
