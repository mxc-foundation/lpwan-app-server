package m2m_ui

import (
	"context"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
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

	gwProfile := api.GetGatewayListResponse.GetGwProfile(&resp.GwProfile)

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetGatewayListResponse{}, err
	}

	userProfile := api.GetGatewayListResponse.GetUserProfile(getUserProfile)

	return &api.GetGatewayListResponse{
		GwProfile:   gwProfile,
		Count:       resp.Count,
		UserProfile: userProfile,
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

	gwProfile := api.GetGatewayProfileResponse.GetGwProfile(&resp.GwProfile)

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetGatewayProfileResponse{}, err
	}

	userProfile := api.GetGatewayProfileResponse.GetUserProfile(getUserProfile)

	return &api.GetGatewayProfileResponse{
		GwProfile:   gwProfile,
		UserProfile: userProfile,
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.GetGatewayHistoryResponse{}, err
	}

	userProfile := api.GetGatewayListResponse.GetUserProfile(getUserProfile)

	return &api.GetGatewayHistoryResponse{
		GwHistory:   resp.GwHistory,
		UserProfile: userProfile,
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

	getUserProfile, err := external.InternalUserAPI{}.Profile(ctx, nil)
	if err != nil {
		log.WithError(err).Error("Cannot get userprofile")
		return &api.SetGatewayModeResponse{}, err
	}

	userProfile := api.SetGatewayModeResponse.GetUserProfile(getUserProfile)

	return &api.SetGatewayModeResponse{
		Status:      resp.Status,
		UserProfile: userProfile,
	}, nil
}
