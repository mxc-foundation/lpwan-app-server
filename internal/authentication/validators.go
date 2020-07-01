package authentication

import (
	"strings"

	"github.com/gofrs/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

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

func ExecuteQuery(db sqlx.Queryer, query string, where [][]string, args ...interface{}) (bool, error) {
	var ors []string
	for _, ands := range where {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	query = "select count(*) from (" + query + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.Get(db, &count, query, args...); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// ValidateMulticastGroupAccess validates if the client has access to the given
// multicast-group.
func ValidateMulticastGroupAccess(flag Flag, multicastGroupID uuid.UUID) ValidatorFunc {
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

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join organization o
			on ak.organization_id = o.id
		left join service_profile sp
			on o.id = sp.organization_id
		left join multicast_group mg
			on sp.service_profile_id = mg.service_profile_id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization users
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "mg.id = $2"},
		}

		// admin api key
		// org api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "mg.id = $2"},
		}
	case Update, Delete:
		// global admin
		// organization admin users
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "mg.id = $2"},
		}

		// admin api key
		// org api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "mg.id = $2"},
		}
	}

	return func(db sqlx.Queryer, claims *Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return ExecuteQuery(db, userQuery, userWhere, claims.Username, multicastGroupID, claims.UserID)
		case SubjectAPIKey:
			return ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, multicastGroupID)
		default:
			return false, nil
		}
	}
}

// ValidateServiceProfileAccess validates if the client has access to the
// given service-profile.
func ValidateServiceProfileAccess(flag Flag, id uuid.UUID) ValidatorFunc {
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

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join service_profile sp
			on ak.organization_id = sp.organization_id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization users to which the service-profile is linked
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "sp.service_profile_id = $2"},
		}

		// admin api key
		// org api key to which the service-profile is linked
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "sp.service_profile_id = $2"},
		}
	case Update, Delete:
		// global admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
		}

		// admin api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true", "$2 = $2"},
		}
	}

	return func(db sqlx.Queryer, claims *Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return ExecuteQuery(db, userQuery, userWhere, claims.Username, id, claims.UserID)
		case SubjectAPIKey:
			return ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, id)
		default:
			return false, nil
		}
	}
}
