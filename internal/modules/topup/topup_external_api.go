package topup

import (
	"context"
	"strconv"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

// TopUpServerAPI defines the topup server api structure
type TopUpServerAPI struct {
	Validator *validator
}

// NewTopUpServerAPI validates the topup server api
func NewTopUpServerAPI(api TopUpServerAPI) *TopUpServerAPI {
	topupServerAPI = TopUpServerAPI{
		Validator: api.Validator,
	}
	return &topupServerAPI
}

var (
	topupServerAPI TopUpServerAPI
)

// GetTopUpHistory defines the topup history request and response
func (s *TopUpServerAPI) GetTopUpHistory(ctx context.Context, req *api.GetTopUpHistoryRequest) (*api.GetTopUpHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetTopUpHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, "fail to connect to m2m server: %s", err.Error())
	}

	topupClient := m2mServer.NewTopUpServiceClient(m2mClient)

	resp, err := topupClient.GetTopUpHistory(ctx, &m2mServer.GetTopUpHistoryRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpHistoryResponse{}, status.Errorf(codes.Unavailable, "call API in m2m failed: &s", err.Error())
	}

	var topUpHistoryList []*api.TopUpHistory
	for _, item := range resp.TopupHistory {
		topUpHistory := &api.TopUpHistory{
			Amount:    item.Amount,
			CreatedAt: item.CreatedAt,
			TxHash:    item.TxHash,
		}

		topUpHistoryList = append(topUpHistoryList, topUpHistory)
	}

	return &api.GetTopUpHistoryResponse{
		Count:        resp.Count,
		TopupHistory: topUpHistoryList,
	}, status.Error(codes.OK, "")
}

// GetTopUpDestination defines the topup destination request and response
func (s *TopUpServerAPI) GetTopUpDestination(ctx context.Context, req *api.GetTopUpDestinationRequest) (*api.GetTopUpDestinationResponse, error) {
	logInfo := "api/appserver_serves_ui/GetTopUpDestination org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := user.GetUserAPI().Validator.GetIsAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if userIsAdmin == false {
		if err := s.Validator.otpValidator.JwtValidator.Validate(ctx, organization.ValidateOrganizationAccess(organization.Read, req.OrgId)); err != nil {
			log.WithError(err).Error(logInfo)
			return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, "fail to connect to m2m server: %s", err.Error())
	}

	topupClient := m2mServer.NewTopUpServiceClient(m2mClient)

	resp, err := topupClient.GetTopUpDestination(ctx, &m2mServer.GetTopUpDestinationRequest{
		OrgId: req.OrgId,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetTopUpDestinationResponse{}, status.Errorf(codes.Unavailable, "call API in m2m failed: &s", err.Error())
	}

	return &api.GetTopUpDestinationResponse{
		ActiveAccount: resp.ActiveAccount,
	}, status.Error(codes.OK, "")
}
