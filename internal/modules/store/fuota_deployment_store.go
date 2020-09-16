package store

import (
	"context"

	"github.com/brocaar/lorawan"

	"github.com/gofrs/uuid"
)

type FUOTADeploymentStore interface {
	CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error)
	CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error)
}

func (h *Handler) CheckReadFUOTADeploymentAccess(ctx context.Context, username string, id uuid.UUID, userID int64) (bool, error) {
	return h.store.CheckReadFUOTADeploymentAccess(ctx, username, id, userID)
}

func (h *Handler) CheckCreateFUOTADeploymentsAccess(ctx context.Context, username string, applicationID int64, devEUI lorawan.EUI64, userID int64) (bool, error) {
	return h.store.CheckCreateFUOTADeploymentsAccess(ctx, username, applicationID, devEUI, userID)
}
