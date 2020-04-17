package external

import (
	"context"
	"strconv"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m_serves_appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WalletServerAPI is the structure that contains the validator
type WalletServerAPI struct {
	validator auth.Validator
}

// NewWalletServerAPI validates the new wallet server api
func NewWalletServerAPI(validator auth.Validator) *WalletServerAPI {
	return &WalletServerAPI{
		validator: validator,
	}
}

// GetWalletBalance gets the wallet balance
func (s *WalletServerAPI) GetWalletBalance(ctx context.Context, req *api.GetWalletBalanceRequest) (*api.GetWalletBalanceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletBalance org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetWalletBalance(ctx, &m2mServer.GetWalletBalanceRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWalletBalanceResponse{
		Balance: resp.Balance,
	}, status.Error(codes.OK, "")
}

func (s *WalletServerAPI) GetWalletMiningIncome(ctx context.Context, req *api.GetWalletMiningIncomeRequest) (*api.GetWalletMiningIncomeResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletMiningIncome org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletMiningIncomeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetWalletMiningIncomeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletMiningIncomeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetWalletMiningIncome(ctx, &m2mServer.GetWalletMiningIncomeRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletMiningIncomeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWalletMiningIncomeResponse{
		MiningIncome: resp.MiningIncome,
	}, status.Error(codes.OK, "")
}

func (s *WalletServerAPI) GetMiningInfo(ctx context.Context, req *api.GetMiningInfoRequest) (*api.GetMiningInfoResponse, error) {
	logInfo := "api/appserver_serves_ui/GetMiningInfo org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMiningInfoResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetMiningInfoResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMiningInfoResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetMiningInfo(ctx, &m2mServer.GetMiningInfoRequest{
		UserId: req.UserId,
		OrgId:  req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMiningInfoResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	/*miningInfo := &api.GetMiningInfoResponse{}
	for _, v := range resp.MiningData {
		miningInfo.MiningData = append(miningInfo.MiningData, v)
	}*/

	var miningData []*api.MiningData
	for _, item := range resp.MiningData {
		miningInfo := &api.MiningData{
			Month: item.Month,
			Amount: item.Amount,
		}

		miningData = append(miningData, miningInfo)
	}

	return &api.GetMiningInfoResponse{
		MiningIncome: resp.MiningIncome,
		MiningData: miningData,
	}, status.Error(codes.OK, "")
}

// GetVmxcTxHistory gets virtual MXC transaction history
func (s *WalletServerAPI) GetVmxcTxHistory(ctx context.Context, req *api.GetVmxcTxHistoryRequest) (*api.GetVmxcTxHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetVmxcTxHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetVmxcTxHistory(ctx, &m2mServer.GetVmxcTxHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var vmxcTxHistoryList []*api.VmxcTxHistory
	for _, item := range resp.TxHistory {
		vmxcTxHistory := &api.VmxcTxHistory{
			From:      item.From,
			To:        item.To,
			TxType:    item.TxType,
			Amount:    item.Amount,
			CreatedAt: item.CreatedAt,
		}

		vmxcTxHistoryList = append(vmxcTxHistoryList, vmxcTxHistory)
	}

	return &api.GetVmxcTxHistoryResponse{
		Count:     resp.Count,
		TxHistory: vmxcTxHistoryList,
	}, status.Error(codes.OK, "")
}

// GetWalletUsageHist gets the walllet usage history
func (s *WalletServerAPI) GetWalletUsageHist(ctx context.Context, req *api.GetWalletUsageHistRequest) (*api.GetWalletUsageHistResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletUsageHist org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetWalletUsageHist(ctx, &m2mServer.GetWalletUsageHistRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var walletUsageHistoryList []*api.GetWalletUsageHist
	for _, item := range resp.WalletUsageHis {
		walletUsageHist := &api.GetWalletUsageHist{
			StartAt:         item.StartAt,
			DurationMinutes: item.DurationMinutes,
			DlCntDv:         item.DlCntDv,
			DlCntDvFree:     item.DlCntDvFree,
			UlCntDv:         item.UlCntDv,
			UlCntDvFree:     item.UlCntDvFree,
			DlCntGw:         item.DlCntGw,
			DlCntGwFree:     item.DlCntDvFree,
			UlCntGw:         item.UlCntGw,
			UlCntGwFree:     item.UlCntGwFree,
			Spend:           item.Spend,
			Income:          item.Income,
			UpdatedBalance:  item.UpdatedBalance,
		}

		walletUsageHistoryList = append(walletUsageHistoryList, walletUsageHist)
	}

	return &api.GetWalletUsageHistResponse{
		WalletUsageHis: walletUsageHistoryList,
		Count:          resp.Count,
	}, status.Error(codes.OK, "")
}

// GetDlPrice gets downlink price from m2m wallet
func (s *WalletServerAPI) GetDlPrice(ctx context.Context, req *api.GetDownLinkPriceRequest) (*api.GetDownLinkPriceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDlPrice org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletClient := m2mServer.NewWalletServiceClient(m2mClient)

	resp, err := walletClient.GetDlPrice(ctx, &m2mServer.GetDownLinkPriceRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDownLinkPriceResponse{
		DownLinkPrice: resp.DownLinkPrice,
	}, status.Error(codes.OK, "")
}
