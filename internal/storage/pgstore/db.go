package pgstore

import (
	"time"

	// register postgresql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Config contains postgres configuration
type Config struct {
	DSN                string `mapstructure:"dsn"`
	Automigrate        bool
	MaxOpenConnections int `mapstructure:"max_open_connections"`
	MaxIdleConnections int `mapstructure:"max_idle_connections"`
}

type controller struct {
	db *sqlx.DB
}

var ctrl *controller

// Setup establishes connection with postgresql server and returns PgStore
// object.
func Setup(cfg Config) (*PgStore, error) {
	log.Info("storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "storage: PostgreSQL connection error")
	}
	d.SetMaxOpenConns(cfg.MaxOpenConnections)
	d.SetMaxIdleConns(cfg.MaxIdleConnections)
	for {
		if err = d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	ctrl = &controller{db: d}

	return &PgStore{db: d}, nil
}
