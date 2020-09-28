package monitoring

import (
	"net/http"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type MonitoringStruct struct {
	Bind                         string `mapstructure:"bind"`
	PrometheusEndpoint           bool   `mapstructure:"prometheus_endpoint"`
	PrometheusAPITimingHistogram bool   `mapstructure:"prometheus_api_timing_histogram"`
	HealthcheckEndpoint          bool   `mapstructure:"healthcheck_endpoint"`
}
type controller struct {
	s MonitoringStruct
}

var ctrl *controller

func SettingsSerup(s MonitoringStruct) error {
	ctrl = &controller{
		s: s,
	}
	return nil
}

// Setup setsup the metrics server.
func Setup() error {
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
	metricsStruct := storage.GetMetricsSettings()
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
