package provisionserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
	"path/filepath"
)


// Setup configures the provision-server package.
func Setup(conf config.Config) error {
	supernode, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "Failed to set up provisioning server config.")
	}

	c := cron.New()
	err = c.AddFunc(conf.ProvisionServer.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := storage.GetGatewayFirmwareList(storage.DB())
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := CreateClientWithCert(conf.ProvisionServer.ProvisionServer,
			conf.ProvisionServer.CACert,
			conf.ProvisionServer.TLSCert,
			conf.ProvisionServer.TLSKey)
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &api.GetUpdateRequest{
				Model:                v.Model,
				SuperNodeAddr:        supernode,
				PortOldGateway:       conf.ApplicationServer.APIForGateway.OldGateway.Bind,
				PortNewGateway:       conf.ApplicationServer.APIForGateway.NewGateway.Bind,
			})
			if err != nil {
				log.WithError(err).Errorf("Failed to get update for gateway model: %s", v.Model)
				continue
			}

			gatewayFw := storage.GatewayFirmware{
				Model:        v.Model,
				ResourceLink: res.ResourceLink,
				Updated:      res.NewFirmware,
			}

			model, err := storage.UpdateGatewayFirmware(storage.DB(), &gatewayFw)
			if model == "" {
				log.Warnf("No fow updated for gateway_firmware at model=%s", v.Model)
			}

		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go c.Start()

	return nil
}

func CreateClientWithCert(hostname, CACert, TLSCert, TLSKey string) (api.ProvisionClient, error) {
	var grpcClient *grpc.ClientConn
	var opts []grpc.DialOption

	if hostname == "" || CACert == "" || TLSCert == "" || TLSKey == "" {
		return nil, errors.New("invalid credentials")
	}

	cert, err := tls.LoadX509KeyPair(TLSCert, TLSKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls key-pair error")
	}

	var caCertPool *x509.CertPool
	rawCaCert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(CACert)))
	if err != nil {
		return nil, errors.Wrap(err, "load ca certificate error")
	}

	caCertPool = x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(rawCaCert) {
		return nil, fmt.Errorf("append ca certificate error: %s", CACert)
	}

	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})))

	grpcClient, err = grpc.Dial(hostname, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial error")
	}

	client := api.NewProvisionClient(grpcClient)

	return client, nil
}
