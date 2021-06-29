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

type PgStore struct {
	db   DB
	txDB txDB
}

// New returns a new database access layer for other stores
func New() *PgStore {
	return &PgStore{
		db: ctrl.db,
	}
}

// Tx starts transaction and executes the function passing to it Handler
// using this transaction. It automatically rolls the transaction back if
// function returns an error. If the error has been caused by serialization
// error, it calls the function again. In order for serialization errors
// handling to work, the function should return Handler errors
// unchanged, or wrap them using %w.
func (ps *PgStore) Tx(ctx context.Context, f func(context.Context, interface{}) error) error {
	for {
		pst, err := ps.TxBegin(ctx)
		if err != nil {
			return err
		}
		err = f(ctx, pst)
		if err == nil {
			if err = pst.TxCommit(ctx); err == nil {
				return nil
			}
		}
		_ = pst.TxRollback(ctx)
		if ps.IsErrorRepeat(err) {
			continue
		}
		return err
	}
}

// InTx returns true if the object is in transaction
func (ps *PgStore) InTx() bool {
	return ps.txDB != nil
}

func (ps *PgStore) TxBegin(ctx context.Context) (*PgStore, error) {
	tx, err := ps.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, err
	}

	return &PgStore{
		db:   ltx{Tx: tx},
		txDB: tx,
	}, nil
}

func (ps *PgStore) TxCommit(ctx context.Context) error {
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

func (ps *PgStore) TxRollback(ctx context.Context) error {
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

func (ps *PgStore) IsErrorRepeat(err error) bool {
	var e pq.Error
	if errors.As(err, &e) {
		if e.Code == "40001" {
			return true
		}
	}
	return false
}
