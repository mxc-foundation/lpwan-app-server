package storage

import (
	"context"

	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"

	mcss "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// RemoteMulticastClassCSession defines a remote multicast-setup Class-C session record.
type RemoteMulticastClassCSession mcss.RemoteMulticastClassCSession

// CreateRemoteMulticastClassCSession creates the given multicast Class-C session.
func CreateRemoteMulticastClassCSession(ctx context.Context, handler *store.Handler, sess *RemoteMulticastClassCSession) error {
	return handler.CreateRemoteMulticastClassCSession(ctx, (*mcss.RemoteMulticastClassCSession)(sess))
}

// GetRemoteMulticastClassCSession returns the multicast Class-C session given
// a DevEUI and multicast-group ID.
func GetRemoteMulticastClassCSession(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, multicastGroupID uuid.UUID, forUpdate bool) (RemoteMulticastClassCSession, error) {
	res, err := handler.GetRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID, forUpdate)
	return RemoteMulticastClassCSession(res), err
}

// GetRemoteMulticastClassCSessionByGroupID returns the multicast Class-C session given
// a DevEUI and McGroupID.
func GetRemoteMulticastClassCSessionByGroupID(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, mcGroupID int, forUpdate bool) (RemoteMulticastClassCSession, error) {
	res, err := handler.GetRemoteMulticastClassCSessionByGroupID(ctx, devEUI, mcGroupID, forUpdate)
	return RemoteMulticastClassCSession(res), err
}

// GetPendingRemoteMulticastClassCSessions returns a slice of pending remote
// multicast Class-C sessions.
func GetPendingRemoteMulticastClassCSessions(ctx context.Context, handler *store.Handler, limit, maxRetryCount int) ([]RemoteMulticastClassCSession, error) {
	res, err := handler.GetPendingRemoteMulticastClassCSessions(ctx, limit, maxRetryCount)
	if err != nil {
		return nil, err
	}

	var resultList []RemoteMulticastClassCSession
	for _, v := range res {
		resultItem := RemoteMulticastClassCSession(v)
		resultList = append(resultList, resultItem)
	}

	return resultList, nil
}

// UpdateRemoteMulticastClassCSession updates the given remote multicast
// Class-C session.
func UpdateRemoteMulticastClassCSession(ctx context.Context, handler *store.Handler, sess *RemoteMulticastClassCSession) error {
	return handler.UpdateRemoteMulticastClassCSession(ctx, (*mcss.RemoteMulticastClassCSession)(sess))
}

// DeleteRemoteMulticastClassCSession deletes the multicast Class-C session
// given a DevEUI and multicast-group ID.
func DeleteRemoteMulticastClassCSession(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, multicastGroupID uuid.UUID) error {
	return handler.DeleteRemoteMulticastClassCSession(ctx, devEUI, multicastGroupID)
}
