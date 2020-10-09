package data

import (
	amqp "github.com/mxc-foundation/lpwan-app-server/internal/integration/amqp/data"
	awssns "github.com/mxc-foundation/lpwan-app-server/internal/integration/awssns/data"
	azureservicebus "github.com/mxc-foundation/lpwan-app-server/internal/integration/azureservicebus/data"
	gcppubsub "github.com/mxc-foundation/lpwan-app-server/internal/integration/gcppubsub/data"
	kafka "github.com/mxc-foundation/lpwan-app-server/internal/integration/kafka/data"
	mqtt "github.com/mxc-foundation/lpwan-app-server/internal/integration/mqtt/data"
	postgresql "github.com/mxc-foundation/lpwan-app-server/internal/integration/postgresql/data"
)

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
