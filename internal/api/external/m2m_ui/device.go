package m2m_ui

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	m2m_api "github.com/mxc-foundation/lpwan-app-server/api/m2m_server"
	api "github.com/mxc-foundation/lpwan-app-server/api/m2m_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceServerAPI struct {
	validator auth.Validator
}

func NewDeviceServerAPI(validator auth.Validator) *DeviceServerAPI {
	return &DeviceServerAPI{
		validator: validator,
	}
}

func (s *DeviceServerAPI) GetDeviceList(ctx context.Context, req *api.GetDeviceListRequest) (*api.GetDeviceListResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceList")

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceList(ctx, &m2m_api.GetDeviceListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetDeviceListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	//devProfile := api.GetDeviceListResponse.GetDevProfile(resp.DevProfile)
	dvProfiles := api.GetDeviceListResponse{}.DevProfile

	for _, v := range resp.DevProfile {
		dvProfile  := api.DeviceProfile{}
		dvProfile.Id = v.Id
		dvProfile.DevEui = v.DevEui
		dvProfile.FkWallet = v.FkWallet
		dvMode := api.DeviceMode(api.DeviceMode_value[string(v.Mode)])
		dvProfile.Mode = dvMode
		dvProfile.CreatedAt = v.CreatedAt
		dvProfile.LastSeenAt = v.LastSeenAt
		dvProfile.ApplicationId = v.ApplicationId
		dvProfile.Name = v.Name

		dvProfiles = append(dvProfiles, &dvProfile)
	}

	return &api.GetDeviceListResponse{
		DevProfile:  dvProfiles,
		Count:       resp.Count,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceProfile(ctx context.Context, req *api.GetDeviceProfileRequest) (*api.GetDeviceProfileResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceProfile")

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceProfile(ctx, &m2m_api.GetDeviceProfileRequest{
		OrgId: req.OrgId,
		DevId: req.DevId,
	})
	if err != nil {
		return &api.GetDeviceProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	//devProfile := api.GetDeviceProfileResponse.GetDevProfile(&resp.DevProfile)
	dvProfile := api.GetDeviceProfileResponse{}.DevProfile
	dvProfile.Id = resp.DevProfile.Id
	dvProfile.ApplicationId = resp.DevProfile.ApplicationId
	dvProfile.CreatedAt = resp.DevProfile.CreatedAt
	dvProfile.DevEui = resp.DevProfile.DevEui
	dvProfile.FkWallet = resp.DevProfile.FkWallet
	dvProfile.LastSeenAt = resp.DevProfile.LastSeenAt
	dvProfile.Mode = api.DeviceMode(api.DeviceMode_value[string(resp.DevProfile.Mode)])
	dvProfile.Name = resp.DevProfile.Name

	return &api.GetDeviceProfileResponse{
		DevProfile:  dvProfile,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) GetDeviceHistory(ctx context.Context, req *api.GetDeviceHistoryRequest) (*api.GetDeviceHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetDeviceHistory")

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetDeviceHistory(ctx, &m2m_api.GetDeviceHistoryRequest{
		OrgId:  req.OrgId,
		DevId:  req.DevId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		return &api.GetDeviceHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetDeviceHistoryResponse{
		DevHistory:  resp.DevHistory,
		UserProfile: &prof,
	}, nil
}

func (s *DeviceServerAPI) SetDeviceMode(ctx context.Context, req *api.SetDeviceModeRequest) (*api.SetDeviceModeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/SetDeviceMode")

	prof, err := getUserProfileByJwt(ctx, req.OrgId)
	if err != nil{
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unauthenticated, err.Error())
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	devMode := m2m_api.DeviceMode(req.DevMode)

	resp, err := m2mClient.SetDeviceMode(ctx, &m2m_api.SetDeviceModeRequest{
		OrgId:   req.OrgId,
		DevId:   req.DevId,
		DevMode: devMode,
	})
	if err != nil {
		return &api.SetDeviceModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetDeviceModeResponse{
		Status:      resp.Status,
		UserProfile: &prof,
	}, nil
}
