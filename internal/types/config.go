package types

import (
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	amqp "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	awssns "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	azureservicebus "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	gcppubsub "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	kafka "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	mqtt "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
	postgresql "github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
)

type GeneralSettingsStruct struct {
	LogLevel               int    `mapstructure:"log_level"`
	PasswordHashIterations int    `mapstructure:"password_hash_iterations"`
	Enable2FALogin         bool   `mapstructure:"enable_2fa_login"`
	DefaultLanguage        string `mapstructure:"defualt_language"`
	ServerAddr             string `mapstructure:"server_addr"`
	ServerRegion           string `mapstructure:"server_region"`
	EnableSTC              bool   `mapstructure:"enable_stc"`
}

// Config contains postgres configuration
type DBConfig struct {
	DSN                string `mapstructure:"dsn"`
	Automigrate        bool
	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
}

type RedisStruct struct {
	URL        string   `mapstructure:"url"` // deprecated
	Servers    []string `mapstructure:"servers"`
	Cluster    bool     `mapstructure:"cluster"`
	MasterName string   `mapstructure:"master_name"`
	PoolSize   int      `mapstructure:"pool_size"`
	Password   string   `mapstructure:"password"`
	Database   int      `mapstructure:"database"`
}

// Operator defines basic settings of operator of this supernode
type Operator struct {
	Operator           string `mapstructure:"name"`
	PrimaryColor       string `mapstructure:"primary_color"`
	SecondaryColor     string `mapstructure:"secondary_color"`
	DownloadAppStore   string `mapstructure:"download_appstore"`
	DownloadGoogle     string `mapstructure:"download_google"`
	DownloadTestFlight string `mapstructure:"download_testflight"`
	DownloadAPK        string `mapstructure:"download_apk"`
	OperatorAddress    string `mapstructure:"operator_address"`
	OperatorLegal      string `mapstructure:"operator_legal_name"`
	OperatorLogo       string `mapstructure:"operator_logo"`
	OperatorContact    string `mapstructure:"operator_contact"`
	OperatorSupport    string `mapstructure:"operator_support"`
}

// SMTPConfig defines smtp service settings
type SMTPConfig struct {
	Email       string `mapstructure:"email"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	AuthType    string `mapstructure:"auth_type"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	TLSRequired bool   `mapstructure:"tls_required"`
}

// ProvisioningServerStruct defines credentails to connect to provisioning-server
type ProvisioningServerStruct struct {
	ServerConifig  grpccli.ConnectionOpts `mapstructure:"grpc_connection"`
	UpdateSchedule string                 `mapstructure:"update_schedule"`
}

type RecaptchaConfig struct {
	HostServer string `mapstructure:"host_server"`
	Secret     string `mapstructure:"secret"`
}

// ExternalAuthentication defines configuration for external_auth section
type ExternalAuthentication struct {
	WechatAuth      auth.WeChatAuthentication `mapstructure:"wechat_auth"`
	DebugWechatAuth auth.WeChatAuthentication `mapstructure:"debug_wechat_auth"`
}

// ShopifyAdminAPI defines shopify admin api configuration
type ShopifyAdminAPI struct {
	Hostname   string `mapstructure:"hostname"`
	APIKey     string `mapstructure:"api_key"`
	Secret     string `mapstructure:"secret"`
	APIVersion string `mapstructure:"api_version"`
	StoreName  string `mapstructure:"store_name"`
}

// BonusSettings defines settings of shopify related bonus
type BonusSettings struct {
	Enable    bool  `mapstructure:"enable"`
	ValueUSD  int64 `mapstructure:"value_usd"`
	ProductID int64 `mapstructure:"product_id"`
}

// Shopify defines full shopify service settings
type Shopify struct {
	AdminAPI ShopifyAdminAPI `mapstructure:"shopify_admin_api"`
	Bonus    BonusSettings   `mapstructure:"bonus"`
}

type UserAuthenticationStruct struct {
	OpenIDConnect struct {
		Enabled                 bool   `mapstructure:"enabled"`
		RegistrationEnabled     bool   `mapstructure:"registration_enabled"`
		RegistrationCallbackURL string `mapstructure:"registration_callback_url"`
		ProviderURL             string `mapstructure:"provider_url"`
		ClientID                string `mapstructure:"client_id"`
		ClientSecret            string `mapstructure:"client_secret"`
		RedirectURL             string `mapstructure:"redirect_url"`
		LogoutURL               string `mapstructure:"logout_url"`
		LoginLabel              string `mapstructure:"login_label"`
	} `mapstructure:"openid_connect"`
}

type CodecStruct struct {
	JS struct {
		MaxExecutionTime time.Duration `mapstructure:"max_execution_time"`
	} `mapstructure:"js"`
}

type IntegrationStruct struct {
	Marshaler       string                                 `mapstructure:"marshaler"`
	Backend         string                                 `mapstructure:"backend"` // deprecated
	Enabled         []string                               `mapstructure:"enabled"`
	AWSSNS          awssns.IntegrationAWSSNSConfig         `mapstructure:"aws_sns"`
	AzureServiceBus azureservicebus.IntegrationAzureConfig `mapstructure:"azure_service_bus"`
	MQTT            mqtt.IntegrationMQTTConfig             `mapstructure:"mqtt"`
	GCPPubSub       gcppubsub.IntegrationGCPConfig         `mapstructure:"gcp_pub_sub"`
	Kafka           kafka.IntegrationKafkaConfig           `mapstructure:"kafka"`
	PostgreSQL      postgresql.IntegrationPostgreSQLConfig `mapstructure:"postgresql"`
	AMQP            amqp.IntegrationAMQPConfig             `mapstructure:"amqp"`
}

