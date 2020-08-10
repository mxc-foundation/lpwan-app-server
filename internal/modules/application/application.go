package application

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"golang.org/x/net/context"
)

type Controller struct {
	St        store.ApplicationStore
	Validator Validator
}

var Service *Controller

func Setup(store store.ApplicationStore) error {
	Service.St = store
	return nil
}

func (c *Controller) GetApplication(ctx context.Context, applicationID int64) (Application, error) {
	return c.St.GetApplication(ctx, applicationID)
}
