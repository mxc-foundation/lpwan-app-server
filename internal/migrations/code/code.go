package code

import (
	"context"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "migrations"

type controller struct {
	name        string
	autoMigrate bool

	st       *store.Handler
	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

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

	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl.st = h
	if ctrl.autoMigrate {
		if err := Migrate("migrate_gorp_migrations", ctrl.st, MigrateGorpMigrations); err != nil {
			log.Fatal(err, "fix gorp_migrations table error")
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

	if err := Migrate("migrate_gw_stats", ctrl.st, MigrateGatewayStats); err != nil {
		log.Fatal(errors.Wrap(err, "migration error"))
	}

	if err := Migrate("migrate_to_cluster_keys", ctrl.st, func(handler *store.Handler) error {
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
			return err
		}

		return f(hander)
	})
}
