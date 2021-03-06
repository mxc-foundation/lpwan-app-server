package staking

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

// Server defines the Staking Server API structure
type Server struct {
	stakingCli pb.StakingServiceClient
	auth       auth.Authenticator
}

// NewServer creates a new staking server
func NewServer(stakingCli pb.StakingServiceClient, auth auth.Authenticator) *Server {
	return &Server{
		stakingCli: stakingCli,
		auth:       auth,
	}
}

// GetStakingPercentage defines the request and response to get staking percentage
func (s *Server) GetStakingPercentage(ctx context.Context, req *api.StakingPercentageRequest) (*api.StakingPercentageResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/GetStakingPercentage")

	resp, err := s.stakingCli.GetStakingPercentage(ctx, &pb.StakingPercentageRequest{
		Currency: req.Currency,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	spr := &api.StakingPercentageResponse{
		StakingInterest: resp.StakingInterest,
	}
	for _, boost := range resp.LockBoosts {
		spr.LockBoosts = append(spr.LockBoosts, &api.Boost{
			LockPeriods: boost.LockPeriods,
			Boost:       boost.Boost,
		})
	}
	return spr, nil
}

// Stake defines the request and response for staking
func (s *Server) Stake(ctx context.Context, req *api.StakeRequest) (*api.StakeResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/Stake org=%d", req.OrgId)
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.Stake(ctx, &pb.StakeRequest{
		OrgId:       req.OrgId,
		Currency:    req.Currency,
		Amount:      req.Amount,
		LockPeriods: req.LockPeriods,
		Boost:       req.Boost,
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
func (s *Server) Unstake(ctx context.Context, req *api.UnstakeRequest) (*api.UnstakeResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/Unstake org=%d", req.OrgId)

	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.Unstake(ctx, &pb.UnstakeRequest{
		OrgId:   req.OrgId,
		StakeId: req.StakeId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.UnstakeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}

// GetActiveStakes defines the request and response to get active stakes
func (s *Server) GetActiveStakes(ctx context.Context, req *api.GetActiveStakesRequest) (*api.GetActiveStakesResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/GetActiveStakes org=%d", req.OrgId)

	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.GetActiveStakes(ctx, &pb.GetActiveStakesRequest{
		OrgId:    req.OrgId,
		Currency: req.Currency,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	gasr := &api.GetActiveStakesResponse{}
	for _, stake := range resp.ActStake {
		gasr.ActStake = append(gasr.ActStake,
			&api.Stake{
				Id:        stake.Id,
				StartTime: stake.StartTime,
				EndTime:   stake.EndTime,
				Amount:    stake.Amount,
				Active:    stake.Active,
				LockTill:  stake.LockTill,
				Boost:     stake.Boost,
				Revenue:   stake.Revenue,
			})
	}
	return gasr, nil
}

// StakeInfo returns information about a particular stake, including all
// received stake revenues
func (s *Server) StakeInfo(ctx context.Context, req *api.StakeInfoRequest) (*api.StakeInfoResponse, error) {
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.StakeInfo(ctx, &pb.StakeInfoRequest{
		OrgId:   req.OrgId,
		StakeId: req.StakeId,
	})
	if err != nil {
		log.WithError(err).Error("m2m StakeInfo returned an error")
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	stakeInfo := &api.StakeInfoResponse{}
	stakeInfo.Stake = &api.Stake{
		Id:        resp.Stake.Id,
		StartTime: resp.Stake.StartTime,
		EndTime:   resp.Stake.EndTime,
		Amount:    resp.Stake.Amount,
		Active:    resp.Stake.Active,
		LockTill:  resp.Stake.LockTill,
		Boost:     resp.Stake.Boost,
		Revenue:   resp.Stake.Revenue,
	}
	for _, revenue := range resp.Revenues {
		stakeInfo.Revenues = append(stakeInfo.Revenues, &api.StakeRevenue{
			Time:   revenue.Time,
			Amount: revenue.Amount,
		},
		)
	}
	return stakeInfo, nil
}

// GetStakingRevenue returns the amount earned from staking during the specified period
func (s *Server) GetStakingRevenue(ctx context.Context, req *api.StakingRevenueRequest) (*api.StakingRevenueResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/GetStakingRevenue org=%d", req.OrgId)
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.GetStakingRevenue(ctx, &pb.StakingRevenueRequest{
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
func (s *Server) GetStakingHistory(ctx context.Context, req *api.StakingHistoryRequest) (*api.StakingHistoryResponse, error) {
	logInfo := fmt.Sprintf("api/appserver_serves_ui/GetStakingHistory org=%d", req.OrgId)
	cred, err := s.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	resp, err := s.stakingCli.GetStakingHistory(ctx, &pb.StakingHistoryRequest{
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
				LockTill:  st.LockTill,
				Boost:     st.Boost,
				Revenue:   st.Revenue,
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
