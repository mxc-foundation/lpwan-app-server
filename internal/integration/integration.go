package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/amqp"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/awssns"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/azureservicebus"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/gcppubsub"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/http"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/influxdb"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/kafka"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/logger"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/loracloud"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/marshaler"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mqtt"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/multi"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mydevices"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/postgresql"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/thingsboard"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// Handler kinds
const (
	HTTP            = "HTTP"
	InfluxDB        = "INFLUXDB"
	ThingsBoard     = "THINGSBOARD"
	MyDevices       = "MYDEVICES"
	LoRaCloud       = "LORACLOUD"
	GCPPubSub       = "GCP_PUBSUB"
	AWSSNS          = "AWS_SNS"
	AzureServiceBus = "AZURE_SERVICE_BUS"
)

var (
	mockIntegration    models.Integration
	marshalType        marshaler.Type
	globalIntegrations []models.IntegrationHandler
)

// AzurePublishMode defines the publish-mode type.
type AzurePublishMode string

// Publish modes.
const (
	AzurePublishModeTopic AzurePublishMode = "topic"
	AzurePublishModeQueue AzurePublishMode = "queue"
)

// IntegrationAWSSNSConfig holds the AWS SNS integration configuration.
type IntegrationAWSSNSConfig struct {
	Marshaler          string `mapstructure:"marshaler" json:"marshaler"`
	AWSRegion          string `mapstructure:"aws_region" json:"region"`
	AWSAccessKeyID     string `mapstructure:"aws_access_key_id" json:"accessKeyID"`
	AWSSecretAccessKey string `mapstructure:"aws_secret_access_key" json:"secretAccessKey"`
	TopicARN           string `mapstructure:"topic_arn" json:"topicARN"`
}

// IntegrationAzureConfig holds the Azure Service-Bus integration configuration.
type IntegrationAzureConfig struct {
	Marshaler        string           `mapstructure:"marshaler" json:"marshaler"`
	ConnectionString string           `mapstructure:"connection_string" json:"connectionString"`
	PublishMode      AzurePublishMode `mapstructure:"publish_mode" json:"-"`
	PublishName      string           `mapstructure:"publish_name" json:"publishName"`
}

// IntegrationGCPConfig holds the GCP Pub/Sub integration configuration.
type IntegrationGCPConfig struct {
	Marshaler            string `mapstructure:"marshaler" json:"marshaler"`
	CredentialsFile      string `mapstructure:"credentials_file" json:"-"`
	CredentialsFileBytes []byte `mapstructure:"-" json:"credentialsFile"`
	ProjectID            string `mapstructure:"project_id" json:"projectID"`
	TopicName            string `mapstructure:"topic_name" json:"topicName"`
}

