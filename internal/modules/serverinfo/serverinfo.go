package serverinfo

import (
	"errors"
	"fmt"

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
}

var ctrl *controller

func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl.st = h
	return nil
}

func GetSettings() GeneralSettingsStruct {
	return ctrl.general
}

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) (err error) {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name:    moduleName,
		general: conf.General,
	}

	return nil
}
