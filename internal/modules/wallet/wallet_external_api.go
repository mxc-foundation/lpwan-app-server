package wallet

import (
	"context"
	"strconv"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/coingecko"
)

// WalletServerAPI is the structure that contains the validator
type WalletServerAPI struct {
}

// NewWalletServerAPI validates the new wallet server api
func NewWalletServerAPI() *WalletServerAPI {
	return &WalletServerAPI{}
}

// GetWalletBalance gets the wallet balance
func (s *WalletServerAPI) GetWalletBalance(ctx context.Context, req *api.GetWalletBalanceRequest) (*api.GetWalletBalanceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletBalance org=" + strconv.FormatInt(req.OrgId, 10)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetWalletBalance(ctx, &pb.GetWalletBalanceRequest{
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

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletMiningIncomeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetWalletMiningIncome(ctx, &pb.GetWalletMiningIncomeRequest{
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

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMiningInfoResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetMiningInfo(ctx, &pb.GetMiningInfoRequest{
		OrgId: req.OrgId,
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
	for _, item := range resp.Data {
		miningInfo := &api.MiningData{
			Month:  item.Month,
			Amount: item.Amount,
		}

		miningData = append(miningData, miningInfo)
	}

	return &api.GetMiningInfoResponse{
		TodayRev: resp.TodayRev,
		Data:     miningData,
	}, status.Error(codes.OK, "")
}

// GetVmxcTxHistory gets virtual MXC transaction history
func (s *WalletServerAPI) GetVmxcTxHistory(ctx context.Context, req *api.GetVmxcTxHistoryRequest) (*api.GetVmxcTxHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetVmxcTxHistory org=" + strconv.FormatInt(req.OrgId, 10)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetVmxcTxHistory(ctx, &pb.GetVmxcTxHistoryRequest{
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
func (s *WalletServerAPI) GetNetworkUsageHist(ctx context.Context, req *api.GetNetworkUsageHistRequest) (*api.GetNetworkUsageHistResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletUsageHist org=" + strconv.FormatInt(req.OrgId, 10)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetNetworkUsageHistResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetNetworkUsageHist(ctx, &pb.GetNetworkUsageHistRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
		From:     req.From,
		Till:     req.Till,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	var walletUsageHistoryList []*api.NetworkUsage
	for _, item := range resp.NetworkUsage {
		walletUsageHist := &api.NetworkUsage{
			Timestamp:    item.Timestamp,
			DlCntDev:     item.DlCntDev,
			DlCntDevFree: item.DlCntDevFree,
			UlCntDev:     item.UlCntDev,
			UlCntDevFree: item.UlCntDevFree,
			DlCntGw:      item.DlCntGw,
			DlCntGwFree:  item.DlCntGwFree,
			UlCntGw:      item.UlCntGw,
			UlCntGwFree:  item.UlCntGwFree,
			Amount:       item.Amount,
		}

		walletUsageHistoryList = append(walletUsageHistoryList, walletUsageHist)
	}

	return &api.GetNetworkUsageHistResponse{
		NetworkUsage: walletUsageHistoryList,
	}, status.Error(codes.OK, "")
}

// GetDlPrice gets downlink price from m2m wallet
func (s *WalletServerAPI) GetDlPrice(ctx context.Context, req *api.GetDownLinkPriceRequest) (*api.GetDownLinkPriceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetDlPrice org=" + strconv.FormatInt(req.OrgId, 10)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient, err := m2mcli.GetWalletServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := walletClient.GetDlPrice(ctx, &pb.GetDownLinkPriceRequest{
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

func (s *WalletServerAPI) GetMXCprice(ctx context.Context, req *api.GetMXCpriceRequest) (*api.GetMXCpriceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetMXCprice org=" + strconv.FormatInt(req.OrgId, 10)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	if req.MxcPrice == "" {
		return &api.GetMXCpriceResponse{MxcPrice: "0"}, nil
	}
	mxc, err := decimal.NewFromString(req.MxcPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mxcPrice must be a number")
	}
	price, err := coingecko.New().GetPrice("mxc", "usd")
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMXCpriceResponse{}, status.Errorf(codes.Internal, "unable to get price from CMC")
	}
	rate := decimal.NewFromFloat(price)
	usd := mxc.Mul(rate).Round(18)

	return &api.GetMXCpriceResponse{MxcPrice: usd.String()}, nil
}
