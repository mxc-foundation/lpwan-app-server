package mxp_portal

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
)

var serviceName = "m2m server"

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "mxp_portal"

type controller struct {
	h         *store.Handler
	mxpCli    data.MxprotocolClientStruct
	mxpServer grpccli.ConnectionOpts
	name      string
	m2mconn   *grpc.ClientConn

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{
		mxpCli:    conf.ApplicationServer.APIForM2M,
		name:      moduleName,
		mxpServer: conf.M2MServer,
	}
	return nil
}

// Setup :
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	log.Info("Set up API for m2m server")

	ctrl.h = h
	var err error
	ctrl.m2mconn, err = grpccli.Connect(ctrl.mxpServer)
	if err != nil {
		return fmt.Errorf("couldn't create mxprotocol server client: %v", err)
	}

	if err := listenWithCredentials(ctrl.mxpCli.Bind,
		ctrl.mxpCli.CACert,
		ctrl.mxpCli.TLSCert,
		ctrl.mxpCli.TLSKey); err != nil {
		return err
	}

	return nil
}

func listenWithCredentials(bind, caCert, tlsCert, tlsKey string) error {
	log.WithFields(log.Fields{
		"bind":     bind,
		"ca-cert":  caCert,
		"tls-cert": tlsCert,
		"tls-key":  tlsKey,
	}).Info("listen With Credentials")

	gs, err := tls.NewServerWithTLSCredentials(serviceName, caCert, tlsCert, tlsKey)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: get new server error")
	}

	pb.RegisterDeviceM2MServiceServer(gs, NewDeviceM2MAPI(ctrl.h))
	pb.RegisterGatewayM2MServiceServer(gs, NewGatewayM2MAPI(ctrl.h))
	pb.RegisterNotificationServiceServer(gs, NewNotificationAPI())

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}

	go func() {
		_ = gs.Serve(ln)
	}()

	return nil
}
