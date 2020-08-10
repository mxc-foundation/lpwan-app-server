package device

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        store.DeviceStore
	Validator Validator
}

var Service *Controller

func Setup(store store.DeviceStore) error {
	Service.St = store
	return nil
}
