package gateway

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	ps "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "gateway"

type controller struct {
	name               string
	st                 *store.Handler
	ps                 ps.ProvisioningServerStruct
	bindPortOldGateway string
	bindPortNewGateway string
	serverAddr         string
	conf               GatewayBindStruct

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, s config.Config) error {

	ctrl = &controller{
		name:       moduleName,
		serverAddr: s.General.ServerAddr,
		ps:         s.ProvisionServer,
		conf:       s.ApplicationServer.APIForGateway,
	}

	bindStruct := s.ApplicationServer.APIForGateway
	if strArray := strings.Split(bindStruct.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", bindStruct.OldGateway.Bind))
	} else {
		ctrl.bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(bindStruct.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", bindStruct.NewGateway.Bind))
	} else {
		ctrl.bindPortNewGateway = strArray[1]
	}

	return nil
}

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	if ctrl.bindPortNewGateway == "" || ctrl.bindPortOldGateway == "" {
		return errors.New("bindPortNewGateway and bindPortOldGateway not initiated")
	}
	ctrl.st = h

	log.Info("Set up API for gateway")

	// listen to new gateways
	if err := listenWithCredentials("New Gateway API", ctrl.conf.NewGateway.Bind,
		ctrl.conf.NewGateway.CACert,
		ctrl.conf.NewGateway.TLSCert,
		ctrl.conf.NewGateway.TLSKey); err != nil {
		return err
	}

	// listen to old gateways
	if err := listenWithCredentials("Old Gateway API", ctrl.conf.OldGateway.Bind,
		ctrl.conf.OldGateway.CACert,
		ctrl.conf.OldGateway.TLSCert,
		ctrl.conf.OldGateway.TLSKey); err != nil {
		return err
	}

	if err := ctrl.updateFirmwareFromProvisioningServer(context.Background()); err != nil {
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

func (c *controller) updateFirmwareFromProvisioningServer(ctx context.Context) error {
	log.WithFields(log.Fields{
		"provisioning-server": ctrl.ps.Server,
		"caCert":              ctrl.ps.CACert,
		"tlsCert":             ctrl.ps.TLSCert,
		"tlsKey":              ctrl.ps.TLSKey,
		"schedule":            ctrl.ps.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")

	supernodeAddr := c.serverAddr

	cron := cron.New()
	err := cron.AddFunc(ctrl.ps.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := c.st.GetGatewayFirmwareList(ctx)
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := pscli.GetPServerClient()
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
				Model:          v.Model,
				SuperNodeAddr:  supernodeAddr,
				PortOldGateway: c.bindPortOldGateway,
				PortNewGateway: c.bindPortNewGateway,
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

			model, _ := c.st.UpdateGatewayFirmware(ctx, &gatewayFw)
			if model == "" {
				log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
			}

		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go cron.Start()

	return nil
}

func AddGateway(ctx context.Context, gw *Gateway, createReq ns.CreateGatewayRequest) error {
	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max gateway count.
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		org, err := handler.GetOrganization(ctx, gw.OrganizationID, true)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		// Validate max. gateway count when != 0.
		if org.MaxGatewayCount != 0 {
			count, err := handler.GetGatewayCount(ctx, GatewayFilters{
				OrganizationID: org.ID,
				Search:         "",
			})
			if err != nil {
				return helpers.ErrToRPCError(err)
			}

			if count >= org.MaxGatewayCount {
				return helpers.ErrToRPCError(errHandler.ErrOrganizationMaxGatewayCount)
			}
		}

		err = handler.CreateGateway(ctx, gw)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		timestampCreatedAt, _ := ptypes.TimestampProto(time.Now())
		// add this gateway to m2m server
		gwClient, err := m2mcli.GetM2MGatewayServiceClient()
		if err != nil {
			return status.Errorf(codes.Unavailable, err.Error())
		}

		_, err = gwClient.AddGatewayInM2MServer(context.Background(), &pb.AddGatewayInM2MServerRequest{
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

		n, err := handler.GetNetworkServer(ctx, gw.NetworkServerID)
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
	}); err != nil {
		return status.Errorf(codes.Unknown, err.Error())
	}

	return nil
}

func DeleteGateway(ctx context.Context, mac lorawan.EUI64) error {
	if err := ctrl.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		// if the gateway is MatchX gateway, unregister it from provisioning server
		obj, err := handler.GetGateway(ctx, mac, false)
		if err != nil {
			return errors.Wrap(err, "get gateway error")
		}

		n, err := handler.GetNetworkServerForGatewayMAC(ctx, mac)
		if err != nil {
			return errors.Wrap(err, "get network-server error")
		}

		if err := handler.DeleteGateway(ctx, obj.MAC); err != nil {
			return err
		}

		// delete this gateway from m2m-server
		gwClient, err := m2mcli.GetM2MGatewayServiceClient()
		if err != nil {
			return err
		}

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
			provClient, err := pscli.GetPServerClient()
			if err != nil {
				return errors.Wrap(err, "failed to connect to provisioning server")
			}

			_, err = provClient.UnregisterGw(context.Background(), &psPb.UnregisterGwRequest{
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
	}); err != nil {
		return err
	}

	return nil
}
