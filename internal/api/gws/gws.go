package gws

import (
	"fmt"
	"net"
	"strings"

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

	var bindPortOldGateway string
	var bindPortNewGateway string
	if strArray := strings.Split(bindStruct.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", bindStruct.OldGateway.Bind))
	} else {
		bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(bindStruct.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", bindStruct.NewGateway.Bind))
	} else {
		bindPortNewGateway = strArray[1]
	}

	gateway.SetupFirmware(bindPortOldGateway, bindPortNewGateway)

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
