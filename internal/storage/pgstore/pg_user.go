package pgstore

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// GetUserIDByExternalUserID gets userID with given external userID and external service
func (ps *PgStore) GetUserIDByExternalUserID(ctx context.Context, service string, externalUserID string) (int64, error) {
	var userID int64
	err := sqlx.GetContext(ctx, ps.db, &userID, `
		select user_id from external_login where service=$1 and external_id=$2`,
		service,
		externalUserID,
	)
	if err != nil {
		return 0, handlePSQLError(Select, err, "select error")
	}
	return userID, nil
}

// GetExternalUserByUserIDAndService gets external user with given userID and external service
func (ps *PgStore) GetExternalUserByUserIDAndService(ctx context.Context, service string, userID int64) (user.ExternalUser, error) {
	var externalUser user.ExternalUser
	err := sqlx.GetContext(ctx, ps.db, &externalUser, `
		select user_id, service, external_id, external_username from external_login where service=$1 and user_id=$2`,
		service,
		userID,
	)
	if err != nil {
		return externalUser, handlePSQLError(Select, err, "select error")
	}
	return externalUser, nil
}

// GetExternalUsersByUserID gets external user list with given userID
func (ps *PgStore) GetExternalUsersByUserID(ctx context.Context, userID int64) ([]user.ExternalUser, error) {
	var externalUsers []user.ExternalUser
	err := sqlx.SelectContext(ctx, ps.db, &externalUsers, `
		select user_id, service, external_id, external_username from external_login where user_id=$1`,
		userID)
	if err != nil {
		return externalUsers, handlePSQLError(Select, err, "select error")
	}
	return externalUsers, nil
}

// AddExternalUserLogin adds new external user binding relation
func (ps *PgStore) AddExternalUserLogin(ctx context.Context, service string, userID int64, externalUserID, externalUsername string) error {
	res, err := ps.db.ExecContext(ctx, `
		insert into external_login (user_id , service, external_id, external_username) values ($1, $2, $3, $4)`,
		userID, service, externalUserID, externalUsername,
	)
	if err != nil {
		return handlePSQLError(Insert, err, "insert error")
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("no record affected")
	}

	return nil
}

// SetExternalUsername updates external user username
func (ps *PgStore) SetExternalUsername(ctx context.Context, service, externalUserID, externalUsername string) error {
	res, err := ps.db.ExecContext(ctx, `
		update external_login set external_username = $1 where service = $2 and external_id = $3`,
		externalUsername, service, externalUserID,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("no record affected")
	}
	return err
}

// DeleteExternalUserLogin delete external user binding relation
func (ps *PgStore) DeleteExternalUserLogin(ctx context.Context, userID int64, service string) error {
	res, err := ps.db.ExecContext(ctx, `
		delete from external_login where service = $1 and user_id = $2`,
		service, userID,
	)
	if err != nil {
		return handlePSQLError(Delete, err, "delete error")
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("no record affected")
	}
	return err
}

// GetExternalUserByUsername returns external user with given external username
func (ps *PgStore) GetExternalUserByUsername(ctx context.Context, service, externalUsername string) (user.ExternalUser, error) {
	var user user.ExternalUser
	err := sqlx.GetContext(ctx, ps.db, &user, `
		select * from external_login where service=$1 and external_username=$2
	`, service, externalUsername)
	if err != nil {
		return user, handlePSQLError(Select, err, "select error")
	}

	return user, nil
}

// GetExternalUserByToken returns external user with given security token
func (ps *PgStore) GetExternalUserByToken(ctx context.Context, service, token string) (user.ExternalUser, error) {
	var user user.ExternalUser
	err := sqlx.GetContext(ctx, ps.db, &user, `
		select * from external_login where service=$1 and external_id=$2
	`, service, token)
	if err != nil {
		return user, handlePSQLError(Select, err, "select error")
	}

	return user, nil
}

// SetExternalUserID updates external id of an external user
func (ps *PgStore) SetExternalUserID(ctx context.Context, extUser user.ExternalUser) error {
	res, err := ps.db.ExecContext(ctx, `
		update 
			external_login
		set 
			external_id=$1 
		where 
			user_id=$2 
		and 
			service=$3 
		and 
			external_username=$4
	`, extUser.ExternalUserID, extUser.UserID, extUser.ServiceName, extUser.ExternalUsername)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("no record affected")
	}

	return err
}