type AppserverStruct struct {
	Bind       string `mapstructure:"bind"`
	CACert     string `mapstructure:"ca_cert"`
	TLSCert    string `mapstructure:"tls_cert"`
	TLSKey     string `mapstructure:"tls_key"`
	PublicHost string `mapstructure:"public_host"`
}

// Config contains configuration of the service
type BonusConfig struct {
	// URL of the server that provides the list of bonuses and credentials
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	// How often to retrieve and pay the bonuses
	CheckInterval int64 `mapstructure:"check_interval_sec"`
	// the identificator for this supernode used by remote side
	SNID string `mapstructure:"supernode_id"`
}

// Config contains configuration for gRPC server serving mxp server
type MXPServerConfig struct {
	Bind    string `mapstructure:"bind"`
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`
}

type GatewayBindStruct struct {
	NewGateway struct {
		Bind    string `mapstructure:"new_gateway_bind"`
		CACert  string `mapstructure:"ecc_ca_cert"`
		TLSCert string `mapstructure:"ecc_tls_cert"`
		TLSKey  string `mapstructure:"ecc_tls_key"`
	} `mapstructure:"new_gateway"`

	OldGateway struct {
		Bind    string `mapstructure:"old_gateway_bind"`
		CACert  string `mapstructure:"rsa_ca_cert"`
		TLSCert string `mapstructure:"rsa_tls_cert"`
		TLSKey  string `mapstructure:"rsa_tls_key"`
	} `mapstructure:"old_gateway"`
}

type ExternalAPIStruct struct {
	Bind            string `mapstructure:"bind"`
	TLSCert         string `mapstructure:"tls_cert"`
	TLSKey          string `mapstructure:"tls_key"`
	JWTSecret       string `mapstructure:"jwt_secret"`
	JWTDefaultTTL   int64  `mapstructure:"jwt_default_ttl_sec"`
	OTPSecret       string `mapstructure:"otp_secret"`
	CORSAllowOrigin string `mapstructure:"cors_allow_origin"`
}

type MulticastStruct struct {
	SyncInterval  time.Duration `mapstructure:"sync_interval"`
	SyncRetries   int           `mapstructure:"sync_retries"`
	SyncBatchSize int           `mapstructure:"sync_batch_size"`
}

type FragmentationStruct struct {
	SyncInterval  time.Duration `mapstructure:"sync_interval"`
	SyncRetries   int           `mapstructure:"sync_retries"`
	SyncBatchSize int           `mapstructure:"sync_batch_size"`
}

type FuotaStruct struct {
	McGroupID int `mapstructure:"mc_group_id"`
	FragIndex int `mapstructure:"frag_index"`
}

// MiningConfig contains mining configuration
type MiningConfig struct {
	// If mining is enabled or not
	Enabled bool `mapstructure:"enabled"`
	// If we haven't got heartbeat for HeartbeatOfflineLimit seconds, we
	// consider gateway to be offline
	HeartbeatOfflineLimit int64 `mapstructure:"heartbeat_offline_limit"`
	// Gateway must be online for at leasts GwOnlineLimit seconds to receive mining reward
	GwOnlineLimit int64 `mapstructure:"gw_online_limit"`
	// Period is the length of the mining period in seconds
	Period int64 `mapstructure:"period"`
}
type JoinServerStruct struct {
	Bind    string
	CACert  string `mapstructure:"ca_cert"`
	TLSCert string `mapstructure:"tls_cert"`
	TLSKey  string `mapstructure:"tls_key"`

	KEK struct {
		ASKEKLabel string `mapstructure:"as_kek_label"`

		Set []struct {
			Label string `mapstructure:"label"`
			KEK   string `mapstructure:"kek"`
		}
	} `mapstructure:"kek"`
}
type MetricsStruct struct {
	Timezone string `mapstructure:"timezone"`
	Redis    struct {
		AggregationIntervals []string      `mapstructure:"aggregation_intervals"`
		MinuteAggregationTTL time.Duration `mapstructure:"minute_aggregation_ttl"`
		HourAggregationTTL   time.Duration `mapstructure:"hour_aggregation_ttl"`
		DayAggregationTTL    time.Duration `mapstructure:"day_aggregation_ttl"`
		MonthAggregationTTL  time.Duration `mapstructure:"month_aggregation_ttl"`
	} `mapstructure:"redis"`
	Prometheus struct {
		EndpointEnabled    bool   `mapstructure:"endpoint_enabled"`
		Bind               string `mapstructure:"bind"`
		APITimingHistogram bool   `mapstructure:"api_timing_histogram"`
	} `mapstructure:"prometheus"`
}

type PProfConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Bind    string `mapstructure:"bind"`
}

type MonitoringStruct struct {
	Bind                         string `mapstructure:"bind"`
	PrometheusEndpoint           bool   `mapstructure:"prometheus_endpoint"`
	PrometheusAPITimingHistogram bool   `mapstructure:"prometheus_api_timing_histogram"`
	HealthcheckEndpoint          bool   `mapstructure:"healthcheck_endpoint"`
}
