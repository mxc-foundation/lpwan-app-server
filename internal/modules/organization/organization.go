package organization

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St *store.Handler
}

var Service = &Controller{}

func Setup(h *store.Handler) error {
	Service.St = h
	return nil
}
