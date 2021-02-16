package config

import (
	external "github.com/mxc-foundation/lpwan-app-server/internal/api/external/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
	fragmentation "github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation/data"
	multicastsetup "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/bonus"
	psconn "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn/data"
	js "github.com/mxc-foundation/lpwan-app-server/internal/codec/js/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/dhx"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	fuota "github.com/mxc-foundation/lpwan-app-server/internal/fuota/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	integration "github.com/mxc-foundation/lpwan-app-server/internal/integration/data"
	joinserver "github.com/mxc-foundation/lpwan-app-server/internal/js/data"
	as "github.com/mxc-foundation/lpwan-app-server/internal/modules/as/data"
	gws "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	metrics "github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics/data"
	mining "github.com/mxc-foundation/lpwan-app-server/internal/modules/mining/data"
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis/data"
	serverinfo "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo/data"
	monitoring "github.com/mxc-foundation/lpwan-app-server/internal/monitoring/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpapisrv"
	oidc "github.com/mxc-foundation/lpwan-app-server/internal/oidc/data"
	pprof "github.com/mxc-foundation/lpwan-app-server/internal/pprof/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

var AppserverVersion string

// Config defines the configuration structure.
type Config struct {
	General serverinfo.GeneralSettingsStruct `mapstructure:"general"`

	PostgreSQL pgstore.Config `mapstructure:"postgresql"`

	Redis rs.RedisStruct `mapstructure:"redis"`

	Operator email.Operator `mapstructure:"operator"`

	SMTP map[string]email.SMTPConfig `mapstructure:"smtp"`

	M2MServer grpccli.ConnectionOpts `mapstructure:"m2m_server"`

	DHXCenter dhx.Config `mapstructure:"dhx_center"`

	ProvisionServer psconn.ProvisioningServerStruct `mapstructure:"provision_server"`

	Recaptcha user.RecaptchaConfig `mapstructure:"recaptcha"`

	ExternalAuth user.ExternalAuthentication `mapstructure:"external_auth"`

	ShopifyConfig user.Shopify `mapstructure:"shopify"`

	ApplicationServer struct {
		ID string `mapstructure:"id"`

		UserAuthentication oidc.UserAuthenticationStruct `mapstructure:"user_authentication"`

		Codec js.CodecStruct `mapstructure:"codec"`

		Integration integration.IntegrationStruct `mapstructure:"integration"`

		API as.AppserverStruct `mapstructure:"api"`

		Airdrop bonus.Config `mapstructure:"airdrop"`

		APIForM2M mxpapisrv.Config `mapstructure:"api_for_m2m"`

		APIForGateway gws.GatewayBindStruct `mapstructure:"api_for_gateway"`

		ExternalAPI external.ExternalAPIStruct `mapstructure:"external_api"`

		RemoteMulticastSetup multicastsetup.MulticastStruct `mapstructure:"remote_multicast_setup"`

		FragmentationSession fragmentation.FragmentationStruct `mapstructure:"fragmentation_session"`

		FUOTADeployment fuota.FuotaStruct `mapstructure:"fuota_deployment"`

		MiningSetUp mining.Config `mapstructure:"mining_setup"`
	} `mapstructure:"application_server"`

	JoinServer joinserver.JoinServerStruct `mapstructure:"join_server"`

	Metrics metrics.MetricsStruct `mapstructure:"metrics"`

	PProf pprof.Config `mapstructure:"pprof"`

	Monitoring monitoring.MonitoringStruct `mapstructure:"monitoring"`
}

// C holds the global configuration.
var C Config
