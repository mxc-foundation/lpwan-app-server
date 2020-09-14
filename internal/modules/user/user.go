package user

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
)

type Controller struct {
	St  *store.Handler
	pwh *pwhash.PasswordHasher
}

var Service = &Controller{}

func Setup(s store.Store) (err error) {
	Service.St, _ = store.New(s)

	Service.pwh, err = pwhash.New(16, config.C.General.PasswordHashIterations)
	if err != nil {
		return err
	}

	return nil
}

func SetUserPassword(user *store.User, pw string) error {
	pwHash, err := Service.pwh.HashPassword(pw)
	if err != nil {
		return err
	}

	user.PasswordHash = pwHash
	return nil
}

func VerifyUserPassword(pw string, pwHash string) error {
	return Service.pwh.Validate(pw, pwHash)
}
