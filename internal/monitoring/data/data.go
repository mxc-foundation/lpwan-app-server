package data

type MonitoringStruct struct {
	Bind                         string `mapstructure:"bind"`
	PrometheusEndpoint           bool   `mapstructure:"prometheus_endpoint"`
	PrometheusAPITimingHistogram bool   `mapstructure:"prometheus_api_timing_histogram"`
	HealthcheckEndpoint          bool   `mapstructure:"healthcheck_endpoint"`
}
