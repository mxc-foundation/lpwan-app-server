package types

import (
	"time"
)

// IntegrationStruct holds the integration package configuration
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

// IntegrationAMQPConfig holds the AMQP integration configuration.
type IntegrationAMQPConfig struct {
	URL                     string `mapstructure:"url"`
	EventRoutingKeyTemplate string `mapstructure:"event_routing_key_template"`
}

// IntegrationAWSSNSConfig holds the AWS SNS integration configuration.
type IntegrationAWSSNSConfig struct {
	Marshaler          string `mapstructure:"marshaler" json:"marshaler"`
	AWSRegion          string `mapstructure:"aws_region" json:"region"`
	AWSAccessKeyID     string `mapstructure:"aws_access_key_id" json:"accessKeyID"`
	AWSSecretAccessKey string `mapstructure:"aws_secret_access_key" json:"secretAccessKey"`
	TopicARN           string `mapstructure:"topic_arn" json:"topicARN"`
}

// AzurePublishMode defines the publish-mode type.
type AzurePublishMode string

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
	RetainEvents         bool          `mapstructure:"retain_events"`
}

// IntegrationPostgreSQLConfig holds the PostgreSQL integration configuration.
type IntegrationPostgreSQLConfig struct {
	DSN                string `json:"dsn"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections"`
}
