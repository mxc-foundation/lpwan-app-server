package networkserver_portal

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
)

// Pool defines the network-server client pool.
type Pool interface {
	get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error)
}

type NSStruct struct {
	Server  string
	CACert  string
	TLSCert string
	TLSKey  string
}

type nsClient struct {
	clientConn *grpc.ClientConn
	caCert     []byte
	tlsCert    []byte
	tlsKey     []byte
}

type pool struct {
	sync.RWMutex
	nsClients map[string]nsClient
}

// Get returns a NetworkServerClient for the given server (hostname:ip).
func (p *pool) get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error) {
	defer p.Unlock()
	p.Lock()

	var connect bool
	c, ok := p.nsClients[hostname]
	if !ok {
		connect = true
	}

	// if the connection exists in the map, but when the certificates changed
	// try to cloe the connection and re-connect
	if ok && (!bytes.Equal(c.caCert, caCert) || !bytes.Equal(c.tlsCert, tlsCert) || !bytes.Equal(c.tlsKey, tlsKey)) {
		c.clientConn.Close()
		delete(p.nsClients, hostname)
		connect = true
	}

	if connect {
		clientConn, err := p.createClient(hostname, caCert, tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "create network-server api client error")
		}
		c = nsClient{
			clientConn: clientConn,
			caCert:     caCert,
			tlsCert:    tlsCert,
			tlsKey:     tlsKey,
		}
		p.nsClients[hostname] = c
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
			logging.UnaryClientCtxIDInterceptor,
		),
		grpc.WithStreamInterceptor(
			grpc_logrus.StreamClientInterceptor(logrusEntry, logrusOpts...),
		),
		grpc.WithBalancerName(roundrobin.Name),
	}

	if len(caCert) == 0 && len(tlsCert) == 0 && len(tlsKey) == 0 {
		nsOpts = append(nsOpts, grpc.WithInsecure())
		log.WithField("server", hostname).Warning("creating insecure network-server client")
	} else {
		log.WithField("server", hostname).Info("creating network-server client")
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nsClient, err := grpc.DialContext(ctx, hostname, nsOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial network-server api error")
	}

	return nsClient, nil
}

func (s *NSStruct) GetNetworkServiceClient() (ns.NetworkServerServiceClient, error) {
	nsconn, err := ctrl.p.get(s.Server, []byte(s.CACert), []byte(s.TLSCert), []byte(s.TLSKey))
	if err != nil {
		return nil, err
	}

	return ns.NewNetworkServerServiceClient(nsconn), nil
}
