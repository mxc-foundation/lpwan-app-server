package auth

import (
	"context"
	"fmt"
)

// User contains information about the user
type User struct {
	ID      int64
	IsAdmin bool
}

// OrgUser contains information about the role of the user in organisation
type OrgUser struct {
	IsOrgAdmin     bool
	IsDeviceAdmin  bool
	IsGatewayAdmin bool
}

// Store provides access to information about users and their roles
type Store interface {
	// GetUser returns user's information given that there is an active user
	// with the given username
	GetUser(ctx context.Context, username string) (*User, error)
	// GetOrgUser returns user's role in the listed organization
	GetOrgUser(ctx context.Context, userID int64, orgID int64) (OrgUser, error)
}

// Credentials provides methods to assert the user's credentials
type Credentials interface {
	// Username returns user's username
	Username() string
	// UserID returns id of the user
	UserID() int64
	// IsGlobalAdmin returns an error if user is not global admin
	IsGlobalAdmin(context.Context) error
	// IsOrgUser returns an error if user does not belong to the organization
	IsOrgUser(context.Context, int64) error
	// IsOrgAdmin returns an error if user is not admin of the organization
	IsOrgAdmin(context.Context, int64) error
	// IsDeviceAdmin returns an error if user is not device admin of the organization
	IsDeviceAdmin(context.Context, int64) error
	// IsGatewayAdmin returns an error if user is not gateway admin of the organization
	IsGatewayAdmin(context.Context, int64) error
}

type credentials struct {
	st            Store
	id            int64
	username      string
	isGlobalAdmin bool
}

// GetLimitedCredentials return credentials that only contain username. All
// other checks will fail and possibly panic.
//
// Deprecated: this is only should be used for the user registration process,
// and user registration process should be fixed to not require this hack
func GetLimitedCredentials(ctx context.Context, st Store, username string) (Credentials, error) {
	return &credentials{
		id:            -1,
		username:      username,
		isGlobalAdmin: false,
	}, nil
}

// GetCredentials returns a new credentials object for the user, assuming that
// the user exists and active
func GetCredentials(ctx context.Context, st Store, username string) (Credentials, error) {
	u, err := st.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &credentials{
		st:            st,
		id:            u.ID,
		username:      username,
		isGlobalAdmin: u.IsAdmin,
	}, nil
}

// Username returns the name of the user
func (c *credentials) Username() string {
	return c.username
}

// UserID returns user id of the user
func (c *credentials) UserID() int64 {
	return c.id
}

// IsGlobalAdmin checks that the user is a global admin and returns an error if
// he's not
func (c *credentials) IsGlobalAdmin(ctx context.Context) error {
	if c.isGlobalAdmin {
		return nil
	}
	return fmt.Errorf("user is not admin")
}

// IsOrgUser checks that the user belongs to the organisation, if not it
// returns an error
func (c *credentials) IsOrgUser(ctx context.Context, orgID int64) error {
	if c.isGlobalAdmin {
		return nil
	}
	_, err := c.st.GetOrgUser(ctx, c.id, orgID)
	if err != nil {
		return err
	}
	return nil
}

// IsOrgAdmin checks that the user is admin for the organisation, if not it
// returns an error
func (c *credentials) IsOrgAdmin(ctx context.Context, orgID int64) error {
	if c.isGlobalAdmin {
		return nil
	}
	ou, err := c.st.GetOrgUser(ctx, c.id, orgID)
	if err != nil {
		return err
	}

	if !ou.IsOrgAdmin {
		return fmt.Errorf("user is not admin for organization %d", orgID)
	}
	return nil
}

// IsDeviceAdmin checks that the user is device admin for the organisation, if
// not it returns an error
func (c *credentials) IsDeviceAdmin(ctx context.Context, orgID int64) error {
	if c.isGlobalAdmin {
		return nil
	}
	ou, err := c.st.GetOrgUser(ctx, c.id, orgID)
	if err != nil {
		return err
	}
	if !(ou.IsDeviceAdmin || ou.IsOrgAdmin) {
		return fmt.Errorf("user is not device admin for organization %d", orgID)
	}
	return nil
}

// IsGatewayAdmin checks that the user is gateway admin for the organisation,
// if not it returns an error
func (c *credentials) IsGatewayAdmin(ctx context.Context, orgID int64) error {
	if c.isGlobalAdmin {
		return nil
	}
	ou, err := c.st.GetOrgUser(ctx, c.id, orgID)
	if err != nil {
		return err
	}
	if !(ou.IsGatewayAdmin || ou.IsOrgAdmin) {
		return fmt.Errorf("user is not gateway admin for organization %d", orgID)
	}
	return nil
}
