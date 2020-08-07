package miningstore

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (st *Store) GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error) {
	return storage.GetGatewayMiningList(ctx, st.db, time, limit)
}

func (st *Store) UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error {
	return storage.UpdateFirstHeartbeatToZero(ctx, st.db, mac)
}
