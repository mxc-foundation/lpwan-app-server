package code

import (
	"context"
	"fmt"
	"log"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "migrations"

func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
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
			return err
		}

		return f(hander)
	})
}
