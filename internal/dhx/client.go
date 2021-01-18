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
	Enable      bool
	SupernodeID string
	DHXServer   grpccli.ConnectionOpts
}

type controller struct {
	dhxcconn *grpc.ClientConn
}

var ctrl *controller

// Setup prepares dhx service module
func Setup(config Config) (err error) {
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
		"supernodeID": config.SupernodeID,
		"server":      config.DHXServer.Server,
		"caCert":      config.DHXServer.CACert,
		"tlsCert":     config.DHXServer.TLSCert,
		"tlsKey":      config.DHXServer.TLSKey,
	}).Infof("register supernode id in dhx-center")

	go func() {
		_, err = GetDHXCenterServerClient().UpdateSupernode(context.Background(),
			&api.UpdateSupernodeRequest{DomainName: config.SupernodeID})

		for err != nil {
			log.Warnf("failed to update supernode in dhx center: %s", err.Error())
			time.Sleep(10 * time.Second)

			_, err = GetDHXCenterServerClient().UpdateSupernode(context.Background(),
				&api.UpdateSupernodeRequest{DomainName: config.SupernodeID})
		}
	}()

	return nil
}

// GetDHXCenterServerClient returns API client for DHXSettingsService
func GetDHXCenterServerClient() api.DHXSettingsServiceClient {
	return api.NewDHXSettingsServiceClient(ctrl.dhxcconn)
}
