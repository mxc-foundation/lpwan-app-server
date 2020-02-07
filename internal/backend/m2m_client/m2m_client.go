package m2m_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"sync"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

var p Pool

// Pool defines the m2m server client pool.
// Actually this connection pool is not necessary now for m2m client, because there is only one m2m server.
type Pool interface {
	Get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error)
}

type m2mServiceClient struct {
	clientConn *grpc.ClientConn
	caCert     []byte
	tlsCert    []byte
	tlsKey     []byte
}

// Setup configures m2mServiceClients
func Setup(conf config.Config) error {
	p = &pool{
		m2mServiceClients: make(map[string]m2mServiceClient),
	}
	return nil
}

// GetPool returns p value
func GetPool() Pool {
	return p
}

// SetPool sets p to pp
func SetPool(pp Pool) {
	p = pp
}

type pool struct {
	sync.RWMutex
	m2mServiceClients map[string]m2mServiceClient
}

// Get returns a M2MServerServiceClient for the given server (hostname:ip).
func (p *pool) Get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error) {
	defer p.Unlock()
	p.Lock()

	var connect bool
	c, ok := p.m2mServiceClients[hostname]
	if !ok {
		connect = true
	}

	// if the connection exists in the map, but when the certificates changed
	// try to cloe the connection and re-connect
	if ok && (!bytes.Equal(c.caCert, caCert) || !bytes.Equal(c.tlsCert, tlsCert) || !bytes.Equal(c.tlsKey, tlsKey)) {
		c.clientConn.Close()
		delete(p.m2mServiceClients, hostname)
		connect = true
	}

	if connect {
		clientConn, err := p.createClient(hostname, caCert, tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "create m2m-server api client error")
		}
		c = m2mServiceClient{
			clientConn: clientConn,
			caCert:     caCert,
			tlsCert:    tlsCert,
			tlsKey:     tlsKey,
		}
		p.m2mServiceClients[hostname] = c
	}

	return c.clientConn, nil
}

func (p *pool) createClient(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error) {
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

	if len(caCert) == 0 && len(tlsCert) == 0 && len(tlsKey) == 0 {
		nsOpts = append(nsOpts, grpc.WithInsecure())
		log.WithField("server", hostname).Warning("creating insecure m2m-server device service client")
	} else {
		log.WithField("server", hostname).Info("creating m2m-server device service client")
		cert, err := tls.X509KeyPair(tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "load x509 keypair error")
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.Wrap(err, "append ca cert to pool error")
		}

		nsOpts = append(nsOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	devServiceClient, err := grpc.DialContext(ctx, hostname, nsOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial m2m-server device service api error")
	}

	return devServiceClient, nil
}
