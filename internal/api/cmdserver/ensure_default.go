package cmdserver

import (
	"context"
	"fmt"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/gofrs/uuid"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/cmdserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/gp"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	spd "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/tls"
)

// CMDServer represents gRPC server serving command server
type CMDServer struct {
	st                       Store
	nsSt                     nsd.Store
	gpSt                     gp.Store
	gs                       *grpc.Server
	nsCli                    *nscli.Client
	applicationServerID      uuid.UUID
	applicationServerPubHost string
}

// Start starts gRPC server that serves mxp server
func Start(st Store, nsSt nsd.Store, gpSt gp.Store, nsCli *nscli.Client,
	applicationServerID uuid.UUID, applicationServerPubHost string) (*CMDServer, error) {
	srv := &CMDServer{
		st:                       st,
		nsSt:                     nsSt,
		gpSt:                     gpSt,
		nsCli:                    nsCli,
		applicationServerID:      applicationServerID,
		applicationServerPubHost: applicationServerPubHost,
	}
	if err := srv.listenWithCredentials("0.0.0.0:1000", "", "", ""); err != nil {
		return nil, err
	}
	return srv, nil
}

func (srv *CMDServer) listenWithCredentials(bind, caCert, tlsCert, tlsKey string) error {
	gs, err := tls.NewServerWithTLSCredentials("command server", caCert, tlsCert, tlsKey)
	if err != nil {
		return fmt.Errorf("listenWithCredentials: get new server error: %v", err)
	}
	srv.gs = gs

	pb.RegisterEnsureDefaultServiceServer(gs, NewServer(srv.st, srv.nsSt, srv.gpSt, srv.nsCli,
		srv.applicationServerID, srv.applicationServerPubHost))

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("listenWithCredentials: start api listener error: %v", err)
	}

	go func() {
		_ = gs.Serve(ln)
	}()

	logrus.Info("start command line server")
	return nil
}

// Stop gracefully stops gRPC server
func (srv *CMDServer) Stop() {
	srv.gs.GracefulStop()
}

// Server defines cmdserver
type Server struct {
	st                       Store
	nsSt                     nsd.Store
	gpSt                     gp.Store
	nsCli                    *nscli.Client
	applicationServerID      uuid.UUID
	applicationServerPubHost string
}

// NewServer returns a new instance of cmdserver
func NewServer(st Store, nsSt nsd.Store, gpSt gp.Store, nsCli *nscli.Client,
	applicationServerID uuid.UUID, applicationServerPubHost string) *Server {
	return &Server{
		st:                       st,
		nsSt:                     nsSt,
		gpSt:                     gpSt,
		nsCli:                    nsCli,
		applicationServerID:      applicationServerID,
		applicationServerPubHost: applicationServerPubHost,
	}
}

// Store defines db APIs used by this package
type Store interface {
	GetNetworkServers(ctx context.Context, filters nsd.NetworkServerFilters) ([]nsd.NetworkServer, error)
	GetGatewayProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]gp.GatewayProfileMeta, error)
	GetServiceProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]spd.ServiceProfileMeta, error)
	GetDeviceProfilesForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]dp.DeviceProfileMeta, error)
	GetGatewaysForNetworkServerID(ctx context.Context, networkServerID int64, limit, offset int) ([]gwd.Gateway, error)
	GetDefaultGatewayProfile(ctx context.Context) (*uuid.UUID, int64, error)
	GetGatewayProfile(ctx context.Context, id uuid.UUID) (gp.GatewayProfile, error)
	GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error)
	UpdateNetworkServer(ctx context.Context, n *nsd.NetworkServer) error
	BatchSetNetworkServerIDForDeviceProfile(ctx context.Context, nsIDBefore, nsIDAfter int64) (int64, error)
	BatchSetNetworkServerIDForServiceProfile(ctx context.Context, nsIDBefore, nsIDAfter int64) (int64, error)
	BatchSetNetworkServerIDAndGatewayProfileIDForGateways(ctx context.Context, nsIDBefore,
		nsIDAfter int64, gpIDBefore, gpIDAfter uuid.UUID) (int64, error)
	UpdateGatewayProfile(ctx context.Context, gp *gp.GatewayProfile) error

	Tx(ctx context.Context, f func(context.Context, Store) error) error
}

type gpObject map[uuid.UUID]string
type spObject struct {
	count int64
}
type dpObject struct {
	count int64
}
type gwObject map[string]int
type nsObject struct {
	name  string
	gpMap gpObject
	spMap spObject
	dpMap dpObject
	gwMap gwObject
}

