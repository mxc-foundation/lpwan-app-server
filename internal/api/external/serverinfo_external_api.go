package external

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
)

// ServerInfoAPI defines the Server Info API structure
type ServerInfoAPI struct {
	serverRegion string
}

// NewServerInfoAPI defines the Server Info API
func NewServerInfoAPI(serverRegion string) *ServerInfoAPI {
	return &ServerInfoAPI{
		serverRegion: serverRegion,
	}
}

// GetAppserverVersion defines the Appserver Version response
func (s *ServerInfoAPI) GetAppserverVersion(ctx context.Context, req *empty.Empty) (*pb.GetAppserverVersionResponse, error) {
	return &pb.GetAppserverVersionResponse{Version: config.AppserverVersion}, nil
}

func (s *ServerInfoAPI) GetServerRegion(ctx context.Context, req *empty.Empty) (*pb.GetServerRegionResponse, error) {
	region := s.serverRegion

	if region == pb.ServerRegion_name[int32(pb.ServerRegion_AVERAGE)] {
		return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_AVERAGE)]}, nil
	}

	if region == pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)] {
		return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_RESTRICTED)]}, nil
	}

	// pb.ServerRegion_NOT_DEFINED same as default region: pb.ServerRegion_AVERAGE
	return &pb.GetServerRegionResponse{ServerRegion: pb.ServerRegion_name[int32(pb.ServerRegion_NOT_DEFINED)]}, nil
}

// GetVersion defines the version of the machine to machine server API
func (s *ServerInfoAPI) GetMxprotocolServerVersion(ctx context.Context, req *empty.Empty) (*pb.GetMxprotocolServerVersionResponse, error) {
	log.WithField("", "").Info("grpc_api/GetVersion")

	verClient := mxpcli.Global.GetServerServiceClient()

	resp, err := verClient.GetVersion(ctx, req)
	if err != nil {
		return &pb.GetMxprotocolServerVersionResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetMxprotocolServerVersionResponse{
		Version: resp.Version,
	}, nil
}
