package gateway

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	"google.golang.org/grpc"

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
	nets "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	org "github.com/mxc-foundation/lpwan-app-server/internal/api/external/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	gw "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

// Server represents gRPC server serving gateway
type Server struct {
	bind       string
	gs         *grpc.Server
	psCli      *pscli.Client
	serverAddr string
	st         Store
}

type controller struct {
	bindOld    string
	bindNew    string
	psCli      psPb.ProvisionClient
	serverAddr string
	st         Store
}

// Store defines db API used by gateway server
type Store interface {
	GetGatewayFirmwareList(ctx context.Context) (list []gw.GatewayFirmware, err error)
	UpdateGatewayFirmware(ctx context.Context, gwFw *gw.GatewayFirmware) (model string, err error)
	GetOrganization(ctx context.Context, id int64, forUpdate bool) (org.Organization, error)
	GetGatewayCount(ctx context.Context, filters gw.GatewayFilters) (int, error)
	CreateGateway(ctx context.Context, gw *gw.Gateway) error
	GetNetworkServer(ctx context.Context, id int64) (nets.NetworkServer, error)
	GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (nets.NetworkServer, error)
	DeleteGateway(ctx context.Context, mac lorawan.EUI64) error
	GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gw.Gateway, error)

	UpdateLastHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error
	UpdateFirstHeartbeat(ctx context.Context, mac lorawan.EUI64, time int64) error
	GetFirstHeartbeat(ctx context.Context, mac lorawan.EUI64) (int64, error)
	GetGatewayFirmware(ctx context.Context, model string, forUpdate bool) (gwFw gw.GatewayFirmware, err error)
	UpdateGatewayAttributes(ctx context.Context, mac lorawan.EUI64, firmware types.MD5SUM,
		osVersion, statistics string) error

	Tx(ctx context.Context, f func(context.Context, *pgstore.PgStore) error) error

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

// Start starts gRPC server of gateway server
func Start(st Store, ServerAddr string, psCli *pscli.Client,
	conf gw.GatewayBindStruct, updateSchedule string) (old *Server, new *Server, err error) {
	log.Info("Set up API for gateway")

	// listen to new gateways
	new = &Server{
		bind:       conf.NewGateway.Bind,
		psCli:      psCli,
		serverAddr: ServerAddr,
		st:         st,
	}
	err = new.listenWithCredentials("New Gateway API", conf.NewGateway.Bind,
		conf.NewGateway.CACert,
		conf.NewGateway.TLSCert,
		conf.NewGateway.TLSKey)
	if err != nil {
		return nil, nil, err
	}
	// listen to old gateways
	old = &Server{
		bind:       conf.OldGateway.Bind,
		psCli:      psCli,
		serverAddr: ServerAddr,
		st:         st,
	}
	err = old.listenWithCredentials("Old Gateway API", conf.OldGateway.Bind,
		conf.OldGateway.CACert,
		conf.OldGateway.TLSCert,
		conf.OldGateway.TLSKey)
	if err != nil {
		return nil, nil, err
	}

	ctrl := &controller{
		bindOld: conf.OldGateway.Bind,
		bindNew: conf.NewGateway.Bind,
		psCli:   psCli.GetPServerClient(),
		st:      st,
	}
	if err := ctrl.scheduleUpdateFirmwareFromProvisioningServer(context.Background(), updateSchedule); err != nil {
		return nil, nil, err
	}

	return old, new, nil
}

// Stop gracefully stops gRPC server
func (s *Server) Stop() {
	s.gs.GracefulStop()
}

func (s *Server) listenWithCredentials(service, bind, caCert, tlsCert, tlsKey string) error {
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

	gwpb.RegisterHeartbeatServiceServer(gs, NewHeartbeatAPI(bind, s.st, s.psCli))

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}

	go func() {
		_ = gs.Serve(ln)
	}()
	s.gs = gs

	return nil
}

// UpdateFirmware respond to manual operation to update gateway firmware from provisioining server
func UpdateFirmware(ctx context.Context, st Store, bindOldGateway, bindNewGateway, serverAddr string, psCli psPb.ProvisionClient) error {
	ctrl := &controller{
		bindOld:    bindOldGateway,
		bindNew:    bindNewGateway,
		psCli:      psCli,
		serverAddr: serverAddr,
		st:         st,
	}
	return ctrl.updateFirmwareFromProvisioningServer(ctx)
}

