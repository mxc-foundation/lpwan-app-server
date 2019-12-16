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

type WithdrawServerAPI struct {
	validator auth.Validator
}

func NewWithdrawServerAPI(validator auth.Validator) *WithdrawServerAPI {
	return &WithdrawServerAPI{
		validator: validator,
	}
}

func (s *WithdrawServerAPI) ModifyWithdrawFee(ctx context.Context, req *api.ModifyWithdrawFeeRequest) (*api.ModifyWithdrawFeeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/ModifyWithdrawFee")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.ModifyWithdrawFee(ctx, &m2m_api.ModifyWithdrawFeeRequest{
		MoneyAbbr:   moneyAbbr,
		WithdrawFee: req.WithdrawFee,
		OrgId:       req.OrgId,
	})
	if err != nil {
		return &api.ModifyWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
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

	return &api.ModifyWithdrawFeeResponse{
		Status:      resp.Status,
		UserProfile: &userProfile,
	}, nil
}

func (s *WithdrawServerAPI) GetWithdrawFee(ctx context.Context, req *api.GetWithdrawFeeRequest) (*api.GetWithdrawFeeResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWithdrawFee")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.GetWithdrawFee(ctx, &m2m_api.GetWithdrawFeeRequest{
		MoneyAbbr: moneyAbbr,
		OrgId:     req.OrgId,
	})
	if err != nil {
		return &api.GetWithdrawFeeResponse{}, status.Errorf(codes.Unavailable, err.Error())
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

	return &api.GetWithdrawFeeResponse{
		WithdrawFee: resp.WithdrawFee,
		UserProfile: &userProfile,
	}, nil
}

func (s *WithdrawServerAPI) GetWithdrawHistory(ctx context.Context, req *api.GetWithdrawHistoryRequest) (*api.GetWithdrawHistoryResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/GetWithdrawHistory")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.GetWithdrawHistory(ctx, &m2m_api.GetWithdrawHistoryRequest{
		OrgId:     req.OrgId,
		Offset:    req.Offset,
		Limit:     req.Limit,
		MoneyAbbr: moneyAbbr,
	})
	if err != nil {
		return &api.GetWithdrawHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	withdrawHist := api.GetWithdrawHistoryResponse.GetWithdrawHistory(resp.WithdrawHistory)

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

	return &api.GetWithdrawHistoryResponse{
		Count:           resp.Count,
		WithdrawHistory: withdrawHist,
		UserProfile:     &userProfile,
	}, nil
}

func (s *WithdrawServerAPI) WithdrawReq(ctx context.Context, req *api.WithdrawReqRequest) (*api.WithdrawReqResponse, error) {
	log.WithField("orgId", req.OrgId).Info("grpc_api/WithdrawReq")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	moneyAbbr := m2m_api.Money(req.MoneyAbbr)

	resp, err := m2mClient.WithdrawReq(ctx, &m2m_api.WithdrawReqRequest{
		OrgId:     req.OrgId,
		MoneyAbbr: moneyAbbr,
		Amount:    req.Amount,
	})
	if err != nil {
		return &api.WithdrawReqResponse{}, status.Errorf(codes.Unavailable, err.Error())
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

	return &api.WithdrawReqResponse{
		Status:      resp.Status,
		UserProfile: &userProfile,
	}, nil
}
