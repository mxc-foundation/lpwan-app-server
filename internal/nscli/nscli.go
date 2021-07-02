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
func (c *Client) Connect(nscfg []NetworkServerConfig) error {
	if c.nsConn == nil {
		c.nsConn = make(map[int64]*grpc.ClientConn)
	}

	for _, v := range nscfg {
		if c.nsConn[v.NetworkServerID] != nil {
			// skip if this network server already exists
			continue
		}
		nsConn, err := grpccli.Connect(v.ConnOptions)
		if err != nil {
			return fmt.Errorf("couldn't create network server client: %v", err)
		}
		c.nsConn[v.NetworkServerID] = nsConn
	}
	return nil
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
	v, ok := cli.nsConn[networkServerID]
	if !ok {
		return nil, fmt.Errorf("no such connection for network server id= %d", networkServerID)
	}
	return ns.NewNetworkServerServiceClient(v), nil
}

// GetNetworkServerExtraServiceClient returns a new NetworkServerExtraServiceClient instance
func (cli *Client) GetNetworkServerExtraServiceClient(networkServerID int64) (nsextra.NetworkServerExtraServiceClient, error) {
	v, ok := cli.nsConn[networkServerID]
	if !ok {
		return nil, fmt.Errorf("no such connection for network server id= %d", networkServerID)
	}
	return nsextra.NewNetworkServerExtraServiceClient(v), nil
}
