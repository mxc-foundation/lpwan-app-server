package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/integration/amqp"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/awssns"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/azureservicebus"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/gcppubsub"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/http"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/influxdb"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/kafka"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/logger"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/marshaler"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mqtt"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/multi"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/postgresql"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/types"
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

var marshalType marshaler.Type

// SetupGlobalIntegrations configures global integration and return instance
func SetupGlobalIntegrations(config types.IntegrationStruct) ([]models.IntegrationHandler, error) {
	log.Info("integration: configuring global integrations")
	var ints []models.IntegrationHandler

	// setup marshaler
	switch config.Marshaler {
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
		return nil, errors.Wrap(err, "new logger integration error")
	}
	ints = append(ints, i)

	// setup global integrations, to be used by all applications
	for _, name := range config.Enabled {
		var i models.IntegrationHandler
		var err error

		switch name {
		case "aws_sns":
			i, err = awssns.New(marshalType, config.AWSSNS)
		case "azure_service_bus":
			i, err = azureservicebus.New(marshalType, config.AzureServiceBus)
		case "mqtt":
			i, err = mqtt.New(marshalType, config.MQTT)
		case "gcp_pub_sub":
			i, err = gcppubsub.New(marshalType, config.GCPPubSub)
		case "kafka":
			i, err = kafka.New(marshalType, config.Kafka)
		case "postgresql":
			i, err = postgresql.New(config.PostgreSQL)
		case "amqp":
			i, err = amqp.New(marshalType, config.AMQP)
		default:
			return nil, fmt.Errorf("unknonwn integration type: %s", name)
		}

		if err != nil {
			return nil, errors.Wrap(err, "new integration error")
		}

		ints = append(ints, i)
	}
	return ints, nil
}

// ForApplicationID returns the integration handler for the given application ID.
// The returned handler will be a "multi-handler", containing both the global
// integrations and the integrations setup specifically for the given
// application ID.
// When the given application ID equals 0, only the global integrations are
// returned.
func ForApplicationID(id int64, gIntegrations []models.IntegrationHandler) models.Integration {
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

	return multi.New(gIntegrations, ints)
}
