package code

import (
	"context"
	"sort"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	pgerr "github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

// MigrateGorpMigrations must be called before db migration in pgstore/db.go
func MigrateGorpMigrations(handler *store.Handler) error {
	ctx := context.Background()
	itemList, err := handler.GetAllFromGorpMigrations(ctx)
	if err != nil {
		if err == pgerr.ErrEmptyGorpMigrations {
			log.Info("gorp_migrations table does not exist, no need to run migrate_gorp_migrations")
			return nil
		}

		return errors.Wrap(err, "failed to get item list of gorp_migrations table")
	}

	newItemList := migrations.AssetNames()
	// must sort out newItemList and itemList
	sort.Strings(newItemList)
	sort.Strings(itemList)

	if len(itemList) > len(newItemList) {
		return errors.New("new list is shorter than existing list, migration must be done manually")
	}

	if len(itemList) <= len(newItemList) {
		for i, v := range itemList {
			if newItemList[i] == v {
				continue
			}

			if err := handler.FixGorpMigrationsItemID(ctx, v, newItemList[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
