package ns

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	pb "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/gp"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
)

// NetworkServer defines the information to connect to a network-server.
type NetworkServer struct {
	ID                          int64     `db:"id"`
	CreatedAt                   time.Time `db:"created_at"`
	UpdatedAt                   time.Time `db:"updated_at"`
	Name                        string    `db:"name"`
	Server                      string    `db:"server"`
	CACert                      string    `db:"ca_cert"`
	TLSCert                     string    `db:"tls_cert"`
	TLSKey                      string    `db:"tls_key"`
	RoutingProfileCACert        string    `db:"routing_profile_ca_cert"`
	RoutingProfileTLSCert       string    `db:"routing_profile_tls_cert"`
	RoutingProfileTLSKey        string    `db:"routing_profile_tls_key"`
	GatewayDiscoveryEnabled     bool      `db:"gateway_discovery_enabled"`
	GatewayDiscoveryInterval    int       `db:"gateway_discovery_interval"`
	GatewayDiscoveryTXFrequency int       `db:"gateway_discovery_tx_frequency"`
	GatewayDiscoveryDR          int       `db:"gateway_discovery_dr"`
	Version                     string    `db:"version"`
	Region                      string    `db:"region"`
}

// Validate validates the network-server data.
func (ns NetworkServer) Validate() error {
	if strings.TrimSpace(ns.Name) == "" || len(ns.Name) > 100 {
		return errHandler.ErrNetworkServerInvalidName
	}

	if ns.GatewayDiscoveryEnabled && ns.GatewayDiscoveryInterval <= 0 {
		return errHandler.ErrInvalidGatewayDiscoveryInterval
	}
	return nil
}

