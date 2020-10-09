package store

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	"github.com/brocaar/lorawan"
)

func NewStore(pg pgstore.PgStore) *mcs {
	return &mcs{
		pg: pg,
	}
}

type mcs struct {
	pg pgstore.MigrateCodePgStore
}

type MigrateCodeStore interface {
	Migrate(ctx context.Context, name string) error
	GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error)
}

func (h *mcs) Migrate(ctx context.Context, name string) error {
	return h.pg.Migrate(ctx, name)
}

func (h *mcs) GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error) {
	return h.pg.GetAllGatewayIDs(ctx)
}
