package devprovision

import (
	"context"

	"github.com/brocaar/lorawan"
	"github.com/golang/protobuf/ptypes/empty"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/ns-extra"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
)

//
type handlerMock struct {
}

func (h *handlerMock) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (gwd.Gateway, error) {
	return gwd.Gateway{Name: "MockGateway"}, nil
}

func (h *handlerMock) GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error) {
	return nsd.NetworkServer{Name: "MockNetworkServer", Server: "0.0.0.0"}, nil
}

//
type networkServerMock struct {
	request *nsextra.SendDelayedProprietaryPayloadRequest
}

// GetPool returns the networkserver pool.
func (n *networkServerMock) SendDelayedProprietaryPayload(ctx context.Context,
	in *nsextra.SendDelayedProprietaryPayloadRequest) (*empty.Empty, error) {
	n.request = in

	return nil, nil
}
