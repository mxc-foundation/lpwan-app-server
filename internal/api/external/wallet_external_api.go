package external

import (
	"context"
	"strconv"

	"github.com/brocaar/lorawan"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/auth"

	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/coingecko"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/wallet"
)

// Pricer allows you to get the price of crypto currency
type Pricer interface {
	GetPrice(crypto, fiat string) (float64, error)
}

// WalletServerAPI is the structure that contains the validator
type WalletServerAPI struct {
	pricer    Pricer
	st        *store.Handler
	auth      auth.Authenticator
	enableSTC bool
}

// NewWalletServerAPI validates the new wallet server api
func NewWalletServerAPI(h *store.Handler, auth auth.Authenticator, enableSTC bool) *WalletServerAPI {
	return &WalletServerAPI{
		pricer:    coingecko.New(),
		st:        h,
		auth:      auth,
		enableSTC: enableSTC,
	}
}

// GetWalletBalance gets the wallet balance
func (s *WalletServerAPI) GetWalletBalance(ctx context.Context, req *api.GetWalletBalanceRequest) (*api.GetWalletBalanceResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletBalance org=" + strconv.FormatInt(req.OrgId, 10)

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

	resp, err := walletClient.GetWalletBalance(ctx, &pb.GetWalletBalanceRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWalletBalanceResponse{
		Balance: resp.Balance,
	}, status.Error(codes.OK, "")
}

// GetGatewayMiningHealth returns information about health of the organization's gateways
func (s *WalletServerAPI) GetGatewayMiningHealth(ctx context.Context, req *api.GetGatewayMiningHealthRequest) (*api.GetGatewayMiningHealthResponse, error) {
	// req.OrgId should be the id of the org that user is making request with
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGatewayAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	gws, err := s.st.GetOrgGateways(ctx, cred.OrgID, req.GatewayMac)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var mreq pb.GetGatewayMiningHealthRequest
	for _, gw := range gws {
		var stcOrgID int64
		if gw.STCOrgID != nil {
			stcOrgID = *gw.STCOrgID
		}
		mreq.Gateway = append(mreq.Gateway, &pb.GatewayMining{
			GatewayMac: gw.MAC.String(),
			OwnerOrgId: gw.OrganizationID,
			StcOrgId:   stcOrgID,
		})
	}
	walletClient := mxpcli.Global.GetMiningServiceClient()
	mresp, err := walletClient.GetGatewayMiningHealth(ctx, &mreq)
	if err != nil {
		return nil, err
	}
	var resp api.GetGatewayMiningHealthResponse
	n := float32(len(mresp.GatewayHealth))
	resp.MiningHealthAverage = &api.MiningHealthAverage{}
	for _, gw := range mresp.GatewayHealth {
		resp.GatewayHealth = append(resp.GatewayHealth, &api.GatewayMiningHealth{
			Id:               gw.GatewayMac,
			OrgId:            gw.OrgId,
			Health:           gw.Health,
			MiningFuel:       gw.MiningFuel,
			MiningFuelMax:    gw.MiningFuelMax,
			MiningFuelHealth: gw.MiningFuelHealth,
			AgeSeconds:       gw.AgeSeconds,
			TotalMined:       gw.TotalMined,
			UptimeHealth:     gw.UptimeHealth,
		})
		resp.MiningHealthAverage.Overall += gw.Health / n
		resp.MiningHealthAverage.MiningFuelHealth += gw.MiningFuelHealth / n
		resp.MiningHealthAverage.UptimeHealth += gw.UptimeHealth / n
	}
	return &resp, nil
}

// TopUpGatewayMiningFuel adds mining fuel to the gateway
func (s *WalletServerAPI) TopUpGatewayMiningFuel(ctx context.Context, req *api.TopUpGatewayMiningFuelRequest) (*api.TopUpGatewayMiningFuelResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	var gwMAC lorawan.EUI64
	if err := gwMAC.UnmarshalText([]byte(req.GatewayMac)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid MAC: %s: %v", req.GatewayMac, err)
	}
	gw, err := s.st.GetGateway(ctx, gwMAC, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if gw.OrganizationID != req.OrgId || !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	walletClient := mxpcli.Global.GetMiningServiceClient()
	mresp, err := walletClient.TopUpGatewayMiningFuel(ctx, &pb.TopUpGatewayMiningFuelRequest{
		OrgId:      req.OrgId,
		GatewayMac: req.GatewayMac,
		Amount:     req.Amount,
		Currency:   req.Currency,
	})
	if err != nil {
		return nil, err
	}
	return &api.TopUpGatewayMiningFuelResponse{
		OrgId:      mresp.OrgId,
		GatewayMac: mresp.GatewayMac,
		Amount:     mresp.Amount,
		Currency:   mresp.Currency,
	}, nil
}

// WithdrawGatewayMiningFuel withdraws gateway's mining fuel
func (s *WalletServerAPI) WithdrawGatewayMiningFuel(ctx context.Context, req *api.WithdrawGatewayMiningFuelRequest) (*api.WithdrawGatewayMiningFuelResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	var gwMAC lorawan.EUI64
	if err := gwMAC.UnmarshalText([]byte(req.GatewayMac)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid MAC: %s: %v", req.GatewayMac, err)
	}
	gw, err := s.st.GetGateway(ctx, gwMAC, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if gw.OrganizationID != req.OrgId || !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	walletClient := mxpcli.Global.GetMiningServiceClient()
	mresp, err := walletClient.WithdrawGatewayMiningFuel(ctx, &pb.WithdrawGatewayMiningFuelRequest{
		OrgId:      req.OrgId,
		GatewayMac: req.GatewayMac,
		Amount:     req.Amount,
		Currency:   req.Currency,
	})
	if err != nil {
		return nil, err
	}
	return &api.WithdrawGatewayMiningFuelResponse{
		OrgId:      mresp.OrgId,
		GatewayMac: mresp.GatewayMac,
		Amount:     mresp.Amount,
		Currency:   mresp.Currency,
	}, nil
}

func (s *WalletServerAPI) GetGatewayMiningIncome(ctx context.Context, req *api.GetGatewayMiningIncomeRequest) (*api.GetGatewayMiningIncomeResponse, error) {
	// req.OrgId should be the id of the org that user is making request with
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayMac)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	// get gateway information
	item, err := s.st.GetGateway(ctx, mac, false)
	if err != nil {
		if err != errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.NotFound, "gateway with mac %s does not exist", req.GatewayMac)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if !cred.IsGlobalAdmin {
		if !cred.IsOrgAdmin {
			// user is neither global admin nor organization admin, return permission denied
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		if req.OrgId != item.OrganizationID {
			// gateway does not belong to req.OrgId, check wether gateway's reseller org id equals to req.OrgId
			if !s.enableSTC || item.STCOrgID == nil || *item.STCOrgID != req.OrgId {
				return nil, status.Errorf(codes.PermissionDenied, "permission denied")
			}
		}
	}

	logInfo := "api/appserver_serves_ui/GetGatewayMiningIncome org=" + strconv.FormatInt(req.OrgId, 10)
	walletClient := mxpcli.Global.GetMiningServiceClient()

	resp, err := walletClient.MiningStats(ctx, &pb.MiningStatsRequest{
		GatewayMac:     req.GatewayMac,
		OrganizationId: req.OrgId,
		FromDate:       req.FromDate,
		TillDate:       req.TillDate,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	stats := &api.GetGatewayMiningIncomeResponse{
		Total: resp.Total,
	}
	for _, ds := range resp.DailyStats {
		stats.DailyStats = append(stats.DailyStats, &api.MiningStats{
			Date:          ds.Date,
			Amount:        ds.Amount,
			OnlineSeconds: ds.OnlineSeconds,
			Health:        ds.Health,
		})
	}
	return stats, nil
}

func (s *WalletServerAPI) GetWalletMiningIncome(ctx context.Context, req *api.GetWalletMiningIncomeRequest) (*api.GetWalletMiningIncomeResponse, error) {
	logInfo := "api/appserver_serves_ui/GetWalletMiningIncome org=" + strconv.FormatInt(req.OrgId, 10)

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

	resp, err := walletClient.GetWalletMiningIncome(ctx, &pb.GetWalletMiningIncomeRequest{
		OrgId:    req.OrgId,
		From:     req.From,
		Till:     req.Till,
		Currency: req.Currency,
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

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

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

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

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

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

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

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()

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

	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	if req.MxcPrice == "" {
		return &api.GetMXCpriceResponse{MxcPrice: "0"}, nil
	}
	mxc, err := decimal.NewFromString(req.MxcPrice)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mxcPrice must be a number")
	}
	price, err := s.pricer.GetPrice("mxc", "usd")
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetMXCpriceResponse{}, status.Errorf(codes.Internal, "unable to get price from CMC")
	}
	rate := decimal.NewFromFloat(price)
	usd := mxc.Mul(rate).Round(18)

	return &api.GetMXCpriceResponse{MxcPrice: usd.String()}, nil
}

// GetTransactionHistory returns the history of transactions of the specified
// type over the specified period
func (s *WalletServerAPI) GetTransactionHistory(ctx context.Context, req *api.GetTransactionHistoryRequest) (*api.GetTransactionHistoryResponse, error) {
	if err := wallet.NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	walletClient := mxpcli.Global.GetWalletServiceClient()
	wcResp, err := walletClient.GetTransactionHistory(ctx, &pb.GetTransactionHistoryRequest{
		OrgId:       req.OrgId,
		Currency:    req.Currency,
		From:        req.From,
		Till:        req.Till,
		PaymentType: req.PaymentType,
	})
	if err != nil {
		return nil, err
	}
	resp := &api.GetTransactionHistoryResponse{}
	for _, tx := range wcResp.Tx {
		resp.Tx = append(resp.Tx, &api.Transaction{
			Id:          tx.Id,
			Timestamp:   tx.Timestamp,
			Amount:      tx.Amount,
			PaymentType: tx.PaymentType,
			DetailsJson: tx.DetailsJson,
		})
	}
	return resp, nil
}
