package gateway

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var configTemplate = &template.Template{}

func LoadTemplates() error {
	// load gateway config templates
	configTemplate = template.Must(template.New("gateway-config/global_conf.json.sx1250.MX.CN490").Parse(
		string(static.MustAsset("gateway-config/global_conf.json.sx1250.MX.CN490"))))

	return nil
}

func GetDefaultGatewayConfig(ctx context.Context, gw *storage.Gateway, networkServerID int64) error {
	if strings.HasPrefix(gw.Model, "MX19") == false {
		return nil
	}

	n, err := storage.GetNetworkServer(ctx, storage.DB(), networkServerID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get network server %d", networkServerID)
		return errors.Wrapf(err, "GetDefaultGatewayConfig")
	}

	var region string

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err == nil {
		resp, err := nsClient.GetVersion(ctx, &empty.Empty{})
		if err == nil {
			region = resp.Region.String()
		}
	}

	defaultGatewayConfig := storage.DefaultGatewayConfig{
		Model:         gw.Model,
		Region:        region,
	}

	err = storage.GetDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to get default gateway config for model= %s, region= %s", defaultGatewayConfig.Model, defaultGatewayConfig.Region)
	}

	gw.Config = strings.Replace(defaultGatewayConfig.DefaultConfig, "{{ .GatewayID }}", fmt.Sprintf("%s", gw.MAC.String()), -1)
	return  nil
}
