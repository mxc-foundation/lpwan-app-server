package multicast

import (
	"context"
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	api "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Enqueue adds the given payload to the multicast-group queue.
func Enqueue(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, fPort uint8,
	data []byte, nsCli *nscli.Client) (uint32, error) {
	fCnts, err := EnqueueMultiple(ctx, handler, multicastGroupID, fPort, [][]byte{data}, nsCli)
	if err != nil {
		return 0, err
	}

	if len(fCnts) != 1 {
		return 0, fmt.Errorf("expected 1 frame-counter, got: %d", len(fCnts))
	}

	return fCnts[0], nil
}

// EnqueueMultiple adds the given payloads to the multicast-group queue.
func EnqueueMultiple(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, fPort uint8,
	payloads [][]byte, nsCli *nscli.Client) ([]uint32, error) {
	// Get and lock multicast-group, the lock is to make sure there are no
	// concurrent enqueue actions for the same multicast-group, which would
	// result in the re-use of the same frame-counter.
	mg, err := GetMulticastGroup(ctx, multicastGroupID, true, false, nsCli)
	if err != nil {
		return nil, errors.Wrap(err, "get multicast-group error")
	}

	// get network-server / client
	n, err := handler.GetNetworkServerForMulticastGroupID(ctx, multicastGroupID)
	if err != nil {
		return nil, errors.Wrap(err, "get network-server error")
	}
	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, errors.Wrap(err, "get network-server client error")
	}

	var out []uint32
	var devAddr lorawan.DevAddr
	copy(devAddr[:], mg.MulticastGroup.McAddr)
	fCnt := mg.MulticastGroup.FCnt

	for _, pl := range payloads {
		// encrypt payload
		b, err := lorawan.EncryptFRMPayload(mg.MCAppSKey, false, devAddr, fCnt, pl)
		if err != nil {
			return nil, errors.Wrap(err, "encrypt frmpayload error")
		}

		_, err = nsClient.EnqueueMulticastQueueItem(ctx, &ns.EnqueueMulticastQueueItemRequest{
			MulticastQueueItem: &ns.MulticastQueueItem{
				MulticastGroupId: multicastGroupID.Bytes(),
				FrmPayload:       b,
				FCnt:             fCnt,
				FPort:            uint32(fPort),
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "enqueue multicast-queue item error")
		}

		out = append(out, fCnt)
		fCnt++
	}

	return out, nil
}

// ListQueue lists the items in the multicast-group queue.
func ListQueue(ctx context.Context, handler *store.Handler, multicastGroupID uuid.UUID, nsCli *nscli.Client) ([]api.MulticastQueueItem, error) {

	mg, err := GetMulticastGroup(ctx, multicastGroupID, false, false, nsCli)
	if err != nil {
		return nil, errors.Wrap(err, "get multicast-group error")
	}

	n, err := handler.GetNetworkServerForMulticastGroupID(ctx, multicastGroupID)
	if err != nil {
		return nil, errors.Wrap(err, "get network-server for multicast-group error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, errors.Wrap(err, "get network-server client error")
	}

	resp, err := nsClient.GetMulticastQueueItemsForMulticastGroup(ctx, &ns.GetMulticastQueueItemsForMulticastGroupRequest{
		MulticastGroupId: multicastGroupID.Bytes(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "get multicast queue-items error")
	}

	var out []api.MulticastQueueItem
	var devAddr lorawan.DevAddr
	copy(devAddr[:], mg.MulticastGroup.McAddr)

	for _, qi := range resp.MulticastQueueItems {
		b, err := lorawan.EncryptFRMPayload(mg.MCAppSKey, false, devAddr, qi.FCnt, qi.FrmPayload)
		if err != nil {
			return nil, errors.Wrap(err, "decrypt multicast queue-item error")
		}

		out = append(out, api.MulticastQueueItem{
			MulticastGroupId: multicastGroupID.String(),
			FCnt:             qi.FCnt,
			FPort:            qi.FPort,
			Data:             b,
		})
	}

	return out, nil
}
