package store

import (
	"context"
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore() *Handler {
	pgstore := pgstore.New()
	return &Handler{
		PgStore: pgstore,
	}
}

type Handler struct {
	pgstore.PgStore
	inTX bool
}

// txBegin creates a transaction and returns a new instance of Handler that
// will either commit or rollback all the changes that done using this
// instance.
func (s *Handler) txBegin(ctx context.Context) (*Handler, error) {
	if s.inTX {
		return nil, fmt.Errorf("already in transaction")
	}
	store, err := s.TxBegin(ctx)
	if err != nil {
		return nil, err
	}
	return &Handler{
		PgStore: store,
		inTX:    true,
	}, nil
}

// Tx starts transaction and executes the function passing to it Handler
// using this transaction. It automatically rolls the transaction back if
// function returns an error. If the error has been caused by serialization
// error, it calls the function again. In order for serialization errors
// handling to work, the function should return Handler errors
// unchanged, or wrap them using %w.
func (h *Handler) Tx(ctx context.Context, f func(context.Context, *Handler) error) error {
	for {
		tx, err := h.txBegin(ctx)
		if err != nil {
			return err
		}
		err = f(ctx, tx)
		if err == nil {
			if err = tx.TxCommit(ctx); err == nil {
				return nil
			}
		}
		_ = tx.TxRollback(ctx)
		if h.IsErrorRepeat(err) {
			// failed due to the serialization error, try again
			continue
		}
		return err
	}
}
