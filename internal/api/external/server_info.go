package external

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// ServerInfoAPI defines the Server Info API structure
type ServerInfoAPI struct {
}

// NewServerInfoAPI defines the Server Info API
func NewServerInfoAPI() *ServerInfoAPI {
	return &ServerInfoAPI{}
}

// GetAppserverVersion defines the Appserver Version response
func (s *ServerInfoAPI) GetAppserverVersion(ctx context.Context, req *empty.Empty) (*pb.GetAppserverVersionResponse, error) {
	return &pb.GetAppserverVersionResponse{Version: config.AppserverVersion}, nil
}