// IntegrationPostgreSQLConfig holds the PostgreSQL integration configuration.
type IntegrationPostgreSQLConfig struct {
	DSN                string `json:"dsn"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}

// IntegrationAMQPConfig holds the AMQP integration configuration.
type IntegrationAMQPConfig struct {
	URL                     string `mapstructure:"url"`
	EventRoutingKeyTemplate string `mapstructure:"event_routing_key_template"`
}

// IntegrationKafkaConfig holds the Kafka integration configuration.
type IntegrationKafkaConfig struct {
	Brokers          []string `mapstructure:"brokers"`
	TLS              bool     `mapstructure:"tls"`
	Topic            string   `mapstructure:"topic"`
	EventKeyTemplate string   `mapstructure:"event_key_template"`
	Username         string   `mapstructure:"username"`
	Password         string   `mapstructure:"password"`
}

// IntegrationMQTTConfig holds the configuration for the MQTT integration.
type IntegrationMQTTConfig struct {
	Server               string        `mapstructure:"server"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	MaxReconnectInterval time.Duration `mapstructure:"max_reconnect_interval"`
	QOS                  uint8         `mapstructure:"qos"`
	CleanSession         bool          `mapstructure:"clean_session"`
	ClientID             string        `mapstructure:"client_id"`
	CACert               string        `mapstructure:"ca_cert"`
	TLSCert              string        `mapstructure:"tls_cert"`
	TLSKey               string        `mapstructure:"tls_key"`
	EventTopicTemplate   string        `mapstructure:"event_topic_template"`
	CommandTopicTemplate string        `mapstructure:"command_topic_template"`
	RetainEvents         bool          `mapstructure:"retain_events"`

	// For backards compatibility
	UplinkTopicTemplate        string `mapstructure:"uplink_topic_template"`
	DownlinkTopicTemplate      string `mapstructure:"downlink_topic_template"`
	JoinTopicTemplate          string `mapstructure:"join_topic_template"`
	AckTopicTemplate           string `mapstructure:"ack_topic_template"`
	ErrorTopicTemplate         string `mapstructure:"error_topic_template"`
	StatusTopicTemplate        string `mapstructure:"status_topic_template"`
	LocationTopicTemplate      string `mapstructure:"location_topic_template"`
	TxAckTopicTemplate         string `mapstructure:"tx_ack_topic_template"`
	IntegrationTopicTemplate   string `mapstructure:"integration_topic_template"`
	UplinkRetainedMessage      bool   `mapstructure:"uplink_retained_message"`
	JoinRetainedMessage        bool   `mapstructure:"join_retained_message"`
	AckRetainedMessage         bool   `mapstructure:"ack_retained_message"`
	ErrorRetainedMessage       bool   `mapstructure:"error_retained_message"`
	StatusRetainedMessage      bool   `mapstructure:"status_retained_message"`
	LocationRetainedMessage    bool   `mapstructure:"location_retained_message"`
	TxAckRetainedMessage       bool   `mapstructure:"tx_ack_retained_message"`
	IntegrationRetainedMessage bool   `mapstructure:"integration_retained_message"`
}

type IntegrationStruct struct {
	Marshaler       string                      `mapstructure:"marshaler"`
	Backend         string                      `mapstructure:"backend"` // deprecated
	Enabled         []string                    `mapstructure:"enabled"`
	AWSSNS          IntegrationAWSSNSConfig     `mapstructure:"aws_sns"`
	AzureServiceBus IntegrationAzureConfig      `mapstructure:"azure_service_bus"`
	MQTT            IntegrationMQTTConfig       `mapstructure:"mqtt"`
	GCPPubSub       IntegrationGCPConfig        `mapstructure:"gcp_pub_sub"`
	Kafka           IntegrationKafkaConfig      `mapstructure:"kafka"`
	PostgreSQL      IntegrationPostgreSQLConfig `mapstructure:"postgresql"`
	AMQP            IntegrationAMQPConfig       `mapstructure:"amqp"`
}

// Setup configures the integration package.
func Setup(conf config.Config) error {
	log.Info("integration: configuring global integrations")

	var ints []models.IntegrationHandler

	// setup marshaler
	switch conf.ApplicationServer.Integration.Marshaler {
	case "protobuf":
		marshalType = marshaler.Protobuf
	case "json":
		marshalType = marshaler.ProtobufJSON
	case "json_v3":
		marshalType = marshaler.JSONV3
	}

	// configure logger integration (for device events in web-interface)
	i, err := logger.New(logger.Config{})
	if err != nil {
		return errors.Wrap(err, "new logger integration error")
	}
	ints = append(ints, i)

	// setup global integrations, to be used by all applications
	for _, name := range conf.ApplicationServer.Integration.Enabled {
		var i models.IntegrationHandler
		var err error

		switch name {
		case "aws_sns":
			i, err = awssns.New(marshalType, conf.ApplicationServer.Integration.AWSSNS)
		case "azure_service_bus":
			i, err = azureservicebus.New(marshalType, conf.ApplicationServer.Integration.AzureServiceBus)
		case "mqtt":
			i, err = mqtt.New(marshalType, conf.ApplicationServer.Integration.MQTT)
		case "gcp_pub_sub":
			i, err = gcppubsub.New(marshalType, conf.ApplicationServer.Integration.GCPPubSub)
		case "kafka":
			i, err = kafka.New(marshalType, conf.ApplicationServer.Integration.Kafka)
		case "postgresql":
			i, err = postgresql.New(conf.ApplicationServer.Integration.PostgreSQL)
		case "amqp":
			i, err = amqp.New(marshalType, conf.ApplicationServer.Integration.AMQP)
		default:
			return fmt.Errorf("unknonwn integration type: %s", name)
		}

		if err != nil {
			return errors.Wrap(err, "new integration error")
		}

		ints = append(ints, i)
	}
	globalIntegrations = ints

	return nil
}

