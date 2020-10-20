package mxp_portal

import (
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
)

// GetM2MDeviceServiceClient returns a new DSDeviceServiceClient of mxprotocol-server
func GetM2MDeviceServiceClient() (pb.DSDeviceServiceClient, error) {
	return pb.NewDSDeviceServiceClient(ctrl.m2mconn), nil
}

// GetM2MGatewayServiceClient returns a new GSGatewayServiceClient of mxprotocol-server
func GetM2MGatewayServiceClient() (pb.GSGatewayServiceClient, error) {
	return pb.NewGSGatewayServiceClient(ctrl.m2mconn), nil
}

// GetMiningServiceClient returns a new MiningServiceClient of mxprotocol-server
func GetMiningServiceClient() (pb.MiningServiceClient, error) {
	return pb.NewMiningServiceClient(ctrl.m2mconn), nil
}

// GetServerServiceClient returns a new M2MServerInfoServiceClient of mxprotocol-server
func GetServerServiceClient() (pb.M2MServerInfoServiceClient, error) {
	return pb.NewM2MServerInfoServiceClient(ctrl.m2mconn), nil
}

// GetSettingsServiceClient returns a new SettingsServiceClient of mxprotocol-server
func GetSettingsServiceClient() (pb.SettingsServiceClient, error) {
	return pb.NewSettingsServiceClient(ctrl.m2mconn), nil
}

// GetStakingServiceClient returns a new StakingServiceClient of mxprotocol-server
func GetStakingServiceClient() (pb.StakingServiceClient, error) {
	return pb.NewStakingServiceClient(ctrl.m2mconn), nil
}

// GetTopupServiceClient returns a new TopUpServiceClient of mxprotocol-server
func GetTopupServiceClient() (pb.TopUpServiceClient, error) {
	return pb.NewTopUpServiceClient(ctrl.m2mconn), nil
}

// GetWalletServiceClient returns a new WalletServiceClient( of mxprotocol-server
func GetWalletServiceClient() (pb.WalletServiceClient, error) {
	return pb.NewWalletServiceClient(ctrl.m2mconn), nil
}

// GetWithdrawServiceClient returns a new WithdrawServiceClient of mxprotocol-server
func GetWithdrawServiceClient() (pb.WithdrawServiceClient, error) {
	return pb.NewWithdrawServiceClient(ctrl.m2mconn), nil
}
