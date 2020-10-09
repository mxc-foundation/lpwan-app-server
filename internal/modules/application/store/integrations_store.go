package store

import (
	"context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewIntegrationStore(pg pgstore.PgStore) *intgs {
	return &intgs{
		pg: pg,
	}
}

type intgs struct {
	pg pgstore.IntegrationPgStore
}

type IntegrationsStore interface {
	CreateIntegration(ctx context.Context, i *Integration) error
	GetIntegration(ctx context.Context, id int64) (Integration, error)
	GetIntegrationByApplicationID(ctx context.Context, applicationID int64, kind string) (Integration, error)
	GetIntegrationsForApplicationID(ctx context.Context, applicationID int64) ([]Integration, error)
	UpdateIntegration(ctx context.Context, i *Integration) error
	DeleteIntegration(ctx context.Context, id int64) error
}

func (h *intgs) CreateIntegration(ctx context.Context, i *Integration) error {
	return h.pg.CreateIntegration(ctx, i)
}

func (h *intgs) GetIntegration(ctx context.Context, id int64) (Integration, error) {
	return h.pg.GetIntegration(ctx, id)
}

func (h *intgs) GetIntegrationByApplicationID(ctx context.Context, applicationID int64, kind string) (Integration, error) {
	return h.pg.GetIntegrationByApplicationID(ctx, applicationID, kind)
}

func (h *intgs) GetIntegrationsForApplicationID(ctx context.Context, applicationID int64) ([]Integration, error) {
	return h.pg.GetIntegrationsForApplicationID(ctx, applicationID)
}

func (h *intgs) UpdateIntegration(ctx context.Context, i *Integration) error {
	return h.pg.UpdateIntegration(ctx, i)
}

func (h *intgs) DeleteIntegration(ctx context.Context, id int64) error {
	return h.pg.DeleteIntegration(ctx, id)
}
