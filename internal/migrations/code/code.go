package code

import (
	"context"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	pgerr "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "migrations"
const (
	migrateGatewayStats  = "migrate_gw_stats"
	migrateToClusterKeys = "migrate_to_cluster_keys"
)

type controller struct {
	name        string
	autoMigrate bool

	st       *store.Handler
	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{
		name:        moduleName,
		autoMigrate: conf.PostgreSQL.Automigrate,
	}

	return nil
}

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h
	if ctrl.autoMigrate {

		if err := MigrateGorpMigrations(ctrl.st); err != nil {
			log.Fatal(err, " fix gorp_migrations table error")
		}

		log.Info("applying PostgreSQL data migrations")
		m := &migrate.AssetMigrationSource{
			Asset:    migrations.Asset,
			AssetDir: migrations.AssetDir,
			Dir:      "",
		}

		if err := ctrl.st.ExecuteMigrateUp(m); err != nil {
			log.Fatal(err, "faile to migrate postgresql data")
		}
	}

	if err := networkserver.Setup(); err != nil {
		return err
	}

	if err := Migrate(migrateGatewayStats, ctrl.st, MigrateGatewayStats); err != nil {
		log.Fatal(errors.Wrap(err, "migration error"))
	}

	if err := Migrate(migrateToClusterKeys, ctrl.st, func(handler *store.Handler) error {
		return MigrateToClusterKeys(config.C)
	}); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Migrate checks if the given function code has been applied and if not
// it will execute the given function.
func Migrate(name string, hander *store.Handler, f func(handler *store.Handler) error) error {
	return hander.Tx(context.Background(), func(ctx context.Context, handler *store.Handler) error {
		err := handler.Migrate(ctx, name)
		if err != nil {
			if err == pgerr.ErrAlreadyExists {
				return nil
			}
			return err
		}

		return f(hander)
	})
}
