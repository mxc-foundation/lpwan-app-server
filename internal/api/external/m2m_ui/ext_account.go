package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	log "github.com/sirupsen/logrus"
)

type ExtAccountServerAPI struct {
	validator auth.Validator
}

func NewMoneyServerAPI(validator auth.Validator) *ExtAccountServerAPI {
	return &ExtAccountServerAPI{
		validator: validator,
	}
}

func (s *ExtAccountServerAPI) ModifyMoneyAccount(ctx context.Context, req *api.ModifyMoneyAccountRequest) (*api.ModifyMoneyAccountResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/ModifyMoneyAccount")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifyMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.ModifyMoneyAccount(ctx, &m2m_api.ModifyMoneyAccountRequest{
		OrgId:          req.OrgId,
		MoneyAbbr:      req.MoneyAbbr,
		CurrentAccount: req.CurrentAccount,
	})
	if err != nil {
		return &api.ModifyMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifyMoneyAccountResponse{
		Status:      resp.Status,
		UserProfile: resp.UserProfile,
	}, nil
}

func (s *ExtAccountServerAPI) GetChangeMoneyAccountHistory(ctx context.Context, req *api.GetMoneyAccountChangeHistoryRequest) (*api.GetMoneyAccountChangeHistoryResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetChangeMoneyAccountHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetMoneyAccountChangeHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetChangeMoneyAccountHistory(ctx, &m2m_api.GetMoneyAccountChangeHistoryRequest{
		OrgId:     req.OrgId,
		Offset:    req.Offset,
		Limit:     req.Limit,
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetMoneyAccountChangeHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetMoneyAccountChangeHistoryResponse{
		Count:         resp.Count,
		ChangeHistory: resp.ChangeHistory,
		UserProfile:   resp.UserProfile,
	}, nil
}

func (s *ExtAccountServerAPI) GetActiveMoneyAccount(ctx context.Context, req *api.GetActiveMoneyAccountRequest) (*api.GetActiveMoneyAccountResponse, error) {
	if err := s.validator.Validate(ctx,
		auth.ValidateActiveUser()); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	log.WithField("orgId", req.OrgId).Info("grpc_api/GetActiveMoneyAccount")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetActiveMoneyAccount(ctx, &m2m_api.GetActiveMoneyAccountRequest{
		OrgId:     req.OrgId,
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetActiveMoneyAccountResponse{
		ActiveAccount: resp.ActiveAccount,
		UserProfile:   resp.UserProfile,
	}, nil
}