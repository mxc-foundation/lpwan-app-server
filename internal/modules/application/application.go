package application

import (
	"errors"
	"fmt"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"golang.org/x/net/context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "application"

type controller struct {
	st *store.Handler
}

var ctrl *controller

func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		st: h,
	}
	ctrl.st = h

	return nil
}

func GetApplication(ctx context.Context, applicationID int64) (Application, error) {
	return ctrl.st.GetApplication(ctx, applicationID)
}
