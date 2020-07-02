package gateway

import (
	"bytes"
	"context"

	/* #nosec */
	"crypto/md5"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/brocaar/lorawan"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	pspb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

// HeartbeatAPI exports the HeartbeatAPI related functions.
type HeartbeatAPI struct {
	BindPort string
	Store    GatewayStore
}

// NewGatewayAPI creates new HeartbeatAPI
func NewHeartbeatAPI(api HeartbeatAPI) *HeartbeatAPI {
	return &HeartbeatAPI{
		BindPort: api.BindPort,
		Store:    api.Store,
	}
}

func (obj *HeartbeatAPI) Heartbeat(ctx context.Context, req *gwpb.HeartbeatRequest) (*gwpb.HeartbeatResponse, error) {
	if obj.BindPort == config.C.ApplicationServer.APIForGateway.OldGateway.Bind {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	response := gwpb.HeartbeatResponse{}

	// check if gateway exists
	var gatewayEUI = lorawan.EUI64{}
	if err := gatewayEUI.UnmarshalText([]byte(req.GatewayMac)); err != nil {
		log.WithError(err).Error("api/Heartbeat: Failed to convert gateway mac address")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid gateway mac format: %s", req.GatewayMac)
	}

	/*	tx, err := obj.Store.
		if err != nil {
			log.WithError(err).Error("Failed to start transaction")
			return nil, status.Errorf(codes.Unknown, "Failed to start transaction: %v", err)
		}
		defer tx.Rollback()*/

	gw, err := obj.Store.GetGateway(ctx, gatewayEUI, true)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			return nil, status.Errorf(codes.Unauthenticated, "Object does not exist: %s", gatewayEUI.String())
		}
		log.WithError(err).Errorf("Failed to select gateway by mac: %s", gatewayEUI.String())
		return nil, status.Errorf(codes.Unknown, "Failed to select gateway by mac: %s", gatewayEUI.String())
	}

	log.Infof("got heartbeat from %s, model %s, last heartbeat %d", gw.MAC.MarshalText, gw.Model, gw.LastHeartbeat)

	// verify gateway model
	if gw.Model != req.Model {
		log.Errorf("Request model does not match saved gateway.")
		return nil, status.Errorf(codes.Unauthenticated, "Request model does not match saved gateway.")
	}

	// important: do this before config and firmware update
	// mining : update heartbeat only for new gateways
	if strings.HasPrefix(gw.Model, "MX19") {
		log.Info("processing MX19 gateway")
		currentHeartbeat := time.Now().Unix()
		lastHeartbeat := gw.LastHeartbeat

		// if last heartbeat == 0 is a new gateway
		if gw.LastHeartbeat == 0 {
			log.Infof("updating heartbeat for the new gw")
			err := obj.Store.UpdateLastHeartbeat(ctx, gatewayEUI, currentHeartbeat)
			if err != nil {
				log.WithError(err).Error("Heartbeat/Update last heartbeat error")
				return nil, status.Errorf(codes.Unimplemented, "Update last heartbeat error")
			}

			err = obj.Store.UpdateFirstHeartbeat(ctx, gatewayEUI, currentHeartbeat)
			if err != nil {
				log.WithError(err).Error("Heartbeat/Update first heartbeat error")
				return nil, status.Errorf(codes.Unimplemented, "Update first heartbeat error")
			}

			goto Next
		}

		// if offline longer than 10 mins, last heartbeat and first heartbeat = current heartbeat
		//if current_heartbeat-last_heartbeat > 600 {
		if currentHeartbeat-lastHeartbeat > config.C.ApplicationServer.MiningSetUp.HeartbeatOfflineLimit {
			err := obj.Store.UpdateLastHeartbeat(ctx, gatewayEUI, currentHeartbeat)
			if err != nil {
				log.WithError(err).Error("Heartbeat/Update last heartbeat error")
				return nil, status.Errorf(codes.Unimplemented, "Update last heartbeat error")
			}

			//err = storage.UpdateFirstHeartbeat(ctx, storage.DB(), mac, current_heartbeat)
			err = obj.Store.UpdateFirstHeartbeatToZero(ctx, gatewayEUI)
			if err != nil {
				log.WithError(err).Error("Heartbeat/Update first heartbeat to zero error")
				return nil, status.Errorf(codes.Unimplemented, "Update first heartbeat to zero error")
			}
			goto Next
		}

		// if first heartbeat != 0 and currentHeartbeat - lastHeart !> 600
		firstHeartbeat, err := obj.Store.GetFirstHeartbeat(ctx, gatewayEUI)
		if err != nil {
			log.WithError(err).Error("Heartbeat/Get first heartbeat error")
			return nil, status.Errorf(codes.DataLoss, "Get firstHeartbeat from DB error")
		}

		if firstHeartbeat == 0 {
			err = obj.Store.UpdateFirstHeartbeat(ctx, gatewayEUI, currentHeartbeat)
			if err != nil {
				log.WithError(err).Error("Heartbeat/Update first heartbeat error")
				return nil, status.Errorf(codes.Unimplemented, "Update first heartbeat error")
			}
		}

		err = obj.Store.UpdateLastHeartbeat(ctx, gatewayEUI, currentHeartbeat)
		if err != nil {
			log.WithError(err).Error("Heartbeat/Update last heartbeat error")
			return nil, status.Errorf(codes.Unimplemented, "Update last heartbeat error")
		}
	}

Next:

	// compare config hash
	/* #nosec */
	configHash := md5.Sum([]byte(gw.Config))
	b := types.MD5SUM{}
	if err := b.UnmarshalText([]byte(req.ConfigHash)); err != nil {
		log.WithError(err).Errorf("Failed to unmarshal config hash: %s", req.ConfigHash)
		return nil, status.Errorf(codes.DataLoss, "Failed to unmarshal config hash: %s", req.ConfigHash)
	}

	if bytes.Equal(configHash[:], b[:]) == false {
		response.Config = gw.Config
	}

	// check if firmware updated
	if gw.AutoUpdateFirmware {
		firmware, err := obj.Store.GetGatewayFirmware(gw.Model, false)
		if err != nil {
			if err == storage.ErrDoesNotExist {
				return nil, status.Errorf(codes.NotFound, "Firmware not found for model: %s", gw.Model)
			}
			log.WithError(err).Errorf("Failed to get firmware information for model: %s", gw.Model)
			return nil, status.Errorf(codes.Unknown, "Failed to get firmware information for model: %s", gw.Model)
		}

		if bytes.Equal(firmware.FirmwareHash[:], gw.FirmwareHash[:]) == false {
			response.NewFirmwareLink = firmware.ResourceLink
			// update gateway firmware hash as well
			copy(gw.FirmwareHash[:], firmware.FirmwareHash[:])
		}
	}

	// update gateway with osVersion and statistics
	if gw.OsVersion != req.OsVersion {
		// update provisioning server
		client, err := provisionserver.CreateClientWithCert(config.C.ProvisionServer.ProvisionServer,
			config.C.ProvisionServer.CACert,
			config.C.ProvisionServer.TLSCert,
			config.C.ProvisionServer.TLSKey)
		if err == nil {
			_, err := client.UpdateGateway(context.Background(), &pspb.UpdateGatewayRequest{
				Sn:        gw.SerialNumber,
				Mac:       gw.MAC.String(),
				OsVersion: req.OsVersion,
			})
			if err != nil {
				log.WithError(err).Error("Failed to call HeartbeatAPI: UpdateGateway")
			}
		} else {
			log.WithError(err).Error("Failed to create provisioning server client.")
		}
	}

	gw.OsVersion = req.OsVersion
	gw.Statistics = req.Statistics
	if err := obj.Store.UpdateGateway(ctx, &gw); err != nil {
		log.WithError(err).Errorf("Failed to update gateway: %s", gw.MAC.String())
	}
	/*	else {
		if err := tx.Commit(); err != nil {
			log.WithError(err).Errorf("Failed to update gateway: %s", gw.MAC.String())
		}
	}*/

	return &response, status.Error(codes.OK, "")
}
