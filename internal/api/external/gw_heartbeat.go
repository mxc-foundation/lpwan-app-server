package external

import (
	"context"
	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mxc-foundation/lpwan-app-server/api/gw_appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

// HeartbeatAPI exports the Node related functions.
type HeartbeatAPI struct {
	validator auth.Validator
}

// NewHeartbeatAPI creates a new NodeAPI.
func NewHeartbeatAPI() *HeartbeatAPI {
	return &HeartbeatAPI{
	}
}

func (a *HeartbeatAPI) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*empty.Empty, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayMac)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	current_heartbeat := time.Now().Unix()

	// if the gateway is new gateway
	if req.Model == "MX1901" || req.Model == "MX1902" || req.Model == "MX1903" {
		last_heartbeat, err := storage.GetLastHeartbeat(ctx, storage.DB(), mac)
		if err != nil {
			// if last heartbeat is nil, update first and last heartbeat
			if err == storage.ErrDoesNotExist {
				err := storage.UpdateLastHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
				if err != nil {
					log.WithError(err).Error("Update last heartbeat error")
					return nil, err
				}

				err = storage.UpdateFirstHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
				if err != nil {
					log.WithError(err).Error("Update first heartbeat error")
					return nil, err
				}
			}
			log.WithError(err).Error("Cannot get last heartbeat from DB.")

			return nil, err
		}

		// if offline longer than 30 mins, last heartbeat and first heartbeat = current heartbeat
		if current_heartbeat-last_heartbeat > 1800 {
			err := storage.UpdateLastHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
			if err != nil {
				log.WithError(err).Error("Update last heartbeat error")
				return &empty.Empty{}, err
			}

			err = storage.UpdateFirstHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
			if err != nil {
				log.WithError(err).Error("Update first heartbeat error")
				return &empty.Empty{}, err
			}
		}

		err = storage.UpdateLastHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
		if err != nil {
			log.WithError(err).Error("Update last heartbeat error")
			return &empty.Empty{}, err
		}
	}

	return &empty.Empty{}, nil
}
