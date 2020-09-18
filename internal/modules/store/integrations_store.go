package store

import (
	"context"
	"encoding/json"
	"time"
)

type IntegrationsStore interface {
	CreateIntegration(ctx context.Context, i *Integration) error
	GetIntegration(ctx context.Context, id int64) (Integration, error)
	GetIntegrationByApplicationID(ctx context.Context, applicationID int64, kind string) (Integration, error)
	GetIntegrationsForApplicationID(ctx context.Context, applicationID int64) ([]Integration, error)
	UpdateIntegration(ctx context.Context, i *Integration) error
	DeleteIntegration(ctx context.Context, id int64) error
}

func (h *Handler) CreateIntegration(ctx context.Context, i *Integration) error {
	return h.store.CreateIntegration(ctx, i)
}

func (h *Handler) GetIntegration(ctx context.Context, id int64) (Integration, error) {
	return h.store.GetIntegration(ctx, id)
}

func (h *Handler) GetIntegrationByApplicationID(ctx context.Context, applicationID int64, kind string) (Integration, error) {
	return h.store.GetIntegrationByApplicationID(ctx, applicationID, kind)
}

func (h *Handler) GetIntegrationsForApplicationID(ctx context.Context, applicationID int64) ([]Integration, error) {
	return h.store.GetIntegrationsForApplicationID(ctx, applicationID)
}

func (h *Handler) UpdateIntegration(ctx context.Context, i *Integration) error {
	return h.store.UpdateIntegration(ctx, i)
}

func (h *Handler) DeleteIntegration(ctx context.Context, id int64) error {
	return h.store.DeleteIntegration(ctx, id)
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
