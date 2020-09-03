package application

import (
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St *store.Handler
}

var Service = &Controller{}

func Setup(s store.Store) error {
	st, _ := store.New(s)
	Service.St = st

	return nil
}

func (c *Controller) GetApplication(ctx context.Context, applicationID int64) (store.Application, error) {
	return c.St.GetApplication(ctx, applicationID)
}
