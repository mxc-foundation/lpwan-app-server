package code

import (
	"context"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

// Migrate checks if the given function code has been applied and if not
// it will execute the given function.
func Migrate(name string, f func(handler *store.Handler) error) error {
	return storage.Transaction(func(ctx context.Context, handler *store.Handler) error {
		_, err := tx.Exec(`lock table code_migration`)
		if err != nil {
			return errors.Wrap(err, "lock code migration table error")
		}

		res, err := tx.Exec(`
			insert into code_migration (
				id,
				applied_at
			) values ($1, $2)
			on conflict
				do nothing
		`, name, time.Now())
		if err != nil {
			switch err := err.(type) {
			case *pq.Error:
				switch err.Code.Name() {
				case "unique_violation":
					return nil
				}
			}

			return err
		}

		ra, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if ra == 0 {
			return nil
		}

		return f(tx)
	})
}
