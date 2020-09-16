package serverinfo

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type controller struct {
	st      *store.Handler
	general ServerSettingsStruct
}

var ctrl *controller

func Setup(s store.Store) error {
	ctrl.st, _ = store.New(s)
	return nil
}

type ServerSettingsStruct struct {
	LogLevel               int    `mapstructure:"log_level"`
	LogToSyslog            bool   `mapstructure:"log_to_syslog"`
	PasswordHashIterations int    `mapstructure:"password_hash_iterations"`
	Enable2FALogin         bool   `mapstructure:"enable_2fa_login"`
	DefaultLanguage        string `mapstructure:"defualt_language"`
	ServerAddr             string `mapstructure:"server_addr"`
	ServerRegion           string `mapstructure:"server_region"`
}

// SettingsSetup init settings extracted values from toml file then assign each modules
func SettingsSetup(conf config.Config) error {
	ctrl = &controller{
		general: conf.General,
	}

	if err := storage.SettingsSetup(storage.SettingStruct{
		Db:                  conf.PostgreSQL,
		Metrics:             conf.Metrics,
		JWTSecret:           conf.ApplicationServer.ExternalAPI.JWTSecret,
		ApplicationServerID: conf.ApplicationServer.ID,
	}); err != nil {
		return err
	}

	if err := redis.SettingsSetup(conf.Redis); err != nil {
		return err
	}

	if err := email.SettingsSetup(conf.SMTP, conf.Operator); err != nil {
		return err
	}

	if err := m2mcli.SettingsSetup(conf.M2MServer); err != nil {
		return err
	}

	if err := gws.SettingsSetup(conf.ApplicationServer.APIForGateway); err != nil {
		return err
	}

	if err := psconn.SettingsSetup(conf.ProvisionServer); err != nil {
		return err
	}

	if err := user.SettingsSetup(conf.Recaptcha); err != nil {
		return err
	}

	if err := external.SettingsSetup(conf.ApplicationServer.ExternalAPI, conf.ApplicationServer.ID); err != nil {
		return err
	}

	if err := oidc.SettingsSetup(conf.ApplicationServer.UserAuthentication); err != nil {
		return err
	}

	return nil
}

func GetSettings() ServerSettingsStruct {
	return ctrl.general
}
