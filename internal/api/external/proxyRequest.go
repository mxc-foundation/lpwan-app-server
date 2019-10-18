package external

import (
	pb "github.com/brocaar/lora-app-server/api"
	m2m_api "github.com/brocaar/lora-app-server/api/m2m_server"
	"github.com/brocaar/lora-app-server/internal/api/external/auth"
	"github.com/brocaar/lora-app-server/internal/backend/m2m_client"
	"github.com/brocaar/lora-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if err := a.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWalletBalance")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &pb.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetWalletBalance(ctx, &m2m_api.GetWalletBalanceRequest{OrgId: req.OrgId})
	if err != nil {
		return &pb.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetWalletBalanceResponse{Balance: resp.Balance}, nil
}
