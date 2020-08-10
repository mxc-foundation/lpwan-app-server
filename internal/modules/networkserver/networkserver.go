package networkserver

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        store.NetworkServerStore
	Validator Validator
}

var Service *Controller

func Setup(store store.NetworkServerStore) error {
	Service.St = store
	return nil
}
