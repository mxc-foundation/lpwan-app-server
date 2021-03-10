package dfi

import (
	"context"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
)

// Server defines DFI service server API structure
type Server struct {
	store  Store
	mxpCli *mxpcli.Client
}

// NewServer creates a new DFI service server
func NewServer(store Store, mxpCli *mxpcli.Client) *Server {
	return &Server{
		store:  store,
		mxpCli: mxpCli,
	}
}

// Store defines db APIs for DFI service
type Store interface {
}

// AuthenticateUser authenticates user with given jwt, return necessary user info for DFI service
func (s Server) AuthenticateUser(ctx context.Context, req *api.DFIAuthenticateUserRequest) (*api.DFIAuthenticateUserResponse, error) {
	return &api.DFIAuthenticateUserResponse{
		UserEmail:      "test@mxc.org",
		OrganizationId: "1",
		MxcBalance:     "1000000",
	}, nil
}

// TopUp allows user to top up DFI margin wallet from DD wallet/supernode wallet
func (s Server) TopUp(ctx context.Context, req *api.TopUpRequest) (*api.TopUpResponse, error) {
	return &api.TopUpResponse{}, nil
}

// Withdraw allows user to withdraw from DFI margin wallet to DD wallet/supernode wallet
func (s Server) Withdraw(ctx context.Context, req *api.WithdrawRequest) (*api.WithdrawResponse, error) {
	return &api.WithdrawResponse{Msg: "success"}, nil
}
