package pgstore

import (
	"time"

	// register postgresql driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
	. "github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore/data"
)

type controller struct {
	db *sqlx.DB
	s  PostgreSQLStruct
	c  Config
}

var ctrl *controller

func SettingsSetup(conf config.Config) error {
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

	ctrl.db = d

	return nil
}