// SetUserLastLogin updates user display name and last login service type
func (ps *PgStore) SetUserLastLogin(ctx context.Context, userID int64, displayName, service string) error {
	res, err := ps.db.ExecContext(ctx, `
		UPDATE 
			"user" 
		SET 
			display_name = $1, updated_at = NOW() , last_login_service = $2 
		WHERE id = $3`,
		displayName, service, userID,
	)
	if err != nil {
		return handlePSQLError(Update, err, "update error")
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return fmt.Errorf("no record affected")
	}

	return err
}

// CreateUser creates the given user.
func (ps *PgStore) CreateUser(ctx context.Context, u user.User, ou []user.OrganizationUser) (user.User, error) {
	err := ps.Tx(ctx, func(ctx context.Context, ps *PgStore) error {
		nu, err := ps.createUser(ctx, u)
		if err != nil {
			return err
		}
		u.ID = nu.ID
		for _, org := range ou {
			org.UserID = u.ID
			if err := ps.createOrganizationUser(ctx, org); err != nil {
				return err
			}
		}
		return nil
	})
	return u, err
}

func (ps *PgStore) createUser(ctx context.Context, u user.User) (user.User, error) {
	var id int64
	err := ps.db.QueryRowContext(ctx,
		`INSERT INTO "user"
	 	(
			is_admin,
			is_active,
			created_at,
			updated_at,
			password_hash,
			email,
			email_verified,
			security_token
		)
		VALUES ($1, $2, NOW(), NOW(), $3, $4, $5, $6)
		RETURNING id`,
		u.IsAdmin,
		u.IsActive,
		u.PasswordHash,
		u.Email,
		u.EmailVerified,
		u.SecurityToken,
	).Scan(&id)
	if err != nil {
		return u, handlePSQLError(Insert, err, "couldn't insert user")
	}
	u.ID = id
	return u, nil
}

func (ps *PgStore) createOrganizationUser(ctx context.Context, ou user.OrganizationUser) error {
	_, err := ps.db.ExecContext(ctx,
		`INSERT INTO organization_user (
			organization_id,
			user_id,
			is_admin,
			is_device_admin,
			is_gateway_admin,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
		ou.OrganizationID,
		ou.UserID,
		ou.IsOrgAdmin,
		ou.IsDeviceAdmin,
		ou.IsGatewayAdmin,
	)
	return err
}

func (ps *PgStore) getUser(ctx context.Context, condition string, args ...interface{}) (user.User, error) {
	var u user.User
	var pass, token sql.NullString
	// considering that this is a private function and condition passed into it
	// is always a literal and doesn't come from some external sources I think
	// it should be ok to concatenate
	// nolint: gosec
	err := ps.db.QueryRowContext(ctx,
		`SELECT id, created_at, updated_at, email, password_hash,
			is_active, is_admin, security_token, email_verified, display_name
		 FROM "user"
		 WHERE `+condition, args...).Scan(
		&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Email, &pass,
		&u.IsActive, &u.IsAdmin, &token, &u.EmailVerified, &u.DisplayName,
	)
	if err == sql.ErrNoRows {
		err = errHandler.ErrDoesNotExist
	}
	u.PasswordHash = pass.String
	u.SecurityToken = token.String
	return u, err
}

// GetUserByID returns the User for the given id.
func (ps *PgStore) GetUserByID(ctx context.Context, userID int64) (user.User, error) {
	return ps.getUser(ctx, "id = $1", userID)
}

// GetUserByEmail returns the User for the given email.
func (ps *PgStore) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	return ps.getUser(ctx, "email = $1", email)
}

func (ps *PgStore) GetUserByToken(ctx context.Context, token string) (user.User, error) {
	return ps.getUser(ctx, "security_token = $1", token)
}

func (ps *PgStore) GetUserOrganizations(ctx context.Context, userID int64) ([]user.OrganizationUser, error) {
	query := `SELECT ou.user_id, ou.organization_id, ou.created_at, ou.updated_at,
				ou.is_admin, ou.is_device_admin, ou.is_gateway_admin, o.name
			  FROM organization_user ou
			    JOIN organization o ON (o.id = ou.organization_id)
			  WHERE ou.user_id = $1`
	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []user.OrganizationUser
	for rows.Next() {
		var ou user.OrganizationUser
		err := rows.Scan(&ou.UserID, &ou.OrganizationID, &ou.CreatedAt, &ou.UpdatedAt,
			&ou.IsOrgAdmin, &ou.IsDeviceAdmin, &ou.IsGatewayAdmin, &ou.OrganizationName)
		if err != nil {
			return nil, err
		}
		res = append(res, ou)
	}
	return res, rows.Err()
}

// GetUserCount returns the total number of users.
func (ps *PgStore) GetUserCount(ctx context.Context) (int64, error) {
	var count int64
	err := ps.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "user"`).Scan(&count)
	return count, err
}

