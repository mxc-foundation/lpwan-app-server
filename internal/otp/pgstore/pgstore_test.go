package pgstore

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

func TestInterface(t *testing.T) {
	var st otp.Store = &OTPPgStore{}
	_ = st
}

func TestGetTOTPInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// check if there's no any configuration for the user
	mock.ExpectQuery("SELECT utc.user_id.*").WithArgs("alice").WillReturnRows(sqlmock.NewRows([]string{"user_id", "is_enabled", "secret", "last_time_slot"}))

	// there is config for the user
	mock.ExpectQuery("SELECT utc.user_id.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "is_enabled", "secret", "last_time_slot"}).
			AddRow(123, true, "secret", 52000000),
	)
	mock.ExpectQuery("SELECT id, code.*").WithArgs(123).WillReturnRows(
		sqlmock.NewRows([]string{"id", "code"}).
			AddRow(12, "code_a").AddRow(13, "code_b").AddRow(14, "code_c"),
	)

	st := New(db)
	ctx := context.Background()
	ti, err := st.GetTOTPInfo(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if ti.Secret != "" || ti.Enabled || len(ti.RecoveryCodes) > 0 {
		t.Fatalf("alice is not as expected: %#v", ti)
	}

	ti, err = st.GetTOTPInfo(ctx, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if ti.Secret != "secret" || ti.LastTimeSlot != 52000000 || !ti.Enabled ||
		len(ti.RecoveryCodes) != 3 || ti.RecoveryCodes[13] != "code_b" {
		t.Fatalf("bob is not as expected: %#v", ti)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEnable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	// user without totp configuration
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("alice").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(21),
	)
	mock.ExpectExec("UPDATE totp_configuration SET is_enabled = true.*").WithArgs(21).WillReturnResult(sqlmock.NewResult(0, 0))

	// non-existing user
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id"}),
	)

	// successful path
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(42),
	)
	mock.ExpectExec("UPDATE totp_configuration SET is_enabled = true.*").WithArgs(42).WillReturnResult(sqlmock.NewResult(0, 1))

	st := New(db)
	ctx := context.Background()
	if err := st.Enable(ctx, "alice"); err == nil {
		t.Fatal("enabled totp for alice")
	}
	if err := st.Enable(ctx, "mallory"); err == nil {
		t.Fatal("enabled totp for mallory")
	}
	if err := st.Enable(ctx, "bob"); err != nil {
		t.Fatalf("couldn't enable totp for bob: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// user without totp configuration
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("alice").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(21),
	)
	mock.ExpectExec("DELETE FROM totp_configuration.*").WithArgs(21).WillReturnResult(sqlmock.NewResult(0, 0))

	// non-existing user
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id"}),
	)

	// user with configuration
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(42),
	)
	mock.ExpectExec("DELETE FROM totp_configuration.*").WithArgs(42).WillReturnResult(sqlmock.NewResult(0, 1))

	st := New(db)
	ctx := context.Background()
	if err := st.Delete(ctx, "alice"); err != nil {
		t.Fatalf("couldn't delete totp for alice: %v", err)
	}
	if err := st.Delete(ctx, "mallory"); err == nil {
		t.Fatal("deleted totp for mallory")
	}
	if err := st.Delete(ctx, "bob"); err != nil {
		t.Fatalf("couldn't delete totp for bob: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestStoreNewSecret(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// user without totp configuration
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.is_enabled.*").WithArgs("alice").WillReturnRows(
		sqlmock.NewRows([]string{"id", "is_enabled"}).AddRow(21, nil),
	)
	mock.ExpectExec("INSERT INTO totp_configuration.*").WithArgs(21, "secret_a").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// non-existent user
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.is_enabled.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id", "is_enabled"}),
	)
	mock.ExpectRollback()

	// user with configuration, disabled
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.is_enabled.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id", "is_enabled"}).AddRow(42, false),
	)
	mock.ExpectExec("UPDATE totp_configuration.*").WithArgs("secret_b", 42).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// user with configuration, enabled already
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.is_enabled.*").WithArgs("carol").WillReturnRows(
		sqlmock.NewRows([]string{"id", "is_enabled"}).AddRow(52, true),
	)
	mock.ExpectRollback()

	st := New(db)
	ctx := context.Background()
	if err := st.StoreNewSecret(ctx, "alice", "secret_a"); err != nil {
		t.Fatalf("couldn't store the secret for alice: %v", err)
	}
	if err := st.StoreNewSecret(ctx, "mallory", "secret_m"); err == nil {
		t.Fatal("stored the secret for mallory")
	}
	if err := st.StoreNewSecret(ctx, "bob", "secret_b"); err != nil {
		t.Fatalf("couldn't store the secret for bob: %v", err)
	}
	if err := st.StoreNewSecret(ctx, "carol", "secret_c"); err == nil {
		t.Fatal("stored the secret for carol")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteRecoveryCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// non-existing user
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id"}),
	)

	// existing user
	mock.ExpectQuery("SELECT id FROM \"user\" WHERE username.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(42),
	)
	mock.ExpectExec("DELETE FROM totp_recovery_codes.*").WithArgs(42, 321).WillReturnResult(sqlmock.NewResult(0, 1))

	st := New(db)
	ctx := context.Background()
	if err := st.DeleteRecoveryCode(ctx, "mallory", 432); err == nil {
		t.Fatal("deleted recovery code for mallory")
	}
	if err := st.DeleteRecoveryCode(ctx, "bob", 321); err != nil {
		t.Fatalf("couldn't delete recovery code for bob: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAddRecoveryCodes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// user w/o totp configuration
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("alice").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(21, nil),
	)
	mock.ExpectRollback()

	// non-existing user
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}),
	)
	mock.ExpectRollback()

	// existing user with 3 recovery codes
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(42, 52000000),
	)
	mock.ExpectQuery("SELECT count... FROM totp_recovery_codes.*").WithArgs(42).WillReturnRows(
		sqlmock.NewRows([]string{"count"}).AddRow(3),
	)
	prep := mock.ExpectPrepare("INSERT INTO totp_recovery_codes.*")
	prep.ExpectExec().WithArgs(42, "bob-1").WillReturnResult(sqlmock.NewResult(31, 1))
	prep.ExpectExec().WithArgs(42, "bob-2").WillReturnResult(sqlmock.NewResult(32, 1))
	prep.ExpectExec().WithArgs(42, "bob-3").WillReturnResult(sqlmock.NewResult(33, 1))
	mock.ExpectCommit()

	// existing user with 10 recovery codes
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("carol").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(52, 52000000),
	)
	mock.ExpectQuery("SELECT count... FROM totp_recovery_codes.*").WithArgs(52).WillReturnRows(
		sqlmock.NewRows([]string{"count"}).AddRow(10),
	)
	mock.ExpectRollback()

	st := New(db)
	ctx := context.Background()
	if err := st.AddRecoveryCodes(ctx, "alice", []string{"alice-1", "alice-2"}); err == nil {
		t.Fatal("added recovery codes for alice")
	}
	if err := st.AddRecoveryCodes(ctx, "mallory", []string{"mallory-1", "mallory-2"}); err == nil {
		t.Fatal("added recovery codes for mallory")
	}
	if err := st.AddRecoveryCodes(ctx, "bob", []string{"bob-1", "bob-2", "bob-3"}); err != nil {
		t.Fatalf("couldn't add recovery codes for bob: %v", err)
	}
	if err := st.AddRecoveryCodes(ctx, "carol", []string{"carol-1", "carol-2", "carol-3"}); err == nil {
		t.Fatal("added recovery codes for carol")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateLastTimeSlot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	// user w/o totp configuration
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("alice").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(21, nil),
	)
	mock.ExpectRollback()

	// non-existing user
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("mallory").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}),
	)
	mock.ExpectRollback()

	// existing user
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("bob").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(42, 52000000),
	)
	mock.ExpectExec("UPDATE totp_configuration SET last_time_slot.*").WithArgs(53000000, 42).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// existing user with changed last_time_slot
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT u.id, utc.last_time_slot.*").WithArgs("carol").WillReturnRows(
		sqlmock.NewRows([]string{"id", "last_ts"}).AddRow(52, 52100000),
	)
	mock.ExpectRollback()

	st := New(db)
	ctx := context.Background()
	if err := st.UpdateLastTimeSlot(ctx, "alice", 52000000, 53000000); err == nil {
		t.Fatal("updated last ts for alice")
	}
	if err := st.UpdateLastTimeSlot(ctx, "mallory", 52000000, 53000000); err == nil {
		t.Fatal("updated last ts for mallory")
	}
	if err := st.UpdateLastTimeSlot(ctx, "bob", 52000000, 53000000); err != nil {
		t.Fatalf("couldn't update last ts for bob: %v", err)
	}
	if err := st.UpdateLastTimeSlot(ctx, "carol", 52000000, 53000000); err == nil {
		t.Fatal("updated last ts for carol")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
