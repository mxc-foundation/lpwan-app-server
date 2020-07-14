package networkserver

import (
	"context"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Credentials *authcus.Credentials
}

type Validate interface {
	ValidateGlobalNetworkServersAccess(ctx context.Context, flag authcus.Flag, orginizationID int64) (bool, error)
	ValidateNetworkServerAccess(ctx context.Context, flag authcus.Flag, id int64) (bool, error)
	GetUser(ctx context.Context) (authcus.User, error)
}

func NewValidator() Validate {
	return &Validator{
		Credentials: authcus.NewCredentials(),
	}
}

func (v *Validator) GetUser(ctx context.Context) (authcus.User, error) {
	return v.Credentials.GetUser(ctx)
}

// validateNetworkServersAccess validates if the client has access to the
// network-servers.
func (v *Validator) ValidateGlobalNetworkServersAccess(ctx context.Context, flag authcus.Flag, orginizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateGlobalNetworkServersAccess")
	}

	switch flag {
	case authcus.Create:
		return Service.St.CheckCreateNetworkServersAccess(u.Username, orginizationID, u.ID)
	case authcus.List:
		return Service.St.CheckListNetworkServersAccess(u.Username, orginizationID, u.ID)
	default:
		panic("ValidateGlobalNetworkServersAccess: unsupported flag")
	}
}

// ValidateNetworkServerAccess validates if the client has access to the
// given network-server.
func (v *Validator) ValidateNetworkServerAccess(ctx context.Context, flag authcus.Flag, networkserverID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNetworkServerAccess")
	}

	switch flag {
	case authcus.Read:
		return Service.St.CheckReadNetworkServerAccess(u.Username, networkserverID, u.ID)
	case authcus.Update, authcus.Delete:
		return Service.St.CheckUpdateDeleteNetworkServerAccess(u.Username, networkserverID, u.ID)
	default:
		panic("ValidateNetworkServerAccess: unsupported flag")
	}
}
