package external

import (
	"context"
	"strconv"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"

	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/withdraw"
)

// WithdrawServerAPI validates the withdraw server api
type WithdrawServerAPI struct {
	auth auth.Authenticator
}

// NewWithdrawServerAPI defines the withdraw server api
func NewWithdrawServerAPI(auth auth.Authenticator) *WithdrawServerAPI {
	return &WithdrawServerAPI{
		auth: auth,
	}
}

// ModifyWithdrawFee modifies the withdraw fee
func (s *WithdrawServerAPI) ModifyWithdrawFee(ctx context.Context, req *api.ModifyWithdrawFeeRequest) (*api.ModifyWithdrawFeeResponse, error) {
	logInfo := "api/appserver_serves_ui/ModifyWithdrawFee"

	if err := withdraw.NewValidator().IsGlobalAdmin(ctx); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	withdrawClient := mxpcli.Global.GetWithdrawServiceClient()

	resp, err := withdrawClient.ModifyWithdrawFee(ctx, &pb.ModifyWithdrawFeeRequest{
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

	withdrawClient := mxpcli.Global.GetWithdrawServiceClient()
	resp, err := withdrawClient.GetWithdrawFee(ctx, &pb.GetWithdrawFeeRequest{
		Currency: req.Currency,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawFeeResponse{
		WithdrawFee: resp.WithdrawFee,
		Currency:    resp.Currency,
	}, status.Error(codes.OK, "")
}

// GetWithdrawHistory gets the withdraw history
func (s *WithdrawServerAPI) GetWithdrawHistory(ctx context.Context, req *api.GetWithdrawHistoryRequest) (*api.GetWithdrawHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdrawHistory org=" + strconv.FormatInt(req.OrgId, 10)

	if err := withdraw.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	withdrawClient := mxpcli.Global.GetWithdrawServiceClient()

	resp, err := withdrawClient.GetWithdrawHistory(ctx, &pb.GetWithdrawHistoryRequest{
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

	options := auth.NewOptions()
	options.WithOrgID(req.OrgId)
	options.WithRequireOTP()
	cred, err := s.auth.GetCredentials(ctx, options)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed:%v", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	withdrawClient := mxpcli.Global.GetWithdrawServiceClient()

	resp, err := withdrawClient.GetWithdraw(ctx, &pb.GetWithdrawRequest{
		OrgId:      req.OrgId,
		Currency:   req.Currency,
		Amount:     req.Amount,
		EthAddress: req.EthAddress,
		Email:      cred.Username,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
