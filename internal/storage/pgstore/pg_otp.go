package pgstore

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type OTPPgStore interface {
	getUserID(ctx context.Context, username string) (int64, error)
	GetTOTPInfo(ctx context.Context, username string) (otp.TOTPInfo, error)
	Enable(ctx context.Context, username string) error
	Delete(ctx context.Context, username string) error
	StoreNewSecret(ctx context.Context, username, secret string) error
	storeNewSecret(ctx context.Context, tx *sql.Tx, username, secret string) error
	DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error
	AddRecoveryCodes(ctx context.Context, username string, codes []string) error
	addRecoveryCodes(ctx context.Context, tx *sql.Tx, username string, codes []string) error
	UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error
	updateLastTimeSlot(ctx context.Context, tx *sql.Tx, username string, previousValue, newValue int64) error
}

func (ps *pgstore) getUserID(ctx context.Context, username string) (int64, error) {
	row := ps.db.QueryRowContext(ctx, "SELECT id FROM \"user\" WHERE email = $1", username)
	var userID int64
	err := row.Scan(&userID)
	return userID, err
}

// GetTOTPInfo returns TOTP configuration info for the user
func (ps *pgstore) GetTOTPInfo(ctx context.Context, username string) (otp.TOTPInfo, error) {
	totpConfigQuery := `
		SELECT utc.user_id, utc.is_enabled, utc.secret, utc.last_time_slot
		FROM "user" u JOIN totp_configuration utc ON (u.id = utc.user_id)
		WHERE u.email = $1
	`
	row := ps.db.QueryRowContext(ctx, totpConfigQuery, username)
	var userID int64
	var ti otp.TOTPInfo
	if err := row.Scan(&userID, &ti.Enabled, &ti.Secret, &ti.LastTimeSlot); err != nil {
		if err == sql.ErrNoRows {
			return ti, nil
		}
		return ti, err
	}
	rows, err := ps.db.QueryContext(ctx,
		`SELECT id, code FROM totp_recovery_codes WHERE user_id = $1`, userID)
	if err != nil {
		return ti, err
	}
	defer rows.Close()
	ti.RecoveryCodes = make(map[int64]string)
	for rows.Next() {
		var codeID int64
		var code string
		if err := rows.Scan(&codeID, &code); err != nil {
			return ti, err
		}
		ti.RecoveryCodes[codeID] = code
	}
	return ti, rows.Err()
}

// Enable enables TOTP for the user
func (ps *pgstore) Enable(ctx context.Context, username string) error {
	userID, err := ps.getUserID(ctx, username)
	if err != nil {
		return err
	}
	res, err := ps.db.ExecContext(
		ctx, "UPDATE totp_configuration SET is_enabled = true WHERE user_id = $1", userID,
	)
	if err != nil {
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if num != 1 {
		return fmt.Errorf("user doesn't have totp configuration")
	}
	return nil
}

// Delete removes TOTP configuration for the user
func (ps *pgstore) Delete(ctx context.Context, username string) error {
	userID, err := ps.getUserID(ctx, username)
	if err != nil {
		return err
	}
	_, err = ps.db.ExecContext(
		ctx, "DELETE FROM totp_configuration WHERE user_id = $1", userID,
	)
	return err
}

// StoreNewSecret stores new secret for user in the database
func (ps *pgstore) StoreNewSecret(ctx context.Context, username, secret string) error {
	tx, err := ps.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	if err := ps.storeNewSecret(ctx, tx.Tx, username, secret); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			ctxlogrus.Extract(ctx).WithError(txErr).Error("couldn't rollback transaction")
		}
		return fmt.Errorf("couldn't store the secret: %v", err)
	}
	return tx.Commit()
}

