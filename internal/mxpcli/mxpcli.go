// Package mxpcli handles connections to mxprotocol server
package mxpcli

import (
	"fmt"
	"google.golang.org/grpc"

	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
)

// Client represents mxprotocol server client
type Client struct {
	mxpConn *grpc.ClientConn
}

// Global mxprotocol server client, it exists only to keep existing code
// working and must not be used in any new code. It should be removed as soon
// as no other module uses it.
var Global *Client

// Connect connects to mxprotocol server and returns the client
func Connect(cfg grpccli.ConnectionOpts) (*Client, error) {
	mxpConn, err := grpccli.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't create mxprotocol server client: %v", err)
	}
	return &Client{
		mxpConn: mxpConn,
	}, nil
}

// Close closes connection to mxprotocol server
func (c *Client) Close() error {
	return c.mxpConn.Close()
}

// GetDHXServiceClient returns a new GetDHXServiceClient of mxprotocol-server
func (c *Client) GetDHXServiceClient() pb.DHXServiceClient {
	return pb.NewDHXServiceClient(c.mxpConn)
}

// GetDistributeBonusServiceClient returns a new DistributeBonusServiceClient instance
func (c *Client) GetDistributeBonusServiceClient() pb.DistributeBonusServiceClient {
	return pb.NewDistributeBonusServiceClient(c.mxpConn)
}

// GetM2MDeviceServiceClient returns a new DSDeviceServiceClient of mxprotocol-server
func (c *Client) GetM2MDeviceServiceClient() pb.DSDeviceServiceClient {
	return pb.NewDSDeviceServiceClient(c.mxpConn)
}

// GetM2MGatewayServiceClient returns a new GSGatewayServiceClient of mxprotocol-server
func (c *Client) GetM2MGatewayServiceClient() pb.GSGatewayServiceClient {
	return pb.NewGSGatewayServiceClient(c.mxpConn)
}

// GetMiningServiceClient returns a new MiningServiceClient of mxprotocol-server
func (c *Client) GetMiningServiceClient() pb.MiningServiceClient {
	return pb.NewMiningServiceClient(c.mxpConn)
}

// GetServerServiceClient returns a new M2MServerInfoServiceClient of mxprotocol-server
func (c *Client) GetServerServiceClient() pb.M2MServerInfoServiceClient {
	return pb.NewM2MServerInfoServiceClient(c.mxpConn)
}

// GetSettingsServiceClient returns a new SettingsServiceClient of mxprotocol-server
func (c *Client) GetSettingsServiceClient() pb.SettingsServiceClient {
	return pb.NewSettingsServiceClient(c.mxpConn)
}

// GetStakingServiceClient returns a new StakingServiceClient of mxprotocol-server
func (c *Client) GetStakingServiceClient() pb.StakingServiceClient {
	return pb.NewStakingServiceClient(c.mxpConn)
}

// GetTopupServiceClient returns a new TopUpServiceClient of mxprotocol-server
func (c *Client) GetTopupServiceClient() pb.TopUpServiceClient {
	return pb.NewTopUpServiceClient(c.mxpConn)
}

// GetWalletServiceClient returns a new WalletServiceClient( of mxprotocol-server
func (c *Client) GetWalletServiceClient() pb.WalletServiceClient {
	return pb.NewWalletServiceClient(c.mxpConn)
}

// GetWithdrawServiceClient returns a new WithdrawServiceClient of mxprotocol-server
func (c *Client) GetWithdrawServiceClient() pb.WithdrawServiceClient {
	return pb.NewWithdrawServiceClient(c.mxpConn)
}

// GetFianceReportClient returns a new FinanceReportServiceClient of mxprotocol-server
func (c *Client) GetFianceReportClient() pb.FinanceReportServiceClient {
	return pb.NewFinanceReportServiceClient(c.mxpConn)
}
