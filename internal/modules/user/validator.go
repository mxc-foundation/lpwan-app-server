package user

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *validator {
	return &validator{otpValidator: otpValidator}
}

// GetIsAdmin returns if the authenticated user is a global amin.
func (v *validator) GetIsAdmin(ctx context.Context) (bool, error) {
	claims, err := v.otpValidator.JwtValidator.GetClaims(ctx)
	if err != nil {
		return false, err
	}

	user, err := userAPI.Store.GetUserByEmail(ctx, claims.Username)
	if err != nil {
		return false, errors.Wrap(err, "get user by username error")
	}

	return user.IsAdmin, nil
}

// GetUser returns the user object.
func (v *validator) GetUser(ctx context.Context) (User, error) {
	claims, err := v.otpValidator.JwtValidator.GetClaims(ctx)
	if err != nil {
		return User{}, err
	}

	if claims.Subject != "user" {
		return User{}, errors.New("subject must be user")
	}

	if claims.UserID != 0 {
		return userAPI.Store.GetUser(ctx, claims.UserID)
	}

	if claims.Username != "" {
		return userAPI.Store.GetUserByEmail(ctx, claims.Username)
	}

	return User{}, errors.New("no username or user_id in claims")
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

// ValidateActiveUser validates if the user in the JWT claim is active.
func validateActiveUser() authcus.ValidatorFunc {
	query := `
		select
			1
		from
			"user" u
	`

	where := [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true"},
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, query, where, claims.Username, claims.UserID)
		case SubjectAPIKey:
			return false, nil
		default:
			return false, nil
		}
	}
}

// ValidateUsersAccess validates if the client has access to the global users
// resource.
func validateUsersAccess(flag Flag) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
	`

	var userWhere [][]string
	var apiKeyWhere [][]string

	switch flag {
	case Create:
		// global admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
		}

		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
		}
	case List:
		// global admin users
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
		}

		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID)
		default:
			return false, nil
		}
	}
}

// ValidateUserAccess validates if the client has access to the given user
// resource.
func validateUserAccess(userID int64, flag Flag) authcus.ValidatorFunc {
	userQuery := `
		select
			1
		from
			"user" u
	`

	apiKeyQuery := `
		select
			1
		from
			api_key ak
	`

	var userWhere [][]string
	var apiKeyWhere [][]string

	switch flag {
	case Read:
		// global admin
		// user itself
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
		}

		// admin token
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
		}
	case Update, Delete:
		// global admin
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
		}

		// admin token
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
		}
	case UpdateProfile:
		// global admin
		// user itself
		userWhere = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
		}

		// admin token
		apiKeyWhere = [][]string{
			{"ak.id = $1", "ak.is_admin = true"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		switch claims.Subject {
		case SubjectUser:
			return authcus.ExecuteQuery(db, userQuery, userWhere, claims.Username, userID, claims.UserID)
		case SubjectAPIKey:
			return authcus.ExecuteQuery(db, apiKeyQuery, apiKeyWhere, claims.APIKeyID)
		default:
			return false, nil
		}
	}
}

// ValidateAPIKeysAccess validates if the client has access to the global
// API key resource.
func validateAPIKeysAccess(flag Flag, organizationID int64, applicationID int64) authcus.ValidatorFunc {
	query := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on ou.organization_id = o.id
		left join application a
			on o.id = a.organization_id
	`

	var where [][]string

	switch flag {
	case Create:
		// global admin
		// organization admin of given org id
		// organization admin of given app id
		where = [][]string{
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 > 0", "$3 = 0", "o.id = $2"},
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$3 > 0", "$2 = 0", "a.id = $3"},
		}

	case List:
		// global admin
		// organization admin of given org id (api key filtered by org in api)
		// organization admin of given app id (api key filtered by app in api)
		where = [][]string{
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$2 > 0", "$3 = 0", "o.id = $2"},
			{"(u.email = $1 or u.id = $4)", "u.is_active = true", "ou.is_admin = true", "$3 > 0", "$2 = 0", "a.id = $3"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		return authcus.ExecuteQuery(db, query, where, claims.Username, organizationID, applicationID, claims.UserID)
	}
}

// ValidateAPIKeyAccess validates if the client has access to the given API
// key.
func validateAPIKeyAccess(flag Flag, id uuid.UUID) authcus.ValidatorFunc {
	query := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
		left join organization o
			on ou.organization_id = o.id
		left join application a
			on o.id = a.organization_id
		left join api_key ak
			on a.id = ak.application_id or o.id = ak.organization_id or u.is_admin
	`

	var where [][]string
	switch flag {
	case Delete:
		// global admin
		// organization admin
		where = [][]string{
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
			{"(u.email = $1 or u.id = $3)", "u.is_active = true", "ou.is_admin = true", "ak.id = $2"},
		}
	default:
		panic("unsupported flag")
	}

	return func(db sqlx.Queryer, claims *authcus.Claims) (bool, error) {
		return authcus.ExecuteQuery(db, query, where, claims.Username, id, claims.UserID)
	}
}
