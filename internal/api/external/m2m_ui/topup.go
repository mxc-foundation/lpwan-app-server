package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TopUpServerAPI struct {
	validator auth.Validator
}

func NewTopUpServerAPI(validator auth.Validator) *TopUpServerAPI {
	return &TopUpServerAPI{
		validator: validator,
	}
}

func (s *TopUpServerAPI) GetTransactionsHistory(ctx context.Context, req *api.GetTransactionsHistoryRequest) (*api.GetTransactionsHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetTransactionsHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetTransactionsHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetTransactionsHistory(ctx, &m2m_api.GetTransactionsHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetTransactionsHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	txHist := api.GetTransactionsHistoryResponse.GetTransactionHistory(&resp.TransactionHistory)

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetTransactionsHistoryResponse{}, err
	}

	userProfile := api.GetTransactionsHistoryResponse.GetUserProfile(getUserProfile)

	return &api.GetTransactionsHistoryResponse{
		Count:              resp.Count,
		TransactionHistory: txHist,
		UserProfile:        userProfile,
	}, nil
}

func (s *TopUpServerAPI) GetTopUpHistory(ctx context.Context, req *api.GetTopUpHistoryRequest) (*api.GetTopUpHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetTopUpHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetTopUpHistory(ctx, &m2m_api.GetTopUpHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	topupHist := api.GetTopUpHistoryResponse.GetTopupHistory(&resp.TopupHistory)

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetTopUpHistoryResponse{}, err
	}

	userProfile := api.GetTopUpHistoryResponse.GetUserProfile(getUserProfile)

	return &api.GetTopUpHistoryResponse{
		Count:        resp.Count,
		TopupHistory: topupHist,
		UserProfile:  userProfile,
	}, nil
}

func (s *TopUpServerAPI) GetTopUpDestination(ctx context.Context, req *api.GetTopUpDestinationRequest) (*api.GetTopUpDestinationResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetTopUpDestination")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.GetTopUpDestination(ctx, &m2m_api.GetTopUpDestinationRequest{
		OrgId:     req.OrgId,
		MoneyAbbr: moneyAbbr,
	})
	if err != nil {
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetTopUpDestinationResponse{}, err
	}

	userProfile := api.GetTopUpDestinationResponse.GetUserProfile(getUserProfile)

	return &api.GetTopUpDestinationResponse{
		ActiveAccount: resp.ActiveAccount,
		UserProfile:   userProfile,
	}, nil
}

func (s *TopUpServerAPI) GetIncome(ctx context.Context, req *api.GetIncomeRequest) (*api.GetIncomeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetIncome")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetIncomeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetIncome(ctx, &m2m_api.GetIncomeRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.GetIncomeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetIncomeResponse{
		Amount: resp.Amount,
	}, nil
}
