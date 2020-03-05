package external

import (
	"context"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SupernodeServerAPI defines the supernode server api structure
type SupernodeServerAPI struct {
	validator auth.Validator
}

// NewSupernodeServerAPI validates the new supernode server api
func NewSupernodeServerAPI(validator auth.Validator) *SupernodeServerAPI {
	return &SupernodeServerAPI{
		validator: validator,
	}
}

// AddSuperNodeMoneyAccount defines the money account addition request and response for the supernode
func (s *SupernodeServerAPI) AddSuperNodeMoneyAccount(ctx context.Context, req *api.AddSuperNodeMoneyAccountRequest) (*api.AddSuperNodeMoneyAccountResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/AddSuperNodeMoneyAccount")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.OrgId)
	if err != nil {
		return &api.AddSuperNodeMoneyAccountResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.AddSuperNodeMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	supernodeClient := api.NewSuperNodeServiceClient(m2mClient)

	resp, err := supernodeClient.AddSuperNodeMoneyAccount(ctx, &api.AddSuperNodeMoneyAccountRequest{
		MoneyAbbr:   req.MoneyAbbr,
		AccountAddr: req.AccountAddr,
		OrgId:       req.OrgId,
	})
	if err != nil {
		return &api.AddSuperNodeMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.AddSuperNodeMoneyAccountResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}

// GetSuperNodeActiveMoneyAccount defines the active money account request and response for the supernode
func (s *SupernodeServerAPI) GetSuperNodeActiveMoneyAccount(ctx context.Context, req *api.GetSuperNodeActiveMoneyAccountRequest) (*api.GetSuperNodeActiveMoneyAccountResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetSuperNodeActiveMoneyAccount")

	prof, err := getUserProfileByJwt(ctx, s.validator, req.OrgId)
	if err != nil {
		return &api.GetSuperNodeActiveMoneyAccountResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetSuperNodeActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	supernodeClient := api.NewSuperNodeServiceClient(m2mClient)

	resp, err := supernodeClient.GetSuperNodeActiveMoneyAccount(ctx, &api.GetSuperNodeActiveMoneyAccountRequest{
		MoneyAbbr: req.MoneyAbbr,
		OrgId:     req.OrgId,
	})
	if err != nil {
		return &api.GetSuperNodeActiveMoneyAccountResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetSuperNodeActiveMoneyAccountResponse{
		SupernodeActiveAccount: resp.SupernodeActiveAccount,
		UserProfile:            &prof,
	}, nil
}
