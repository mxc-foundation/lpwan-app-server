package organization

import (
	"fmt"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "organization"

type controller struct {
	st *store.Handler
}

var ctrl *controller

func Setup(name string, h *store.Handler) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		st: h,
	}

	return nil
}

// GetOrganizationCount :
func GetOrganizationCount(ctx context.Context, filters OrganizationFilters) (int, error) {
	return ctrl.st.GetOrganizationCount(ctx, filters)
}

// GetOrganizationIDList :
func GetOrganizationIDList(ctx context.Context, limit, offset int, search string) ([]int, error) {
	return ctrl.st.GetOrganizationIDList(ctx, limit, offset, search)
}

// GetOrganizationUsers :
func GetOrganizationUsers(ctx context.Context, organizationID int64, limit, offset int) ([]OrganizationUser, error) {
	return ctrl.st.GetOrganizationUsers(ctx, organizationID, limit, offset)
}
