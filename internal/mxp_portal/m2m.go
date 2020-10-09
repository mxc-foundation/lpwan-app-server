package mxp_portal

import (
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal/data"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

var serviceName = "m2m server"

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "mxp_portal"

type controller struct {
	h         *store.Handler
	mxpCli    MxprotocolClientStruct
	mxpServer MxprotocolServerStruct
	name      string
	p         Pool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		mxpCli:    conf.ApplicationServer.APIForM2M,
		name:      moduleName,
		mxpServer: conf.M2MServer,
	}
	return nil
}

// Setup :
func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	log.Info("Set up API for m2m server")

	ctrl.h = h
	ctrl.p = &pool{
		mxprotocolServiceClients: make(map[string]mxprotocolServiceClient),
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
