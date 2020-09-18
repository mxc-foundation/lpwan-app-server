package store

import (
	"context"
	"github.com/brocaar/lorawan"
)

type MigrateCodeStore interface {
	Migrate(ctx context.Context, name string) error
	GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error)
}

func (h *Handler) Migrate(ctx context.Context, name string) error {
	return h.store.Migrate(ctx, name)
}

func (h *Handler) GetAllGatewayIDs(ctx context.Context) ([]lorawan.EUI64, error) {
	return h.store.GetAllGatewayIDs(ctx)
}
