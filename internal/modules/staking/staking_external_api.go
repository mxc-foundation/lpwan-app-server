package staking

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
)

// StakingServerAPI defines the Staking Server API structure
type StakingServerAPI struct{}

// NewStakingServerAPI defines the Stagking Server API Validator
func NewStakingServerAPI() *StakingServerAPI {
	return &StakingServerAPI{}
}

// GetStakingPercentage defines the request and response to get staking percentage
func (s *StakingServerAPI) GetStakingPercentage(ctx context.Context, req *api.StakingPercentageRequest) (*api.StakingPercentageResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetStakingPercentage org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.GetStakingPercentage(ctx, &pb.StakingPercentageRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakingPercentageResponse{
		StakingPercentage: resp.StakingPercentage,
	}, status.Error(codes.OK, "")
}

// Stake defines the request and response for staking
func (s *StakingServerAPI) Stake(ctx context.Context, req *api.StakeRequest) (*api.StakeResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/Stake org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.Stake(ctx, &pb.StakeRequest{
		OrgId:  req.OrgId,
		Amount: req.Amount,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.StakeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// Unstake defines the request and response to unstake
func (s *StakingServerAPI) Unstake(ctx context.Context, req *api.UnstakeRequest) (*api.UnstakeResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/Unstake org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.Unstake(ctx, &pb.UnstakeRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.UnstakeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// GetActiveStakes defines the request and response to get active stakes
func (s *StakingServerAPI) GetActiveStakes(ctx context.Context, req *api.GetActiveStakesRequest) (*api.GetActiveStakesResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetActiveStakes org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.GetActiveStakes(ctx, &pb.GetActiveStakesRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unknown, err.Error())
	}

	if resp.ActStake == nil {
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.OK, "")
	}

	return &api.GetActiveStakesResponse{
		ActStake: &api.ActiveStake{
			Id:          resp.ActStake.Id,
			Amount:      resp.ActStake.Amount,
			StakeStatus: resp.ActStake.StakeStatus,
			StartTime:   resp.ActStake.StartTime,
			EndTime:     resp.ActStake.EndTime,
		},
	}, status.Error(codes.OK, "")
}

// GetStakingRevenue returns the amount earned from staking during the specified period
func (s *StakingServerAPI) GetStakingRevenue(ctx context.Context, req *api.StakingRevenueRequest) (*api.StakingRevenueResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetStakingRevenue org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingRevenueResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.GetStakingRevenue(ctx, &pb.StakingRevenueRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
		From:     req.From,
		Till:     req.Till,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Error(codes.Internal, "couldn't get response from m2m")
	}
	return &api.StakingRevenueResponse{Amount: resp.Amount}, nil
}

// GetStakingHistory defines the request and response to get staking history
func (s *StakingServerAPI) GetStakingHistory(ctx context.Context, req *api.StakingHistoryRequest) (*api.StakingHistoryResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetStakingHistory org=%d", req.OrgId)

	if err := NewValidator().IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	stakeClient, err := m2mcli.GetStakingServiceClient()
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := stakeClient.GetStakingHistory(ctx, &pb.StakingHistoryRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
		From:     req.From,
		Till:     req.Till,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var stakeHistoryList []*api.StakingHistory
	for _, item := range resp.StakingHist {
		var stake *api.Stake
		if st := item.Stake; st != nil {
			stake = &api.Stake{
				Id:        st.Id,
				Amount:    st.Amount,
				Active:    st.Active,
				StartTime: st.StartTime,
				EndTime:   st.EndTime,
			}
		}
		stakeHistory := &api.StakingHistory{
			Timestamp: item.Timestamp,
			Amount:    item.Amount,
			Type:      item.Type,
			Stake:     stake,
		}

		stakeHistoryList = append(stakeHistoryList, stakeHistory)
	}

	return &api.StakingHistoryResponse{
		StakingHist: stakeHistoryList,
	}, status.Error(codes.OK, "")
}
