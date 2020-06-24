package mining

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// MiningServerAPI defines the Mining Server API structure
type MiningServerAPI struct {
	validator auth.Validator
}

// NewMiningServerAPI defines the Mining Server API validator
func NewMiningServerAPI(validator auth.Validator) *MiningServerAPI {
	return &MiningServerAPI{
		validator: validator,
	}
}

// Mining defines the request to give m2m the gateway list that should receive the mining tokens
func (s *MiningServerAPI) Mining(ctx context.Context, req *api.MiningRequest) (*empty.Empty, error) {

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error("create m2mClient for mining error")
		return &empty.Empty{}, status.Errorf(codes.Unavailable, err.Error())
	}

	miningClient := api.NewMiningServiceClient(m2mClient)

	_, err = miningClient.Mining(ctx, &api.MiningRequest{
		GatewayMac: req.GatewayMac,
	})
	if err != nil {
		log.WithError(err).Error("Mining API request error")
		return &empty.Empty{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &empty.Empty{}, nil
}
