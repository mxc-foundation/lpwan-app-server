package organization

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Controller struct {
	St        store.OrganizationStore
	Validator Validator
}

var Service *Controller

func Setup(store store.OrganizationStore) error {
	Service.St = store
	return nil
}
