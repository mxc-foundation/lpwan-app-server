package nscli

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/brocaar/chirpstack-api/go/v3/ns"

	nsextra "github.com/mxc-foundation/lpwan-app-server/api/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

// Client represents network server client
type Client struct {
	nsConn map[int64]*grpc.ClientConn
}

// NetworkServerConfig defines data structure for creating ns clients
type NetworkServerConfig struct {
	NetworkServerID int64
	ConnOptions     grpccli.ConnectionOpts
}

// Connect connects to given network server config and add connection to pool
func Connect(nscfg []NetworkServerConfig) (*Client, error) {
	var cli Client

	if cli.nsConn == nil {
		cli.nsConn = make(map[int64]*grpc.ClientConn)
	}

	for _, v := range nscfg {
		nsConn, err := grpccli.Connect(v.ConnOptions)
		if err != nil {
			return nil, fmt.Errorf("couldn't create network server client: %v", err)
		}
		cli.nsConn[v.NetworkServerID] = nsConn
	}
	return nil, nil
}

// Close closes connection to network server
func (cli *Client) Close() error {
	for _, v := range cli.nsConn {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}

// GetNumberOfNetworkServerClients returns number of network server clients, used for set default network server
func (cli *Client) GetNumberOfNetworkServerClients() int {
	return len(cli.nsConn)
}

// GetNetworkServerServiceClient returns a new NetworkServerServiceClient instance
func (cli *Client) GetNetworkServerServiceClient(networkServerID int64) (ns.NetworkServerServiceClient, error) {
	if v, ok := cli.nsConn[networkServerID]; !ok {
		return nil, fmt.Errorf("no such connection for network server id= %d", networkServerID)
	} else {
		return ns.NewNetworkServerServiceClient(v), nil
	}
}

// GetNetworkServerExtraServiceClient returns a new NetworkServerExtraServiceClient instance
func (cli *Client) GetNetworkServerExtraServiceClient(networkServerID int64) (nsextra.NetworkServerExtraServiceClient, error) {
	if v, ok := cli.nsConn[networkServerID]; !ok {
		return nil, fmt.Errorf("no such connection for network server id= %d", networkServerID)
	} else {
		return nsextra.NewNetworkServerExtraServiceClient(v), nil
	}
}
