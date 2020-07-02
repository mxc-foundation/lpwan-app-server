package m2m

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/tls"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	devPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	gwPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

var serviceName = "m2m server"

// Setup :
func Setup(conf config.Config) error {
	log.Info("Set up API for m2m server")

	if err := listenWithCredentials(conf.ApplicationServer.APIForM2M.Bind,
		conf.ApplicationServer.APIForM2M.CACert,
		conf.ApplicationServer.APIForM2M.TLSCert,
		conf.ApplicationServer.APIForM2M.TLSKey); err != nil {
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

	tx, err := storage.DB().Beginx()
	if err != nil {
		return err
	}

	pb.RegisterDeviceM2MServiceServer(gs, device.NewDeviceM2MAPI(device.DeviceM2MAPI{
		Store: devPg.New(tx.Tx, storage.DB().DB),
	}))

	pb.RegisterGatewayM2MServiceServer(gs, gateway.NewGatewayM2MAPI(gateway.GatewayM2MAPI{
		Store: gwPg.New(tx.Tx, storage.DB().DB),
	}))

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return errors.Wrap(err, "listenWithCredentials: start api listener error")
	}
	go gs.Serve(ln)

	return nil
}
