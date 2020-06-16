package gateway

import (
	"context"
	"fmt"
	"strings"
	"text/template"

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

func GetDefaultGatewayConfig(ctx context.Context, gw *storage.Gateway) error {
	if strings.HasPrefix(gw.Model, "MX19") == false {
		return nil
	}

	n, err := storage.GetNetworkServer(ctx, storage.DB(), gw.NetworkServerID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get network server %d", gw.NetworkServerID)
		return errors.Wrapf(err, "GetDefaultGatewayConfig")
	}

	defaultGatewayConfig := storage.DefaultGatewayConfig{
		Model:         gw.Model,
		Region:        n.Region,
	}

	err = storage.GetDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to get default gateway config for model= %s, region= %s", defaultGatewayConfig.Model, defaultGatewayConfig.Region)
	}

	gw.Config = strings.Replace(defaultGatewayConfig.DefaultConfig, "{{ .GatewayID }}", fmt.Sprintf("%s", gw.MAC.String()), -1)
	return  nil
}
