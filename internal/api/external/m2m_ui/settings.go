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

type SettingsServerAPI struct {
	validator auth.Validator
}

func NewSettingsServerAPI(validator auth.Validator) *SettingsServerAPI {
	return &SettingsServerAPI{
		validator: validator,
	}
}

func (s *SettingsServerAPI) GetSettings(ctx context.Context, req *api.GetSettingsRequest) (*api.GetSettingsResponse, error) {
	log.WithField("", "").Info("grpc_api/GetSettings")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.GetSettings(ctx, &m2m_api.GetSettingsRequest{})
	if err != nil {
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetSettingsResponse{
		LowBalanceWarning:          resp.LowBalanceWarning,
		DownlinkFee:                resp.DownlinkFee,
		TransactionPercentageShare: resp.TransactionPercentageShare,
	}, nil
}

func (s *SettingsServerAPI) ModifySettings(ctx context.Context, req *api.ModifySettingsRequest) (*api.ModifySettingsResponse, error) {
	log.WithField("", "").Info("grpc_api/ModifySettings")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := m2mClient.ModifySettings(ctx, &m2m_api.ModifySettingsRequest{
		LowBalanceWarning:          req.LowBalanceWarning,
		DownlinkFee:                req.DownlinkFee,
		TransactionPercentageShare: req.TransactionPercentageShare,
	})
	if err != nil {
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.ModifySettingsResponse{
		Status: resp.Status,
	}, nil
}
