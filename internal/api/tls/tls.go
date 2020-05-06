package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"path/filepath"
)

func getTransportCredentials(caCert, tlsCert, tlsKey string, verifyClientCert bool) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls key-pair error")
	}

	var caCertPool *x509.CertPool
	if caCert != "" {
		rawCaCert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(caCert)))
		if err != nil {
			return nil, errors.Wrap(err, "load ca certificate error")
		}

		caCertPool = x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(rawCaCert) {
			return nil, fmt.Errorf("append ca certificate error: %s", caCert)
		}
	}

	if verifyClientCert {
		return credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}), nil
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}), nil
}

func serverOptions() []grpc.ServerOption {
	logrusEntry := log.NewEntry(log.StandardLogger())
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}

	return []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
			grpc_prometheus.StreamServerInterceptor,
		),
	}
}

// NewServerWithTLSCredentials :
func NewServerWithTLSCredentials(service, caCert, tlsCert, tlsKey string) (*grpc.Server, error) {
	opts := serverOptions()

	if caCert != "" && tlsCert != "" && tlsKey != "" {
		creds, err := getTransportCredentials(caCert, tlsCert, tlsKey, true)
		if err != nil {
			return nil, errors.Wrap(err, "get transport credentials error")
		}

		opts = append(opts, grpc.Creds(creds))
	} else {
		log.WithField("service name", service).Warn("Transport credentials not set, create new server with insecure connection.")
	}

	s := grpc.NewServer(opts...)

	return s, nil
}
