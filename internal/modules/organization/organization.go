package organization

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

// GetOrganizationCount :
func GetOrganizationCount(ctx context.Context, filters store.OrganizationFilters) (int, error) {
	return ctrl.st.GetOrganizationCount(ctx, filters)
}

// GetOrganizationIDList :
func GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error) {
	return ctrl.st.GetOrganizationIDList(ctx, limit, offset, search)
}

// GetOrganizationUsers :
func GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]store.OrganizationUser, error) {
	return ctrl.st.GetOrganizationUsers(ctx, organizationID, limit, offset)
}
