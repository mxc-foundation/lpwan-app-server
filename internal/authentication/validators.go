package authentication

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
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
	UpdatePassword
	FinishRegistration
)

// ValidateMulticastGroupAccess validates if the client has access to the given
// multicast-group.
func (c *Credentials) ValidateMulticastGroupAccess(ctx context.Context, flag Flag, multicastGroupID uuid.UUID) (bool, error) {
	u, err := c.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateMulticastGroupAccess")
	}

	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
		left join multicast_group mg
			on sp.service_profile_id = mg.service_profile_id
	`

	var userWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization users
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "mg.id = $2"},
		}
	case Update, Delete:
		// global admin
		// organization admin users
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "mg.id = $2"},
		}
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(storage.DB().DB, &count, userQuery, u.UserEmail, multicastGroupID, u.ID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// ValidateServiceProfileAccess validates if the client has access to the
// given service-profile.
func (c *Credentials) ValidateServiceProfileAccess(ctx context.Context, flag Flag, id uuid.UUID) (bool, error) {
	u, err := c.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateServiceProfileAccess")
	}

	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join service_profile sp
			on sp.organization_id = ou.organization_id
	`

	var userWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization users to which the service-profile is linked
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "sp.service_profile_id = $2"},
		}
	case Update, Delete:
		// global admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
		}
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(storage.DB().DB, &count, userQuery, u.UserEmail, id, u.ID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}