// InspectNetworkServerSettings inspects all existsing network servers in db together with all other settings which are
// referring to network server id
func (a *Server) InspectNetworkServerSettings(ctx context.Context,
	req *pb.InspectNetworkServerSettingsRequest) (*pb.InspectNetworkServerSettingsResponse, error) {
	_, result, err := a.inspectNetworkServerSettings(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.InspectNetworkServerSettingsResponse{InspectResult: result}, nil
}

func (a *Server) inspectNetworkServerSettings(ctx context.Context) (map[int64]nsObject, []string, error) {
	var nsList []nsd.NetworkServer
	var err error
	nsMap := make(map[int64]nsObject)
	response := []string{""}

	nsList, err = a.st.GetNetworkServers(ctx, nsd.NetworkServerFilters{
		Limit:  999,
		Offset: 0,
	})
	if err != nil {
		return nsMap, response, err
	}

	if len(nsList) == 0 {
		return nil, []string{"no netowrk servers set yet"}, nil
	}

	for _, v := range nsList {
		nsObj := nsObject{
			name: v.Name,
		}
		nsObj.gpMap = make(map[uuid.UUID]string)
		nsObj.gwMap = make(map[string]int)
		nsMap[v.ID] = nsObj

		result := fmt.Sprintf("network server: id=%d, name=%s \n", v.ID, v.Name)

		// get gateway profiles
		result += fmt.Sprintf("gateway profiles: \n")
		gpList, err := a.st.GetGatewayProfilesForNetworkServerID(ctx, v.ID, 999, 0)
		if err != nil {
			return nsMap, response, err
		}

		for _, v := range gpList {
			result += fmt.Sprintf("    nsID=%d, nsName=%s, gateway_profile_id=%s, gateway_profile_name=%s \n",
				v.NetworkServerID, v.NetworkServerName, v.GatewayProfileID.String(), v.Name)
			nsObj.gpMap[v.GatewayProfileID] = v.Name
		}

		// get service profiles
		result += fmt.Sprintf("service profiles: \n")
		spList, err := a.st.GetServiceProfilesForNetworkServerID(ctx, v.ID, 999, 0)
		if err != nil {
			return nsMap, response, err
		}
		result += fmt.Sprintf("    nsID=%d, nsName=%s, count=%d \n", v.ID, v.Name, len(spList))
		nsObj.spMap.count = int64(len(spList))

		// get device profiles
		result += fmt.Sprintf("device profiles: \n")
		dpList, err := a.st.GetDeviceProfilesForNetworkServerID(ctx, v.ID, 999, 0)
		if err != nil {
			return nsMap, response, err
		}
		result += fmt.Sprintf("    nsID=%d, nsName=%s, count=%d \n", v.ID, v.Name, len(dpList))
		nsObj.dpMap.count = int64(len(dpList))

		// get gateways
		result += fmt.Sprintf("gateways: \n")
		gwList, err := a.st.GetGatewaysForNetworkServerID(ctx, v.ID, 999, 0)
		if err != nil {
			return nsMap, response, err
		}
		for _, v := range gwList {
			nsObj.gwMap[*v.GatewayProfileID]++
		}
		for gwpID, count := range nsObj.gwMap {
			result += fmt.Sprintf("    nsID=%d, gateway_profile_id=%s, count=%d", v.ID, gwpID, count)
		}

		response = append(response, result)
	}

	return nsMap, response, nil
}

// CorrectNetworkServerSettings
// - removes network server and all related settings from DB except for the given network server id
// - rename the only network server to default_network_server
// - ensure default_gateway_profile is set
func (a *Server) CorrectNetworkServerSettings(ctx context.Context,
	req *pb.CorrectNetworkServerSettingsRequest) (*pb.CorrectNetworkServerSettingsResponse, error) {
	nsMap, _, err := a.inspectNetworkServerSettings(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if len(nsMap) == 0 {
		return &pb.CorrectNetworkServerSettingsResponse{
			Report: "no network servers set yet",
		}, nil
	}

	/*	if len(nsMap) == 1 {
		for nsID := range nsMap {
			// ensure default_network_server
			if err = a.ensureDefaultNetworkServerName(ctx, nsID); err != nil {
				return nil, status.Errorf(codes.Internal, "update network server name error: %v", err)
			}
			// ensure defaul_gateway_profile
			if err = a.ensureDefaultGatewayProfile(ctx, nsID); err != nil {
				return nil, status.Errorf(codes.Internal, "ensure default gateway profile error: %v", err)
			}
		}
	}*/

	// clean up redundant network server and all related settings
	report := ""

	defaultGpID, err := a.ensureDefaultGatewayProfile(ctx, req.NetworkServerId, &report)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	defaultNsID := req.NetworkServerId

	for nsID, nsObj := range nsMap {
		if nsID == defaultNsID {
			continue
		}
		// reset network server id for device profile and service profile
		// use transaction to guarantee device profiles and service profiles which are using same network server will
		//  use same network server no matter what
		if err = a.st.Tx(ctx, func(ctx context.Context, store Store) error {
			// reset device profiles' network server id
			if err = correctDeviceProfiles(ctx, store, defaultNsID, nsID, nsObj.dpMap.count, &report); err != nil {
				return status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			// reset service profiles' network server id
			if err = correctServiceProfiles(ctx, store, defaultNsID, nsID, nsObj.spMap.count, &report); err != nil {
				return status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			return nil
		}); err != nil {
			return nil, err
		}

		// reset gateways' network server id and gateway profile
		for gpIDStr, count := range nsObj.gwMap {
			gpID := uuid.UUID{}
			if err = gpID.UnmarshalText([]byte(gpIDStr)); err != nil {
				return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
			}

			a.st.BatchSetNetworkServerIDAndGatewayProfileIDForGateways(ctx, nsID, defaultNsID, gpID, *defaultGpID)
		}

		if err = a.st.Tx(ctx, func(ctx context.Context, store Store) error {
			store.BatchSetNetworkServerIDAndGatewayProfileIDForGateways(ctx, nsID, defaultNsID, nsMap[defaultNsID].gpMap, *defaultGpID)
		}); err != nil {
			return nil, err
		}

	}

	return &pb.CorrectNetworkServerSettingsResponse{}, nil
}

func correctServiceProfiles(ctx context.Context, st Store, nsIDKept, nsID int64, count int64, report *string) error {
	var reportStr string
	ra, err := st.BatchSetNetworkServerIDForServiceProfile(ctx, nsID, nsIDKept)
	if err != nil {
		return err
	}
	if ra != count {
		reportStr = *report + fmt.Sprintf("WARNING: reset service profiles' network server id from %d to %d, expect done %d, actually done %d \n",
			nsID, nsIDKept, count, ra)
		*report = reportStr
	}
	return nil
}

func correctDeviceProfiles(ctx context.Context, st Store, nsIDKept, nsID int64, count int64, report *string) error {
	var reportStr string
	ra, err := st.BatchSetNetworkServerIDForDeviceProfile(ctx, nsID, nsIDKept)
	if err != nil {
		return err
	}
	if ra != count {
		reportStr = *report + fmt.Sprintf("WARNING: reset device profiles' network server id from %d to %d, expect done %d, actually done %d \n",
			nsID, nsIDKept, count, ra)
		*report = reportStr
	}
	return nil
}

func (a *Server) ensureDefaultGatewayProfile(ctx context.Context, defaultNsID int64, report *string) (gpID *uuid.UUID, err error) {
	gpID, _, err = a.st.GetDefaultGatewayProfile(ctx)
	if err != nil && err != errHandler.ErrDoesNotExist {
		return gpID, err
	} else if err == errHandler.ErrDoesNotExist {
		// no default gateway profile exists, create one with default network server id
		newGpID, err := gp.CreateGatewayProfile(ctx, a.gpSt, a.nsCli, &gp.GatewayProfile{
			NetworkServerID: defaultNsID,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
			Name:            nsd.DefaultGatewayProfileName,
			GatewayProfile: ns.GatewayProfile{
				Channels: []uint32{0, 1, 2},
			},
		})
		if err != nil {
			return gpID, err
		}
		gpID = newGpID
		return gpID, nil
	}

	gatewayProfile, err := a.st.GetGatewayProfile(ctx, *gpID)
	if err != nil {
		return gpID, err
	}
	if defaultNsID != gatewayProfile.NetworkServerID {
		// existing default gatway profile is not using default network server id, update it
		if err = a.st.UpdateGatewayProfile(ctx, &gp.GatewayProfile{
			NetworkServerID: defaultNsID,
			CreatedAt:       gatewayProfile.CreatedAt,
			UpdatedAt:       time.Now().UTC(),
			Name:            nsd.DefaultGatewayProfileName,
			GatewayProfile:  gatewayProfile.GatewayProfile,
		}); err != nil {
			return gpID, err
		}
	}
	return gpID, nil
}

func (a *Server) ensureDefaultNetworkServerName(ctx context.Context, nsID int64) error {
	networkServer, err := a.st.GetNetworkServer(ctx, nsID)
	if err != nil {
		return err
	}

	if networkServer.Name != nsd.DefaultNetworkServerName {
		if err := a.st.UpdateNetworkServer(ctx, &nsd.NetworkServer{
			ID:                          networkServer.ID,
			CreatedAt:                   networkServer.CreatedAt,
			UpdatedAt:                   time.Now().UTC(),
			Name:                        nsd.DefaultNetworkServerName,
			Server:                      networkServer.Server,
			CACert:                      networkServer.CACert,
			TLSCert:                     networkServer.TLSCert,
			TLSKey:                      networkServer.TLSKey,
			RoutingProfileCACert:        networkServer.RoutingProfileCACert,
			RoutingProfileTLSCert:       networkServer.RoutingProfileTLSCert,
			RoutingProfileTLSKey:        networkServer.RoutingProfileTLSKey,
			GatewayDiscoveryEnabled:     networkServer.GatewayDiscoveryEnabled,
			GatewayDiscoveryInterval:    networkServer.GatewayDiscoveryInterval,
			GatewayDiscoveryTXFrequency: networkServer.GatewayDiscoveryTXFrequency,
			GatewayDiscoveryDR:          networkServer.GatewayDiscoveryDR,
			Version:                     networkServer.Version,
			Region:                      networkServer.Region,
		}); err != nil {
			return err
		}
	}
	return nil
}
