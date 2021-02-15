package mxpapisrv

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
)

var serviceName = "m2m server"

// Config contains configuration for gRPC server serving mxp server
type Config struct {
	Bind    string `mapstructure:"bind"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

// MXPAPIServer represents gRPC server serving mxp server
type MXPAPIServer struct {
	h  *pgstore.PgStore
	gs *grpc.Server
}

// Start starts gRPC server that serves mxp server
func Start(h *pgstore.PgStore, cfg Config) (*MXPAPIServer, error) {
	log.Info("Starting API for m2m server")

	srv := &MXPAPIServer{
		h: h,
	}

	if err := srv.listenWithCredentials(
		cfg.Bind,
		cfg.CACert,
		cfg.TLSCert,
		cfg.TLSKey); err != nil {
		return nil, err
	}

	return srv, nil
}

// Stop gracefully stops gRPC server
func (srv *MXPAPIServer) Stop() {
	srv.gs.GracefulStop()
}

func (srv *MXPAPIServer) listenWithCredentials(bind, caCert, tlsCert, tlsKey string) error {
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
	srv.gs = gs

	pb.RegisterDeviceM2MServiceServer(gs, NewDeviceM2MAPI(srv.h))
	pb.RegisterGatewayM2MServiceServer(gs, NewGatewayM2MAPI(srv.h))
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
