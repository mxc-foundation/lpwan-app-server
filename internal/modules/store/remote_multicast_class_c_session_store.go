package store

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
)

type RemoteMulticastClassCSessionStore interface {
	CreateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error
	GetRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastClassCSession, error)
	GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastClassCSession, error)
	GetPendingRemoteMulticastClassCSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastClassCSession, error)
	UpdateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error
	DeleteRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error
}

func (h *Handler) CreateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error {
	return h.store.CreateRemoteMulticastClassCSession(ctx, sess)
}
func (h *Handler) GetRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastClassCSession, error) {
	return h.store.GetRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID, forUpdate)
}
func (h *Handler) GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastClassCSession, error) {
	return h.store.GetRemoteMulticastClassCSessionByGroupID(ctx, devEUI, mcGroupID, forUpdate)
}
func (h *Handler) GetPendingRemoteMulticastClassCSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastClassCSession, error) {
	return h.store.GetPendingRemoteMulticastClassCSessions(ctx, limit, maxRetryCount)
}
func (h *Handler) UpdateRemoteMulticastClassCSession(ctx context.Context, sess *RemoteMulticastClassCSession) error {
	return h.store.UpdateRemoteMulticastClassCSession(ctx, sess)
}
func (h *Handler) DeleteRemoteMulticastClassCSession(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return h.store.DeleteRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID)
}

// RemoteMulticastClassCSession defines a remote multicast-setup Class-C session record.
type RemoteMulticastClassCSession struct {
	DevEUI           lorawan.EUI64 `db:"dev_eui"`
	MulticastGroupID uuid.UUID     `db:"multicast_group_id"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
	McGroupID        int           `db:"mc_group_id"`
	SessionTime      time.Time     `db:"session_time"`
	SessionTimeOut   int           `db:"session_time_out"`
	DLFrequency      int           `db:"dl_frequency"`
	DR               int           `db:"dr"`
	StateProvisioned bool          `db:"state_provisioned"`
	RetryAfter       time.Time     `db:"retry_after"`
	RetryCount       int           `db:"retry_count"`
	RetryInterval    time.Duration `db:"retry_interval"`
}
