package setdefault

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// Setup add default network server and related configurations if none exists
func Setup(ctx context.Context, st *store.Handler, applicationServerID uuid.UUID,
	applicationPublicHost string, nsCli *nscli.Client) error {
	if 0 == nsCli.GetNumberOfNetworkServerClients() {
		// create default network server
		if err := ns.CreateNetworkServer(ctx, &nsd.NetworkServer{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      "default_network_server",
			Server:    "network-server:8000",
		}, st, nsCli, applicationServerID, applicationPublicHost); err != nil {
			return err
		}
	}
	return nil
}
