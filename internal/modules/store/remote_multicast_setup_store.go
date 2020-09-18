package store

import (
	"context"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
)

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

func (h *Handler) GetDevEUIsWithFragmentationSessionSetup(ctx context.Context, id *uuid.UUID, fragIdx int) ([]lorawan.EUI64, error) {
	return h.store.GetDevEUIsWithFragmentationSessionSetup(ctx, id, fragIdx)
}

func (h *Handler) GetDevEUIsWithMulticastSetup(ctx context.Context, id *uuid.UUID) ([]lorawan.EUI64, error) {
	return h.store.GetDevEUIsWithMulticastSetup(ctx, id)
}
func (h *Handler) CreateRemoteMulticastSetup(ctx context.Context, dms *RemoteMulticastSetup) error {
	return h.store.CreateRemoteMulticastSetup(ctx, dms)
}
func (h *Handler) GetRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastSetup, error) {
	return h.store.GetRemoteMulticastSetup(ctx, devEUI, multicastGroupID, forUpdate)
}
func (h *Handler) GetRemoteMulticastSetupByGroupID(ctx context.Context, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastSetup, error) {
	return h.store.GetRemoteMulticastSetupByGroupID(ctx, devEUI, mcGroupID, forUpdate)
}
func (h *Handler) GetPendingRemoteMulticastSetupItems(ctx context.Context, limit, maxRetryCount int) ([]RemoteMulticastSetup, error) {
	return h.store.GetPendingRemoteMulticastSetupItems(ctx, limit, maxRetryCount)
}
func (h *Handler) UpdateRemoteMulticastSetup(ctx context.Context, dmg *RemoteMulticastSetup) error {
	return h.store.UpdateRemoteMulticastSetup(ctx, dmg)
}
func (h *Handler) DeleteRemoteMulticastSetup(ctx context.Context, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return h.store.DeleteRemoteMulticastSetup(ctx, devEUI, multicastGroupID)
}

// RemoteMulticastSetupState defines the state type.
type RemoteMulticastSetupState string

// Possible states
const (
	RemoteMulticastSetupSetup  RemoteMulticastSetupState = "SETUP"
	RemoteMulticastSetupDelete RemoteMulticastSetupState = "DELETE"
)

// RemoteMulticastSetup defines a remote multicast-setup record.
type RemoteMulticastSetup struct {
	DevEUI           lorawan.EUI64             `db:"dev_eui"`
	MulticastGroupID uuid.UUID                 `db:"multicast_group_id"`
	CreatedAt        time.Time                 `db:"created_at"`
	UpdatedAt        time.Time                 `db:"updated_at"`
	McGroupID        int                       `db:"mc_group_id"`
	McAddr           lorawan.DevAddr           `db:"mc_addr"`
	McKeyEncrypted   lorawan.AES128Key         `db:"mc_key_encrypted"`
	MinMcFCnt        uint32                    `db:"min_mc_f_cnt"`
	MaxMcFCnt        uint32                    `db:"max_mc_f_cnt"`
	State            RemoteMulticastSetupState `db:"state"`
	StateProvisioned bool                      `db:"state_provisioned"`
	RetryInterval    time.Duration             `db:"retry_interval"`
	RetryAfter       time.Time                 `db:"retry_after"`
	RetryCount       int                       `db:"retry_count"`
}
