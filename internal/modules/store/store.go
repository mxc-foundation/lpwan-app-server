package store

import (
	"context"
	"fmt"
)

func New(store Store) (*Handler, error) {
	return &Handler{store: store}, nil
}

type Store interface {
	// TxBegin starts a new transaction and returns a new ApplicationStore instance
	TxBegin(ctx context.Context) (Store, error)
	// TxCommit commits the transaction, store is not usable after this call
	TxCommit(ctx context.Context) error
	// TxRollback rolls the transaction back, store is not usable after this call
	TxRollback(ctx context.Context) error

	// IsErrorRepeat returns true if the error indicates that the action failed
	// because of a conflict with another transaction and that it should be
	// repeated
	IsErrorRepeat(err error) bool

	ApplicationStore
	DeviceStore
	GatewayStore
	GatewayProfileStore
	NetworkServerStore
	OrganizationStore
	UserStore
	ServiceProfileStore
	DeviceProfileStore
	MulticastGroupStore
}

type Handler struct {
	store Store
	inTX  bool
}

// txBegin creates a transaction and returns a new instance of Handler that
// will either commit or rollback all the changes that done using this
// instance.
func (s *Handler) txBegin(ctx context.Context) (*Handler, error) {
	if s.inTX {
		return nil, fmt.Errorf("already in transaction")
	}
	store, err := s.store.TxBegin(ctx)
	if err != nil {
		return nil, err
	}
	btx := *s
	btx.store = store
	btx.inTX = true
	return &btx, nil
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
			if err = tx.store.TxCommit(ctx); err == nil {
				return nil
			}
		}
		_ = tx.store.TxRollback(ctx)
		if h.IsErrorRepeat(err) {
			// failed due to the serialization error, try again
			continue
		}
		return err
	}
}

// IsErrorRepeat returns true if the error indicates that the action has failed
// because of the conflict with another transaction and that the application
// should try to repeat the action
func (h *Handler) IsErrorRepeat(err error) bool {
	return h.store.IsErrorRepeat(err)
}
