package provisionserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"path/filepath"
)

func CreateClientWithCert(hostname, CACert, TLSCert, TLSKey string) (api.ProvisionClient, error) {
	var grpcClient *grpc.ClientConn
	var opts []grpc.DialOption

	if hostname == "" || CACert == "" || TLSCert == "" || TLSKey == "" {
		return nil, errors.New("invalid credentials")
	}

	cert, err := tls.LoadX509KeyPair(TLSCert, TLSKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls key-pair error")
	}

	var caCertPool *x509.CertPool
	rawCaCert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(CACert)))
	if err != nil {
		return nil, errors.Wrap(err, "load ca certificate error")
	}

	caCertPool = x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(rawCaCert) {
		return nil, fmt.Errorf("append ca certificate error: %s", CACert)
	}

	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})))

	grpcClient, err = grpc.Dial(hostname, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial error")
	}

	client := api.NewProvisionClient(grpcClient)

	return client, nil
}
