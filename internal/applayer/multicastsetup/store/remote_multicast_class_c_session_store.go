package store

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	. "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewClassCSessionStore(pg pgstore.PgStore) *mcClassCss {
	return &mcClassCss{
		pg: pg,
	}
}

type mcClassCss struct {
	pg pgstore.MulticastClassCSessionPgStore
}

type RemoteMulticastClassCSessionStore interface {
	CreateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error
	GetRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastClassCSession, error)
	GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastClassCSession, error)
	GetPendingRemoteMulticastClassCSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastClassCSession, error)
	UpdateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error
	DeleteRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error
}

func (h *mcClassCss) CreateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error {
	return h.pg.CreateRemoteMulticastClassCSession(ctx, sess)
}
func (h *mcClassCss) GetRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastClassCSession, error) {
	return h.pg.GetRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID, forUpdate)
}
func (h *mcClassCss) GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastClassCSession, error) {
	return h.pg.GetRemoteMulticastClassCSessionByGroupID(ctx, devEUI, mcGroupID, forUpdate)
}
func (h *mcClassCss) GetPendingRemoteMulticastClassCSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastClassCSession, error) {
	return h.pg.GetPendingRemoteMulticastClassCSessions(ctx, limit, maxRetryCount)
}
func (h *mcClassCss) UpdateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error {
	return h.pg.UpdateRemoteMulticastClassCSession(ctx, sess)
}
func (h *mcClassCss) DeleteRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return h.pg.DeleteRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID)
}
