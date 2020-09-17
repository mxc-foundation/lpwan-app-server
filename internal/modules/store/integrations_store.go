package store

import (
	"context"
	"encoding/json"
	"time"
)

type IntegrationsStore interface {
	CreateIntegration(ctx context.Context, i *Integration) error
}

func (h *Handler) CreateIntegration(ctx context.Context, i *Integration) error {
	return h.store.CreateIntegration(ctx, i)
}

// Integration represents an integration.
type Integration struct {
	ID            int64           `db:"id"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	ApplicationID int64           `db:"application_id"`
	Kind          string          `db:"kind"`
	Settings      json.RawMessage `db:"settings"`
}
