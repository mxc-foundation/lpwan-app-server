package code

import (
	"context"
	"sort"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// MigrateGorpMigrations must be called before db migration in pgstore/db.go
func MigrateGorpMigrations(handler *store.Handler) error {
	ctx := context.Background()

	itemList, err := handler.GetAllFromGorpMigrations(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get item list of gorp_migrations table")
	}

	newItemList := migrations.AssetNames()
	// must sort out newItemList
	sort.Strings(newItemList)

	if len(itemList) > len(newItemList) {
		return errors.New("new list is shorter than existing list, migration must be done manually")
	}

	if len(itemList) <= len(newItemList) {
		for i, v := range itemList {
			if newItemList[i] == v {
				continue
			}

			if err := handler.FixGorpMigrationsItemId(ctx, v, newItemList[i]); err != nil {
				return err
			}
		}
	}

	return nil
}