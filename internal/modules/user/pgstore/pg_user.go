package pgstore

import (
	"context"
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	umod "github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

const externalUserFields = "id, username, is_admin, is_active, session_ttl, created_at, updated_at, email, note, security_token"
const internalUserFields = "*"

type UserPgHandler struct {
	db *storage.DBLogger
}

func New(db *storage.DBLogger) *UserPgHandler {
	return &UserPgHandler{
		db: db,
	}
}

// CreateUser creates the given user.
func (h *UserPgHandler) CreateUser(ctx context.Context, user *umod.User) error {
	if err := user.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := sqlx.Get(h.db, &user.ID, `
		insert into "user" (
			is_admin,
			is_active,
			session_ttl,
			created_at,
			updated_at,
			password_hash,
			email,
			email_verified,
			note,
			external_id
		)
		values (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning
			id`,
		user.IsAdmin,
		user.IsActive,
		user.SessionTTL,
		user.CreatedAt,
		user.UpdatedAt,
		user.PasswordHash,
		user.Email,
		user.EmailVerified,
		user.Note,
		user.ExternalID,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	var externalID string
	if user.ExternalID != nil {
		externalID = *user.ExternalID
	}

	log.WithFields(log.Fields{
		"id":          user.ID,
		"external_id": externalID,
		"email":       user.Email,
		"ctx_id":      ctx.Value(logging.ContextIDKey),
	}).Info("storage: user created")

	return nil
}

// GetUser returns the User for the given id.
func (h *UserPgHandler) GetUser(ctx context.Context, id int64) (umod.User, error) {
	var user umod.User

	err := sqlx.Get(h.db, &user, `
		select
			*
		from
			"user"
		where
			id = $1
	`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("not exist")
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByExternalID returns the User for the given ext. ID.
func (h *UserPgHandler) GetUserByExternalID(ctx context.Context, externalID string) (umod.User, error) {
	var user umod.User

	err := sqlx.Get(h.db, &user, `
		select
			*
		from
			"user"
		where
			external_id = $1
	`, externalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("not exist")
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByEmail returns the User for the given email.
func (h *UserPgHandler) GetUserByEmail(ctx context.Context, email string) (umod.User, error) {
	var user umod.User

	err := sqlx.Get(h.db, &user, `
		select
			*
		from
			"user"
		where
			email = $1
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("not exist")
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserCount returns the total number of users.
func (h *UserPgHandler) GetUserCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.Get(h.db, &count, `
		select
			count(*)
		from "user"
	`)
	if err != nil {
		return 0, errors.Wrap(err, "select error")
	}
	return count, nil
}

// GetUsers returns a slice of users, respecting the given limit and offset.
func (h *UserPgHandler) GetUsers(ctx context.Context, limit, offset int) ([]umod.User, error) {
	var users []umod.User

	err := sqlx.Select(h.db, &users, `
		select
			*
		from
			"user"
		order by
			email
		limit $1
		offset $2
	`, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return users, nil
}

// UpdateUser updates the given User.
func (h *UserPgHandler) UpdateUser(ctx context.Context, u *umod.User) error {
	if err := u.Validate(); err != nil {
		return errors.Wrap(err, "validate user error")
	}

	u.UpdatedAt = time.Now()

	res, err := h.db.Exec(`
		update "user"
		set
			updated_at = $2,
			is_admin = $3,
			is_active = $4,
			session_ttl = $5,
			email = $6,
			email_verified = $7,
			note = $8,
			external_id = $9,
			password_hash = $10
		where
			id = $1`,
		u.ID,
		u.UpdatedAt,
		u.IsAdmin,
		u.IsActive,
		u.SessionTTL,
		u.Email,
		u.EmailVerified,
		u.Note,
		u.ExternalID,
		u.PasswordHash,
	)
	if err != nil {
		return errors.Wrap(err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	var extUser string
	if u.ExternalID != nil {
		extUser = *u.ExternalID
	}

	log.WithFields(log.Fields{
		"id":          u.ID,
		"external_id": extUser,
		"ctx_id":      ctx.Value(logging.ContextIDKey),
	}).Info("storage: user updated")

	return nil
}

// DeleteUser deletes the User record matching the given ID.
func (h *UserPgHandler) DeleteUser(ctx context.Context, id int64) error {
	res, err := h.db.Exec(`
		delete from
			"user"
		where
			id = $1
	`, id)
	if err != nil {
		return errors.Wrap(err, "delete error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return errors.New("not exist")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("storage: user deleted")
	return nil
}

// LoginUserByPassword returns a JWT token for the user matching the given email
// and password combination.
func (h *UserPgHandler) LoginUserByPassword(ctx context.Context, email string, password string) (string, error) {
	// get the user by email
	var user umod.User
	err := sqlx.Get(h.db, &user, `
		select
			*
		from
			"user"
		where
			email = $1
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid username or password")
		}
		return "", errors.Wrap(err, "select error")
	}

	// Compare the passed in password with the hash in the database.
	if !user.HashCompare(password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	return h.GetUserToken(user)
}

// GetProfile returns the user profile (user, applications and organizations
// to which the user is linked).
func (h *UserPgHandler) GetProfile(ctx context.Context, id int64) (umod.UserProfile, error) {
	var prof umod.UserProfile

	user, err := h.GetUser(ctx, id)
	if err != nil {
		return prof, errors.Wrap(err, "get user error")
	}
	prof.User = umod.UserProfileUser{
		ID:         user.ID,
		Email:      user.Email,
		SessionTTL: user.SessionTTL,
		IsAdmin:    user.IsAdmin,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	err = sqlx.Select(h.db, &prof.Organizations, `
		select
			ou.organization_id as organization_id,
			o.name as organization_name,
			ou.is_admin as is_admin,
			ou.is_device_admin as is_device_admin,
			ou.is_gateway_admin as is_gateway_admin,
			ou.created_at as created_at,
			ou.updated_at as updated_at
		from
			organization_user ou,
			organization o
		where
			ou.user_id = $1
			and ou.organization_id = o.id`,
		id,
	)
	if err != nil {
		return prof, errors.Wrap(err, "select error")
	}

	return prof, nil
}

func (h *UserPgHandler) GetUserToken(u umod.User) (string, error) {
	// Generate the token.
	now := time.Now()
	nowSecondsSinceEpoch := now.Unix()
	var expSecondsSinceEpoch int64
	if u.SessionTTL > 0 {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + (60 * int64(u.SessionTTL))
	} else {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + int64(umod.DefaultSessionTTL/time.Second)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "lpwan-app-server",
		"aud":      "lpwan-app-server",
		"nbf":      nowSecondsSinceEpoch,
		"exp":      expSecondsSinceEpoch,
		"sub":      "user",
		"id":       u.ID,
		"username": u.Email, // backwards compatibility
	})

	jwt, err := token.SignedString([]byte(config.C.ApplicationServer.ExternalAPI.JWTSecret))
	if err != nil {
		return jwt, errors.Wrap(err, "get jwt signed string error")
	}
	return jwt, err
}

// RegisterUser ...
func (h *UserPgHandler) RegisterUser(user *umod.User, token string) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Add the new user.
	err := sqlx.Get(h.db, &user.ID, `
		insert into "user" (
			email,
			is_admin,
			is_active,
			session_ttl,
			created_at,
			updated_at,
	        security_token)
		values (
			$1, $2, $3, $4, $5, $6, $7) returning id`,
		user.Email,
		user.IsAdmin,
		user.IsActive,
		user.SessionTTL,
		user.CreatedAt,
		user.UpdatedAt,
		token,
	)
	if err != nil {
		return errors.Wrap(err, "insert error")
	}

	log.WithFields(log.Fields{
		"email":       user.Email,
		"session_ttl": user.SessionTTL,
		"is_admin":    user.IsAdmin,
	}).Info("Registration: user created")
	return nil
}

// GetUserByToken ...
func (h *UserPgHandler) GetUserByToken(token string) (umod.User, error) {
	var user umod.User
	err := sqlx.Get(h.db, &user, "select "+externalUserFields+" from \"user\" where security_token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("not exist")
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetTokenByUsername ...
func (h *UserPgHandler) GetTokenByUsername(ctx context.Context, username string) (string, error) {
	//var user User
	var otp string
	err := sqlx.Get(h.db, &otp, "select security_token from \"user\" where username = $1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return otp, errors.New("not exist")
		}
		return otp, errors.Wrap(err, "select error")
	}

	return otp, nil
}

// FinishRegistration ...
func (h *UserPgHandler) FinishRegistration(userID int64, pwdHash string) error {
	_, err := h.db.Exec(`
		update "user"
		set
			password_hash = $1,
			is_active = true,
			security_token = null,
			updated_at = now()
		where id = $2`,
		pwdHash,
		userID,
	)
	if err != nil {
		return errors.Wrap(err, "update error")
	}

	log.WithFields(log.Fields{
		"id": userID,
	}).Info("user password updated")

	return nil
}
