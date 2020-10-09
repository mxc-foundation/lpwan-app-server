package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

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
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/loracloud"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/marshaler"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mqtt"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/multi"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/mydevices"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/postgresql"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/thingsboard"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	. "github.com/mxc-foundation/lpwan-app-server/internal/integration/awssns/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/integration/azureservicebus/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/integration/data"
	. "github.com/mxc-foundation/lpwan-app-server/internal/integration/gcppubsub/data"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "integration"

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

type controller struct {
	h                  *store.Handler
	name               string
	mockIntegration    models.Integration
	marshalType        marshaler.Type
	globalIntegrations []models.IntegrationHandler
	s                  IntegrationStruct
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name: moduleName,
		s:    conf.ApplicationServer.Integration,
	}

	return nil
}

// Setup configures the integration package.
func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	log.Info("integration: configuring global integrations")

	var ints []models.IntegrationHandler

	// setup marshaler
	switch ctrl.s.Marshaler {
	case "protobuf":
		ctrl.marshalType = marshaler.Protobuf
	case "json":
		ctrl.marshalType = marshaler.ProtobufJSON
	case "json_v3":
		ctrl.marshalType = marshaler.JSONV3
	}

	ctrl.h = h
	// configure logger integration (for device events in web-interface)
	i, err := logger.New(logger.Config{})
	if err != nil {
		return errors.Wrap(err, "new logger integration error")
	}
	ints = append(ints, i)

	// setup global integrations, to be used by all applications
	for _, name := range ctrl.s.Enabled {
		var i models.IntegrationHandler
		var err error

		switch name {
		case "aws_sns":
			i, err = awssns.New(ctrl.marshalType, ctrl.s.AWSSNS)
		case "azure_service_bus":
			i, err = azureservicebus.New(ctrl.marshalType, ctrl.s.AzureServiceBus)
		case "mqtt":
			i, err = mqtt.New(ctrl.marshalType, ctrl.s.MQTT)
		case "gcp_pub_sub":
			i, err = gcppubsub.New(ctrl.marshalType, ctrl.s.GCPPubSub)
		case "kafka":
			i, err = kafka.New(ctrl.marshalType, ctrl.s.Kafka)
		case "postgresql":
			i, err = postgresql.New(ctrl.s.PostgreSQL)
		case "amqp":
			i, err = amqp.New(ctrl.marshalType, ctrl.s.AMQP)
		default:
			return fmt.Errorf("unknonwn integration type: %s", name)
		}

		if err != nil {
			return errors.Wrap(err, "new integration error")
		}

		ints = append(ints, i)
	}
	ctrl.globalIntegrations = ints

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
	if ctrl.mockIntegration != nil {
		return ctrl.mockIntegration
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
			i, err = http.New(ctrl.marshalType, conf)
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
			i, err = gcppubsub.New(ctrl.marshalType, conf)
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
			i, err = awssns.New(ctrl.marshalType, conf)
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
			i, err = azureservicebus.New(ctrl.marshalType, conf)
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

	return multi.New(ctrl.globalIntegrations, ints)
}

// SetMockIntegration mocks the integration.
func SetMockIntegration(i models.Integration) {
	ctrl.mockIntegration = i
}
