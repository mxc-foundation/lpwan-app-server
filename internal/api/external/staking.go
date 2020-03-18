package external

import (
	"context"
	"fmt"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m_serves_appserver"
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
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetStakingPercentage org=%d", req.OrgId)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.StakingPercentageResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetStakingPercentage(ctx, &m2mServer.StakingPercentageRequest{
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

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.StakeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	} else {
		return &api.StakeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.Stake(ctx, &m2mServer.StakeRequest{
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

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.UnstakeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.UnstakeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	} else {
		return &api.UnstakeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.UnstakeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.Unstake(ctx, &m2mServer.UnstakeRequest{
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

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetActiveStakes(ctx, &m2mServer.GetActiveStakesRequest{
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
			Id:             resp.ActStake.Id,
			FkWallet:       resp.ActStake.FkWallet,
			Amount:         resp.ActStake.Amount,
			StakeStatus:    resp.ActStake.StakeStatus,
			StartStakeTime: resp.ActStake.StartStakeTime,
			UnstakeTime:    resp.ActStake.UnstakeTime,
		},
	}, status.Error(codes.OK, "")
}

// GetStakingHistory defines the request and response to get staking history
func (s *StakingServerAPI) GetStakingHistory(ctx context.Context, req *api.StakingHistoryRequest) (*api.StakingHistoryResponse, error) {
	logInfo, _ := fmt.Printf("api/appserver_serves_ui/GetStakingHistory org=%d", req.OrgId)

	// verify if user is global admin
	userIsAdmin, err := s.validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.validator.Validate(ctx, auth.ValidateOrganizationAccess(auth.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.StakingHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)

	resp, err := stakeClient.GetStakingHistory(ctx, &m2mServer.StakingHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var stakeHistoryList []*api.GetStakingHistory
	for _, item := range resp.StakingHist {
		stakeHistory := &api.GetStakingHistory{
			StakeAmount:   item.StakeAmount,
			Start:         item.Start,
			End:           item.End,
			RevMonth:      item.RevMonth,
			NetworkIncome: item.NetworkIncome,
			MonthlyRate:   item.MonthlyRate,
			Revenue:       item.Revenue,
			Balance:       item.Balance,
		}

		stakeHistoryList = append(stakeHistoryList, stakeHistory)
	}

	return &api.StakingHistoryResponse{
		StakingHist: stakeHistoryList,
		Count:       resp.Count,
	}, status.Error(codes.OK, "")
}
