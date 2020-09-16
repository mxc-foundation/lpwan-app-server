package config

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/api/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	joinserver "github.com/mxc-foundation/lpwan-app-server/internal/api/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup"
	mxprotocolconn "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/codec/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/fuota"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/monitoring"
	"github.com/mxc-foundation/lpwan-app-server/internal/pprof"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

var AppserverVersion string

// Config defines the configuration structure.
type Config struct {
	General serverinfo.ServerSettingsStruct `mapstructure:"general"`

	PostgreSQL storage.PostgreSQLStruct `mapstructure:"postgresql"`

	Redis rs.RedisStruct `mapstructure:"redis"`

	Operator email.OperatorStruct `mapstructure:"operator"`

	SMTP map[string]email.SMTPStruct `mapstructure:"smtp"`

	M2MServer mxprotocolconn.MxprotocolServerStruct `mapstructure:"m2m_server"`

	ProvisionServer psconn.ProvisioningServerStruct `mapstructure:"provision_server"`

	Recaptcha user.RecaptchaStruct `mapstructure:"recaptcha"`

	ApplicationServer struct {
		ID string `mapstructure:"id"`

		UserAuthentication oidc.UserAuthenticationStruct `mapstructure:"user_authentication"`

		Codec js.CodecStruct `mapstructure:"codec"`

		Integration integration.IntegrationStruct `mapstructure:"integration"`

		API as.AppserverStruct `mapstructure:"api"`

		APIForM2M m2m.M2MStruct `mapstructure:"api_for_m2m"`

		APIForGateway gws.GatewayBindStruct `mapstructure:"api_for_gateway"`

		ExternalAPI external.ExternalAPIStruct `mapstructure:"external_api"`

		RemoteMulticastSetup multicastsetup.MulticastStruct `mapstructure:"remote_multicast_setup"`

		FragmentationSession fragmentation.FragmentationStruct `mapstructure:"fragmentation_session"`

		FUOTADeployment fuota.FuotaStruct `mapstructure:"fuota_deployment"`

		MiningSetUp mining.Config `mapstructure:"mining_setup"`
	} `mapstructure:"application_server"`

	JoinServer joinserver.JoinServerStruct `mapstructure:"join_server"`

	Metrics storage.MetricsStruct `mapstructure:"metrics"`

	PProf pprof.Config `mapstructure:"pprof"`

	Monitoring monitoring.MonitoringStruct `mapstructure:"monitoring"`
}

// C holds the global configuration.
var C Config
