package as

import (
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/brocaar/chirpstack-api/go/v3/as"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/as/data"
)

// NetworkServerAPIServer represents gRPC server serving network server
type NetworkServerAPIServer struct {
	gs *grpc.Server
}

// Start configures the package.
func Start(h *store.Handler, config AppserverStruct, gIntegrations []models.IntegrationHandler) (*NetworkServerAPIServer, error) {
	log.WithFields(log.Fields{
		"bind":     config.Bind,
		"ca_cert":  config.CACert,
		"tls_cert": config.TLSCert,
		"tls_key":  config.TLSKey,
	}).Info("api/as: starting application-server api")

	grpcOpts := helpers.GetgRPCServerOptions()
	if config.CACert != "" && config.TLSCert != "" && config.TLSKey != "" {
		creds, err := helpers.GetTransportCredentials(config.CACert, config.TLSCert, config.TLSKey, true)
		if err != nil {
			return nil, errors.Wrap(err, "get transport credentials error")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	server := grpc.NewServer(grpcOpts...)
	as.RegisterApplicationServerServiceServer(server, NewApplicationServerAPI(h, gIntegrations))

	ln, err := net.Listen("tcp", config.Bind)
	if err != nil {
		return nil, errors.Wrap(err, "start application-server api listener error")
	}
	go func() {
		_ = server.Serve(ln)
	}()

	return &NetworkServerAPIServer{gs: server}, nil
}

// Stop gracefully stops gRPC server
func (srv *NetworkServerAPIServer) Stop() {
	srv.gs.GracefulStop()
}
