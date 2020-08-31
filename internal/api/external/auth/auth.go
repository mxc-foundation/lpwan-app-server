package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"golang.org/x/net/context"
)

// Claims defines the struct containing the token claims.
type Claims struct {
	jwt.StandardClaims

	// Username defines the identity of the user.
	Username string `json:"username"`

	// UserID defines the ID of th user.
	UserID int64 `json:"user_id"`

	// APIKeyID defines the API key ID.
	APIKeyID uuid.UUID `json:"api_key_id"`
}

// Validator defines the interface a validator needs to implement.
type Validator interface {
	// Validate validates the given set of validators against the given context.
	// Must return after the first validator function either returns true or
	// and error. The way how the validation must be seens is:
	//   if validatorFunc1 || validatorFunc2 || validatorFunc3 ...
	// In case multiple validators must validate to true, then a validator
	// func needs to be implemented which validates a given set of funcs as:
	//   if validatorFunc1 && validatorFunc2 && ValidatorFunc3 ...
	Validate(context.Context, ...ValidatorFunc) error

	// GetSubject returns the claim subject.
	GetSubject(context.Context) (string, error)

	// GetUser returns the user object.
	GetUser(context.Context) (storage.User, error)

	// GetAPIKey returns the API key ID.
	GetAPIKeyID(context.Context) (uuid.UUID, error)
}

// ValidatorFunc defines the signature of a claim validator function.
// It returns a bool indicating if the validation passed or failed and an
// error in case an error occurred (e.g. db connectivity).
type ValidatorFunc func(sqlx.Queryer, *Claims) (bool, error)
