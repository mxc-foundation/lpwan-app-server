package device

import (
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/codec"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// DeviceQueueAPI exposes the downlink queue methods.
type DeviceQueueAPI struct {
	st *store.Handler
}

// NewDeviceQueueAPI creates a new DeviceQueueAPI.
func NewDeviceQueueAPI() *DeviceQueueAPI {
	return &DeviceQueueAPI{
		st: Service.St,
	}
}

// Enqueue adds the given item to the device-queue.
func (d *DeviceQueueAPI) Enqueue(ctx context.Context, req *pb.EnqueueDeviceQueueItemRequest) (*pb.EnqueueDeviceQueueItemResponse, error) {
	var fCnt uint32

	if req.DeviceQueueItem == nil {
		return nil, status.Errorf(codes.InvalidArgument, "queue_item must not be nil")
	}

	if req.DeviceQueueItem.FPort == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "f_port must be > 0")
	}

	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DeviceQueueItem.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceQueueAccess(ctx, devEUI, authcus.Create); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := d.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		// Lock the device to avoid concurrent enqueue actions for the same
		// device as this would result in re-use of the same frame-counter.
		dev, err := handler.GetDevice(ctx, devEUI, true)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		// if JSON object is set, try to encode it to bytes
		if req.DeviceQueueItem.JsonObject != "" && req.DeviceQueueItem.JsonObject != "null" {
			app, err := handler.GetApplication(ctx, dev.ApplicationID)
			if err != nil {
				return helpers.ErrToRPCError(err)
			}

			dp, err := handler.GetDeviceProfile(ctx, dev.DeviceProfileID, false)
			if err != nil {
				log.WithError(err).WithField("id", dev.DeviceProfileID).Error("get device-profile error")
				return grpc.Errorf(codes.Internal, "get device-profile error: %s", err)
			}

			// TODO: in the next major release, remove this and always use the
			// device-profile codec fields.
			payloadCodec := app.PayloadCodec
			payloadEncoderScript := app.PayloadEncoderScript

			if dp.PayloadCodec != "" {
				payloadCodec = dp.PayloadCodec
				payloadEncoderScript = dp.PayloadEncoderScript
			}

			req.DeviceQueueItem.Data, err = codec.JSONToBinary(codec.Type(payloadCodec), uint8(req.DeviceQueueItem.FPort), dev.Variables, payloadEncoderScript, []byte(req.DeviceQueueItem.JsonObject))
			if err != nil {
				return helpers.ErrToRPCError(err)
			}
		}

		fCnt, err = handler.EnqueueDownlinkPayload(ctx, devEUI, req.DeviceQueueItem.Confirmed, uint8(req.DeviceQueueItem.FPort), req.DeviceQueueItem.Data)
		if err != nil {
			return status.Errorf(codes.Internal, "enqueue downlink payload error: %s", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.EnqueueDeviceQueueItemResponse{
		FCnt: fCnt,
	}, nil
}

// Flush flushes the downlink device-queue.
func (d *DeviceQueueAPI) Flush(ctx context.Context, req *pb.FlushDeviceQueueRequest) (*empty.Empty, error) {
	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceQueueAccess(ctx, devEUI, authcus.Delete); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := storage.GetNetworkServerForDevEUI(ctx, storage.DB(), devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	// add this device to network server
	nStruct := &nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}

	nsClient, err := nStruct.GetNetworkServiceClient()
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	_, err = nsClient.FlushDeviceQueueForDevEUI(ctx, &ns.FlushDeviceQueueForDevEUIRequest{
		DevEui: devEUI[:],
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// List lists the items in the device-queue.
func (d *DeviceQueueAPI) List(ctx context.Context, req *pb.ListDeviceQueueItemsRequest) (*pb.ListDeviceQueueItemsResponse, error) {
	var devEUI lorawan.EUI64
	if err := devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "devEUI: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceQueueAccess(ctx, devEUI, authcus.List); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	device, err := storage.GetDevice(ctx, storage.DB(), devEUI, false, true)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := storage.GetNetworkServerForDevEUI(ctx, storage.DB(), devEUI)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	queueItemsResp, err := nsClient.GetDeviceQueueItemsForDevEUI(ctx, &ns.GetDeviceQueueItemsForDevEUIRequest{
		DevEui:    devEUI[:],
		CountOnly: req.CountOnly,
	})
	if err != nil {
		return nil, err
	}

	resp := pb.ListDeviceQueueItemsResponse{
		TotalCount: queueItemsResp.TotalCount,
	}
	for _, qi := range queueItemsResp.Items {
		b, err := lorawan.EncryptFRMPayload(device.AppSKey, false, device.DevAddr, qi.FCnt, qi.FrmPayload)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.DeviceQueueItems = append(resp.DeviceQueueItems, &pb.DeviceQueueItem{
			DevEui:    devEUI.String(),
			Confirmed: qi.Confirmed,
			FPort:     qi.FPort,
			FCnt:      qi.FCnt,
			Data:      b,
		})
	}

	return &resp, nil
}
