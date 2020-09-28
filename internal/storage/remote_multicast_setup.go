package storage

import (
	"context"

	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// Possible states
const (
	RemoteMulticastSetupSetup  = store.RemoteMulticastSetupSetup
	RemoteMulticastSetupDelete = store.RemoteMulticastSetupDelete
)

// RemoteMulticastSetup defines a remote multicast-setup record.
type RemoteMulticastSetup store.RemoteMulticastSetup

// CreateRemoteMulticastSetup creates the given multicast-setup.
func CreateRemoteMulticastSetup(ctx context.Context, handler *store.Handler, dms *RemoteMulticastSetup) error {
	return handler.CreateRemoteMulticastSetup(ctx, (*store.RemoteMulticastSetup)(dms))
}

// GetRemoteMulticastSetup returns the multicast-setup given a multicast-group ID and DevEUI.
func GetRemoteMulticastSetup(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastSetup, error) {
	res, err := handler.GetRemoteMulticastSetup(ctx, devEUI, multicastGroupID, forUpdate)
	return RemoteMulticastSetup(res), err
}

// GetRemoteMulticastSetupByGroupID returns the multicast-setup given a DevEUI and McGroupID.
func GetRemoteMulticastSetupByGroupID(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastSetup, error) {
	res, err := handler.GetRemoteMulticastSetupByGroupID(ctx, devEUI, mcGroupID, forUpdate)
	return RemoteMulticastSetup(res), err
}

// GetPendingRemoteMulticastSetupItems returns a slice of pending remote multicast-setup items.
// The selected items will be locked.
func GetPendingRemoteMulticastSetupItems(ctx context.Context, handler *store.Handler, limit, maxRetryCount int) ([]RemoteMulticastSetup, error) {
	res, err := handler.GetPendingRemoteMulticastSetupItems(ctx, limit, maxRetryCount)
	if err != nil {
		return nil, err
	}

	var rmsList []RemoteMulticastSetup
	for _, v := range res {
		rmsItem := RemoteMulticastSetup(v)
		rmsList = append(rmsList, rmsItem)
	}

	return rmsList, nil
}

// UpdateRemoteMulticastSetup updates the given update multicast-group setup.
func UpdateRemoteMulticastSetup(ctx context.Context, handler *store.Handler, dmg *RemoteMulticastSetup) error {
	return handler.UpdateRemoteMulticastSetup(ctx, (*store.RemoteMulticastSetup)(dmg))
}

// DeleteRemoteMulticastSetup deletes the multicast-setup given a multicast-group ID and DevEUI.
func DeleteRemoteMulticastSetup(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return handler.DeleteRemoteMulticastSetup(ctx, devEUI, multicastGroupID)
}
