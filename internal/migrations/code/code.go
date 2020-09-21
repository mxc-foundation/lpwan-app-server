package code

import (
	"context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

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
