package networkserver

import (
	"github.com/jmoiron/sqlx"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *validator {
	return &validator{otpValidator: otpValidator}
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

// validateNetworkServersAccess validates if the client has access to the
// network-servers.
func validateNetworkServersAccess(flag Flag, organizationID int64) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Create:
		// global admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
		}

		// admin api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true", "$2 = $2"},
		}
	case List:
		// global admin
		// organization user
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2"},
		}

		// admin api key
		// org api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "ak.organization_id = $2"},
		}
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, organizationID, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, organizationID)
		default:
			return false, nil
		}
	}
}

// ValidateNetworkServerAccess validates if the client has access to the
// given network-server.
func validateNetworkServerAccess(flag Flag, id int64) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join service_profile sp
			on sp.organization_id = o.id
		left join network_server ns
			on ns.id = sp.network_server_id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join service_profile sp
			on ak.organization_id = sp.organization_id
		left join network_server ns
			on sp.network_server_id = ns.id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization admin
		// organization gateway admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "ns.id = $2"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_gateway_admin = true", "ns.id = $2"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "ns.id = $2"},
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

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, id, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, id)
		default:
			return false, nil
		}
	}
}
