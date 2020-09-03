package serverinfo

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St              *store.Handler
	SupernodeAddr   string
	DefaultLanguage string
	ServerRegion    string
}

var Service = &Controller{}

func Setup(conf config.Config, s store.Store) error {
	Service.St, _ = store.New(s)
	Service.SupernodeAddr = conf.General.ServerAddr
	Service.DefaultLanguage = conf.General.DefaultLanguage
	Service.ServerRegion = conf.General.ServerRegion
	return nil
}
