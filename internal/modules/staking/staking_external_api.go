package staking

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
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

	// verify if user is global admin
	userIsAdmin, err := NewValidator().GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingPercentageResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
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
	userIsAdmin, err := NewValidator().GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
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
	userIsAdmin, err := NewValidator().GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.UnstakeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
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
	userIsAdmin, err := NewValidator().GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetActiveStakesResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
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
	cred, err := s.validator.GetCredentials(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}
	if err := cred.IsOrgAdmin(ctx, req.OrgId); err != nil {
		return nil, status.Error(codes.PermissionDenied, "must be an organization admin")
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	stakeClient := m2mServer.NewStakingServiceClient(m2mClient)
	resp, err := stakeClient.GetStakingRevenue(ctx, &m2mServer.StakingRevenueRequest{
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

	// verify if user is global admin
	userIsAdmin, err := NewValidator().GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.StakingHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
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
