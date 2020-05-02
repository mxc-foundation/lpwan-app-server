package gws

import (
	"context"
	"crypto/md5"
	"github.com/brocaar/lorawan"
	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strings"
	"time"
)

// Setup configures the package.
func Setup(conf config.Config) error {
	log.Info("Set up API for gateway")

	// listen to new gateways
	if err := listenWithCredentials("New Gateway API", conf.ApplicationServer.APIForGateway.NewGateway.Bind,
		conf.ApplicationServer.APIForGateway.NewGateway.CACert,
		conf.ApplicationServer.APIForGateway.NewGateway.TLSCert,
		conf.ApplicationServer.APIForGateway.NewGateway.TLSKey); err != nil {
		return err
	}

	// listen to old gateways
	if err := listenWithCredentials("Old Gateway API", conf.ApplicationServer.APIForGateway.OldGateway.Bind,
		conf.ApplicationServer.APIForGateway.OldGateway.CACert,
		conf.ApplicationServer.APIForGateway.OldGateway.TLSCert,
		conf.ApplicationServer.APIForGateway.OldGateway.TLSKey); err != nil {
		return err
	}

	return nil
}

func listenWithCredentials(service, bind, caCert, tlsCert, tlsKey string) error {
	log.WithFields(log.Fields{
		"bind":     bind,
		"ca-cert":  caCert,
		"tls-cert": tlsCert,
		"tls-key":  tlsKey,
	}).Info("listen With Credentials")

	gs, err := tls.NewServerWithTLSCredentials(service, caCert, tlsCert, tlsKey)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: get new server error")
	}

	gwAPI := GatewayAPI(bind)
	gwpb.RegisterHeartbeatServiceServer(gs, gwAPI)

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}
	go gs.Serve(ln)

	return nil
}

// API exports the API related functions.
type API struct{
	BindPort string
}

// NewGatewayAPI creates new API
func GatewayAPI(bindPort string) *API {
	return &API{
		BindPort: bindPort,
	}
}

func (obj *API) Heartbeat(ctx context.Context, req *gwpb.HeartbeatRequest) (*gwpb.HeartbeatResponse, error) {
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

	gw, err := storage.GetGateway(ctx, storage.DB(), gatewayEUI, true)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			return nil, status.Errorf(codes.Unauthenticated, "Object does not exist: %s", gatewayEUI.String())
		}
		log.WithError(err).Errorf("Failed to select gateway by mac: %s", gatewayEUI.String())
		return nil, status.Errorf(codes.Unknown, "Failed to select gateway by mac: %s", gatewayEUI.String())
	}

	// verify gateway model
	if gw.Model != req.Model {
		log.Errorf("Request model does not match saved gateway.")
		return nil, status.Errorf(codes.Unauthenticated, "Request model does not match saved gateway.")
	}

	// compare config hash
	h := md5.Sum([]byte(gw.Config))
	configHash := string(h[:])
	if configHash != req.ConfigHash {
		response.Config = gw.Config
	}

	// check if firmware updated
	firmware, err := storage.GetGatewayFirmware(storage.DB(), gw.Model, false)
	if err != nil {
		if err == storage.ErrDoesNotExist {
			return nil, status.Errorf(codes.NotFound, "Firmware not found for model: %s", gw.Model)
		}
		log.WithError(err).Errorf("Failed to get firmware information for model: %s", gw.Model)
		return nil, status.Errorf(codes.Unknown, "Failed to get firmware information for model: %s", gw.Model)
	}

	if firmware.Updated == true {
		response.NewFirmwareLink = firmware.ResourceLink
	}

	// mining : update heartbeat only for new gateways
	if strings.HasPrefix(gw.Model, "MX19") {
		currentHeartbeat := time.Now().Unix()
		lastHeartbeat := gw.LastHeartbeat

		// if last heartbeat is 0, update first, last heartbeat
		// if offline longer than 30 mins, last heartbeat and first heartbeat = current heartbeat
		if lastHeartbeat == 0 || currentHeartbeat-lastHeartbeat > 1800 {
			gw.LastHeartbeat = currentHeartbeat
			gw.FirstHeartbeat = currentHeartbeat
		} else {
			gw.LastHeartbeat = currentHeartbeat
		}
	}

	// update gateway with osVersion and statistics
	gw.OsVersion = req.OsVersion
	gw.Statistics = req.Statistics
	if err := storage.UpdateGateway(ctx, storage.DB(), &gw); err != nil {
		log.WithError(err).Errorf("Failed to update gateway: %s", gw.MAC.String())
	}

	return &response, status.Error(codes.OK, "")
}