package pgstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DB represents database interface.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)

	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type ltx struct {
	*sqlx.Tx
}

func (t ltx) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return nil, fmt.Errorf("already in transaction")
}

type txDB interface {
	Commit() error
	Rollback() error
}

type pgstore struct {
	db   DB
	txDB txDB
}

type PgStore interface {
	TxBegin(ctx context.Context) (PgStore, error)
	TxCommit(ctx context.Context) error
	TxRollback(ctx context.Context) error
	IsErrorRepeat(err error) bool

	DevicePgstore
	DeviceProfilePgStore
	ApplicationPgStore
	IntegrationPgStore
	GatewayPgStore
	GatewayDefaultConfigPgStore
	OrganizationPgStore
	GatewayProfilePgStore
	FuotaDeploymentPgStore
	MigrateCodePgStore
	MulticastClassCSessionPgStore
	MulticastGroupPgStore
	MulticastSetupPgStore
	ServiceProfilePgStore
	NetworkServerPgStore
	UserPgStore
	FragmentationSessionPgStore
	OTPPgStore
	AuthenticationPgStore
}

// New returns a new database access layer for other stores
func New() PgStore {
	return &pgstore{
		db: ctrl.db,
	}
}

func (ps *pgstore) TxBegin(ctx context.Context) (PgStore, error) {
	tx, err := ps.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, err
	}

	return &pgstore{
		db:   ltx{Tx: tx},
		txDB: tx,
	}, nil
}

func (ps *pgstore) TxCommit(ctx context.Context) error {
	if ps.txDB == nil {
		return fmt.Errorf("not in transaction")
	}
	err := ps.txDB.Commit()
	if err == nil {
		ps.db = nil
		ps.txDB = nil
	}
	return err
}

func (ps *pgstore) TxRollback(ctx context.Context) error {
	if ps.txDB == nil {
		return fmt.Errorf("not in transaction")
	}
	err := ps.txDB.Rollback()
	if err == nil {
		ps.db = nil
		ps.txDB = nil
	}
	return err
}

func (ps *pgstore) IsErrorRepeat(err error) bool {
	var e pq.Error
	if errors.As(err, &e) {
		if e.Code == "40001" {
			return true
		}
	}
	return false
}
