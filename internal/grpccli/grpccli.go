package grpccli

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ConnectionOpts struct {
	Server  string `mapstructure:"server"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

func Connect(co ConnectionOpts) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	var tlsConf *tls.Config
	if co.CACert != "" {
		rawCACert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(co.CACert)))
		if err != nil {
			return nil, fmt.Errorf("couldn't load ca cert: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(rawCACert) {
			return nil, fmt.Errorf("no certificates were added from %s", co.CACert)
		}
		tlsConf = &tls.Config{
			RootCAs: caCertPool,
		}

		if co.TLSCert != "" {
			cert, err := tls.LoadX509KeyPair(co.TLSCert, co.TLSKey)
			if err != nil {
				return nil, fmt.Errorf("couldn't load client certificate and key: %v", err)
			}
			tlsConf.Certificates = []tls.Certificate{cert}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(co.Server, opts...)
}
