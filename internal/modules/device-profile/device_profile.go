package devprofile

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type controller struct {
	st *store.Handler
}

var ctrl *controller

func Setup(h *store.Handler) error {
	ctrl = &controller{
		st: h,
	}

	return nil
}
