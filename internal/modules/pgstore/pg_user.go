package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

const externalUserFields = "id, is_admin, is_active, session_ttl, created_at, updated_at, email, note, security_token"
const internalUserFields = "*"

func (ps *pgstore) CheckActiveUser(ctx context.Context, userEmail string, userID int64) (bool, error) {
	var userQuery = "select count(*) from " +
		"(select 1 from public.user u where (u.email = $1 or u.id = $2) " +
		"and u.is_active = true limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckCreateUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
	`
	// global admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckListUserAcess(ctx context.Context, userEmail string, userID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
		left join organization_user ou
			on u.id = ou.user_id
	`
	// global admin users
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $2)", "u.is_active = true", "u.is_admin = true"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckReadUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// global admin
	// user itself
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdateUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// user itself
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckDeleteUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// global admin
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true", "$2 = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdatePasswordUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`

	// user itself
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

func (ps *pgstore) CheckUpdateProfileUserAccess(ctx context.Context, userEmail string, userID, operatorUserID int64) (bool, error) {
	userQuery := `
		select
			1
		from
			"user" u
	`
	// global admin
	// user itself
	var userWhere = [][]string{
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.is_admin = true"},
		{"(u.email = $1 or u.id = $3)", "u.is_active = true", "u.id = $2"},
	}

	var ors []string
	for _, ands := range userWhere {
		ors = append(ors, "(("+strings.Join(ands, ") and (")+"))")
	}
	whereStr := strings.Join(ors, " or ")
	userQuery = "select count(*) from (" + userQuery + " where " + whereStr + " limit 1) count_only"

	var count int64
	if err := sqlx.GetContext(ctx, ps.db, &count, userQuery, userEmail, userID, operatorUserID); err != nil {
		return false, errors.Wrap(err, "select error")
	}
	return count > 0, nil
}

// CreateUser creates the given user.
func (ps *pgstore) CreateUser(ctx context.Context, user *store.User) error {
	if err := user.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}

	err := ps.SetUserPassword(user, user.Password)
	if err != nil {
		return err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = sqlx.GetContext(ctx, ps.db, &user.ID, `
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
		return handlePSQLError(Insert, err, "insert error")
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
func (ps *pgstore) GetUser(ctx context.Context, id int64) (store.User, error) {
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+externalUserFields+" from \"user\" where id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, store.ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByExternalID returns the User for the given ext. ID.
func (ps *pgstore) GetUserByExternalID(ctx context.Context, externalID string) (store.User, error) {
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+externalUserFields+" from \"user\" external_id id = $1", externalID)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, store.ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByUsername returns the User for the given email.
func (ps *pgstore) GetUserByUsername(ctx context.Context, userEmail string) (store.User, error) {
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+externalUserFields+" from \"user\" where email = $1", userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, store.ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserByEmail returns the User for the given userEmail.
func (ps *pgstore) GetUserByEmail(ctx context.Context, userEmail string) (store.User, error) {
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+externalUserFields+" from \"user\" where email = $1", userEmail)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, store.ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetUserCount returns the total number of users.
func (ps *pgstore) GetUserCount(ctx context.Context) (int, error) {
	var count int
	err := sqlx.GetContext(ctx, ps.db, &count, `
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
func (ps *pgstore) GetUsers(ctx context.Context, limit, offset int) ([]store.User, error) {
	var users []store.User
	err := sqlx.SelectContext(ctx, ps.db, &users, "select "+externalUserFields+
		" from public.user order by email limit $1 offset $2", limit, offset)

	if err != nil {
		return nil, errors.Wrap(err, "select error")
	}
	return users, nil
}

// UpdateUser updates the given User.
func (ps *pgstore) UpdateUser(ctx context.Context, u *store.User) error {
	if err := u.Validate(); err != nil {
		return errors.Wrap(err, "validate user error")
	}

	u.UpdatedAt = time.Now()

	res, err := ps.db.ExecContext(ctx, `
		update "user"
		set
			updated_at = $2,
			is_admin = $3,
			is_active = $4,
			session_ttl = $5,
			email = $6,
			email_verified = $7,
			note = $8,
			external_id = $9
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
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "get rows affected error")
	}
	if ra == 0 {
		return store.ErrDoesNotExist
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
func (ps *pgstore) DeleteUser(ctx context.Context, id int64) error {
	res, err := ps.db.ExecContext(ctx, `
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
		return store.ErrDoesNotExist
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("storage: user deleted")
	return nil
}

// UpdatePassword updates the user with the new password.
func (ps *pgstore) UpdatePassword(ctx context.Context, id int64, newpassword string) error {
	if err := store.ValidatePassword(newpassword); err != nil {
		return errors.Wrap(err, "validation error")
	}

	pwHash, err := ps.s.PWH.HashPassword(newpassword)
	if err != nil {
		return err
	}

	// update password
	_, err = ps.db.ExecContext(ctx, "update \"user\" set password_hash = $1, updated_at = now() where id = $2",
		pwHash, id)
	if err != nil {
		return errors.Wrap(err, "update error")
	}

	log.WithFields(log.Fields{
		"id":     id,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("user password updated")
	return nil

}

// LoginUserByPassword checks the password for the user matching the given email
func (ps *pgstore) LoginUserByPassword(ctx context.Context, userEmail string, password string) error {
	// get the user by userEmail
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+internalUserFields+" from \"user\" where email = $1", userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrDoesNotExist
		}
		return errors.Wrap(err, "select error")
	}

	// Compare the passed in password with the hash in the database.
	if err := ps.VerifyUserPassword(password, user.PasswordHash); err != nil {
		return errors.Wrap(err, "password doesn't match email")
	}

	return nil
}

// GetProfile returns the user profile (user, applications and organizations
// to which the user is linked).
func (ps *pgstore) GetProfile(ctx context.Context, id int64) (store.UserProfile, error) {
	var prof store.UserProfile

	user, err := ps.GetUser(ctx, id)
	if err != nil {
		return prof, errors.Wrap(err, "get user error")
	}
	prof.User = store.UserProfileUser{
		ID:         user.ID,
		Email:      user.Email,
		SessionTTL: user.SessionTTL,
		IsAdmin:    user.IsAdmin,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	err = sqlx.SelectContext(ctx, ps.db, &prof.Organizations, `
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

func (ps *pgstore) GetUserToken(ctx context.Context, u store.User) (string, error) {
	// Generate the token.
	now := time.Now()
	nowSecondsSinceEpoch := now.Unix()
	var expSecondsSinceEpoch int64
	if u.SessionTTL > 0 {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + (60 * int64(u.SessionTTL))
	} else {
		expSecondsSinceEpoch = nowSecondsSinceEpoch + int64(store.DefaultSessionTTL/time.Second)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "lpwan-app-server",
		"aud":   "lpwan-app-server",
		"nbf":   nowSecondsSinceEpoch,
		"exp":   expSecondsSinceEpoch,
		"sub":   "user",
		"id":    u.ID,
		"email": u.Email, // backwards compatibility
	})

	jwt, err := token.SignedString([]byte(ps.s.JWTSecret))
	if err != nil {
		return jwt, errors.Wrap(err, "get jwt signed string error")
	}
	return jwt, err
}

// RegisterUser ...
func (ps *pgstore) RegisterUser(ctx context.Context, user *store.User, token string) error {
	if err := user.Validate(); err != nil {
		return errors.Wrap(err, "validation error")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Add the new user.
	err := sqlx.GetContext(ctx, ps.db, &user.ID, `
		insert into "user"(
			email,
			is_admin,
			is_active,
			session_ttl,
			created_at,
			updated_at,
			email_verified,
			note,
			external_id,
			security_token)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)	returning id`,
		user.Email,
		user.IsAdmin,
		user.IsActive,
		user.SessionTTL,
		user.CreatedAt,
		user.UpdatedAt,
		user.EmailVerified,
		user.Note,
		user.ExternalID,
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
func (ps *pgstore) GetUserByToken(ctx context.Context, token string) (store.User, error) {
	var user store.User
	err := sqlx.GetContext(ctx, ps.db, &user, "select "+externalUserFields+" from \"user\" where security_token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, store.ErrDoesNotExist
		}
		return user, errors.Wrap(err, "select error")
	}

	return user, nil
}

// GetTokenByUsername ...
func (ps *pgstore) GetTokenByUsername(ctx context.Context, userEmail string) (string, error) {
	//var user User
	var otp string
	err := sqlx.GetContext(ctx, ps.db, &otp, "select security_token from \"user\" where email = $1", userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return otp, store.ErrDoesNotExist
		}
		return otp, errors.Wrap(err, "select error")
	}

	return otp, nil
}

// FinishRegistration ...
func (ps *pgstore) FinishRegistration(ctx context.Context, userID int64, password string) error {
	if err := store.ValidatePassword(password); err != nil {
		return errors.Wrap(err, "validation error")
	}

	pwdHash, err := ps.s.PWH.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = ps.db.ExecContext(ctx, `
	update
	"user"
	set
	password_hash = $1,
		is_active = true,
		security_token = null,
		updated_at = now()
	where
	id = $2
	`,
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

func (ps *pgstore) SetOTP(ctx context.Context, pr *store.PasswordResetRecord) error {
	res, err := ps.db.ExecContext(ctx, `
	UPDATE
	password_reset
	SET
	otp = $1, generated_at = $2, attempts_left = $3
	WHERE
	user_id = $4
	`,
		pr.OTP, pr.GeneratedAt, pr.AttemptsLeft, pr.UserID)
	if err != nil {
		return err
	}
	// we need to make sure that we've updated exactly one row
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return fmt.Errorf("expected to update 1 row, but updated %d", rowCnt)
	}
	return nil
}

func (ps *pgstore) ReduceAttempts(ctx context.Context, pr *store.PasswordResetRecord) error {
	pr.AttemptsLeft--
	res, err := ps.db.ExecContext(ctx, `
	UPDATE
	password_reset
	SET
	attempts_left = $1
	WHERE
	user_id = $2
	`,
		pr.AttemptsLeft, pr.UserID)
	if err != nil {
		return err
	}
	// we need to make sure that we've updated exactly one row
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt != 1 {
		return fmt.Errorf("expected to update 1 row, but updated %d", rowCnt)
	}
	return nil
}

func (ps *pgstore) GetPasswordResetRecord(ctx context.Context, userID int64) (*store.PasswordResetRecord, error) {
	query := `
	SELECT
	otp, generated_at, attempts_left
	FROM
	password_reset
	WHERE
	user_id = $1
	`
	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	res := &store.PasswordResetRecord{UserID: userID}
	var count int
	defer rows.Close()
	for rows.Next() {
		if count > 0 {
			return nil, fmt.Errorf("got multiple reset password rows for %d", userID)
		}
		count++
		if err := rows.Scan(&res.OTP, &res.GeneratedAt, &res.AttemptsLeft); err != nil {
			return nil, fmt.Errorf("scan has failed: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if count == 0 {
		_, err := ps.db.ExecContext(ctx, `
	INSERT
	INTO
	password_reset(user_id, otp, generated_at, attempts_left)
	VALUES($1, $2, $3, $4)`, res.UserID, res.OTP, res.GeneratedAt, res.AttemptsLeft)
		if err != nil {
			return nil, fmt.Errorf("failed to add a new password reset record: %v", err)
		}
	}
	return res, nil
}

func (ps *pgstore) SetUserPassword(user *store.User, pw string) error {
	pwHash, err := ps.s.PWH.HashPassword(pw)
	if err != nil {
		return err
	}

	user.PasswordHash = pwHash
	return nil
}

func (ps *pgstore) VerifyUserPassword(pw string, pwHash string) error {
	return ps.s.PWH.Validate(pw, pwHash)
}
