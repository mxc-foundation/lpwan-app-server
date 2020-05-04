package external

import (
	"context"
	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
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

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Internal, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.ModifyWithdrawFee(ctx, &m2mServer.ModifyWithdrawFeeRequest{
		MoneyAbbr:   m2mServer.Money(req.MoneyAbbr),
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
		MoneyAbbr: m2mServer.Money(req.MoneyAbbr),
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

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawHistory(ctx, &m2mServer.GetWithdrawHistoryRequest{
		OrgId:     req.OrgId,
		Offset:    req.Offset,
		Limit:     req.Limit,
		MoneyAbbr: m2mServer.Money(req.MoneyAbbr),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var withdrawHistoryList []*api.WithdrawHistory
	for _, item := range resp.WithdrawHistory {
		withdrawHistory := &api.WithdrawHistory{
			Amount:      item.Amount,
			TxSentTime:  item.TxSentTime,
			TxStatus:    item.TxStatus,
			TxHash:      item.TxHash,
			DenyComment: item.DenyComment,
		}

		withdrawHistoryList = append(withdrawHistoryList, withdrawHistory)
	}

	return &api.GetWithdrawHistoryResponse{
		Count:           resp.Count,
		WithdrawHistory: withdrawHistoryList,
	}, status.Error(codes.OK, "")
}

// GetWithdraw sends the requests to cobo directly
func (s *WithdrawServerAPI) GetWithdraw (ctx context.Context, req *api.GetWithdrawRequest) (*api.GetWithdrawResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdraw org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// if user is not global admin, user must have access to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetWithdrawResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdraw(ctx, &m2mServer.GetWithdrawRequest{
		OrgId:            req.OrgId,
		Amount:           req.Amount,
		EthAddress:       req.EthAddress,
		AvailableBalance: req.AvailableBalance,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWithdrawResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// WithdrawReq defines request for withdraw
func (s *WithdrawServerAPI) WithdrawReq(ctx context.Context, req *api.WithdrawReqRequest) (*api.WithdrawReqResponse, error) {
	logInfo := "api/appserver_serves_ui/WithdrawReq org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// if user is not global admin, user must have access to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.WithdrawReqResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.WithdrawReq(ctx, &m2mServer.WithdrawReqRequest{
		OrgId:            req.OrgId,
		Amount:           req.Amount,
		EthAddress:       req.EthAddress,
		AvailableBalance: req.AvailableBalance,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.WithdrawReqResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

func (s *WithdrawServerAPI) ConfirmWithdraw(ctx context.Context, req *api.ConfirmWithdrawRequest) (*api.ConfirmWithdrawResponse, error) {
	logInfo := "api/appserver_serves_ui/ConfirmWithdraw org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ConfirmWithdrawResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.ConfirmWithdrawResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ConfirmWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.ConfirmWithdraw(ctx, &m2mServer.ConfirmWithdrawRequest{
		OrgId:         req.OrgId,
		ConfirmStatus: req.ConfirmStatus,
		DenyComment:   req.DenyComment,
		WithdrawId:    req.WithdrawId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.ConfirmWithdrawResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ConfirmWithdrawResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// GetWithdrawRequestList returns all users withdrawal requests to the front-end
func (s *WithdrawServerAPI) GetWithdrawRequestList(ctx context.Context, req *api.GetWithdrawRequestListRequest) (*api.GetWithdrawRequestListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdrawRequestList"
	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawRequestListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.GetWithdrawRequestListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawRequestListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawRequestList(ctx, &m2mServer.GetWithdrawRequestListRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawRequestListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	// ToDo: Get the user name from DB
	/*for _, v := range resp.WithdrawRequest {
		v.UserName = ""
	}*/

	var withdrawRequestList []*api.WithdrawRequest
	for _, item := range resp.WithdrawRequest {
		withdrawRequest := &api.WithdrawRequest{
			UserName:       item.UserName,
			AvailableToken: item.AvailableToken,
			Amount:         item.Amount,
			WithdrawId:     item.WithdrawId,
		}

		withdrawRequestList = append(withdrawRequestList, withdrawRequest)
	}

	return &api.GetWithdrawRequestListResponse{
		Count:           resp.Count,
		WithdrawRequest: withdrawRequestList,
	}, status.Error(codes.OK, "")
}

func (s *WithdrawServerAPI) GetWithdrawMonthly(ctx context.Context, req *api.GetWithdrawMonthlyRequest) (*api.GetWithdrawMonthlyResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWithdrawMonthly"
	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawMonthlyResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		return &api.GetWithdrawMonthlyResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawMonthlyResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawClient := m2mServer.NewWithdrawServiceClient(m2mClient)

	resp, err := withdrawClient.GetWithdrawMonthly (ctx, &m2mServer.GetWithdrawMonthlyRequest{
		UserId: req.UserId,
		OrgId:  req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWithdrawMonthlyResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var withdrawData []*api.MonthlyData
	for _, item := range resp.MonthlyData {
		monthlyWithdraw := &api.MonthlyData{
			Month: item.Month,
			Amount: item.Amount,
		}

		withdrawData = append(withdrawData, monthlyWithdraw)
	}

	return &api.GetWithdrawMonthlyResponse{
		MonthlyData: withdrawData,
	}, status.Error(codes.OK, "")
}