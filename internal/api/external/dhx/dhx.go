package dhx

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
)

// Server defines the DHX service Server API structure
type Server struct {
	store  Store
	dhxCli pb.DHXServiceClient
	auth   auth.Authenticator
}

// NewServer creates a new DHX service server
func NewServer(cli pb.DHXServiceClient, auth auth.Authenticator, store Store) *Server {
	return &Server{
		dhxCli: cli,
		auth:   auth,
		store:  store,
	}
}

// Store defines db APIs for dhx service
type Store interface {
	// GetOnlineGatewayCount returns count of gateways that meet certain requirements:
	// 1. online (last_seen_at is not earlier than 10 mins ago)
	// 2. must be matchx new model (sn and modle are not empty string)
	GetOnlineGatewayCount(ctx context.Context, orgID int64) (int, error)
	// GetOrganizationName returns the name of the organization with the given ID
	GetOrganizationName(ctx context.Context, orgID int64) (string, error)
}

// DHXCreateStake creates new dhx stake
func (a *Server) DHXCreateStake(ctx context.Context, req *api.DHXCreateStakeRequest) (*api.DHXCreateStakeResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	res, err := a.dhxCli.DHXCreateStake(ctx, &pb.DHXCreateStakeRequest{
		CouncilId:  req.CouncilId,
		Amount:     req.Amount,
		Currency:   req.Currency,
		LockMonths: req.LockMonths,
		Boost:      req.Boost,
		OrgId:      req.OrganizationId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.DHXCreateStakeResponse{StakeId: res.StakeId}, nil
}

// DHXCreateCouncil creates new council chair
func (a *Server) DHXCreateCouncil(ctx context.Context, req *api.DHXCreateCouncilRequest) (*api.DHXCreateCouncilResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// check whether council has 5+ gateways
	count, err := a.store.GetOnlineGatewayCount(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if count < 5 {
		return nil, status.Error(codes.FailedPrecondition,
			"at least 5 online gateways are required to be registered for creating council chair")
	}

	name, err := a.store.GetOrganizationName(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't get organization name: %v", err)
	}

	res, err := a.dhxCli.DHXCreateCouncil(ctx, &pb.DHXCreateCouncilRequest{
		OrgId:      req.OrganizationId,
		Amount:     req.Amount,
		Currency:   req.Currency,
		LockMonths: req.LockMonths,
		Boost:      req.Boost,
		Name:       name,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.DHXCreateCouncilResponse{CouncilId: res.CouncilId, StakeId: res.StakeId}, nil
}

// DHXListCouncils lists all existing councils
func (a *Server) DHXListCouncils(ctx context.Context, req *api.DHXListCouncilsRequest) (*api.DHXListCouncilsResponse, error) {
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	res, err := a.dhxCli.DHXListCouncils(ctx, &pb.DHXListCouncilsRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	var response []*api.Council
	for _, v := range res.Council {
		item := api.Council{
			Id:             v.Id,
			ChairOrgId:     v.ChairOrgId,
			Name:           v.Name,
			LastPaidDate:   v.LastPaidDate,
			LastDhxRevenue: v.LastDhxRevenue,
			LastMpower:     v.LastMpower,
		}

		response = append(response, &item)
	}

	return &api.DHXListCouncilsResponse{Council: response}, nil
}

// DHXListStakes lists all dhx stake records
func (a *Server) DHXListStakes(ctx context.Context, req *api.DHXListStakesRequest) (*api.DHXListStakesResponse, error) {
	if req.OrganizationId != 0 {
		cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
		if !cred.IsOrgUser && !cred.IsGlobalAdmin {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

	} else if req.ChairOrgId != 0 {
		cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.ChairOrgId))
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
		if !cred.IsOrgUser && !cred.IsGlobalAdmin {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
	} else {
		// both req.OrganizationId and req.ChairOrgId are 0, only global admin is allowed to make the request
		cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
		if !cred.IsGlobalAdmin {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
	}

	res, err := a.dhxCli.DHXListStakes(ctx, &pb.DHXListStakesRequest{
		ChairOrgId: req.ChairOrgId,
		OrgId:      req.OrganizationId,
	})
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}

	var response []*api.DHXStake
	for _, v := range res.Stake {
		item := api.DHXStake{
			Id:             v.Id,
			OrganizationId: v.OrgId,
			Amount:         v.Amount,
			Currency:       v.Currency,
			CouncilId:      v.CouncilId,
			CouncilName:    v.CouncilName,
			Created:        v.Created,
			LockTill:       v.LockTill,
			Boost:          v.Boost,
			Closed:         v.Closed,
			DhxMined:       v.DhxMined,
		}

		response = append(response, &item)
	}

	return &api.DHXListStakesResponse{Stake: response}, nil
}

// DHXGetLastMining returns info about the last paid day of DHX mining
func (a *Server) DHXGetLastMining(ctx context.Context, req *api.DHXGetLastMiningRequest) (*api.DHXGetLastMiningResponse, error) {
	authOpts := auth.NewOptions()
	if req.OrgId != 0 {
		authOpts = authOpts.WithOrgID(req.OrgId)
	}
	cred, err := a.auth.GetCredentials(ctx, authOpts)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if req.OrgId != 0 && !cred.IsOrgUser {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	res, err := a.dhxCli.DHXGetLastMining(ctx, &pb.DHXGetLastMiningRequest{OrgId: req.OrgId})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.DHXGetLastMiningResponse{
		Date:               res.Date,
		MiningPower:        res.MiningPower,
		DhxAmount:          res.DhxAmount,
		OrgId:              res.OrgId,
		OrgMiningPower:     res.OrgMiningPower,
		OrgDhxLimit:        res.OrgDhxLimit,
		OrgDhxAmount:       res.OrgDhxAmount,
		CouncilId:          res.CouncilId,
		CouncilName:        res.CouncilName,
		CouncilMiningPower: res.CouncilMiningPower,
		CouncilDhxAmount:   res.CouncilDhxAmount,
	}, nil
}

// DHXBond bonds the specified amount of DHX for DHX mining
func (a *Server) DHXBond(ctx context.Context, req *api.DHXBondRequest) (*api.DHXBondResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsOrgAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	_, err = a.dhxCli.DHXBond(ctx, &pb.DHXBondRequest{
		OrgId:  req.OrgId,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.DHXBondResponse{}, nil
}

// DHXUnbond unbonds the specified amount of DHX
func (a *Server) DHXUnbond(ctx context.Context, req *api.DHXUnbondRequest) (*api.DHXUnbondResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsOrgAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	_, err = a.dhxCli.DHXUnbond(ctx, &pb.DHXUnbondRequest{
		OrgId:  req.OrgId,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.DHXUnbondResponse{}, nil
}

// DHXBondInfo returns info about amounts of DHX bonded, unbonding and cooling off
func (a *Server) DHXBondInfo(ctx context.Context, req *api.DHXBondInfoRequest) (*api.DHXBondInfoResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsOrgAdmin {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	res, err := a.dhxCli.DHXBondInfo(ctx, &pb.DHXBondInfoRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	var dhxUnbonding []*api.DHXUnbonding
	for _, ub := range res.DhxUnbonding {
		dhxUnbonding = append(dhxUnbonding, &api.DHXUnbonding{
			Amount:  ub.Amount,
			Created: ub.Created,
		})
	}
	var dhxCoolingOff []*api.DHXCoolingOff
	for _, co := range res.DhxCoolingOff {
		dhxCoolingOff = append(dhxCoolingOff, &api.DHXCoolingOff{
			Amount:  co.Amount,
			Created: co.Created,
		})
	}

	return &api.DHXBondInfoResponse{
		OrgId:              res.OrgId,
		DhxBonded:          res.DhxBonded,
		DhxUnbonding:       dhxUnbonding,
		DhxUnbondingTotal:  res.DhxUnbondingTotal,
		DhxCoolingOff:      dhxCoolingOff,
		DhxCoolingOffTotal: res.DhxCoolingOffTotal,
	}, nil
}
