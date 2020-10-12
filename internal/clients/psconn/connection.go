package psconn

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	. "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "provisioning_server_portal"

type controller struct {
	name               string
	p                  Pool
	provisioningServer ProvisioningServerStruct

	moduleUp bool
}

// Pool defines a set of interfaces to operate with internal resource
type Pool interface {
	get() (*grpc.ClientConn, error)
}

type provisioningServerClient struct {
	clientConn *grpc.ClientConn
}

type pool struct {
	sync.RWMutex
	provisioningServerClients map[string]provisioningServerClient
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {

	ctrl = &controller{
		name:               moduleName,
		provisioningServer: conf.ProvisionServer,
	}

	return nil
}

// Setup initialize module on start
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.p = &pool{
		provisioningServerClients: make(map[string]provisioningServerClient),
	}

	if _, err := ctrl.p.get(); err != nil {
		return err
	}

	return nil
}

// get returns a M2MServerServiceClient for the given server (hostname:ip).
func (p *pool) get() (*grpc.ClientConn, error) {
	defer p.Unlock()
	p.Lock()

	hostname := ctrl.provisioningServer.Server
	caCert := ctrl.provisioningServer.CACert
	tlsCert := ctrl.provisioningServer.TLSCert
	tlsKey := ctrl.provisioningServer.TLSKey

	var connect bool
	c, ok := p.provisioningServerClients[hostname]
	if !ok {
		connect = true
	}

	if connect {
		clientConn, err := p.createClient(hostname, caCert, tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "create provisioning server client error")
		}
		c = provisioningServerClient{
			clientConn: clientConn,
		}
		p.provisioningServerClients[hostname] = c
	}

	return c.clientConn, nil
}

func (p *pool) createClient(hostname, caCert, tlsCert, tlsKey string) (*grpc.ClientConn, error) {
	logrusEntry := log.NewEntry(log.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}

	nsOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_logrus.UnaryClientInterceptor(logrusEntry, logrusOpts...),
		),
		grpc.WithStreamInterceptor(
			grpc_logrus.StreamClientInterceptor(logrusEntry, logrusOpts...),
		),
	}

	if hostname == "" || caCert == "" || tlsCert == "" || tlsKey == "" {
		return nil, errors.New("invalid credentials")
	}

	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls key-pair error")
	}

	var caCertPool *x509.CertPool
	rawCaCert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(caCert)))
	if err != nil {
		return nil, errors.Wrap(err, "load ca certificate error")
	}

	caCertPool = x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(rawCaCert) {
		return nil, fmt.Errorf("append ca certificate error: %s", ctrl.provisioningServer.CACert)
	}

	nsOpts = append(nsOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})))

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	grpcClient, err := grpc.DialContext(ctx, hostname, nsOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial provisioning server error")
	}

	return grpcClient, nil
}

// GetPServerClient returns a new ProvisionClient of provisioning server
func GetPServerClient() (api.ProvisionClient, error) {
	grpcClient, err := ctrl.p.get()
	if err != nil {
		return nil, err
	}

	client := api.NewProvisionClient(grpcClient)

	return client, nil
}
