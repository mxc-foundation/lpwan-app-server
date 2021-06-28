package gateway

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/lorawan"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	org "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	nets "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type Server struct {
	bindOldGateway string
	bindNewGateway string

	psCli      psPb.ProvisionClient
	serverAddr string
	st         Store
}

type Store interface {
	GetGatewayFirmwareList(ctx context.Context) (list []GatewayFirmware, err error)
	UpdateGatewayFirmware(ctx context.Context, gwFw *GatewayFirmware) (model string, err error)
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (org.Organization, error)
	GetGatewayCount(ctx context.Context, filters GatewayFilters) (int, error)
	CreateGateway(ctx context.Context, gw *Gateway) error
	GetNetworkServer(ctx context.Context, id int64) (nets.NetworkServer, error)
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (Gateway, error)
	GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (nets.NetworkServer, error)
	DeleteGateway(ctx context.Context, mac lorawan.EUI64) error

	InTx() bool
}

func extractPort(bindStr string) (string, error) {
	var strArray []string
	if strArray = strings.Split(bindStr, ":"); len(strArray) != 2 {
		return "", fmt.Errorf("cannot parse port from %s", bindStr)
	}
	if strArray[1] == "" {
		return "", fmt.Errorf("no valid bind port defined")
	}
	return strArray[1], nil
}

func Start(st Store, ServerAddr string, PSCli psPb.ProvisionClient, conf GatewayBindStruct, updateSchedule string) error {
	server := &Server{
		st:             st,
		psCli:          PSCli,
		serverAddr:     ServerAddr,
		bindOldGateway: conf.OldGateway.Bind,
		bindNewGateway: conf.NewGateway.Bind,
	}

	log.Info("Set up API for gateway")

	// listen to new gateways
	if err := listenWithCredentials("New Gateway API", conf.NewGateway.Bind,
		conf.NewGateway.CACert,
		conf.NewGateway.TLSCert,
		conf.NewGateway.TLSKey); err != nil {
		return err
	}

	// listen to old gateways
	if err := listenWithCredentials("Old Gateway API", conf.OldGateway.Bind,
		conf.OldGateway.CACert,
		conf.OldGateway.TLSCert,
		conf.OldGateway.TLSKey); err != nil {
		return err
	}

	if err := server.scheduleUpdateFirmwareFromProvisioningServer(context.Background(), updateSchedule); err != nil {
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

	gwpb.RegisterHeartbeatServiceServer(gs, NewHeartbeatAPI(bind))

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}

	go func() {
		_ = gs.Serve(ln)
	}()

	return nil
}

// UpdateFirmware respond to manual operation to update gateway firmware from provisioining server
func UpdateFirmware(ctx context.Context, st Store, bindOldGateway, bindNewGateway, serverAddr string, psCli psPb.ProvisionClient) error {
	server := &Server{
		bindOldGateway: bindOldGateway,
		bindNewGateway: bindNewGateway,
		psCli:          psCli,
		serverAddr:     serverAddr,
		st:             st,
	}
	return server.updateFirmwareFromProvisioningServer(ctx)
}

func (s *Server) updateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.Info("Check firmware update...")
	gwFwList, err := s.st.GetGatewayFirmwareList(ctx)
	if err != nil {
		log.WithError(err).Errorf("Failed to get gateway firmware list.")
		return err
	}

	bindPortOld, err := extractPort(s.bindOldGateway)
	if err != nil {
		return err
	}
	bindPortNew, err := extractPort(s.bindNewGateway)
	if err != nil {
		return err
	}

	// send update
	for _, v := range gwFwList {
		res, err := s.psCli.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
			Model:          v.Model,
			SuperNodeAddr:  s.serverAddr,
			PortOldGateway: bindPortOld,
			PortNewGateway: bindPortNew,
		})
		if err != nil {
			log.WithError(err).Errorf("Failed to get update for gateway model: %s", v.Model)
			continue
		}

		var md5sum types.MD5SUM
		if err := md5sum.UnmarshalText([]byte(res.FirmwareHash)); err != nil {
			log.WithError(err).Errorf("Failed to unmarshal firmware hash: %s", res.FirmwareHash)
			continue
		}

		gatewayFw := GatewayFirmware{
			Model:        v.Model,
			ResourceLink: res.ResourceLink,
			FirmwareHash: md5sum,
		}

		model, _ := s.st.UpdateGatewayFirmware(ctx, &gatewayFw)
		if model == "" {
			log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
		}

	}

	return nil
}

