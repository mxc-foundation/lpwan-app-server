package user

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

type RecaptchaStruct struct {
	HostServer string `mapstructure:"host_server"`
	Secret     string `mapstructure:"secret"`
}

type controller struct {
	s   RecaptchaStruct
	st  *store.Handler
	pwh *pwhash.PasswordHasher
}

var ctrl *controller

func SettingsSetup(s RecaptchaStruct) error {
	ctrl = &controller{
		s: s,
	}

	return nil
}
func GetSettings() RecaptchaStruct {
	return ctrl.s
}

func Setup(s store.Store) (err error) {
	ctrl.st, _ = store.New(s)

	ctrl.pwh, err = pwhash.New(16, serverinfo.GetSettings().PasswordHashIterations)
	if err != nil {
		return err
	}

	return nil
}

func SetUserPassword(user *store.User, pw string) error {
	pwHash, err := ctrl.pwh.HashPassword(pw)
	if err != nil {
		return err
	}

	user.PasswordHash = pwHash
	return nil
}

func VerifyUserPassword(pw string, pwHash string) error {
	return ctrl.pwh.Validate(pw, pwHash)
}
