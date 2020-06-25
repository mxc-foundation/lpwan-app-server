package user

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// Action defines the action type.
type Action int

// Possible actions
const (
	Select Action = iota
	Insert
	Update
	Delete
	Scan
)

// errors
var (
	ErrAlreadyExists             = errors.New("object already exists")
	ErrDoesNotExist              = errors.New("object does not exist")
	ErrUsedByOtherObjects        = errors.New("this object is used by other objects, remove them first")
	ErrUserPasswordLength        = errors.New("passwords must be at least 6 characters long")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrInvalidEmail              = errors.New("invalid e-mail")
)

func handlePSQLError(action Action, err error, description string) error {
	if err == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	switch err := err.(type) {
	case *pq.Error:
		switch err.Code.Name() {
		case "unique_violation":
			return ErrAlreadyExists
		case "foreign_key_violation":
			switch action {
			case Delete:
				return ErrUsedByOtherObjects
			default:
				return ErrDoesNotExist
			}
		}
	}

	return errors.Wrap(err, description)
}
