package pgstore

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Action store.Action

// Possible actions
const (
	Select = Action(store.Select)
	Insert = Action(store.Insert)
	Update = Action(store.Update)
	Delete = Action(store.Delete)
	Scan   = Action(store.Scan)
)

func handlePSQLError(action Action, err error, description string) error {
	return store.HandlePSQLError(store.Action(action), err, description)
}
