package external

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
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
	logInfo := "api/appserver_serves_ui/ModifyWithdrawFee"

	cred, err := s.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsGlobalAdmin(ctx); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be a global admin")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.ModifyWithdrawFee(ctx, &m2mServer.ModifyWithdrawFeeRequest{
		Currency:    req.Currency,
		WithdrawFee: req.WithdrawFee,
		Password:    req.Password,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifyWithdrawFeeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// GetWithdrawFee gets the withdraw fee
func (s *WithdrawServerAPI) GetWithdrawFee(ctx context.Context, req *api.GetWithdrawFeeRequest) (*api.GetWithdrawFeeResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdrawFee"

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawFee(ctx, &m2mServer.GetWithdrawFeeRequest{
		Currency: req.Currency,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawFeeResponse{
		WithdrawFee: resp.WithdrawFee,
	}, status.Error(codes.OK, "")
}

// GetWithdrawHistory gets the withdraw history
func (s *WithdrawServerAPI) GetWithdrawHistory(ctx context.Context, req *api.GetWithdrawHistoryRequest) (*api.GetWithdrawHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdrawHistory org=" + strconv.FormatInt(req.OrgId, 10)

	cred, err := s.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgAdmin(ctx, req.OrgId); err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.PermissionDenied, "must be an organization admin")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawHistory(ctx, &m2mServer.GetWithdrawHistoryRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
		From:     req.From,
		Till:     req.Till,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var withdrawHistoryList []*api.WithdrawHistory
	for _, item := range resp.WithdrawHistory {
		withdrawHistory := &api.WithdrawHistory{
			Amount:      item.Amount,
			Timestamp:   item.Timestamp,
			TxStatus:    item.TxStatus,
			TxHash:      item.TxHash,
			DenyComment: item.DenyComment,
			WithdrawFee: item.WithdrawFee,
		}

		withdrawHistoryList = append(withdrawHistoryList, withdrawHistory)
	}

	return &api.GetWithdrawHistoryResponse{
		WithdrawHistory: withdrawHistoryList,
	}, status.Error(codes.OK, "")
}

// GetWithdraw sends the requests to cobo directly
func (s *WithdrawServerAPI) GetWithdraw(ctx context.Context, req *api.GetWithdrawRequest) (*api.GetWithdrawResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdraw org=" + strconv.FormatInt(req.OrgId, 10)
	cred, err := s.validator.GetCredentials(ctx, auth.WithValidOTP())
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	// user must be organization admin to withdraw
	if err := cred.IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "must be an organization admin")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdraw(ctx, &m2mServer.GetWithdrawRequest{
		OrgId:      req.OrgId,
		Currency:   req.Currency,
		Amount:     req.Amount,
		EthAddress: req.EthAddress,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
