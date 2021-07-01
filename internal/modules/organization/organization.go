package organization

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"golang.org/x/net/context"

	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "organization"

type controller struct {
	st *store.Handler

	moduleUp bool
}

var ctrl *controller

func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{}
	return nil
}
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h

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
