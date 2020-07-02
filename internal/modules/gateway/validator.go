package gateway

import (
	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type Validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *Validator {
	return &Validator{otpValidator: otpValidator}
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

// ValidateGatewaysAccess validates if the client has access to the gateways.
func validateGatewaysAccess(flag Flag, organizationID int64) authcus.ValidatorFunc {
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
		left join organization o
			on ak.organization_id = o.id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Create:
		// global admin
		// organization admin
		// gateway admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_admin = true", "o.can_have_gateways = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "o.id = $2", "ou.is_gateway_admin = true", "o.can_have_gateways = true"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "o.id = $2", "o.can_have_gateways = true"},
		}
	case List:
		// global admin
		// organization user
		// any active user (result filtered on user)
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 > 0", "o.id = $2"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "$2 = 0"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "o.id = $2"},
		}

	default:
		panic("unsupported flag")
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

// ValidateGatewayAccess validates if the client has access to the given gateway.
func validateGatewayAccess(flag Flag, mac lorawan.EUI64) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join gateway g
			on o.id = g.organization_id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join gateway g
			on ak.organization_id = g.organization_id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization user
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "g.mac = $2"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "g.mac = $2"},
		}
	case Update, Delete:
		userWhere = [][]string{
			// global admin
			// organization admin
			// organization gateway admin
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "g.mac = $2", "ou.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "g.mac = $2", "ou.is_gateway_admin = true"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "g.mac = $2"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, mac[:], claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, mac[:])
		default:
			return false, nil
		}
	}
}

// ValidateOrganizationNetworkServerAccess validates if the given client has
// access to the given organization id / network server id combination.
func validateOrganizationNetworkServerAccess(flag Flag, organizationID, networkServerID int64) authcus.ValidatorFunc {
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
		left join device_profile dp
			on dp.organization_id = o.id
		left join network_server ns
			on ns.id = sp.network_server_id or ns.id = dp.network_server_id
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
		// organization user
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "o.id = $2", "ns.id = $3"},
		}

		// admin api key
		// organization api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "ak.organization_id = $2", "ns.id = $3"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, organizationID, networkServerID, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, organizationID, networkServerID)
		default:
			return false, nil
		}
	}
}
