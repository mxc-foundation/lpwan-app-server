package pgstore

import (
	"fmt"
	"time"

	// register postgresql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
	. "github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore/data"
)

const moduleName = "storage"

type dbsql struct {
	*sqlx.DB
}

type controller struct {
	db dbsql
	s  PostgreSQLStruct
	c  Config
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}
	var err error

	// set up pgstore settings
	pgStoreConfig := Config{}
	pgStoreConfig.PWH, err = pwhash.New(16, conf.General.PasswordHashIterations)
	if err != nil {
		return err
	}
	if err = pgStoreConfig.ApplicationServerID.UnmarshalText([]byte(conf.ApplicationServer.ID)); err != nil {
		return errors.Wrap(err, "decode application_server.id error")
	}
	pgStoreConfig.JWTSecret = conf.ApplicationServer.ExternalAPI.JWTSecret
	pgStoreConfig.ApplicationServerPublicHost = conf.ApplicationServer.API.PublicHost

	ctrl = &controller{
		s: conf.PostgreSQL,
		c: pgStoreConfig,
	}

	return nil
}

func Setup() error {
	log.Info("storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", ctrl.s.DSN)
	if err != nil {
		return errors.Wrap(err, "storage: PostgreSQL connection error")
	}
	d.SetMaxOpenConns(ctrl.s.MaxOpenConnections)
	d.SetMaxIdleConns(ctrl.s.MaxIdleConnections)
	for {
		if err = d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	ctrl.db = dbsql{d}

	if ctrl.s.Automigrate {
		log.Info("storage: applying PostgreSQL data migrations")
		m := &migrate.AssetMigrationSource{
			Asset:    migrations.Asset,
			AssetDir: migrations.AssetDir,
			Dir:      "",
		}
		n, err := migrate.Exec(ctrl.db.DB.DB, "postgres", m, migrate.Up)
		if err != nil {
			return errors.Wrap(err, "storage: applying PostgreSQL data migrations error")
		}
		log.WithField("count", n).Info("storage: PostgreSQL data migrations applied")
	}

	return nil
}
