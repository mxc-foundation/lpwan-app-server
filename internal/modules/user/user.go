package user

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        store.UserStore
	Validator Validator
}

var Service *Controller

func Setup(store store.UserStore) error {
	Service.St = store
	return nil
}
