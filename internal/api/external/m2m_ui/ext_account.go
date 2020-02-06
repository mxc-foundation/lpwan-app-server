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

// ExtAccountServerAPI defines the ExtAccountServerAPI structure
type ExtAccountServerAPI struct {
	validator auth.Validator
}

// NewMoneyServerAPI defines the NewMoneyServerAPI validator
func NewMoneyServerAPI(validator auth.Validator) *ExtAccountServerAPI {
	return &ExtAccountServerAPI{
		validator: validator,
	}
}

// ModifyMoneyAccount defines the modify money account request and respone
func (s *ExtAccountServerAPI) ModifyMoneyAccount(ctx context.Context, req *api.ModifyMoneyAccountRequest) (*api.ModifyMoneyAccountResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/ModifyMoneyAccount")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.OrgId)
	if err != nil {
		return &api.ModifyMoneyAccountResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifyMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyClient := api.NewMoneyServiceClient(m2mClient)

	resp, err := moneyClient.ModifyMoneyAccount(ctx, &api.ModifyMoneyAccountRequest{
		OrgId:          req.OrgId,
		MoneyAbbr:      req.MoneyAbbr,
		CurrentAccount: req.CurrentAccount,
	})
	if err != nil {
		return &api.ModifyMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifyMoneyAccountResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

// GetChangeMoneyAccountHistory defines the get change money account history request and response
func (s *ExtAccountServerAPI) GetChangeMoneyAccountHistory(ctx context.Context, req *api.GetMoneyAccountChangeHistoryRequest) (*api.GetMoneyAccountChangeHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetChangeMoneyAccountHistory")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.OrgId)
	if err != nil {
		return &api.GetMoneyAccountChangeHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetMoneyAccountChangeHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyClient := api.NewMoneyServiceClient(m2mClient)

	resp, err := moneyClient.GetChangeMoneyAccountHistory(ctx, &api.GetMoneyAccountChangeHistoryRequest{
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
		UserProfile:   &prof,
	}, nil
}

// GetActiveMoneyAccount defines the get active money account request and response
func (s *ExtAccountServerAPI) GetActiveMoneyAccount(ctx context.Context, req *api.GetActiveMoneyAccountRequest) (*api.GetActiveMoneyAccountResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetActiveMoneyAccount")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.OrgId)
	if err != nil {
		return &api.GetActiveMoneyAccountResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyClient := api.NewMoneyServiceClient(m2mClient)

	resp, err := moneyClient.GetActiveMoneyAccount(ctx, &api.GetActiveMoneyAccountRequest{
		OrgId:     req.OrgId,
		MoneyAbbr: req.MoneyAbbr,
	})
	if err != nil {
		return &api.GetActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetActiveMoneyAccountResponse{
		ActiveAccount: resp.ActiveAccount,
		UserProfile:   &prof,
	}, nil
}