// GetUsers returns a slice of users, respecting the given limit and offset.
func (ps *PgStore) GetUsers(ctx context.Context, limit, offset int) ([]user.User, error) {
	query := `SELECT id, created_at, updated_at, email, password_hash,
			    is_active, is_admin, security_token, email_verified, display_name
		 	  FROM "user"
			  ORDER BY email
			  LIMIT $1
			  OFFSET $2`
	rows, err := ps.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []user.User
	for rows.Next() {
		var u user.User
		var pass, token sql.NullString
		err := rows.Scan(
			&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.Email, &pass,
			&u.IsActive, &u.IsAdmin, &token, &u.EmailVerified, &u.DisplayName)
		if err != nil {
			return nil, err
		}
		u.PasswordHash = pass.String
		u.SecurityToken = token.String
		users = append(users, u)
	}
	return users, rows.Err()
}

// SetUserDisplayName updates display name of the user
func (ps *PgStore) SetUserDisplayName(ctx context.Context, displayName string, userID int64) error {
	query := `UPDATE "user"
			  SET display_name = $1, updated_at = NOW()
			  WHERE id = $2`
	_, err := ps.db.ExecContext(ctx, query, displayName, userID)
	return err
}

// SetUserEmail changes the email address of the user
func (ps *PgStore) SetUserEmail(ctx context.Context, userID int64, email string) error {
	query := `UPDATE "user"
			  SET email = $1, updated_at = NOW()
			  WHERE id = $2`
	_, err := ps.db.ExecContext(ctx, query, email, userID)
	return err
}

// SetUserActiveStatus disables or enables the user
func (ps *PgStore) SetUserActiveStatus(ctx context.Context, userID int64, isActive bool) error {
	query := `UPDATE "user"
			  SET is_active = $1, updated_at = NOW()
			  WHERE id = $2`
	_, err := ps.db.ExecContext(ctx, query, isActive, userID)
	return err
}

// SetUserPasswordHash sets the password hash for the user
func (ps *PgStore) SetUserPasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	query := `UPDATE "user"
			  SET password_hash = $1, updated_at = NOW()
			  WHERE id = $2`
	_, err := ps.db.ExecContext(ctx, query, passwordHash, userID)
	return err
}

func (ps *PgStore) unsetUserSecurityToken(ctx context.Context, userID int64) error {
	query := `UPDATE "user"
			  SET security_token = null
			  WHERE id = $1`
	_, err := ps.db.ExecContext(ctx, query, userID)
	return err
}

type prRecord struct {
	otp          string
	generatedAt  time.Time
	attemptsLeft int64
}

func (ps *PgStore) getPasswordResetRecord(ctx context.Context, userID int64) (prRecord, error) {
	query := `SELECT otp, generated_at, attempts_left
				 FROM password_reset
				 WHERE user_id = $1`
	var pr prRecord
	err := ps.db.QueryRowContext(ctx, query, userID).
		Scan(&pr.otp, &pr.generatedAt, &pr.attemptsLeft)
	return pr, err
}

func (ps *PgStore) setPasswordResetRecord(ctx context.Context, userID int64, pr prRecord) error {
	query := `INSERT INTO password_reset (
				user_id, otp, generated_at, attempts_left)
			  VALUES ($1, $2, $3, $4)
			  ON CONFLICT (user_id) DO UPDATE SET
			 	otp = $2, generated_at = $3, attempts_left = $4`
	_, err := ps.db.ExecContext(ctx, query, userID, pr.otp, pr.generatedAt, pr.attemptsLeft)
	return err
}

