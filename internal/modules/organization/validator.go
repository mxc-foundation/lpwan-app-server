package organization

import (
	"context"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
)

type Validator struct {
	Store       OrganizationStore
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

// ValidateOrganizationAccess validates if the client has access to the
// given organization.
func (v *Validator) ValidateOrganizationAccess(ctx context.Context, flag Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationAccess")
	}

	switch flag {
	case Read:
		return v.Store.CheckReadOrganizationAccess(u.Username, u.ID, organizationID)
	case Update:
		return v.Store.CheckUpdateOrganizationAccess(u.Username, u.ID, organizationID)
	case Delete:
		return v.Store.CheckDeleteOrganizationAccess(u.Username, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationsAccess validates if the client has access to the
// organizations.
func (v *Validator) ValidateOrganizationsAccess(ctx context.Context, flag Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationsAccess")
	}

	switch flag {
	case Create:
		return v.Store.CheckCreateOrganizationAccess(u.Username, u.ID)
	case List:
		return v.Store.CheckListOrganizationAccess(u.Username, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUsersAccess validates if the client has access to
// the organization users.
func (v *Validator) ValidateOrganizationUsersAccess(ctx context.Context, flag Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case Create:
		return v.Store.CheckCreateOrganizationUserAccess(u.Username, u.ID, organizationID)
	case List:
		return v.Store.CheckListOrganizationUserAccess(u.Username, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUserAccess validates if the client has access to the
// given user of the given organization.
func (v *Validator) ValidateOrganizationUserAccess(ctx context.Context, flag Flag, organizationID, userID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case Read:
		return v.Store.CheckReadOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	case Update:
		return v.Store.CheckUpdateOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	case Delete:
		return v.Store.CheckDeleteOrganizationUserAccess(u.Username, organizationID, userID, u.ID)
	default:
		panic("unsupported flag")
	}
}
