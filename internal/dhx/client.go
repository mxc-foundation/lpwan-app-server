// Package dhx registers the appserver with dhx-center server
package dhx

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	api "github.com/mxc-foundation/lpwan-app-server/api/dhxapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

// Config defines configuration struct for dhx service module
type Config struct {
	Enable    bool                   `mapstructure:"enable"`
	DHXServer grpccli.ConnectionOpts `mapstructure:"dhx_server"`
}

type controller struct {
	dhxcconn *grpc.ClientConn
}

var ctrl *controller

// Register registers appserver on dhx-center server
func Register(supernodeID string, config Config) (err error) {
	if !config.Enable {
		log.Info("dhx center service disabled")
		return nil
	}

	ctrl = &controller{}
	ctrl.dhxcconn, err = grpccli.Connect(config.DHXServer)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"supernodeID": supernodeID,
		"server":      config.DHXServer.Server,
		"caCert":      config.DHXServer.CACert,
		"tlsCert":     config.DHXServer.TLSCert,
		"tlsKey":      config.DHXServer.TLSKey,
	}).Infof("register supernode id in dhx-center")

	go func() {
		cli := api.NewDHXSettingsServiceClient(ctrl.dhxcconn)

		for {
			_, err = cli.UpdateSupernode(context.Background(),
				&api.UpdateSupernodeRequest{DomainName: supernodeID})
			if err == nil {
				return
			}
			log.Warnf("failed to update supernode in dhx center: %s", err.Error())
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}
