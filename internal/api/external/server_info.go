package external

import (
	"context"
	pb "github.com/mxc-foundation/lpwan-app-server/api"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/golang/protobuf/ptypes/empty"
)

type ServerInfoAPI struct {
}

func NewServerInfoAPI() *ServerInfoAPI {
	return &ServerInfoAPI{}
}

func (s *ServerInfoAPI) GetAppserverVersion(ctx context.Context, req *empty.Empty) (*pb.GetAppserverVersionResponse, error) {
	return &pb.GetAppserverVersionResponse{Version:config.AppserverVersion}, nil
}