// ForApplicationID returns the integration handler for the given application ID.
// The returned handler will be a "multi-handler", containing both the global
// integrations and the integrations setup specifically for the given
// application ID.
// When the given application ID equals 0, only the global integrations are
// returned.
func ForApplicationID(id int64) models.Integration {
	// for testing, return mock integration
	if mockIntegration != nil {
		return mockIntegration
	}

	var appints []storage.Integration
	var err error

	// retrieve application integrations when ID != 0
	if id != 0 {
		appints, err = storage.GetIntegrationsForApplicationID(context.TODO(), storage.DB(), id)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"application_id": id,
			}).Error("integrations: get application integrations error")
		}
	}

	// parse integration configs and setup integrations
	var ints []models.IntegrationHandler
	for _, appint := range appints {
		var i models.IntegrationHandler
		var err error

		switch appint.Kind {
		case HTTP:
			// read config
			var conf http.Config
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read http configuration error")
				continue
			}

			// create new http integration
			i, err = http.New(marshalType, conf)
		case InfluxDB:
			// read config
			var conf influxdb.Config
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read influxdb configuration error")
				continue
			}

			// create new influxdb integration
			i, err = influxdb.New(conf)
		case ThingsBoard:
			// read config
			var conf thingsboard.Config
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read thingsboard configuration error")
				continue
			}

			// create new thingsboard integration
			i, err = thingsboard.New(conf)
		case MyDevices:
			// read config
			var conf mydevices.Config
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read mydevices configuration error")
				continue
			}

			// create new mydevices integration
			i, err = mydevices.New(conf)
		case LoRaCloud:
			// read config
			var conf loracloud.Config
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read loracloud configuration error")
				continue
			}

			// create new loracloud integration
			i, err = loracloud.New(conf)
		case GCPPubSub:
			// read config
			var conf IntegrationGCPConfig
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read gcp pubsub configuration error")
				continue
			}

			// create new gcp pubsub integration
			i, err = gcppubsub.New(marshalType, conf)
		case AWSSNS:
			// read config
			var conf IntegrationAWSSNSConfig
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read aws sns configuration error")
				continue
			}

			// create new aws sns integration
			i, err = awssns.New(marshalType, conf)
		case AzureServiceBus:
			// read config
			var conf IntegrationAzureConfig
			if err := json.NewDecoder(bytes.NewReader(appint.Settings)).Decode(&conf); err != nil {
				log.WithError(err).WithFields(log.Fields{
					"application_id": id,
				}).Error("integrtations: read azure service-bus configuration error")
				continue
			}

			// create new aws sns integration
			i, err = azureservicebus.New(marshalType, conf)
		default:
			log.WithFields(log.Fields{
				"application_id": id,
				"kind":           appint.Kind,
			}).Error("integrations: unknown integration type")
			continue
		}

		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"application_id": id,
				"kind":           appint.Kind,
			}).Error("integrations: new integration error")
			continue
		}

		ints = append(ints, i)
	}

	return multi.New(globalIntegrations, ints)
}

// SetMockIntegration mocks the integration.
func SetMockIntegration(i models.Integration) {
	mockIntegration = i
}
