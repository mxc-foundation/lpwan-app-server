package gatewayprofile

import (
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

func GetGatewayProfileCount(ctx context.Context) (int, error) {
	return ctrl.st.GetGatewayProfileCount(ctx)
}

func GetGatewayProfiles(ctx context.Context, limit, offset int) ([]store.GatewayProfileMeta, error) {
	return ctrl.st.GetGatewayProfiles(ctx, limit, offset)
}
