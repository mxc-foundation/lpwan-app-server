package store

import (
	"context"

	"github.com/brocaar/lorawan"

	. "github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

func NewStore(pg pgstore.PgStore) *fss {
	return &fss{
		pg: pg,
	}
}

type fss struct {
	pg pgstore.FragmentationSessionPgStore
}

type RemoteFragmentaionSessionStore interface {
	CreateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error
	GetRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (RemoteFragmentationSession, error)
	GetPendingRemoteFragmentationSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteFragmentationSession, error)
	UpdateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error
	DeleteRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int) error
}

func (h *fss) CreateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error {
	return h.pg.CreateRemoteFragmentationSession(ctx, sess)
}
func (h *fss) GetRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (RemoteFragmentationSession, error) {
	return h.pg.GetRemoteFragmentationSession(ctx, devEUI, fragIndex, forUpdate)
}
func (h *fss) GetPendingRemoteFragmentationSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteFragmentationSession, error) {
	return h.pg.GetPendingRemoteFragmentationSessions(ctx, limit, maxRetryCount)
}
func (h *fss) UpdateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error {
	return h.pg.UpdateRemoteFragmentationSession(ctx, sess)
}
func (h *fss) DeleteRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int) error {
	return h.pg.DeleteRemoteFragmentationSession(ctx, devEUI, fragIndex)
}
