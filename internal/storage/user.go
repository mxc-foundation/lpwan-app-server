package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// User defines the user structure.
type User struct {
	ID            int64     `db:"id"`
	IsAdmin       bool      `db:"is_admin"`
	IsActive      bool      `db:"is_active"`
	SessionTTL    int32     `db:"session_ttl"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	PasswordHash  string    `db:"password_hash"`
	Email         string    `db:"email"`
	EmailVerified bool      `db:"email_verified"`
	EmailOld      string    `db:"email_old"`
	Note          string    `db:"note"`
	ExternalID    *string   `db:"external_id"` // must be pointer for unique index
	SecurityToken *string   `db:"security_token"`
}

// GetUser returns the User for the given id.
func GetUser(ctx context.Context, db sqlx.Queryer, id int64) (User, error) {
	var user User

	err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByEmail returns the User for the given email.
func GetUserByEmail(ctx context.Context, db sqlx.Queryer, email string) (User, error) {
	var user User

	err := sqlx.Get(db, &user, `
		select
			*
		from
			"user"
		where
			email = $1
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}