// GetOrSetPasswordResetOTP if the password reset OTP has been generated
// already then returns it, otherwise sets the new OTP and returns it
func (ps *PgStore) GetOrSetPasswordResetOTP(ctx context.Context, userID int64, otp string) (string, error) {
	err := ps.Tx(ctx, func(ctx context.Context, ps *PgStore) error {
		pr, err := ps.getPasswordResetRecord(ctx, userID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows ||
			pr.generatedAt.Before(time.Now().Add(-30*24*time.Hour)) {
			// if no password reset record, or the record is old create a new one
			return ps.setPasswordResetRecord(ctx, userID, prRecord{
				otp:          otp,
				generatedAt:  time.Now(),
				attemptsLeft: 3,
			})
		}
		if pr.attemptsLeft == 0 {
			return fmt.Errorf("can not reset password more often than once a month")
		}
		// if password reset record is still valid, just return the existing OTP
		otp = pr.otp
		return nil
	})
	return otp, err
}

// SetUserPasswordIfOTPMatch sets the user's password if the OTP provided is correct
func (ps *PgStore) SetUserPasswordIfOTPMatch(ctx context.Context, userID int64, otp, passwordHash string) error {
	var invalidOTP bool
	err := ps.Tx(ctx, func(ctx context.Context, ps *PgStore) error {
		pr, err := ps.getPasswordResetRecord(ctx, userID)
		if err != nil {
			return err
		}
		if pr.generatedAt.Before(time.Now().Add(-30*24*time.Hour)) || len(pr.otp) < 6 {
			// pr record has expired
			return errHandler.ErrDoesNotExist
		}
		if pr.attemptsLeft < 1 {
			return fmt.Errorf("no attempts left")
		}
		pr.attemptsLeft--
		if subtle.ConstantTimeCompare([]byte(pr.otp), []byte(otp)) == 1 {
			// otp matches, unset otp, update password
			pr.otp = ""
			pr.attemptsLeft = 0
			if err := ps.SetUserPasswordHash(ctx, userID, passwordHash); err != nil {
				return err
			}
		} else {
			invalidOTP = true
		}
		return ps.setPasswordResetRecord(ctx, userID, pr)
	})
	if err == nil && invalidOTP {
		return fmt.Errorf("invalid otp")
	}
	return err
}

// DeleteUser deletes the User record matching the given ID.
func (ps *PgStore) DeleteUser(ctx context.Context, userID int64) error {
	res, err := ps.db.ExecContext(ctx,
		`DELETE FROM
			"user"
		WHERE
			id = $1`,
		userID)
	if err != nil {
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get the number of rows deleted: %w", err)
	}
	if ra == 0 {
		return errHandler.ErrDoesNotExist
	}
	return nil
}

func (ps *PgStore) createOrganization(ctx context.Context, org user.Organization) (user.Organization, error) {
	query := `INSERT INTO organization (
				name,
				display_name,
				created_at,
				updated_at,
				can_have_gateways,
				max_gateway_count,
				max_device_count
			  )
			  VALUES ($1, $2, NOW(), NOW(), $3, $4, $5)
			  RETURNING id, created_at, updated_at`
	err := ps.db.QueryRowContext(ctx, query, org.Name, org.DisplayName, org.CanHaveGateways,
		org.MaxGatewayCount, org.MaxDeviceCount).
		Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)
	return org, err
}

// ActivateUser creates the organization for the new user, adds the user to
// the org and activates the user
func (ps *PgStore) ActivateUser(ctx context.Context, userID int64, passwordHash, orgName, orgDisplayName string) error {
	err := ps.Tx(ctx, func(ctx context.Context, ps *PgStore) error {
		org := user.Organization{
			Name:            orgName,
			DisplayName:     orgDisplayName,
			CanHaveGateways: true,
		}
		var err error
		org, err = ps.createOrganization(ctx, org)
		if err != nil {
			return err
		}
		ou := user.OrganizationUser{
			OrganizationID: org.ID,
			UserID:         userID,
			IsOrgAdmin:     true,
		}
		if err := ps.createOrganizationUser(ctx, ou); err != nil {
			return err
		}
		if err := ps.SetUserActiveStatus(ctx, userID, true); err != nil {
			return err
		}
		if err := ps.SetUserPasswordHash(ctx, userID, passwordHash); err != nil {
			return err
		}
		if err := ps.unsetUserSecurityToken(ctx, userID); err != nil {
			return err
		}
		return nil
	})
	return err
}
