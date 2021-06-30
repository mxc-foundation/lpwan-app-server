package pscli

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

// Client represents provisioning server client
type Client struct {
	psconn *grpc.ClientConn
}

// Connect connects to provisioning server and returns the client
func Connect(config types.ProvisioningServerStruct) (*Client, error) {
	psconn, err := grpccli.Connect(config.ServerConifig)
	if err != nil {
		return nil, fmt.Errorf("failed to create provisioning server client: %v", err)
	}
	return &Client{
		psconn: psconn,
	}, nil
}

// Close closes connection to provisioning server
func (c *Client) Close() error {
	return c.psconn.Close()
}

// GetPServerClient returns a new ProvisionClient of provisioning server
func (c *Client) GetPServerClient() pb.ProvisionClient {
	return pb.NewProvisionClient(c.psconn)
}

// GetDeviceProvisionServiceClient returns a new DeviceProvisionClient of provisioning server
func (c *Client) GetDeviceProvisionServiceClient() pb.DeviceProvisionClient {
	return pb.NewDeviceProvisionClient(c.psconn)
}
