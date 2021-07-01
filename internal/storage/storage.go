package storage

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "storage"

type controller struct {
	name      string
	jwtsecret []byte
	// HashIterations denfines the number of times a password is hashed.
	HashIterations      int
	applicationServerID uuid.UUID

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {

	ctrl = &controller{
		name:           moduleName,
		HashIterations: 100000,
		jwtsecret:      []byte(conf.ApplicationServer.ExternalAPI.JWTSecret),
	}

	if err := ctrl.applicationServerID.UnmarshalText([]byte(conf.ApplicationServer.ID)); err != nil {
		return errors.Wrap(err, "decode application_server.id error")
	}

	return nil
}

// Setup configures the storage package.
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	// redis client is set up in metrics setup function
	if err := metrics.Setup(); err != nil {
		return err
	}

	return nil
}
