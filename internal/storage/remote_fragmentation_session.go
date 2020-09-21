package storage

import (
	"context"

	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// RemoteFragmentationSession defines a remote fragmentation session record.
type RemoteFragmentationSession store.RemoteFragmentationSession

// CreateRemoteFragmentationSession creates the given fragmentation session.
func CreateRemoteFragmentationSession(ctx context.Context, handler *store.Handler, sess *RemoteFragmentationSession) error {
	return handler.CreateRemoteFragmentationSession(ctx, (*store.RemoteFragmentationSession)(sess))
}

// GetRemoteFragmentationSession returns the fragmentation session given a
// DevEUI and fragmentation index.
func GetRemoteFragmentationSession(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, fragIndex int, forUpdate bool) (RemoteFragmentationSession, error) {
	res, err := handler.GetRemoteFragmentationSession(ctx, devEUI, fragIndex, forUpdate)
	return RemoteFragmentationSession(res), err
}

// GetPendingRemoteFragmentationSessions returns a slice of pending remote
// fragmentation sessions.
func GetPendingRemoteFragmentationSessions(ctx context.Context, handler *store.Handler, limit, maxRetryCount int) ([]RemoteFragmentationSession, error) {
	res, err := handler.GetPendingRemoteFragmentationSessions(ctx, limit, maxRetryCount)
	if err != nil {
		return nil, err
	}

	var rfsList []RemoteFragmentationSession
	for _, v := range res {
		rfsItem := RemoteFragmentationSession(v)
		rfsList = append(rfsList, rfsItem)
	}
	return rfsList, nil
}

// UpdateRemoteFragmentationSession updates the given fragmentation session.
func UpdateRemoteFragmentationSession(ctx context.Context, handler *store.Handler, sess *RemoteFragmentationSession) error {
	return handler.UpdateRemoteFragmentationSession(ctx, (*store.RemoteFragmentationSession)(sess))
}

// DeleteRemoteFragmentationSession removes the fragmentation session for the
// given DevEUI / fragmentation index combination.
func DeleteRemoteFragmentationSession(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, fragIndex int) error {
	return handler.DeleteRemoteFragmentationSession(ctx, devEUI, fragIndex)
}
