package pgstore

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/errors"
)

type Action errors.Action

// Possible actions
const (
	Select = Action(errors.Select)
	Insert = Action(errors.Insert)
	Update = Action(errors.Update)
	Delete = Action(errors.Delete)
	Scan   = Action(errors.Scan)
)

func handlePSQLError(action Action, err error, description string) error {
	return errors.HandlePSQLError(errors.Action(action), err, description)
}
