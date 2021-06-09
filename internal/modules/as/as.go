package as

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"net"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/brocaar/chirpstack-api/go/v3/as"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/as/data"
)

// Setup configures the package.
func Setup(h *store.Handler, config AppserverStruct, gIntegrations []models.IntegrationHandler) error {

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
			return errors.Wrap(err, "get transport credentials error")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	server := grpc.NewServer(grpcOpts...)
	as.RegisterApplicationServerServiceServer(server, NewApplicationServerAPI(h, gIntegrations))

	ln, err := net.Listen("tcp", config.Bind)
	if err != nil {
		return errors.Wrap(err, "start application-server api listener error")
	}
	go func() {
		_ = server.Serve(ln)
	}()

	return nil
}
