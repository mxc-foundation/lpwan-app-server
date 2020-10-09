package monitoring

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	. "github.com/mxc-foundation/lpwan-app-server/internal/monitoring/data"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "monitoring"

type controller struct {
	name string
	s    MonitoringStruct
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name: moduleName,
		s:    conf.Monitoring,
	}
	return nil
}

// Setup setsup the metrics server.
func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	if ctrl.s.Bind != "" {
		return setupNew()
	}
	return setupLegacy()
}

func setupNew() error {
	if ctrl.s.Bind == "" {
		return nil
	}

	log.WithFields(log.Fields{
		"bind": ctrl.s.Bind,
	}).Info("monitoring: setting up monitoring endpoint")

	mux := http.NewServeMux()

	if ctrl.s.PrometheusAPITimingHistogram {
		log.Info("monitoring: enabling Prometheus api timing histogram")
		grpc_prometheus.EnableHandlingTimeHistogram()
	}

	if ctrl.s.PrometheusEndpoint {
		log.WithFields(log.Fields{
			"endpoint": "/metrics",
		}).Info("monitoring: registering Prometheus endpoint")
		mux.Handle("/metrics", promhttp.Handler())
	}

	if ctrl.s.HealthcheckEndpoint {
		log.WithFields(log.Fields{
			"endpoint": "/healthcheck",
		}).Info("monitoring: registering healthcheck endpoint")
		mux.HandleFunc("/health", healthCheckHandlerFunc)
	}

	server := http.Server{
		Handler: mux,
		Addr:    ctrl.s.Bind,
	}

	go func() {
		err := server.ListenAndServe()
		log.WithError(err).Error("monitoring: monitoring server error")
	}()

	return nil
}

func setupLegacy() error {
	metricsStruct := metrics.GetMetricsSettings()
	if !metricsStruct.Prometheus.EndpointEnabled {
		return nil
	}

	if metricsStruct.Prometheus.APITimingHistogram {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}

	log.WithFields(log.Fields{
		"bind": metricsStruct.Prometheus.Bind,
	}).Info("metrics: starting prometheus metrics server")

	server := http.Server{
		Handler: promhttp.Handler(),
		Addr:    metricsStruct.Prometheus.Bind,
	}

	go func() {
		err := server.ListenAndServe()
		log.WithError(err).Error("metrics: prometheus metrics server error")
	}()

	return nil
}
