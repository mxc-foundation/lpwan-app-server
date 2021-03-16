package external

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	. "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// NetworkServerAPI exports the NetworkServer related functions.
type NetworkServerAPI struct {
	st *store.Handler
}

// NewNetworkServerAPI creates a new NetworkServerAPI.
func NewNetworkServerAPI(h *store.Handler) *NetworkServerAPI {
	return &NetworkServerAPI{
		st: h,
	}
}

// Create creates the given network-server.
func (a *NetworkServerAPI) Create(ctx context.Context, req *pb.CreateNetworkServerRequest) (*pb.CreateNetworkServerResponse, error) {
	if req.NetworkServer == nil {
		return nil, status.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	if valid, err := nsmod.NewValidator().ValidateGlobalNetworkServersAccess(ctx, auth.Create, 0); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	nStruct := &nsmod.NSStruct{
		Server:  req.NetworkServer.Server,
		CACert:  req.NetworkServer.CaCert,
		TLSCert: req.NetworkServer.TlsCert,
		TLSKey:  req.NetworkServer.TlsKey,
	}
	nsClient, err := nStruct.GetNetworkServiceClient()
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	res, err := nsClient.GetVersion(ctx, &empty.Empty{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	region := res.Region.String()
	version := res.Version

	networkServer := NetworkServer{
		Name:                        req.NetworkServer.Name,
		Server:                      req.NetworkServer.Server,
		CACert:                      req.NetworkServer.CaCert,
		TLSCert:                     req.NetworkServer.TlsCert,
		TLSKey:                      req.NetworkServer.TlsKey,
		RoutingProfileCACert:        req.NetworkServer.RoutingProfileCaCert,
		RoutingProfileTLSCert:       req.NetworkServer.RoutingProfileTlsCert,
		RoutingProfileTLSKey:        req.NetworkServer.RoutingProfileTlsKey,
		GatewayDiscoveryEnabled:     req.NetworkServer.GatewayDiscoveryEnabled,
		GatewayDiscoveryInterval:    int(req.NetworkServer.GatewayDiscoveryInterval),
		GatewayDiscoveryTXFrequency: int(req.NetworkServer.GatewayDiscoveryTxFrequency),
		GatewayDiscoveryDR:          int(req.NetworkServer.GatewayDiscoveryDr),
		Region:                      region,
		Version:                     version,
	}

	if err := nsmod.CreateNetworkServer(ctx, &networkServer, a.st); err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &pb.CreateNetworkServerResponse{
		Id: networkServer.ID,
	}, nil
}

// Get returns the network-server matching the given id.
func (a *NetworkServerAPI) Get(ctx context.Context, req *pb.GetNetworkServerRequest) (*pb.GetNetworkServerResponse, error) {
	if valid, err := nsmod.NewValidator().ValidateNetworkServerAccess(ctx, auth.Read, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := a.st.GetNetworkServer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}
	
	if n.Region == "" || n.Version == "" {
		nStruct := &nsmod.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		nsClient, err := nStruct.GetNetworkServiceClient()
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}

		res, err := nsClient.GetVersion(ctx, &empty.Empty{})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		n.Region = res.Region.String()
		n.Version = res.Version

		_ = a.st.UpdateNetworkServer(ctx, &n)
	}

	response := pb.GetNetworkServerResponse{
		NetworkServer: &pb.NetworkServer{
			Id:                          n.ID,
			Name:                        n.Name,
			Server:                      n.Server,
			CaCert:                      n.CACert,
			TlsCert:                     n.TLSCert,
			RoutingProfileCaCert:        n.RoutingProfileCACert,
			RoutingProfileTlsCert:       n.RoutingProfileTLSCert,
			GatewayDiscoveryEnabled:     n.GatewayDiscoveryEnabled,
			GatewayDiscoveryInterval:    uint32(n.GatewayDiscoveryInterval),
			GatewayDiscoveryTxFrequency: uint32(n.GatewayDiscoveryTXFrequency),
			GatewayDiscoveryDr:          uint32(n.GatewayDiscoveryDR),
		},
		Region:  n.Region,
		Version: n.Version,
	}

	response.CreatedAt, err = ptypes.TimestampProto(n.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}
	response.UpdatedAt, err = ptypes.TimestampProto(n.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &response, nil
}

// Update updates the given network-server.
func (a *NetworkServerAPI) Update(ctx context.Context, req *pb.UpdateNetworkServerRequest) (*empty.Empty, error) {
	if req.NetworkServer == nil {
		return nil, status.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	if valid, err := nsmod.NewValidator().ValidateNetworkServerAccess(ctx, auth.Update, req.NetworkServer.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	networkServer, err := a.st.GetNetworkServer(ctx, req.NetworkServer.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	networkServer.Name = req.NetworkServer.Name
	networkServer.Server = req.NetworkServer.Server
	networkServer.CACert = req.NetworkServer.CaCert
	networkServer.TLSCert = req.NetworkServer.TlsCert
	networkServer.RoutingProfileCACert = req.NetworkServer.RoutingProfileCaCert
	networkServer.RoutingProfileTLSCert = req.NetworkServer.RoutingProfileTlsCert
	networkServer.GatewayDiscoveryEnabled = req.NetworkServer.GatewayDiscoveryEnabled
	networkServer.GatewayDiscoveryInterval = int(req.NetworkServer.GatewayDiscoveryInterval)
	networkServer.GatewayDiscoveryTXFrequency = int(req.NetworkServer.GatewayDiscoveryTxFrequency)
	networkServer.GatewayDiscoveryDR = int(req.NetworkServer.GatewayDiscoveryDr)

	if req.NetworkServer.TlsKey != "" {
		networkServer.TLSKey = req.NetworkServer.TlsKey
	}
	if networkServer.TLSCert == "" {
		networkServer.TLSKey = ""
	}

	if req.NetworkServer.RoutingProfileTlsKey != "" {
		networkServer.RoutingProfileTLSKey = req.NetworkServer.RoutingProfileTlsKey
	}
	if networkServer.RoutingProfileTLSCert == "" {
		networkServer.RoutingProfileTLSKey = ""
	}

	if err := nsmod.UpdateNetworkServer(ctx, &networkServer, a.st); err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the network-server matching the given id.
func (a *NetworkServerAPI) Delete(ctx context.Context, req *pb.DeleteNetworkServerRequest) (*empty.Empty, error) {
	if valid, err := nsmod.NewValidator().ValidateNetworkServerAccess(ctx, auth.Delete, req.Id); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := nsmod.DeleteNetworkServer(ctx, req.Id, a.st); err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &empty.Empty{}, nil
}

// List lists the available network-servers.
func (a *NetworkServerAPI) List(ctx context.Context, req *pb.ListNetworkServerRequest) (*pb.ListNetworkServerResponse, error) {
	if valid, err := nsmod.NewValidator().ValidateGlobalNetworkServersAccess(ctx, auth.List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := nsmod.NewValidator().GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	var count int
	var nss []NetworkServer

	if req.OrganizationId == 0 {
		if u.IsGlobalAdmin {
			count, err = a.st.GetNetworkServerCount(ctx, NetworkServerFilters{})
			if err != nil {
				return nil, status.Errorf(codes.Unknown, "%s", err)
			}
			nss, err = a.st.GetNetworkServers(ctx, NetworkServerFilters{
				Limit:  int(req.Limit),
				Offset: int(req.Offset),
			})
			if err != nil {
				return nil, status.Errorf(codes.Unknown, "%s", err)
			}
		}
	} else {
		count, err = a.st.GetNetworkServerCountForOrganizationID(ctx, req.OrganizationId)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}
		nss, err = a.st.GetNetworkServersForOrganizationID(ctx, req.OrganizationId, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}
	}

	resp := pb.ListNetworkServerResponse{
		TotalCount: int64(count),
	}

	for _, ns := range nss {
		row := pb.NetworkServerListItem{
			Id:     ns.ID,
			Name:   ns.Name,
			Server: ns.Server,
		}

		row.CreatedAt, err = ptypes.TimestampProto(ns.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(ns.UpdatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}
