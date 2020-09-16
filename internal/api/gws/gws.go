package gws

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	gwpb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
)

type GatewayBindStruct struct {
	NewGateway struct {
		Bind    string `mapstructure:"new_gateway_bind"`
		CACert  string `mapstructure:"ecc_ca_cert"`
		TLSCert string `mapstructure:"ecc_tls_cert"`
		TLSKey  string `mapstructure:"ecc_tls_key"`
	} `mapstructure:"new_gateway"`

	OldGateway struct {
		Bind    string `mapstructure:"old_gateway_bind"`
		CACert  string `mapstructure:"rsa_ca_cert"`
		TLSCert string `mapstructure:"rsa_tls_cert"`
		TLSKey  string `mapstructure:"rsa_tls_key"`
	} `mapstructure:"old_gateway"`
}

type controller struct {
	s GatewayBindStruct
}

var ctrl *controller

func SettingsSetup(bindStruct GatewayBindStruct) error {
	ctrl = &controller{
		s: bindStruct,
	}

	return nil
}
func GetSettings() GatewayBindStruct {
	return ctrl.s
}

// Setup configures the package.
func Setup() error {
	log.Info("Set up API for gateway")

	// listen to new gateways
	if err := listenWithCredentials("New Gateway API", ctrl.s.NewGateway.Bind,
		ctrl.s.NewGateway.CACert,
		ctrl.s.NewGateway.TLSCert,
		ctrl.s.NewGateway.TLSKey); err != nil {
		return err
	}

	// listen to old gateways
	if err := listenWithCredentials("Old Gateway API", ctrl.s.OldGateway.Bind,
		ctrl.s.OldGateway.CACert,
		ctrl.s.OldGateway.TLSCert,
		ctrl.s.OldGateway.TLSKey); err != nil {
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
