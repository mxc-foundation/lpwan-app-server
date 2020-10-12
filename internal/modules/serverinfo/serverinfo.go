package serverinfo

import (
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "server_info"

type controller struct {
	name    string
	st      *store.Handler
	general GeneralSettingsStruct

	moduleUp bool
}

var ctrl *controller

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h

	return nil
}

func GetSettings() GeneralSettingsStruct {
	return ctrl.general
}

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) (err error) {

	ctrl = &controller{
		name:    moduleName,
		general: conf.General,
	}

	return nil
}
