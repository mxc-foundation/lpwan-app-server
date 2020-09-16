package storage

import (
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
)

type controller struct {
	Db        PostgreSQLStruct
	Metrics   MetricsStruct
	jwtsecret []byte
	// HashIterations denfines the number of times a password is hashed.
	HashIterations      int
	applicationServerID uuid.UUID
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

type PostgreSQLStruct struct {
	DSN                string `mapstructure:"dsn"`
	Automigrate        bool
	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
}

type SettingStruct struct {
	Db                  PostgreSQLStruct
	Metrics             MetricsStruct
	JWTSecret           string
	ApplicationServerID string
}

// GetMetricsSettings :
func GetMetricsSettings() MetricsStruct {
	return ctrl.Metrics
}

var ctrl *controller

func SettingsSetup(s SettingStruct) error {
	ctrl = &controller{
		Db:             s.Db,
		Metrics:        s.Metrics,
		HashIterations: 100000,
		jwtsecret:      []byte(s.JWTSecret),
	}

	if err := ctrl.applicationServerID.UnmarshalText([]byte(s.ApplicationServerID)); err != nil {
		return errors.Wrap(err, "decode application_server.id error")
	}

	return nil
}

// Setup configures the storage package.
func Setup() error {

	log.Info("storage: setup metrics")
	// setup aggregation intervals
	var intervals []AggregationInterval
	for _, agg := range ctrl.Metrics.Redis.AggregationIntervals {
		intervals = append(intervals, AggregationInterval(strings.ToUpper(agg)))
	}
	if err := SetAggregationIntervals(intervals); err != nil {
		return errors.Wrap(err, "set aggregation intervals error")
	}

	// setup timezone
	if err := SetTimeLocation(ctrl.Metrics.Timezone); err != nil {
		return errors.Wrap(err, "set time location error")
	}

	// setup storage TTL
	SetMetricsTTL(
		ctrl.Metrics.Redis.MinuteAggregationTTL,
		ctrl.Metrics.Redis.HourAggregationTTL,
		ctrl.Metrics.Redis.DayAggregationTTL,
		ctrl.Metrics.Redis.MonthAggregationTTL,
	)

	if err := rs.SetupRedis(); err != nil {
		return errors.Wrap(err, "set up redis error")
	}

	log.Info("storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", ctrl.Db.DSN)
	if err != nil {
		return errors.Wrap(err, "storage: PostgreSQL connection error")
	}
	d.SetMaxOpenConns(ctrl.Db.MaxOpenConnections)
	d.SetMaxIdleConns(ctrl.Db.MaxIdleConnections)
	for {
		if err := d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	db = &DBLogger{d}

	if ctrl.Db.Automigrate {
		log.Info("storage: applying PostgreSQL data migrations")
		m := &migrate.AssetMigrationSource{
			Asset:    migrations.Asset,
			AssetDir: migrations.AssetDir,
			Dir:      "",
		}
		n, err := migrate.Exec(db.DB.DB, "postgres", m, migrate.Up)
		if err != nil {
			return errors.Wrap(err, "storage: applying PostgreSQL data migrations error")
		}
		log.WithField("count", n).Info("storage: PostgreSQL data migrations applied")
	}

	return nil
}
