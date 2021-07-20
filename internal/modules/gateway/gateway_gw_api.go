package gateway

import (
	"bytes"
	"context"

	/* #nosec */
	"crypto/md5"
	"strings"
	"time"

	"github.com/brocaar/lorawan"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	pspb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	gw "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

// HeartbeatAPI exports the HeartbeatAPI related functions.
type HeartbeatAPI struct {
	BindPort string
	st       Store
	psCli    pspb.ProvisionClient
}

// NewHeartbeatAPI creates new HeartbeatAPI
func NewHeartbeatAPI(bind string, st Store, psCli *pscli.Client) *HeartbeatAPI {
	return &HeartbeatAPI{
		BindPort: bind,
		st:       st,
		psCli:    psCli.GetPServerClient(),
	}
}

func (a *HeartbeatAPI) verifyGateway(ctx context.Context, mac string, model string) (*gw.Gateway, error) {
	// check if gateway exists
	var gatewayEUI = lorawan.EUI64{}
	if err := gatewayEUI.UnmarshalText([]byte(mac)); err != nil {
		logrus.WithError(err).Error("api/Heartbeat: Failed to convert gateway mac address")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid gateway mac format: %s", mac)
	}

	gateway, err := a.st.GetGateway(ctx, gatewayEUI, true)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "Failed to select gateway by mac: %s", gatewayEUI.String())
	}

	// verify gateway model
	if gateway.Model != model {
		logrus.Errorf("Request model does not match saved gateway.")
		return nil, status.Errorf(codes.Unauthenticated, "Request model does not match saved gateway.")
	}

	// important: do this before config and firmware update
	// mining : update heartbeat only for new gateways
	if !strings.HasPrefix(gateway.Model, "MX19") {
		// do not proceed, return nil for *gw.Gateway
		// it depends on the gateway which sends heartbeat to server, to prevent gatway from sending continuous request
		//  after receiving response with error, here we return nil for error
		return nil, nil
	}
	return &gateway, nil
}

// Heartbeat receives gateway heartbeat signals, updates heartbeat timestamp, os_version, statistics, firmware_hash,
//  returns new firmware link and config if changed
func (a *HeartbeatAPI) Heartbeat(ctx context.Context, req *gwpb.HeartbeatRequest) (*gwpb.HeartbeatResponse, error) {
	if a.BindPort == "8005" {
		return nil, status.Error(codes.PermissionDenied, "only new model of gateway should be processed")
	}

	gateway, err := a.verifyGateway(ctx, req.GatewayMac, req.Model)
	if err != nil {
		return nil, err
	} else if gateway == nil {
		// it depends on the gateway which sends heartbeat to server, to prevent gatway from sending continuous request
		// after receiving response with error, here we return nil for error
		return &gwpb.HeartbeatResponse{}, nil
	}

	logrus.Info("processing MX19 gateway")
	currentHeartbeat := time.Now().Unix()
	lastHeartbeat := gateway.LastHeartbeat
	firstHeartbeat := gateway.FirstHeartbeat

	if lastHeartbeat == 0 {
		firstHeartbeat = currentHeartbeat
		lastHeartbeat = currentHeartbeat
	} else if currentHeartbeat-lastHeartbeat > mining.GetSettings().HeartbeatOfflineLimit {
		// gateway is considered as went offline before in this case, during offline time, gateway should not be paid
		// set firstHeartbeat to currentHeartbeat
		firstHeartbeat = currentHeartbeat
		lastHeartbeat = currentHeartbeat
	} else {
		// gateway is considered as online all the time, no need to update firstHeartbeat unless firstHeartbeat is 0
		if firstHeartbeat == 0 {
			// TODO: before deploying this fix, there might be firstHeartbeat set to 0 in live servers, can be optimized
			//  off later
			firstHeartbeat = lastHeartbeat
		}
		lastHeartbeat = currentHeartbeat
	}
	err = a.st.UpdateGatewayHeartbeat(ctx, gateway.MAC, firstHeartbeat, lastHeartbeat)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Update last heartbeat error")
	}
	return a.checkStatusAndFirmwareUpdate(ctx, gateway, req)
}

type gatewayAttributes struct {
	firmwareHash types.MD5SUM
	osVersion    string
	statistics   string
}

func (a *HeartbeatAPI) checkStatusAndFirmwareUpdate(ctx context.Context, gateway *gw.Gateway,
	req *gwpb.HeartbeatRequest) (*gwpb.HeartbeatResponse, error) {
	response := gwpb.HeartbeatResponse{}
	updatedGateway := gatewayAttributes{
		firmwareHash: gateway.FirmwareHash,
		osVersion:    gateway.OsVersion,
		statistics:   gateway.Statistics,
	}
	// compare config hash
	/* #nosec */
	configHash := md5.Sum([]byte(gateway.Config))
	b := types.MD5SUM{}
	if err := b.UnmarshalText([]byte(req.ConfigHash)); err != nil {
		logrus.WithError(err).Errorf("Failed to unmarshal config hash: %s", req.ConfigHash)
		return nil, status.Errorf(codes.DataLoss, "Failed to unmarshal config hash: %s", req.ConfigHash)
	}

	if !bytes.Equal(configHash[:], b[:]) {
		response.Config = gateway.Config
	}

	// check if firmware updated
	if gateway.AutoUpdateFirmware {
		firmware, err := a.st.GetGatewayFirmware(ctx, gateway.Model, false)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to get firmware information for model: %s", gateway.Model)
		}

		if !bytes.Equal(firmware.FirmwareHash[:], gateway.FirmwareHash[:]) {
			response.NewFirmwareLink = firmware.ResourceLink
			// update gateway firmware hash as well
			copy(updatedGateway.firmwareHash[:], firmware.FirmwareHash[:])
		}
	}

	updatedGateway.osVersion = req.OsVersion
	updatedGateway.statistics = req.Statistics

	if !bytes.Equal(updatedGateway.firmwareHash[:], gateway.FirmwareHash[:]) ||
		updatedGateway.osVersion != gateway.OsVersion ||
		updatedGateway.statistics != gateway.Statistics {
		if err := a.st.UpdateGatewayAttributes(ctx, gateway.MAC, updatedGateway.firmwareHash,
			updatedGateway.osVersion, updatedGateway.statistics); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	// update gateway osVersion on provisioning server
	if gateway.OsVersion != req.OsVersion {
		// update provisioning server
		_, err := a.psCli.UpdateGateway(context.Background(), &pspb.UpdateGatewayRequest{
			Sn:        gateway.SerialNumber,
			Mac:       gateway.MAC.String(),
			OsVersion: req.OsVersion,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to call HeartbeatAPI: UpdateGateway")
		}
	}
	return &response, nil
}
