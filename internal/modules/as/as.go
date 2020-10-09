package as

import (
	"fmt"
	"net"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/brocaar/chirpstack-api/go/v3/as"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/as/data"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "appserver"

type controller struct {
	name string
	st   *store.Handler
	s    AppserverStruct
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name: moduleName,
		s:    conf.ApplicationServer.API,
	}
	return nil
}

// Setup configures the package.
func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl.st = h

	log.WithFields(log.Fields{
		"bind":     ctrl.s.Bind,
		"ca_cert":  ctrl.s.CACert,
		"tls_cert": ctrl.s.TLSCert,
		"tls_key":  ctrl.s.TLSKey,
	}).Info("api/as: starting application-server api")

	grpcOpts := helpers.GetgRPCServerOptions()
	if ctrl.s.CACert != "" && ctrl.s.TLSCert != "" && ctrl.s.TLSKey != "" {
		creds, err := helpers.GetTransportCredentials(ctrl.s.CACert, ctrl.s.TLSCert, ctrl.s.TLSKey, true)
		if err != nil {
			return errors.Wrap(err, "get transport credentials error")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	server := grpc.NewServer(grpcOpts...)
	as.RegisterApplicationServerServiceServer(server, NewApplicationServerAPI())

	ln, err := net.Listen("tcp", ctrl.s.Bind)
	if err != nil {
		return errors.Wrap(err, "start application-server api listener error")
	}
	go server.Serve(ln)

	return nil
}
