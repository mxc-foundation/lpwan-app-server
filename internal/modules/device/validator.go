package device

import (
	"context"
	"github.com/brocaar/lorawan"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Store       DeviceStore
	Credentials *authcus.Credentials
}

func NewValidator(v Validator) *Validator {
	return &Validator{
		Store:       v.Store,
		Credentials: v.Credentials,
	}
}

// API key subjects.
const (
	SubjectUser   = "user"
	SubjectAPIKey = "api_key"
)

// Flag defines the authorization flag.
type Flag int

// Authorization flags.
const (
	Create Flag = iota
	Read
	Update
	Delete
	List
	UpdateProfile
	FinishRegistration
)

// ValidateNodesAccess validates if the client has access to the global nodes
// resource.
func (v *Validator) ValidateGlobalNodesAccess(ctx context.Context, flag Flag, applicationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateNodesAccess")
	}

	switch flag {
	case Create:
		return v.Store.CheckCreateNodeAccess(u.Username, applicationID, u.ID)
	case List:
		return v.Store.CheckListNodeAccess(u.Username, applicationID, u.ID)
	default:
		panic("ValidateNodesAccess: unsupported flag")
	}

}

// ValidateNodeAccess validates if the client has access to the given node.
func (v *Validator) ValidateNodeAccess(ctx context.Context, flag Flag, devEUI lorawan.EUI64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "validateNodeAccess")
	}

	switch flag {
	case Read:
		return v.Store.CheckReadNodeAccess(u.Username, devEUI, u.ID)
	case Update:
		return v.Store.CheckUpdateNodeAccess(u.Username, devEUI, u.ID)
	case Delete:
		return v.Store.CheckDeleteNodeAccess(u.Username, devEUI, u.ID)
	default:
		panic("validateNodeAccess: unsupported flag")
	}

}