func (ps *pgstore) storeNewSecret(ctx context.Context, tx *sql.Tx, username, secret string) error {
	totpConfigQuery := `
		SELECT u.id, utc.is_enabled
		FROM "user" u LEFT JOIN totp_configuration utc ON (u.id = utc.user_id)
		WHERE u.email = $1
	`
	row := tx.QueryRowContext(ctx, totpConfigQuery, username)
	var userID int64
	var enabled sql.NullBool
	if err := row.Scan(&userID, &enabled); err != nil {
		return err
	}

	if enabled.Valid && enabled.Bool {
		return fmt.Errorf("totp is already enabled")
	}
	if enabled.Valid {
		// we already have the row for the user
		updateQuery := `
			UPDATE totp_configuration
			SET secret = $1, last_time_slot = 0
			WHERE user_id = $2
		`
		_, err := tx.ExecContext(ctx, updateQuery, secret, userID)
		if err != nil {
			return err
		}
	} else {
		insertQuery := `
			INSERT INTO totp_configuration (user_id, secret)
			VALUES ($1, $2)
		`
		_, err := tx.ExecContext(ctx, insertQuery, userID, secret)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteRecoveryCode removes recovery code from the database
func (ps *pgstore) DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error {
	userID, err := ps.getUserID(ctx, username)
	if err != nil {
		return err
	}
	_, err = ps.db.ExecContext(
		ctx, `DELETE FROM totp_recovery_codes WHERE user_id = $1 AND id = $2`, userID, codeID,
	)

	return err
}

// AddRecoveryCodes adds new recovery codes to the database
func (ps *pgstore) AddRecoveryCodes(ctx context.Context, username string, codes []string) error {
	tx, err := ps.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	if err := ps.addRecoveryCodes(ctx, tx.Tx, username, codes); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			ctxlogrus.Extract(ctx).WithError(txErr).Error("couldn't rollback transaction")
		}
		return fmt.Errorf("couldn't add recovery codes for %s: %v", username, err)
	}
	return tx.Commit()
}

func (ps *pgstore) addRecoveryCodes(ctx context.Context, tx *sql.Tx, username string, codes []string) error {
	totpConfigQuery := `
		SELECT u.id, utc.last_time_slot
		FROM "user" u LEFT JOIN totp_configuration utc ON (u.id = utc.user_id)
		WHERE u.email = $1
	`
	row := tx.QueryRowContext(ctx, totpConfigQuery, username)
	var userID int64
	var dbLastTimeSlot sql.NullInt64
	if err := row.Scan(&userID, &dbLastTimeSlot); err != nil {
		return err
	}
	if !dbLastTimeSlot.Valid {
		return fmt.Errorf("user doesn't have totp configured")
	}

	recoveryCodesQuery := `
		SELECT count(*) FROM totp_recovery_codes
		WHERE user_id = $1
	`
	row = tx.QueryRowContext(ctx, recoveryCodesQuery, userID)
	var count int
	if err := row.Scan(&count); err != nil {
		return err
	}
	if count+len(codes) > 10 {
		return fmt.Errorf("can't add %d recovery codes as user already has %d", len(codes), count)
	}

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO totp_recovery_codes (user_id, code) VALUES ($1, $2)`,
	)
	if err != nil {
		return err
	}
	for _, code := range codes {
		if _, err := stmt.ExecContext(ctx, userID, code); err != nil {
			return err
		}
	}

	return nil
}

// UpdateLastTimeSlot updates last time slot value in the database
func (ps *pgstore) UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error {
	tx, err := ps.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	if err := ps.updateLastTimeSlot(ctx, tx.Tx, username, previousValue, newValue); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			ctxlogrus.Extract(ctx).WithError(txErr).Error("couldn't rollback transaction")
		}
		return fmt.Errorf("couldn't update last_time_slot: %v", err)
	}
	return tx.Commit()
}

func (ps *pgstore) updateLastTimeSlot(ctx context.Context, tx *sql.Tx, username string, previousValue, newValue int64) error {
	totpConfigQuery := `
		SELECT u.id, utc.last_time_slot
		FROM "user" u LEFT JOIN totp_configuration utc ON (u.id = utc.user_id)
		WHERE u.email = $1
	`
	row := tx.QueryRowContext(ctx, totpConfigQuery, username)
	var userID int64
	var lastTimeSlot sql.NullInt64
	if err := row.Scan(&userID, &lastTimeSlot); err != nil {
		return err
	}
	if !lastTimeSlot.Valid {
		return fmt.Errorf("user don't have totp configuration")
	}
	_, err := tx.ExecContext(ctx,
		`UPDATE totp_configuration SET last_time_slot=$1 WHERE user_id=$2`,
		newValue, userID,
	)
	return err
}
