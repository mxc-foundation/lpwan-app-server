package store

import (
	"context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/otp/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pg pgstore.PgStore) *otps {
	return &otps{
		pg: pg,
	}
}

type otps struct {
	pg pgstore.OTPPgStore
}

// Store provides methods to access TOTP information for users stored in the database
type Store interface {
	// GetTOTPInfo returns TOTP configuration info for the user
	GetTOTPInfo(ctx context.Context, username string) (TOTPInfo, error)
	// Enable enables TOTP for the user
	Enable(ctx context.Context, username string) error
	// Delete removes TOTP configuration for the user
	Delete(ctx context.Context, username string) error
	// StoreNewSecret stores new secret for user in the database
	StoreNewSecret(ctx context.Context, username, secret string) error
	// DeleteRecoveryCode removes recovery code from the database
	DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error
	// AddRecoveryCodes adds new recovery codes to the database
	AddRecoveryCodes(ctx context.Context, username string, codes []string) error
	// UpdateLastTimeSlot updates last time slot value in the database
	UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error
}

// GetTOTPInfo returns TOTP configuration info for the user
func (s *otps) GetTOTPInfo(ctx context.Context, username string) (TOTPInfo, error) {
	return s.pg.GetTOTPInfo(ctx, username)
}

// Enable enables TOTP for the user
func (s *otps) Enable(ctx context.Context, username string) error {
	return s.pg.Enable(ctx, username)
}

// Delete removes TOTP configuration for the user
func (s *otps) Delete(ctx context.Context, username string) error {
	return s.pg.Delete(ctx, username)
}

// StoreNewSecret stores new secret for user in the database
func (s *otps) StoreNewSecret(ctx context.Context, username, secret string) error {
	return s.pg.StoreNewSecret(ctx, username, secret)
}

// DeleteRecoveryCode removes recovery code from the database
func (s *otps) DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error {
	return s.pg.DeleteRecoveryCode(ctx, username, codeID)
}

// AddRecoveryCodes adds new recovery codes to the database
func (s *otps) AddRecoveryCodes(ctx context.Context, username string, codes []string) error {
	return s.pg.AddRecoveryCodes(ctx, username, codes)
}

// UpdateLastTimeSlot updates last time slot value in the database
func (s *otps) UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error {
	return s.pg.UpdateLastTimeSlot(ctx, username, previousValue, newValue)
}
