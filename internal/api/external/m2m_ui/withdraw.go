package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WithdrawServerAPI struct {
	validator auth.Validator
}

func NewWithdrawServerAPI(validator auth.Validator) *WithdrawServerAPI {
	return &WithdrawServerAPI{
		validator: validator,
	}
}

func (s *WithdrawServerAPI) ModifyWithdrawFee(ctx context.Context, req *api.ModifyWithdrawFeeRequest) (*api.ModifyWithdrawFeeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/ModifyWithdrawFee")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.ModifyWithdrawFee(ctx, &m2m_api.ModifyWithdrawFeeRequest{
		MoneyAbbr:   moneyAbbr,
		WithdrawFee: req.WithdrawFee,
		OrgId:       req.OrgId,
	})
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifyWithdrawFeeResponse{
		Status:      resp.Status,
		UserProfile: UserProfile,
	}, nil
}

func (s *WithdrawServerAPI) GetWithdrawFee(ctx context.Context, req *api.GetWithdrawFeeRequest) (*api.GetWithdrawFeeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWithdrawFee")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.GetWithdrawFee(ctx, &m2m_api.GetWithdrawFeeRequest{
		MoneyAbbr: moneyAbbr,
		OrgId:     req.OrgId,
	})
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawFeeResponse{
		WithdrawFee: resp.WithdrawFee,
		UserProfile: UserProfile,
	}, nil
}

func (s *WithdrawServerAPI) GetWithdrawHistory(ctx context.Context, req *api.GetWithdrawHistoryRequest) (*api.GetWithdrawHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWithdrawHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.GetWithdrawHistory(ctx, &m2m_api.GetWithdrawHistoryRequest{
		OrgId:     req.OrgId,
		Offset:    req.Offset,
		Limit:     req.Limit,
		MoneyAbbr: moneyAbbr,
	})
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawHist := api.GetWithdrawHistoryResponse.GetWithdrawHistory(resp.WithdrawHistory)

	return &api.GetWithdrawHistoryResponse{
		Count:           resp.Count,
		WithdrawHistory: withdrawHist,
		UserProfile:     UserProfile,
	}, nil
}

func (s *WithdrawServerAPI) WithdrawReq(ctx context.Context, req *api.WithdrawReqRequest) (*api.WithdrawReqResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/WithdrawReq")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.WithdrawReq(ctx, &m2m_api.WithdrawReqRequest{
		OrgId:     req.OrgId,
		MoneyAbbr: moneyAbbr,
		Amount:    req.Amount,
	})
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.WithdrawReqResponse{
		Status:      resp.Status,
		UserProfile: UserProfile,
	}, nil
}
