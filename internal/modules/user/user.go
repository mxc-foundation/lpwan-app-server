package user

import (
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "user"

type controller struct {
	s  Config
	st *store.Handler

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{
		s: Config{
			Recaptcha:      conf.Recaptcha,
			Enable2FALogin: conf.General.Enable2FALogin,
		},
	}

	return nil
}

func Setup(name string, h *store.Handler) (err error) {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h

	return nil
}
