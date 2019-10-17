package external

import (
	pb "github.com/brocaar/lora-app-server/api"
	"github.com/brocaar/lora-app-server/internal/api/external/auth"
	"golang.org/x/net/context"
)

type ProxyRequestAPI struct {
	validator auth.Validator
}

func NewProxyRequestAPI(validator auth.Validator) *ProxyRequestAPI {
	return &ProxyRequestAPI{
		validator: validator,
	}
}

func (a *ProxyRequestAPI) GetWalletBalance(ctx context.Context, req *pb.GetWalletBalanceRequest) (*pb.GetWalletBalanceResponse, error) {


	return &pb.GetWalletBalanceResponse{}, nil
}
