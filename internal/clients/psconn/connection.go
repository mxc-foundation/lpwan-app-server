package psconn

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
)

type controller struct {
	provisioningServer ProvisioningServerStruct
}

type ProvisioningServerStruct struct {
	Server         string `mapstructure:"provision_server"`
	CACert         string `mapstructure:"ca_cert"`
	TLSCert        string `mapstructure:"tls_cert"`
	TLSKey         string `mapstructure:"tls_key"`
	UpdateSchedule string `mapstructure:"update_schedule"`
}

var ctrl *controller

func SettingsSetup(s ProvisioningServerStruct) error {
	ctrl = &controller{
		provisioningServer: s,
	}

	return nil
}
func GetSettings() ProvisioningServerStruct {
	return ctrl.provisioningServer
}

func Setup() error {
	return nil
}

func CreateClientWithCert() (api.ProvisionClient, error) {
	var grpcClient *grpc.ClientConn
	var opts []grpc.DialOption

	if ctrl.provisioningServer.Server == "" ||
		ctrl.provisioningServer.CACert == "" ||
		ctrl.provisioningServer.TLSCert == "" ||
		ctrl.provisioningServer.TLSKey == "" {
		return nil, errors.New("invalid credentials")
	}

	cert, err := tls.LoadX509KeyPair(ctrl.provisioningServer.TLSCert, ctrl.provisioningServer.TLSKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls key-pair error")
	}

	var caCertPool *x509.CertPool
	rawCaCert, err := ioutil.ReadFile(filepath.Join(filepath.Clean(ctrl.provisioningServer.CACert)))
	if err != nil {
		return nil, errors.Wrap(err, "load ca certificate error")
	}

	caCertPool = x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(rawCaCert) {
		return nil, fmt.Errorf("append ca certificate error: %s", ctrl.provisioningServer.CACert)
	}

	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})))

	grpcClient, err = grpc.Dial(ctrl.provisioningServer.Server, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "dial error")
	}

	client := api.NewProvisionClient(grpcClient)

	return client, nil
}
