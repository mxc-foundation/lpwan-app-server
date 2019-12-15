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

	username, err := auth.JWTValidator{}.GetUsername(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	//userProfile := api.GetDeviceListResponse.GetUserProfile(prof)
	userProfile := api.ProfileResponse{
		User: &api.User{
			Id:         string(prof.User.ID),
			Username:   prof.User.Username,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &api.ProfileSettings{
			DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
		},
	}

	for _, org := range prof.Organizations {
		row := api.OrganizationLink{
			OrganizationId:   org.ID,
			OrganizationName: org.Name,
			IsAdmin:          org.IsAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		userProfile.Organizations = append(userProfile.Organizations, &row)
	}

	return &api.GetGatewayListResponse{
		GwProfile:   gwProfiles,
		Count:       resp.Count,
		UserProfile: &userProfile,
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

	username, err := auth.JWTValidator{}.GetUsername(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	//userProfile := api.GetDeviceListResponse.GetUserProfile(prof)
	userProfile := api.ProfileResponse{
		User: &api.User{
			Id:         string(prof.User.ID),
			Username:   prof.User.Username,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &api.ProfileSettings{
			DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
		},
	}

	for _, org := range prof.Organizations {
		row := api.OrganizationLink{
			OrganizationId:   org.ID,
			OrganizationName: org.Name,
			IsAdmin:          org.IsAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		userProfile.Organizations = append(userProfile.Organizations, &row)
	}

	return &api.GetGatewayProfileResponse{
		GwProfile:   gwProfile,
		UserProfile: &userProfile,
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

	username, err := auth.JWTValidator{}.GetUsername(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	//userProfile := api.GetDeviceListResponse.GetUserProfile(prof)
	userProfile := api.ProfileResponse{
		User: &api.User{
			Id:         string(prof.User.ID),
			Username:   prof.User.Username,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &api.ProfileSettings{
			DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
		},
	}

	for _, org := range prof.Organizations {
		row := api.OrganizationLink{
			OrganizationId:   org.ID,
			OrganizationName: org.Name,
			IsAdmin:          org.IsAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		userProfile.Organizations = append(userProfile.Organizations, &row)
	}

	return &api.GetGatewayHistoryResponse{
		GwHistory:   resp.GwHistory,
		UserProfile: &userProfile,
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

	username, err := auth.JWTValidator{}.GetUsername(ctx)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	// Get the user id based on the username.
	user, err := storage.GetUserByUsername(ctx, storage.DB(), username)
	if nil != err {
		return nil, helpers.ErrToRPCError(err)
	}

	prof, err := storage.GetProfile(ctx, storage.DB(), user.ID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	//userProfile := api.GetDeviceListResponse.GetUserProfile(prof)
	userProfile := api.ProfileResponse{
		User: &api.User{
			Id:         string(prof.User.ID),
			Username:   prof.User.Username,
			SessionTtl: prof.User.SessionTTL,
			IsAdmin:    prof.User.IsAdmin,
			IsActive:   prof.User.IsActive,
		},
		Settings: &api.ProfileSettings{
			DisableAssignExistingUsers: auth.DisableAssignExistingUsers,
		},
	}

	for _, org := range prof.Organizations {
		row := api.OrganizationLink{
			OrganizationId:   org.ID,
			OrganizationName: org.Name,
			IsAdmin:          org.IsAdmin,
		}

		row.CreatedAt, err = ptypes.TimestampProto(org.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(org.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		userProfile.Organizations = append(userProfile.Organizations, &row)
	}

	return &api.SetGatewayModeResponse{
		Status:      resp.Status,
		UserProfile: &userProfile,
	}, nil
}
