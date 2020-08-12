package user

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        *store.Handler
	Validator Validator
}

var Service = &Controller{}

func Setup(s store.Store, conf config.Config) (err error) {
	Service.St, _ = store.New(s)

	return nil
}