func (c *controller) updateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.Info("Check firmware update...")
	gwFwList, err := c.st.GetGatewayFirmwareList(ctx)
	if err != nil {
		log.WithError(err).Errorf("Failed to get gateway firmware list.")
		return err
	}

	bindPortOld, err := extractPort(c.bindOld)
	if err != nil {
		return err
	}
	bindPortNew, err := extractPort(c.bindNew)
	if err != nil {
		return err
	}

	// send update
	for _, v := range gwFwList {
		res, err := c.psCli.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
			Model:          v.Model,
			SuperNodeAddr:  c.serverAddr,
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

		gatewayFw := gw.GatewayFirmware{
			Model:        v.Model,
			ResourceLink: res.ResourceLink,
			FirmwareHash: md5sum,
		}

		model, _ := c.st.UpdateGatewayFirmware(ctx, &gatewayFw)
		if model == "" {
			log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
		}

	}

	return nil
}

func (c *controller) scheduleUpdateFirmwareFromProvisioningServer(ctx context.Context, updateSchedule string) error {
	log.Info("Start schedule to update gateway firmware...")

	cron := cron.New()
	err := cron.AddFunc(updateSchedule, func() {
		if err := c.updateFirmwareFromProvisioningServer(ctx); err != nil {
			log.WithError(err).Error("update firmware on schdule error")
		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go cron.Start()

	return nil
}

// AddGateway add new gateway and sync across all relevant servers
func AddGateway(ctx context.Context, st Store, gateway *gw.Gateway, createReq ns.CreateGatewayRequest,
	mxpCli pb.GSGatewayServiceClient, nsCli *nscli.Client) error {
	organization, err := st.GetOrganization(ctx, gateway.OrganizationID, true)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	// Validate max. gateway count when != 0.
	if organization.MaxGatewayCount != 0 {
		count, err := st.GetGatewayCount(ctx, gw.GatewayFilters{
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

	timestampCreatedAt := timestamppb.New(time.Now())
	// add this gateway to m2m server
	_, err = mxpCli.AddGatewayInM2MServer(context.Background(), &pb.AddGatewayInM2MServerRequest{
		OrgId: gateway.OrganizationID,
		GwProfile: &pb.AppServerGatewayProfile{
			Mac:         gateway.MAC.String(),
			OrgId:       gateway.OrganizationID,
			Description: gateway.Description,
			Name:        gateway.Name,
			CreatedAt:   timestampCreatedAt,
		},
	})
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return helpers.ErrToRPCError(err)
	}

	n, err := st.GetNetworkServer(ctx, gateway.NetworkServerID)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	client, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	_, err = client.CreateGateway(ctx, &createReq)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}

	err = st.CreateGateway(ctx, gateway)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	return nil
}

// DeleteGateway deletes gateway and sync across all relevant servers. Must be called from within transaction
func DeleteGateway(ctx context.Context, mac lorawan.EUI64, st Store, psCli psPb.ProvisionClient, nsCli *nscli.Client) error {
	// if the gateway is MatchX gateway, unregister it from provisioning server
	obj, err := st.GetGateway(ctx, mac, false)
	if err != nil {
		return errors.Wrap(err, "get gateway error")
	}

	n, err := st.GetNetworkServerForGatewayMAC(ctx, mac)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	// delete this gateway from m2m-server
	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	_, err = gwClient.DeleteGatewayInM2MServer(context.Background(), &pb.DeleteGatewayInM2MServerRequest{
		MacAddress: mac.String(),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		log.WithError(err).Error("delete gateway from m2m-server error")
	}

	client, err := nsCli.GetNetworkServerServiceClient(n.ID)
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
		if err != nil && status.Code(err) != codes.NotFound {
			return errors.Wrap(err, "failed to unregister from provisioning server")
		}
	}

	if err := st.DeleteGateway(ctx, obj.MAC); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"id":     mac,
		"ctx_id": ctx.Value(logging.ContextIDKey),
	}).Info("gateway deleted")

	return nil
}
