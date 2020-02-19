package provisionserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sync"
	"time"
)

var p Pool

// Pool defines the provision-server client pool.
type Pool interface {
	Get(hostname string, caCert, tlsCert, tlsKey []byte) (api.ProvisionServiceClient, error)
}

type client struct {
	client     api.ProvisionServiceClient
	clientConn *grpc.ClientConn
	caCert     []byte
	tlsCert    []byte
	tlsKey     []byte
}

// Setup configures the provision-server package.
func Setup(conf config.Config) error {
	p = &pool{
		clients: make(map[string]client),
	}

	return nil
}

type pool struct {
	sync.RWMutex
	clients map[string]client
}

func GetPool() Pool {
	return p
}

func SetPool(pool Pool) {
	p = pool
}

// Get returns a ProvisionServerClient for the given server (hostname:ip).
func (p *pool) Get(hostname string, caCert, tlsCert, tlsKey []byte) (api.ProvisionServiceClient, error) {
	p.Lock()
	defer p.Unlock()

	var connect bool
	// if there is no client, try to create a new one
	c, ok := p.clients[hostname]
	if !ok {
		connect = true
	}

	// if the connection exists in the map, but when the certificates changed
	// try to cloe the connection and re-connect
	if ok && (!bytes.Equal(c.caCert, caCert) || !bytes.Equal(c.tlsCert, tlsCert) || !bytes.Equal(c.tlsKey, tlsKey)) {
		if err := c.clientConn.Close(); err != nil {
			log.WithError(err).Error("Cannot close the provision client connection")
		}

		delete(p.clients, hostname)
		connect = true
	}

	if connect {
		clientConn, provClient, err := p.createClient(hostname, caCert, tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "create provision-server api client error")
		}
		c = client{
			client:     provClient,
			clientConn: clientConn,
			caCert:     caCert,
			tlsCert:    tlsCert,
			tlsKey:     tlsKey,
		}
		p.clients[hostname] = c
	}

	return c.client, nil
}

func (p *pool) createClient(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, api.ProvisionServiceClient, error) {
	logrusEntry := log.NewEntry(log.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}

	provOps := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			logging.UnaryClientCtxIDInterceptor,
		),
		grpc.WithStreamInterceptor(
			grpc_logrus.StreamClientInterceptor(logrusEntry, logrusOpts...),
		),
	}

	if len(caCert) == 0 && len(tlsCert) == 0 && len(tlsKey) == 0 {
		provOps = append(provOps, grpc.WithInsecure())
		log.WithField("server", hostname).Warning("creating insecure network-server client")
	} else {
		log.WithField("server", hostname).Info("creating network-server client")
		cert, err := tls.X509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return nil, nil, errors.Wrap(err, "load x509 keypair error")
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, nil, errors.Wrap(err, "append ca cert to pool error")
		}

		provOps = append(provOps, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	provClientConn, err := grpc.DialContext(ctx, hostname, provOps...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "dial provision-server api error")
	}

	return provClientConn, api.NewProvisionServiceClient(provClientConn), nil
}