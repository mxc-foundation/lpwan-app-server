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

// SettingsServerAPI defines the settings of the Server API structure
type SettingsServerAPI struct {
	validator auth.Validator
}

// NewSettingsServerAPI defines the SettingsServerAPI validator
func NewSettingsServerAPI(validator auth.Validator) *SettingsServerAPI {
	return &SettingsServerAPI{
		validator: validator,
	}
}

// GetSettings defines the settings of the Server API request and response
func (s *SettingsServerAPI) GetSettings(ctx context.Context, req *api.GetSettingsRequest) (*api.GetSettingsResponse, error) {
	log.WithField("", "").Info("grpc_api/GetSettings")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	settingClient := api.NewSettingsServiceClient(m2mClient)

	resp, err := settingClient.GetSettings(ctx, &api.GetSettingsRequest{})
	if err != nil {
		return &api.GetSettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetSettingsResponse{
		LowBalanceWarning:          resp.LowBalanceWarning,
		DownlinkFee:                resp.DownlinkFee,
		TransactionPercentageShare: resp.TransactionPercentageShare,
	}, nil
}

// ModifySettings defines the modification of the Server API settings
func (s *SettingsServerAPI) ModifySettings(ctx context.Context, req *api.ModifySettingsRequest) (*api.ModifySettingsResponse, error) {
	log.WithField("", "").Info("grpc_api/ModifySettings")

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		return &api.ModifySettingsResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	settingClient := api.NewSettingsServiceClient(m2mClient)

	resp, err := settingClient.ModifySettings(ctx, &api.ModifySettingsRequest{
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
