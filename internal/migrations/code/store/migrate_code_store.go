package store

import (
	"context"
	migrate "github.com/rubenv/sql-migrate"

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
	ExecuteMigrateUp(m migrate.MigrationSource) error
	GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error)
	GetAllFromGorpMigrations(ctx context.Context) ([]string, error)
	FixGorpMigrationsItemId(ctx context.Context, oldID, newID string) error
}

func (h *mcs) Migrate(ctx context.Context, name string) error {
	return h.pg.Migrate(ctx, name)
}

func (h *mcs) GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error) {
	return h.pg.GetAllGatewayIDs(ctx)
}

func (h *mcs) GetAllFromGorpMigrations(ctx context.Context) ([]string, error) {
	return h.pg.GetAllFromGorpMigrations(ctx)
}

func (h *mcs) FixGorpMigrationsItemId(ctx context.Context, oldID, newID string) error {
	return h.pg.FixGorpMigrationsItemId(ctx, oldID, newID)
}

func (h *mcs) ExecuteMigrateUp(m migrate.MigrationSource) error {
	return h.pg.ExecuteMigrateUp(m)
}
