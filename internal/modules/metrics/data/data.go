package data

import (
	"time"
)

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
