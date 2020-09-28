package device

import (
	"github.com/brocaar/lorawan"
	"golang.org/x/net/context"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

type Config struct {
	ApplicationServerID string
}
type controller struct {
	st *store.Handler
	s  Config
}

var ctrl *controller

func SettingsSetup(s Config) error {
	ctrl = &controller{
		s: s,
	}
	return nil
}

func Setup(h *store.Handler) error {
	ctrl.st = h
	return nil
}

func GetDevice(ctx context.Context, devEUI lorawan.EUI64, forUpdate bool) (store.Device, error) {
	return ctrl.st.GetDevice(ctx, devEUI, forUpdate)
}
