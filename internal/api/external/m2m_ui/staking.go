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

type StakingServerAPI struct {
	validator auth.Validator
}

func NewStakingServerAPI(validator auth.Validator) *StakingServerAPI {
	return &StakingServerAPI{
		validator: validator,
	}
}

func (s *StakingServerAPI) GetStakingPercentage(ctx context.Context, req *api.StakingPercentageRequest) (*api.StakingPercentageResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetStakingPercentage")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetStakingPercentage(ctx, &m2m_api.StakingPercentageRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakingPercentageResponse{
		StakingPercentage: resp.StakingPercentage,
	}, nil
}

func (s *StakingServerAPI) Stake(ctx context.Context, req *api.StakeRequest) (*api.StakeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/Stake")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.Stake(ctx, &m2m_api.StakeRequest{
		OrgId:  req.OrgId,
		Amount: req.Amount,
	})
	if err != nil {
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.StakeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.StakeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

func (s *StakingServerAPI) Unstake(ctx context.Context, req *api.UnstakeRequest) (*api.UnstakeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/Unstake")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.Unstake(ctx, &m2m_api.UnstakeRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.UnstakeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.UnstakeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

func (s *StakingServerAPI) GetActiveStakes(ctx context.Context, req *api.GetActiveStakesRequest) (*api.GetActiveStakesResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetActiveStakes")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetActiveStakes(ctx, &m2m_api.GetActiveStakesRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	actStake := api.GetActiveStakesResponse{}.GetActStake()
	actStake.FkWallet = resp.ActStake.FkWallet
	actStake.Id = resp.ActStake.Id
	actStake.Amount = resp.ActStake.Amount
	actStake.StakeStatus = resp.ActStake.StakeStatus
	actStake.StartStakeTime = resp.ActStake.StartStakeTime
	actStake.UnstakeTime = resp.ActStake.UnstakeTime

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.GetActiveStakesResponse{
		ActStake:    actStake,
		UserProfile: &prof,
	}, nil
}

func (s *StakingServerAPI) GetStakingHistory(ctx context.Context, req *api.StakingHistoryRequest) (*api.StakingHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetStakingHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetStakingHistory(ctx, &m2m_api.StakingHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakingHist := api.StakingHistoryResponse.GetStakingHist(&resp.StakingHist)

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.StakingHistoryResponse{
		UserProfile: &prof,
		StakingHist: stakingHist,
		Count:       resp.Count,
	}, nil
}
