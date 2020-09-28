package user

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type RecaptchaStruct struct {
	HostServer string `mapstructure:"host_server"`
	Secret     string `mapstructure:"secret"`
}

type Config struct {
	Recaptcha      RecaptchaStruct
	Enable2FALogin bool
}

type controller struct {
	s  Config
	st *store.Handler
}

var ctrl *controller

func SettingsSetup(s Config) error {
	ctrl = &controller{
		s: s,
	}

	return nil
}

func Setup(h *store.Handler) (err error) {
	ctrl.st = h

	return nil
}
