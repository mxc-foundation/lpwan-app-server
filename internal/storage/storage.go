package storage

import (
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	uuid "github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
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

type SettingStruct struct {
	JWTSecret           string
	ApplicationServerID string
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name:           moduleName,
		HashIterations: 100000,
		jwtsecret:      []byte(conf.ApplicationServer.ExternalAPI.JWTSecret),
	}

	if err := ctrl.applicationServerID.UnmarshalText([]byte(conf.ApplicationServer.ID)); err != nil {
		return errors.Wrap(err, "decode application_server.id error")
	}

	if err := pgstore.SettingsSetup(conf); err != nil {
		return err
	}
	return nil
}

// Setup configures the storage package.
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

	// redis client is set up in metrics setup function
	if err := metrics.Setup(); err != nil {
		return err
	}

	if err := pgstore.Setup(); err != nil {
		return err
	}

	return nil
}
