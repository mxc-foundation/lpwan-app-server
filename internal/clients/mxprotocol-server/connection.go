package mxprotocolconn

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
	"google.golang.org/grpc/credentials"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
)

type controller struct {
	p                Pool
	mxprotocolServer MxprotocolServerStruct
}

type Pool interface {
	get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error)
}

type MxprotocolServerStruct struct {
	Server  string `mapstructure:"m2m_server"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

type mxprotocolServiceClient struct {
	clientConn *grpc.ClientConn
	caCert     []byte
	tlsCert    []byte
	tlsKey     []byte
}

type pool struct {
	sync.RWMutex
	mxprotocolServiceClients map[string]mxprotocolServiceClient
}

var ctrl *controller

func SettingsSetup(s MxprotocolServerStruct) error {
	ctrl = &controller{
		mxprotocolServer: s,
	}

	return nil
}

func Setup() error {
	ctrl = &controller{
		mxprotocolServer: MxprotocolServerStruct{
			Server:  ctrl.mxprotocolServer.Server,
			CACert:  ctrl.mxprotocolServer.CACert,
			TLSCert: ctrl.mxprotocolServer.TLSCert,
			TLSKey:  ctrl.mxprotocolServer.TLSKey,
		},
		p: &pool{
			mxprotocolServiceClients: make(map[string]mxprotocolServiceClient),
		},
	}

	return nil
}

// Get returns a M2MServerServiceClient for the given server (hostname:ip).
func (p *pool) get(hostname string, caCert, tlsCert, tlsKey []byte) (*grpc.ClientConn, error) {
	defer p.Unlock()
	p.Lock()

	var connect bool
	c, ok := p.mxprotocolServiceClients[hostname]
	if !ok {
		connect = true
	}

	// if the connection exists in the map, but when the certificates changed
	// try to cloe the connection and re-connect
	if ok && (!bytes.Equal(c.caCert, caCert) || !bytes.Equal(c.tlsCert, tlsCert) || !bytes.Equal(c.tlsKey, tlsKey)) {
		c.clientConn.Close()
		delete(p.mxprotocolServiceClients, hostname)
		connect = true
	}

	if connect {
		clientConn, err := p.createClient(hostname, caCert, tlsCert, tlsKey)
		if err != nil {
			return nil, errors.Wrap(err, "create m2m-server api client error")
		}
		c = mxprotocolServiceClient{
			clientConn: clientConn,
			caCert:     caCert,
			tlsCert:    tlsCert,
			tlsKey:     tlsKey,
		}
		p.mxprotocolServiceClients[hostname] = c
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

func GetM2MDeviceServiceClient() (pb.DSDeviceServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewDSDeviceServiceClient(m2mServerClient), nil
}

func GetM2MGatewayServiceClient() (pb.GSGatewayServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewGSGatewayServiceClient(m2mServerClient), nil
}

func GetMiningServiceClient() (pb.MiningServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewMiningServiceClient(m2mServerClient), nil
}

func GetServerServiceClient() (pb.M2MServerInfoServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewM2MServerInfoServiceClient(m2mServerClient), nil
}

func GetSettingsServiceClient() (pb.SettingsServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewSettingsServiceClient(m2mServerClient), nil
}

func GetStakingServiceClient() (pb.StakingServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewStakingServiceClient(m2mServerClient), nil
}

func GetTopupServiceClient() (pb.TopUpServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewTopUpServiceClient(m2mServerClient), nil
}

func GetWalletServiceClient() (pb.WalletServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewWalletServiceClient(m2mServerClient), nil
}

func GetWithdrawServiceClient() (pb.WithdrawServiceClient, error) {
	m2mServerClient, err := ctrl.p.get(ctrl.mxprotocolServer.Server, []byte(ctrl.mxprotocolServer.CACert),
		[]byte(ctrl.mxprotocolServer.TLSCert), []byte(ctrl.mxprotocolServer.TLSKey))
	if err != nil {
		return nil, err
	}

	return pb.NewWithdrawServiceClient(m2mServerClient), nil
}
