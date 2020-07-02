package device

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

// ValidateNodesAccess validates if the client has access to the global nodes
// resource.
func validateNodesAccess(applicationID int64, flag Flag) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join application a
			on ak.application_id = a.id or ak.organization_id = a.organization_id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Create:
		// global admin
		// organization admin
		// organization device admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "a.id = $2"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin = true", "a.id = $2"},
		}

		// admin api key
		// organization api key
		// application api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "a.id = $2"}, // application is joined on a.id and a.organization_id
		}
	case List:
		// global admin
		// organization user
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "a.id = $2"},
		}

		// admin api key
		// organization api key
		// application api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "a.id = $2"}, // application is joined on a.id and a.organization_id
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, applicationID, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, applicationID)
		default:
			return false, nil
		}
	}
}

// ValidateNodeAccess validates if the client has access to the given node.
func validateNodeAccess(devEUI lorawan.EUI64, flag Flag) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on o.id = ou.organization_id
		left join application a
			on a.organization_id = o.id
		left join device d
			on a.id = d.application_id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
		left join application a
			on ak.application_id = a.id or ak.organization_id = a.organization_id
		left join device d
			on a.id = d.application_id
	`

	var userWhere = [][]string{}
	var apiKeyWhere = [][]string{}

	switch flag {
	case Read:
		// global admin
		// organization user
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "d.dev_eui = $2"},
		}

		// admin api key
		// organization api key
		// application api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "d.dev_eui = $2"}, // application is joined on a.id and a.organization_id
		}

	case Update:
		// global admin
		// organization admin
		// organization device admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "d.dev_eui = $2"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin = true", "d.dev_eui = $2"},
		}

		// admin api key
		// organization api key
		// application api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "d.dev_eui = $2"}, // application is joined on a.id and a.organization_id
		}
	case Delete:
		// global admin
		// organization admin
		// organization device admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "d.dev_eui = $2"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_device_admin = true", "d.dev_eui = $2"},
		}

		// admin api key
		// organization api key
		// application api key
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
			{"ak.id = $1", "d.dev_eui = $2"}, // application is joined on a.id and a.organization_id
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, devEUI[:], claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID, devEUI[:])
		default:
			return false, nil
		}
	}
}
