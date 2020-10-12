package application

import (
	"errors"
	"fmt"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "application"

type controller struct {
	st *store.Handler

	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{}
	return nil
}
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl.st = h

	return nil
}

func GetApplication(ctx context.Context, applicationID int64) (Application, error) {
	return ctrl.st.GetApplication(ctx, applicationID)
}
