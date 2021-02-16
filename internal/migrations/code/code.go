package code

import (
	"context"

	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	pgerr "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Setup checks migration status and runs migration scripts right after db is ready
func Setup(h *store.Handler, autoMigrate bool) error {
	if autoMigrate {

		if err := MigrateGorpMigrations(h); err != nil {
			log.Fatal(err, " fix gorp_migrations table error")
		}

		log.Info("applying PostgreSQL data migrations")
		m := &migrate.AssetMigrationSource{
			Asset:    migrations.Asset,
			AssetDir: migrations.AssetDir,
			Dir:      "",
		}

		if err := h.ExecuteMigrateUp(m); err != nil {
			log.Fatal(err, "faile to migrate postgresql data")
		}
	}

	if err := networkserver.Setup(); err != nil {
		return err
	}

	if err := Migrate("migrate_gw_stats", h, MigrateGatewayStats); err != nil {
		log.Fatal(errors.Wrap(err, "migration error"))
	}

	if err := Migrate("migrate_to_cluster_keys", h, func(handler *store.Handler) error {
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
