package gateway

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/brocaar/lorawan"
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

func GetDefaultGatewayConfig(ctx context.Context, mac lorawan.EUI64, server string, networkServerID int64) (string, error) {
	n, err := storage.GetNetworkServer(ctx, storage.DB(), networkServerID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get network server %d", networkServerID)
		return "", errors.Wrapf(err, "GetDefaultGatewayConfig")
	}

	var region string
	var version string

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err == nil {
		resp, err := nsClient.GetVersion(ctx, &empty.Empty{})
		if err == nil {
			region = resp.Region.String()
			version = resp.Version
		}
	}

	if strings.HasPrefix(region, "CN") == false {
		log.Warnf("Gateway registered on network region=%s, version=%s doesn not have default config yet", region, version)
		return "", nil
	}

	var defaultGwConfig bytes.Buffer
	if err := configTemplate.Execute(&defaultGwConfig, struct {
		GatewayID, ServerAddr string
	}{
		GatewayID:  mac.String(),
		ServerAddr: server,
	}); err != nil {
		log.WithError(err).Error("Failed to execute gateway config template")
		return "", errors.Wrapf(err, "GetDefaultGatewayConfig for mac=%s", mac.String())
	}

	return defaultGwConfig.String(), nil
}
