package m2m_ui

import (
	"context"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TopUpServerAPI defines the topup server api structure
type TopUpServerAPI struct {
	validator auth.Validator
}

// NewTopUpServerAPI validates the topup server api
func NewTopUpServerAPI(validator auth.Validator) *TopUpServerAPI {
	return &TopUpServerAPI{
		validator: validator,
	}
}

// GetTopUpHistory defines the topup history request and response
func (s *TopUpServerAPI) GetTopUpHistory(ctx context.Context, req *api.GetTopUpHistoryRequest) (*api.GetTopUpHistoryResponse, error) {
	log.WithField("orgId", req.UserId).Info("grpc_api/GetTopUpHistory")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.UserId)
	if err != nil {
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	topupClient := api.NewTopUpServiceClient(m2mClient)

	resp, err := topupClient.GetTopUpHistory(ctx, &api.GetTopUpHistoryRequest{
		UserId:  req.UserId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetTopUpHistoryResponse{
		Count:        resp.Count,
		TopupHistory: resp.TopupHistory,
		UserProfile:  &prof,
	}, nil
}

// GetTopUpDestination defines the topup destination request and response
func (s *TopUpServerAPI) GetTopUpDestination(ctx context.Context, req *api.GetTopUpDestinationRequest) (*api.GetTopUpDestinationResponse, error) {
	log.WithField("orgId", req.UserId).Info("grpc_api/GetTopUpDestination")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.UserId)
	if err != nil {
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	topupClient := api.NewTopUpServiceClient(m2mClient)

	resp, err := topupClient.GetTopUpDestination(ctx, &api.GetTopUpDestinationRequest{
		UserId:     req.UserId,
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetTopUpDestinationResponse{
		ActiveAccount: resp.ActiveAccount,
		UserProfile:   &prof,
	}, nil
}