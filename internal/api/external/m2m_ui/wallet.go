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

type WalletServerAPI struct {
	validator auth.Validator
}

func NewWalletServerAPI(validator auth.Validator) *WalletServerAPI {
	return &WalletServerAPI{
		validator: validator,
	}
}

func (s *WalletServerAPI) GetWalletBalance(ctx context.Context, req *api.GetWalletBalanceRequest) (*api.GetWalletBalanceResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWalletBalance")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetWalletBalance(ctx, &m2m_api.GetWalletBalanceRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.GetWalletBalanceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetWalletBalanceResponse{
		Balance:     resp.Balance,
		UserProfile: UserProfile,
	}, nil
}

func (s *WalletServerAPI) GetVmxcTxHistory(ctx context.Context, req *api.GetVmxcTxHistoryRequest) (*api.GetVmxcTxHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetVmxcTxHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetVmxcTxHistory(ctx, &m2m_api.GetVmxcTxHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetVmxcTxHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	txHist := api.GetVmxcTxHistoryResponse.GetTxHistory(&resp.TxHistory)

	return &api.GetVmxcTxHistoryResponse{
		Count:       resp.Count,
		TxHistory:   txHist,
		UserProfile: UserProfile,
	}, nil
}

func (s *WalletServerAPI) GetWalletUsageHist(ctx context.Context, req *api.GetWalletUsageHistRequest) (*api.GetWalletUsageHistResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWalletUsageHist")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetWalletUsageHist(ctx, &m2m_api.GetWalletUsageHistRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetWalletUsageHistResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	walletUsageHist := api.GetWalletUsageHistResponse.GetWalletUsageHis(&resp.WalletUsageHis)

	return &api.GetWalletUsageHistResponse{
		WalletUsageHis: walletUsageHist,
		UserProfile:    UserProfile,
		Count:          resp.Count,
	}, nil
}

func (s *WalletServerAPI) GetDlPrice(ctx context.Context, req *api.GetDownLinkPriceRequest) (*api.GetDownLinkPriceResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDlPrice")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDlPrice(ctx, &m2m_api.GetDownLinkPriceRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.GetDownLinkPriceResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDownLinkPriceResponse{
		DownLinkPrice: resp.DownLinkPrice,
		UserProfile:   UserProfile,
	}, nil
}
