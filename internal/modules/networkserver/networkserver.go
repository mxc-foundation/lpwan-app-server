package networkserver

import (
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"golang.org/x/net/context"
)

type controller struct {
	st *store.Handler
}

var ctrl *controller

func Setup(h *store.Handler) error {
	ctrl = &controller{
		st: h,
	}

	return nil
}

// GetNetworkServerForDevEUI :
func GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (store.NetworkServer, error) {
	return ctrl.st.GetNetworkServerForDevEUI(ctx, devEUI)
}

// GetNetworkServer :
func GetNetworkServer(ctx context.Context, id int64) (store.NetworkServer, error) {
	return ctrl.st.GetNetworkServer(ctx, id)
}

// GetNetworkServerForGatewayProfileID :
func GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (store.NetworkServer, error) {
	return ctrl.st.GetNetworkServerForGatewayProfileID(ctx, id)
}

// GetNetworkServerForGatewayMAC :
func GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (store.NetworkServer, error) {
	return ctrl.st.GetNetworkServerForGatewayMAC(ctx, mac)
}
