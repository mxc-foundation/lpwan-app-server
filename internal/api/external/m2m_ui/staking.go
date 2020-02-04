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

// StakingServerAPI defines the Staking Server API structure
type StakingServerAPI struct {
	validator auth.Validator
}

// NewStakingServerAPI defines the Stagking Server API validator
func NewStakingServerAPI(validator auth.Validator) *StakingServerAPI {
	return &StakingServerAPI{
		validator: validator,
	}
}

// GetStakingPercentage defines the request and response to get staking percentage
func (s *StakingServerAPI) GetStakingPercentage(ctx context.Context, req *api.StakingPercentageRequest) (*api.StakingPercentageResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetStakingPercentage")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := api.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetStakingPercentage(ctx, &api.StakingPercentageRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakingPercentageResponse{
		StakingPercentage: resp.StakingPercentage,
	}, nil
}

// Stake defines the request and response for staking
func (s *StakingServerAPI) Stake(ctx context.Context, req *api.StakeRequest) (*api.StakeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/Stake")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.StakeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := api.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.Stake(ctx, &api.StakeRequest{
		OrgId:  req.OrgId,
		Amount: req.Amount,
	})
	if err != nil {
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

// Unstake defines the request and response to unstake
func (s *StakingServerAPI) Unstake(ctx context.Context, req *api.UnstakeRequest) (*api.UnstakeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/Unstake")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := api.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.Unstake(ctx, &api.UnstakeRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.UnstakeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

// GetActiveStakes defines the request and response to get active stakes
func (s *StakingServerAPI) GetActiveStakes(ctx context.Context, req *api.GetActiveStakesRequest) (*api.GetActiveStakesResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetActiveStakes")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := api.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetActiveStakes(ctx, &api.GetActiveStakesRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetActiveStakesResponse{
		ActStake:    resp.ActStake,
		UserProfile: &prof,
	}, nil
}

// GetStakingHistory defines the request and response to get staking history
func (s *StakingServerAPI) GetStakingHistory(ctx context.Context, req *api.StakingHistoryRequest) (*api.StakingHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetStakingHistory")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := api.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetStakingHistory(ctx, &api.StakingHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakingHistoryResponse{
		UserProfile: &prof,
		StakingHist: resp.StakingHist,
		Count:       resp.Count,
	}, nil
}
