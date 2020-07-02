package networkserver

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

type NetworkServerStore interface {
	CreateNetworkServer(ctx context.Context, n *NetworkServer) error
	GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error)
	UpdateNetworkServer(ctx context.Context, n *NetworkServer) error
	DeleteNetworkServer(ctx context.Context, id int64) error
	GetNetworkServerCount(ctx context.Context) (int, error)
	GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error)
	GetNetworkServers(ctx context.Context, limit, offset int) ([]NetworkServer, error)
	GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error)
	GetNetworkServerForDevEUI(ctx context.Context, devEUI lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForDeviceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForServiceProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForGatewayMAC(ctx context.Context, mac lorawan.EUI64) (NetworkServer, error)
	GetNetworkServerForGatewayProfileID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetNetworkServerForMulticastGroupID(ctx context.Context, id uuid.UUID) (NetworkServer, error)
	GetDefaultNetworkServer(ctx context.Context) (NetworkServer, error)
}

// NetworkServerAPI exports the NetworkServer related functions.
type NetworkServerAPI struct {
	Validator *validator
	Store     NetworkServerStore
}

// NewNetworkServerAPI creates a new NetworkServerAPI.
func NewNetworkServerAPI(api NetworkServerAPI) *NetworkServerAPI {
	networkServerAPI = NetworkServerAPI{
		Validator: api.Validator,
		Store:     api.Store,
	}

	_ = networkServerAPI.SetupDefault()

	return &networkServerAPI
}

var (
	networkServerAPI NetworkServerAPI
)

func GetNetworkServerAPI() *NetworkServerAPI {
	return &networkServerAPI
}

func (a *NetworkServerAPI) SetupDefault() error {
	ctx := context.Background()
	count, err := storage.GetGatewayProfileCount(ctx, storage.DB())
	if err != nil && err != storage.ErrDoesNotExist {
		return errors.Wrap(err, "Failed to load gateway profiles")
	}

	if count != 0 {
		// check if default gateway profile already exists
		gpList, err := storage.GetGatewayProfiles(ctx, storage.DB(), count, 0)
		if err != nil {
			return errors.Wrap(err, "Failed to load gateway profiles")
		}

		for _, v := range gpList {
			if v.Name == "default_gateway_profile" {
				return nil
			}
		}
	}

	// none default_gateway_profile exists, add one
	var networkServer NetworkServer
	n, err := a.Store.GetNetworkServers(ctx, 1, 0)
	if err != nil && err != storage.ErrDoesNotExist {
		return errors.Wrap(err, "Load network server internal error")
	}

	if len(n) >= 1 {
		networkServer = n[0]
	} else {
		// insert default one
		err := storage.Transaction(func(tx sqlx.Ext) error {
			return storage.CreateNetworkServer(ctx, storage.DB(), &storage.NetworkServer{
				Name:                    "default_network_server",
				Server:                  "network-server:8000",
				GatewayDiscoveryEnabled: false,
			})
		})
		if err != nil {
			return nil
		}

		// get network-server id

		networkServer, err = a.Store.GetDefaultNetworkServer(ctx)
		if err != nil {
			return err
		}
	}

	gp := storage.GatewayProfile{
		NetworkServerID: networkServer.ID,
		Name:            "default_gateway_profile",
		GatewayProfile: ns.GatewayProfile{
			Channels:      []uint32{0, 1, 2},
			ExtraChannels: []*ns.GatewayProfileExtraChannel{},
		},
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		return storage.CreateGatewayProfile(ctx, tx, &gp)
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create default gateway profile")
	}

	return nil
}

// Create creates the given network-server.
func (a *NetworkServerAPI) Create(ctx context.Context, req *pb.CreateNetworkServerRequest) (*pb.CreateNetworkServerResponse, error) {
	if req.NetworkServer == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateNetworkServersAccess(Create, 0)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

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
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		return a.Store.CreateNetworkServer(ctx, &networkServer)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateNetworkServerResponse{
		Id: networkServer.ID,
	}, nil
}

// Get returns the network-server matching the given id.
func (a *NetworkServerAPI) Get(ctx context.Context, req *pb.GetNetworkServerRequest) (*pb.GetNetworkServerResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateNetworkServerAccess(Read, req.Id)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := a.Store.GetNetworkServer(ctx, req.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var region string
	var version string

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err == nil {
		resp, err := nsClient.GetVersion(ctx, &empty.Empty{})
		if err == nil {
			region = resp.Region.String()
			version = resp.Version
		}
	}

	resp := pb.GetNetworkServerResponse{
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
		Region:  region,
		Version: version,
	}

	resp.CreatedAt, err = ptypes.TimestampProto(n.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(n.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &resp, nil
}

// Update updates the given network-server.
func (a *NetworkServerAPI) Update(ctx context.Context, req *pb.UpdateNetworkServerRequest) (*empty.Empty, error) {
	if req.NetworkServer == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateNetworkServerAccess(Update, req.NetworkServer.Id)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	networkServer, err := a.Store.GetNetworkServer(ctx, req.NetworkServer.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
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

	err = storage.Transaction(func(tx sqlx.Ext) error {
		return a.Store.UpdateNetworkServer(ctx, &networkServer)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the network-server matching the given id.
func (a *NetworkServerAPI) Delete(ctx context.Context, req *pb.DeleteNetworkServerRequest) (*empty.Empty, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateNetworkServerAccess(Delete, req.Id)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := storage.Transaction(func(tx sqlx.Ext) error {
		return a.Store.DeleteNetworkServer(ctx, req.Id)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// List lists the available network-servers.
func (a *NetworkServerAPI) List(ctx context.Context, req *pb.ListNetworkServerRequest) (*pb.ListNetworkServerResponse, error) {
	if err := a.Validator.otpValidator.JwtValidator.Validate(ctx, validateNetworkServersAccess(List, req.OrganizationId)); err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	u, err := user.GetUserAPI().Validator.GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	var count int
	var nss []NetworkServer

	if req.OrganizationId == 0 {
		if u.IsAdmin {
			count, err = a.Store.GetNetworkServerCount(ctx)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
			nss, err = a.Store.GetNetworkServers(ctx, int(req.Limit), int(req.Offset))
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}
	} else {
		count, err = a.Store.GetNetworkServerCountForOrganizationID(ctx, req.OrganizationId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		nss, err = a.Store.GetNetworkServersForOrganizationID(ctx, req.OrganizationId, int(req.Limit), int(req.Offset))
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
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
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(ns.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}