// NetworkServerFilters provides filters for filtering network-servers.
type NetworkServerFilters struct {
	OrganizationID int64 `db:"organization_id"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f NetworkServerFilters) SQL() string {
	var filters []string

	if f.OrganizationID != 0 {
		filters = append(filters, "sp.organization_id = :organization_id")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// NetworkServerAPI exports the NetworkServer related functions.
type NetworkServerAPI struct {
	st                          Store
	gpSt                        gp.Store
	auth                        auth.Authenticator
	nsCli                       *nscli.Client
	applicationServerID         uuid.UUID
	applicationServerPublicHost string
}

// Store defines db APIs used by this package
type Store interface {
	CreateNetworkServer(ctx context.Context, n *NetworkServer) error
	UpdateNetworkServer(ctx context.Context, n *NetworkServer) error
	GetNetworkServersForOrganizationID(ctx context.Context, organizationID int64, limit, offset int) ([]NetworkServer, error)
	GetNetworkServerCountForOrganizationID(ctx context.Context, organizationID int64) (int, error)
	GetNetworkServers(ctx context.Context, filters NetworkServerFilters) ([]NetworkServer, error)
	GetNetworkServerCount(ctx context.Context, filters NetworkServerFilters) (int, error)
	GetNetworkServer(ctx context.Context, id int64) (NetworkServer, error)
	DeleteNetworkServer(ctx context.Context, id int64) error
}

// NewNetworkServerAPI creates a new NetworkServerAPI.
func NewNetworkServerAPI(st Store, gpSt gp.Store, nsCli *nscli.Client, auth auth.Authenticator,
	applicationServerID uuid.UUID, applicationServerPublicHost string) *NetworkServerAPI {
	return &NetworkServerAPI{
		st:                          st,
		gpSt:                        gpSt,
		auth:                        auth,
		nsCli:                       nsCli,
		applicationServerID:         applicationServerID,
		applicationServerPublicHost: applicationServerPublicHost,
	}
}

// DefaultNetworkServerName defines name of the default network server
const DefaultNetworkServerName = "default_network_server"

// DefaultNetworkServerAddress defines address of the default network server,
// this value is unique for the same supernode
const DefaultNetworkServerAddress = "network-server:8000"

// DefaultGatewayProfileName defines name of the default gateway profile for default network server,
// this value is unique
const DefaultGatewayProfileName = "default_gateway_profile"

// DefaultServiceProfileName defines name of the default service profile for default network server with given org id
// this value is unique: default_service_profile_ORGID
const DefaultServiceProfileName = "default_service_profile_"

// CreateNetworkServer creates network server config in appserver and network server
func CreateNetworkServer(ctx context.Context, n *NetworkServer, st Store, gpSt gp.Store,
	nsCli *nscli.Client, applicationServerID uuid.UUID, applicationServerPublicHost string) error {
	// adding new network server connection if not exists yet
	conn, err := grpccli.Connect(grpccli.ConnectionOpts{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	})
	if err != nil {
		return err
	}
	nsClient := ns.NewNetworkServerServiceClient(conn)
	// adding application server id of local appserver to network server, this is unique in network server
	// check whether application server id already exists
	_, err = nsClient.GetRoutingProfile(ctx, &ns.GetRoutingProfileRequest{
		Id: applicationServerID.Bytes(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "code = NotFound") {
			// create routing profile if none exists
			if _, err = nsClient.CreateRoutingProfile(ctx, &ns.CreateRoutingProfileRequest{
				RoutingProfile: &ns.RoutingProfile{
					Id:      applicationServerID.Bytes(),
					AsId:    applicationServerPublicHost,
					CaCert:  n.RoutingProfileCACert,
					TlsCert: n.RoutingProfileTLSCert,
					TlsKey:  n.RoutingProfileTLSKey,
				},
			}); err != nil {
				return err
			}
		} else {
			// unknow error
			return err
		}
	} // do not create routing profile if there is one existing already

	res, err := nsClient.GetVersion(ctx, &empty.Empty{})
	if err != nil {
		return err
	}
	n.Region = res.Region.String()
	n.Version = res.Version

	// create network server after appserver routing profile is created in network server
	// when adding routing profile fails, no need to proceed adding network server locally
	// when adding routing profile succeeds, creating network server locally fails, user can try again later
	if err := st.CreateNetworkServer(ctx, n); err != nil {
		return err
	}

	// save network server connection
	if err = nsCli.Save(n.ID, conn); err != nil {
		logrus.WithError(err).Warnf("save network server connection error")
	}

	// create default gatway profile for this network server
	gatewayProfile := gp.GatewayProfile{
		NetworkServerID: n.ID,
		Name:            DefaultGatewayProfileName,
		GatewayProfile: ns.GatewayProfile{
			Channels: []uint32{0, 1, 2},
		},
	}
	_, err = gp.CreateGatewayProfile(ctx, gpSt, nsCli, &gatewayProfile)
	if err != nil {
		return fmt.Errorf("create default gateway profile for network server returns error: %v", err)
	}

	return nil
}

// UpdateNetworkServer updates network server config in appserver and network server
func UpdateNetworkServer(ctx context.Context, n *NetworkServer,
	st Store, nsCli *nscli.Client, applicationServerID uuid.UUID, applicationServerPublicHost string) error {
	nsClient, err := nsCli.GetNetworkServerServiceClient(n.ID)
	if err != nil {
		return fmt.Errorf("failed to get ns client for network server %d: %v", n.ID, err)
	}
	_, err = nsClient.UpdateRoutingProfile(ctx, &ns.UpdateRoutingProfileRequest{
		RoutingProfile: &ns.RoutingProfile{
			Id:      applicationServerID.Bytes(),
			AsId:    applicationServerPublicHost,
			CaCert:  n.RoutingProfileCACert,
			TlsCert: n.RoutingProfileTLSCert,
			TlsKey:  n.RoutingProfileTLSKey,
		},
	})
	if err != nil {
		return errors.Wrap(err, "update routing-profile error")
	}
	// update remote first, so that user can try again if local change fails to be done
	if err := st.UpdateNetworkServer(ctx, n); err != nil {
		return err
	}
	return nil
}

// DeleteNetworkServer deletes network server config from appserver and network server
func DeleteNetworkServer(ctx context.Context, id int64, st Store, nsCli *nscli.Client, applicationServerID uuid.UUID) error {

	nsClient, err := nsCli.GetNetworkServerServiceClient(id)
	if err != nil {
		return fmt.Errorf("failed to get ns client for network server %d: %v", id, err)
	}

	_, err = nsClient.GetRoutingProfile(ctx, &ns.GetRoutingProfileRequest{
		Id: applicationServerID.Bytes(),
	})
	if err == nil {
		_, err = nsClient.DeleteRoutingProfile(ctx, &ns.DeleteRoutingProfileRequest{
			Id: applicationServerID.Bytes(),
		})
		if err != nil {
			return errors.Wrap(err, "delete routing-profile error")
		}
	}

	// delete remote routing profile first, so that user can try again if local change fails to be done
	if err := st.DeleteNetworkServer(ctx, id); err != nil {
		return err
	}

	return nil
}

// Create creates the given network-server.
func (a *NetworkServerAPI) Create(ctx context.Context, req *pb.CreateNetworkServerRequest) (*pb.CreateNetworkServerResponse, error) {
	if req.NetworkServer == nil {
		return nil, status.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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
	if err := CreateNetworkServer(ctx, &networkServer, a.st, a.gpSt, a.nsCli, a.applicationServerID,
		a.applicationServerPublicHost); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateNetworkServerResponse{
		Id: networkServer.ID,
	}, nil
}

// Get returns the network-server matching the given id.
func (a *NetworkServerAPI) Get(ctx context.Context, req *pb.GetNetworkServerRequest) (*pb.GetNetworkServerResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	n, err := a.st.GetNetworkServer(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
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
	response.CreatedAt = timestamppb.New(n.CreatedAt)
	response.UpdatedAt = timestamppb.New(n.UpdatedAt)

	return &response, nil
}

// Update updates the given network-server.
func (a *NetworkServerAPI) Update(ctx context.Context, req *pb.UpdateNetworkServerRequest) (*empty.Empty, error) {
	if req.NetworkServer == nil {
		return nil, status.Errorf(codes.InvalidArgument, "network_server must not be nil")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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

	if err := UpdateNetworkServer(ctx, &networkServer, a.st, a.nsCli, a.applicationServerID,
		a.applicationServerPublicHost); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

// Delete deletes the network-server matching the given id.
func (a *NetworkServerAPI) Delete(ctx context.Context, req *pb.DeleteNetworkServerRequest) (*empty.Empty, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := DeleteNetworkServer(ctx, req.Id, a.st, a.nsCli, a.applicationServerID); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

// List lists the available network-servers.
func (a *NetworkServerAPI) List(ctx context.Context, req *pb.ListNetworkServerRequest) (*pb.ListNetworkServerResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	var count int
	var nss []NetworkServer

	if req.OrganizationId == 0 {
		if cred.IsGlobalAdmin {
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

		row.CreatedAt = timestamppb.New(ns.CreatedAt)
		row.UpdatedAt = timestamppb.New(ns.UpdatedAt)

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}
