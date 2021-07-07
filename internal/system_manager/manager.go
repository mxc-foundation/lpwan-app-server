package system_manager

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

type SettingsSetupFunc func(string, config.Config) error
type ModuleSetupFunc func(string, *store.Handler) error

func RegisterSettingsSetup(name string, f SettingsSetupFunc) {
	if settingsSetupCallbackFunc[name] != nil {
		log.Fatalf("settingsSetupFunc for %s already registered", name)
	}

	settingsSetupCallbackFunc[name] = &f
}

func RegisterModuleSetup(name string, f ModuleSetupFunc) {
	if moduleSetupCallbackFunc[name] != nil {
		log.Fatalf("moduleSetupFunc for %s already registered", name)
	}

	moduleSetupCallbackFunc[name] = &f
}

var settingsSetupCallbackFunc = map[string]*SettingsSetupFunc{}
var moduleSetupCallbackFunc = map[string]*ModuleSetupFunc{}

func SetupSystemSettings(conf config.Config) error {
	for n, f := range settingsSetupCallbackFunc {
		log.Infof("Set up %s settings...", n)
		if err := (*f)(n, conf); err != nil {
			return err
		}
	}

	return nil
}

func SetupSystemModules() error {
	if _, ok := moduleSetupCallbackFunc["storage"]; !ok {
		return errors.New(fmt.Sprintf("setup function is not found for %s", "storage"))
	}

	f := moduleSetupCallbackFunc["storage"]
	if err := (*f)("storage", nil); err != nil {
		return err
	}

	for n, f := range moduleSetupCallbackFunc {
		if err := (*f)(n, store.NewStore()); err != nil {
			return err
		}
	}
	return nil
}
