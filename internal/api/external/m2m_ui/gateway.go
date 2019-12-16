package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatewayServerAPI struct {
	validator auth.Validator
}

func NewGatewayServerAPI(validator auth.Validator) *GatewayServerAPI {
	return &GatewayServerAPI{
		validator: validator,
	}
}

func (s *GatewayServerAPI) GetGatewayList(ctx context.Context, req *api.GetGatewayListRequest) (*api.GetGatewayListResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayList")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetGatewayList(ctx, &m2m_api.GetGatewayListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	//gwProfile := api.GetGatewayListResponse.GetGwProfile(&resp.GwProfile)
	gwProfiles := api.GetGatewayListResponse{}.GwProfile
	for _, v := range resp.GwProfile {
		gwProfile := api.GatewayProfile{}
		gwProfile.Mode = api.GatewayMode(api.DeviceMode_value[string(v.Mode)])
		gwProfile.Name = v.Name
		gwProfile.LastSeenAt = v.LastSeenAt
		gwProfile.FkWallet = v.FkWallet
		gwProfile.Id = v.Id
		gwProfile.CreateAt = v.CreateAt
		gwProfile.Description = v.Description
		gwProfile.FkGwNs = v.FkGwNs
		gwProfile.Mac = v.Mac
		gwProfile.OrgId = v.OrgId

		gwProfiles = append(gwProfiles, &gwProfile)
	}

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.GetGatewayListResponse{
		GwProfile:   gwProfiles,
		Count:       resp.Count,
		UserProfile: &prof,
	}, nil
}

func (s *GatewayServerAPI) GetGatewayProfile(ctx context.Context, req *api.GetGatewayProfileRequest) (*api.GetGatewayProfileResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayProfile")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetGatewayProfile(ctx, &m2m_api.GetGatewayProfileRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwProfile := api.GetGatewayProfileResponse{}.GwProfile
	gwProfile.OrgId = resp.GwProfile.OrgId
	gwProfile.Mac = resp.GwProfile.Mac
	gwProfile.FkGwNs = resp.GwProfile.FkGwNs
	gwProfile.Description = resp.GwProfile.Description
	gwProfile.CreateAt = resp.GwProfile.CreateAt
	gwProfile.Id = resp.GwProfile.Id
	gwProfile.FkWallet = resp.GwProfile.FkWallet
	gwProfile.LastSeenAt = resp.GwProfile.LastSeenAt
	gwProfile.Name = resp.GwProfile.Name
	gwProfile.Mode = api.GatewayMode(api.DeviceMode_value[string(resp.GwProfile.Mode)])

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetGatewayProfileResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.GetGatewayProfileResponse{
		GwProfile:   gwProfile,
		UserProfile: &prof,
	}, nil
}

func (s *GatewayServerAPI) GetGatewayHistory(ctx context.Context, req *api.GetGatewayHistoryRequest) (*api.GetGatewayHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetGatewayHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetGatewayHistory(ctx, &m2m_api.GetGatewayHistoryRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.GetGatewayHistoryResponse{
		GwHistory:   resp.GwHistory,
		UserProfile: &prof,
	}, nil
}

func (s *GatewayServerAPI) SetGatewayMode(ctx context.Context, req *api.SetGatewayModeRequest) (*api.SetGatewayModeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/SetGatewayMode")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwMode := m2m_api.GatewayMode(req.GwMode)

	resp, err := m2mClient.SetGatewayMode(ctx, &m2m_api.SetGatewayModeRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		GwMode: gwMode,
	})
	if err != nil {
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	return &api.SetGatewayModeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}
