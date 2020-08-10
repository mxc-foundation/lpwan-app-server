package gatewayprofile

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        store.GatewayProfileStore
	Validator Validator
}

var Service *Controller

func Setup(store store.GatewayProfileStore) error {
	Service.St = store
	return nil
}
