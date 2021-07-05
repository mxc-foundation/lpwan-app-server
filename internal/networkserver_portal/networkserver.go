package networkserver_portal

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "network_server"

type controller struct {
	st                          *store.Handler
	applicationServerID         uuid.UUID
	p                           Pool
	applicationServerPublicHost string

	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {

	appServerID, err := uuid.FromString(conf.ApplicationServer.ID)
	if err != nil {
		return errors.Wrap(err, "failed to convert applicationserver id from string to uuid")
	}

	ctrl = &controller{
		applicationServerID:         appServerID,
		applicationServerPublicHost: conf.ApplicationServer.API.PublicHost,
	}

	return nil
}

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h
	ctrl.p = &pool{
		nsClients: make(map[string]nsClient),
	}

	return nil
}
