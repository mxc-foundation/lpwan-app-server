package pgstore

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	dbx := sqlx.NewDb(db, "postgresql")

	query := `SELECT id, is_admin FROM "user".*`
	cols := []string{"id", "is_admin"}
	// no user alice
	mock.ExpectQuery(query).WithArgs("alice").WillReturnRows(sqlmock.NewRows(cols))
	// there is admin bob
	mock.ExpectQuery(query).WithArgs("bob").WillReturnRows(
		sqlmock.NewRows(cols).AddRow(42, true))
	// there is user carol
	mock.ExpectQuery(query).WithArgs("carol").WillReturnRows(
		sqlmock.NewRows(cols).AddRow(43, false))

	s := New(dbx)
	ctx := context.Background()
	if alice, err := s.GetUser(ctx, "alice"); err == nil {
		t.Errorf("got success for alice: %#v", alice)
	}

	if bob, err := s.GetUser(ctx, "bob"); err != nil || bob.ID != 42 || !bob.IsGlobalAdmin {
		t.Errorf("failed for bob: %v: %#v", err, bob)
	}

	if carol, err := s.GetUser(ctx, "carol"); err != nil || carol.ID != 43 || carol.IsGlobalAdmin {
		t.Errorf("failed for carol: %v: %#v", err, carol)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectation failed: %v", err)
	}
}

func TestGetOrgUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	dbx := sqlx.NewDb(db, "postgresql")

	query := `SELECT is_admin, is_device_admin, is_gateway_admin.*`
	cols := []string{"is_admin", "is_dev", "is_gw"}
	// user is not in org
	mock.ExpectQuery(query).WithArgs(42, 23).WillReturnRows(sqlmock.NewRows(cols))
	// user is admin
	mock.ExpectQuery(query).WithArgs(43, 43).WillReturnRows(
		sqlmock.NewRows(cols).AddRow(true, true, true))
	// user is device admin
	mock.ExpectQuery(query).WithArgs(43, 100).WillReturnRows(
		sqlmock.NewRows(cols).AddRow(false, true, false))
	// user is not admin
	mock.ExpectQuery(query).WithArgs(44, 100).WillReturnRows(
		sqlmock.NewRows(cols).AddRow(false, false, false))

	s := New(dbx)
	ctx := context.Background()
	if ou, err := s.GetOrgUser(ctx, 42, 23); err == nil {
		t.Errorf("expected an error, got %#v", ou)
	}
	if ou, err := s.GetOrgUser(ctx, 43, 43); err != nil || !ou.IsOrgAdmin || !ou.IsDeviceAdmin || !ou.IsGatewayAdmin {
		t.Errorf("43,43 is not as expected: %v: %#v", err, ou)
	}
	if ou, err := s.GetOrgUser(ctx, 43, 100); err != nil || ou.IsOrgAdmin || !ou.IsDeviceAdmin || ou.IsGatewayAdmin {
		t.Errorf("43,100 is not as expected: %v: %#v", err, ou)
	}
	if ou, err := s.GetOrgUser(ctx, 44, 100); err != nil || ou.IsOrgAdmin || ou.IsDeviceAdmin || ou.IsGatewayAdmin {
		t.Errorf("44,100 is not as expected: %v: %#v", err, ou)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectation failed: %v", err)
	}
}
