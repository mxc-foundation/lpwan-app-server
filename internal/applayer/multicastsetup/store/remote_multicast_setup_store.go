package store

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	. "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewSetupStore(pg pgstore.PgStore) *mcss {
	return &mcss{
		pg: pg,
	}
}

type mcss struct {
	pg pgstore.MulticastSetupPgStore
}

type RemoteMulticastSetupStore interface {
	CreateRemoteMulticastSetup(ctx context.Context, dms *RemoteMulticastSetup) error
	GetRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastSetup, error)
	GetRemoteMulticastSetupByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastSetup, error)
	GetPendingRemoteMulticastSetupItems(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastSetup, error)
	UpdateRemoteMulticastSetup(ctx context.Context, dmg *RemoteMulticastSetup) error
	DeleteRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error
	GetDevEUIsWithMulticastSetup(ctx context.Context, id *uuid.UUID) ([]lorawan.EUI64, error)
	GetDevEUIsWithFragmentationSessionSetup(ctx context.Context, id *uuid.UUID, fragIdx int) ([]lorawan.EUI64, error)
}

func (h *mcss) GetDevEUIsWithFragmentationSessionSetup(ctx context.Context, id *uuid.UUID, fragIdx int) ([]lorawan.EUI64, error) {
	return h.pg.GetDevEUIsWithFragmentationSessionSetup(ctx, id, fragIdx)
}

func (h *mcss) GetDevEUIsWithMulticastSetup(ctx context.Context, id *uuid.UUID) ([]lorawan.EUI64, error) {
	return h.pg.GetDevEUIsWithMulticastSetup(ctx, id)
}
func (h *mcss) CreateRemoteMulticastSetup(ctx context.Context, dms *RemoteMulticastSetup) error {
	return h.pg.CreateRemoteMulticastSetup(ctx, dms)
}
func (h *mcss) GetRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastSetup, error) {
	return h.pg.GetRemoteMulticastSetup(ctx, devEUI, multicastGroupID, forUpdate)
}
func (h *mcss) GetRemoteMulticastSetupByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastSetup, error) {
	return h.pg.GetRemoteMulticastSetupByGroupID(ctx, devEUI, mcGroupID, forUpdate)
}
func (h *mcss) GetPendingRemoteMulticastSetupItems(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastSetup, error) {
	return h.pg.GetPendingRemoteMulticastSetupItems(ctx, limit, maxRetryCount)
}
func (h *mcss) UpdateRemoteMulticastSetup(ctx context.Context, dmg *RemoteMulticastSetup) error {
	return h.pg.UpdateRemoteMulticastSetup(ctx, dmg)
}
func (h *mcss) DeleteRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return h.pg.DeleteRemoteMulticastSetup(ctx, devEUI, multicastGroupID)
}
