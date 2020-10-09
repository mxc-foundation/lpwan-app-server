package storage

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/jmoiron/sqlx"

	// register postgresql driver
	_ "github.com/lib/pq"
)

// db holds the PostgreSQL connection pool.
var db *DBLogger

// DBLogger is a DB wrapper which logs the executed sql queries and their
// duration.
type DBLogger struct {
	*sqlx.DB
}

// DB returns the PostgreSQL database object.
func DB() *store.Handler {
	return store.NewStore()
}

func DBTest() *DBLogger {
	return db
}

// Transaction wraps the given function in a transaction. In case the given
// functions returns an error, the transaction will be rolled back.
func Transaction(f func(ctx context.Context, handler *store.Handler) error) error {
	if err := DB().Tx(context.TODO(), func(ctx context.Context, handler *store.Handler) error {
		if err := f(ctx, handler); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
