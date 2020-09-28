package application

import (
	"golang.org/x/net/context"

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
	ctrl.st = h

	return nil
}

func GetApplication(ctx context.Context, applicationID int64) (store.Application, error) {
	return ctrl.st.GetApplication(ctx, applicationID)
}
