package gws

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
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

	gwpb.RegisterHeartbeatServiceServer(gs, gateway.NewHeartbeatAPI(bind))

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}
	go gs.Serve(ln)

	return nil
}
