package store

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
)

type RemoteFragmentaionSessionStore interface {
	CreateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error
	GetRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (RemoteFragmentationSession, error)
	GetPendingRemoteFragmentationSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteFragmentationSession, error)
	UpdateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error
	DeleteRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int) error
}

func (h *Handler) CreateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error {
	return h.store.CreateRemoteFragmentationSession(ctx, sess)
}
func (h *Handler) GetRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (RemoteFragmentationSession, error) {
	return h.store.GetRemoteFragmentationSession(ctx, devEUI, fragIndex, forUpdate)
}
func (h *Handler) GetPendingRemoteFragmentationSessions(ctx context.Context, limit, maxRetryCount int) ([]RemoteFragmentationSession, error) {
	return h.store.GetPendingRemoteFragmentationSessions(ctx, limit, maxRetryCount)
}
func (h *Handler) UpdateRemoteFragmentationSession(ctx context.Context, sess *RemoteFragmentationSession) error {
	return h.store.UpdateRemoteFragmentationSession(ctx, sess)
}
func (h *Handler) DeleteRemoteFragmentationSession(ctx context.Context, devEUI lorawan.EUI64, fragIndex int) error {
	return h.store.DeleteRemoteFragmentationSession(ctx, devEUI, fragIndex)
}

// RemoteFragmentationSession defines a remote fragmentation session record.
type RemoteFragmentationSession struct {
	DevEUI              lorawan.EUI64             `db:"dev_eui"`
	FragIndex           int                       `db:"frag_index"`
	CreatedAt           time.Time                 `db:"created_at"`
	UpdatedAt           time.Time                 `db:"updated_at"`
	MCGroupIDs          []int                     `db:"mc_group_ids"`
	NbFrag              int                       `db:"nb_frag"`
	FragSize            int                       `db:"frag_size"`
	FragmentationMatrix uint8                     `db:"fragmentation_matrix"`
	BlockAckDelay       int                       `db:"block_ack_delay"`
	Padding             int                       `db:"padding"`
	Descriptor          [4]byte                   `db:"descriptor"`
	State               RemoteMulticastSetupState `db:"state"`
	StateProvisioned    bool                      `db:"state_provisioned"`
	RetryAfter          time.Time                 `db:"retry_after"`
	RetryCount          int                       `db:"retry_count"`
	RetryInterval       time.Duration             `db:"retry_interval"`
}
