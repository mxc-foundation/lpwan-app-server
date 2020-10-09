package fuotamod

import (
	"errors"
	"fmt"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "fuota_deployment"

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

	return nil
}
