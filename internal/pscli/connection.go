package pscli

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
)

type Client struct {
	psconn *grpc.ClientConn
}

func Connect(config types.ProvisioningServerStruct) (*Client, error) {
	psconn, err := grpccli.Connect(config.ServerConifig)
	if err != nil {
		return nil, fmt.Errorf("failed to create provisioning server client: %v", err)
	}
	return &Client{
		psconn: psconn,
	}, nil
}

func (c *Client) Close() error {
	return c.psconn.Close()
}

func (c *Client) GetPServerClient() pb.ProvisionClient {
	return pb.NewProvisionClient(c.psconn)
}

func (c *Client) GetDeviceProvisionServiceClient() pb.DeviceProvisionClient {
	return pb.NewDeviceProvisionClient(c.psconn)
}
