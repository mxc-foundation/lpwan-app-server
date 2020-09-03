package gatewayprofile

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St *store.Handler
}

var Service = &Controller{}

func Setup(s store.Store) error {
	Service.St, _ = store.New(s)
	return nil
}
