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

// GatewayServerAPI defines the Gateway Server API structure
type GatewayServerAPI struct {
	validator auth.Validator
}

// NewGatewayServerAPI defines the Gateway Server API validator
func NewGatewayServerAPI(validator auth.Validator) *GatewayServerAPI {
	return &GatewayServerAPI{
		validator: validator,
	}
}

// GetGatewayList defines the get Gateway list request and response
func (s *GatewayServerAPI) GetGatewayList(ctx context.Context, req *api.GetGatewayListRequest) (*api.GetGatewayListResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayList")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := api.NewGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayList(ctx, &api.GetGatewayListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGatewayListResponse{
		GwProfile:   resp.GwProfile,
		Count:       resp.Count,
		UserProfile: &prof,
	}, nil
}

// GetGatewayProfile defines the get Gateway Profile request and response
func (s *GatewayServerAPI) GetGatewayProfile(ctx context.Context, req *api.GetGatewayProfileRequest) (*api.GetGatewayProfileResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayProfile")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := api.NewGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayProfile(ctx, &api.GetGatewayProfileRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGatewayProfileResponse{
		GwProfile:   resp.GwProfile,
		UserProfile: &prof,
	}, nil
}

// GetGatewayHistory defines the get Gateway History request and response
func (s *GatewayServerAPI) GetGatewayHistory(ctx context.Context, req *api.GetGatewayHistoryRequest) (*api.GetGatewayHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayHistory")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := api.NewGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayHistory(ctx, &api.GetGatewayHistoryRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGatewayHistoryResponse{
		GwHistory:   resp.GwHistory,
		UserProfile: &prof,
	}, nil
}

// SetGatewayMode defines the set Gateway mode request and response
func (s *GatewayServerAPI) SetGatewayMode(ctx context.Context, req *api.SetGatewayModeRequest) (*api.SetGatewayModeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/SetGatewayMode")

	prof, err := getUserProfileByJwt(s.validator, ctx, req.OrgId)
	if err != nil {
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := api.NewGatewayServiceClient(m2mClient)

	resp, err := gwClient.SetGatewayMode(ctx, &api.SetGatewayModeRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		GwMode: req.GwMode,
	})
	if err != nil {
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetGatewayModeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}
