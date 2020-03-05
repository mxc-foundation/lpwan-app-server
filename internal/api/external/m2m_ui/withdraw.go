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

// WithdrawServerAPI validates the withdraw server api
type WithdrawServerAPI struct {
	validator auth.Validator
}

// NewWithdrawServerAPI defines the withdraw server api
func NewWithdrawServerAPI(validator auth.Validator) *WithdrawServerAPI {
	return &WithdrawServerAPI{
		validator: validator,
	}
}

// ModifyWithdrawFee modifies the withdraw fee
func (s *WithdrawServerAPI) ModifyWithdrawFee(ctx context.Context, req *api.ModifyWithdrawFeeRequest) (*api.ModifyWithdrawFeeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/ModifyWithdrawFee")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := api.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.ModifyWithdrawFee(ctx, &api.ModifyWithdrawFeeRequest{
		MoneyAbbr:   req.MoneyAbbr,
		WithdrawFee: req.WithdrawFee,
		OrgId:       req.OrgId,
	})
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifyWithdrawFeeResponse{
		Status:      resp.Status,
	}, nil
}

// GetWithdrawFee gets the withdraw fee
func (s *WithdrawServerAPI) GetWithdrawFee(ctx context.Context, req *api.GetWithdrawFeeRequest) (*api.GetWithdrawFeeResponse, error) {
	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := api.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawFee(ctx, &api.GetWithdrawFeeRequest{
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawFeeResponse{
		WithdrawFee: resp.WithdrawFee,
	}, nil
}

// GetWithdrawHistory gets the withdraw history
func (s *WithdrawServerAPI) GetWithdrawHistory(ctx context.Context, req *api.GetWithdrawHistoryRequest) (*api.GetWithdrawHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWithdrawHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := api.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawHistory(ctx, &api.GetWithdrawHistoryRequest{
		OrgId:     req.OrgId,
		Offset:    req.Offset,
		Limit:     req.Limit,
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawHistoryResponse{
		Count:           resp.Count,
		WithdrawHistory: resp.WithdrawHistory,
	}, nil
}

// WithdrawReq defines request for withdraw
func (s *WithdrawServerAPI) WithdrawReq(ctx context.Context, req *api.WithdrawReqRequest) (*api.WithdrawReqResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/WithdrawReq")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := api.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.WithdrawReq(ctx, &api.WithdrawReqRequest{
		OrgId:      req.OrgId,
		MoneyAbbr:  req.MoneyAbbr,
		Amount:     req.Amount,
		EthAddress: req.EthAddress,
	})
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.WithdrawReqResponse{
		Status:      resp.Status,
	}, nil
}

func (s *WithdrawServerAPI) ConfirmWithdraw (ctx context.Context, req *api.ConfirmWithdrawRequest) (*api.ConfirmWithdrawResponse, error) {
	return &api.ConfirmWithdrawResponse{}, nil
}

func (s *WithdrawServerAPI) GetWithdrawRequestList(ctx context.Context, req *api.GetWithdrawRequestListRequest) (*api.GetWithdrawRequestListResponse, error) {

	return &api.GetWithdrawRequestListResponse{}, nil
}