func (s *Server) scheduleUpdateFirmwareFromProvisioningServer(ctx context.Context, updateSchedule string) error {
	log.Info("Start schedule to update gateway firmware...")

	cron := cron.New()
	err := cron.AddFunc(updateSchedule, func() {
		if err := s.updateFirmwareFromProvisioningServer(ctx); err != nil {
			log.WithError(err).Error("update firmware on schdule error")
		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go cron.Start()

	return nil
}

func AddGateway(ctx context.Context, st Store, gw *Gateway, createReq ns.CreateGatewayRequest,
	mxpCli pb.GSGatewayServiceClient) error {
	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max gateway count.
	if !st.InTx() {
		return fmt.Errorf("AddGateway must be called from within transaction")
	}
	organization, err := st.GetOrganization(ctx, gw.OrganizationID, true)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	// Validate max. gateway count when != 0.
	if organization.MaxGatewayCount != 0 {
		count, err := st.GetGatewayCount(ctx, GatewayFilters{
			OrganizationID: organization.ID,
			Search:         "",
		})
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		if count >= organization.MaxGatewayCount {
			return helpers.ErrToRPCError(errHandler.ErrOrganizationMaxGatewayCount)
		}
	}

	err = st.CreateGateway(ctx, gw)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	timestampCreatedAt := timestamppb.New(time.Now())
	// add this gateway to m2m server
	_, err = mxpCli.AddGatewayInM2MServer(context.Background(), &pb.AddGatewayInM2MServerRequest{
		OrgId: gw.OrganizationID,
		GwProfile: &pb.AppServerGatewayProfile{
			Mac:         gw.MAC.String(),
			OrgId:       gw.OrganizationID,
			Description: gw.Description,
			Name:        gw.Name,
			CreatedAt:   timestampCreatedAt,
		},
	})
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	n, err := st.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	nStruct := &nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nStruct.GetNetworkServiceClient()
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	_, err = client.CreateGateway(ctx, &createReq)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}

	return nil
}

func DeleteGateway(ctx context.Context, mac lorawan.EUI64, st Store, psCli psPb.ProvisionClient) error {
	if !st.InTx() {
		return fmt.Errorf("DeleteGateway must be called from within transaction")
	}
	// if the gateway is MatchX gateway, unregister it from provisioning server
	obj, err := st.GetGateway(ctx, mac, false)
	if err != nil {
		return errors.Wrap(err, "get gateway error")
	}

	n, err := st.GetNetworkServerForGatewayMAC(ctx, mac)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	if err := st.DeleteGateway(ctx, obj.MAC); err != nil {
		return err
	}

	// delete this gateway from m2m-server
	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	_, err = gwClient.DeleteGatewayInM2MServer(context.Background(), &pb.DeleteGatewayInM2MServerRequest{
		MacAddress: mac.String(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		log.WithError(err).Error("delete gateway from m2m-server error")
	}

	nsStruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nsStruct.GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	_, err = client.DeleteGateway(ctx, &ns.DeleteGatewayRequest{
		Id: mac[:],
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return errors.Wrap(err, "delete gateway error")
	}

	if strings.HasPrefix(obj.Model, "MX") {
		_, err = psCli.UnregisterGw(context.Background(), &psPb.UnregisterGwRequest{
			Sn:  obj.SerialNumber,
			Mac: obj.MAC.String(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to unregister from provisioning server")
		}
	}

	log.WithFields(log.Fields{
		"id":     mac,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway deleted")

	return nil
}
