package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func newStore(pg pgstore.PgStore) *st {
	return &st{
		pg: pg,
	}
}

type st struct {
	pg pgstore.PgStore
}

type basicStore interface {
	// TxBegin starts a new transaction and returns a new ApplicationStore instance
	TxBegin(ctx context.Context) (pgstore.PgStore, error)
	// TxCommit commits the transaction, store is not usable after this call
	TxCommit(ctx context.Context) error
	// TxRollback rolls the transaction back, store is not usable after this call
	TxRollback(ctx context.Context) error

	// IsErrorRepeat returns true if the error indicates that the action failed
	// because of a conflict with another transaction and that it should be
	// repeated
	IsErrorRepeat(err error) bool
}

func (s *st) TxBegin(ctx context.Context) (pgstore.PgStore, error) {
	return s.pg.TxBegin(ctx)
}
func (s *st) TxCommit(ctx context.Context) error {
	return s.pg.TxCommit(ctx)
}
func (s *st) TxRollback(ctx context.Context) error {
	return s.pg.TxRollback(ctx)
}
func (s *st) IsErrorRepeat(err error) bool {
	return s.pg.IsErrorRepeat(err)
}
