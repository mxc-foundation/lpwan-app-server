package cmdserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mxc-foundation/lpwan-app-server/api/cmdserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/gp"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	gwd "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	spd "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
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
	GetGatewaysCountForNetworkServerID(ctx context.Context, networkServerID int64) (int64, error)
	GetDefaultGatewayProfile(ctx context.Context) (*uuid.UUID, int64, error)
	GetGatewayProfile(ctx context.Context, id uuid.UUID) (gp.GatewayProfile, error)
	GetNetworkServer(ctx context.Context, id int64) (nsd.NetworkServer, error)
	UpdateNetworkServer(ctx context.Context, n *nsd.NetworkServer) error
	BatchSetNetworkServerIDForDeviceProfileAndServiceProfile(ctx context.Context,
		nsIDBefore, nsIDAfter int64) (int64, int64, error)
	UpdateNetworkServerIDAndGatewayProfileIDForGateway(ctx context.Context,
		nsIDAfter int64, gpIDAfter uuid.UUID, mac lorawan.EUI64) (int64, error)
	UpdateGatewayProfile(ctx context.Context, gp *gp.GatewayProfile) error
	GetGatewayProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error)
	GetServiceProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error)
	GetDeviceProfileCountForNetworkServerID(ctx context.Context, networkServerID int64) (int, error)
	GetGatewaysCountForGatewayProfileID(ctx context.Context, gpID uuid.UUID) (int, error)
	DeleteNetworkServer(ctx context.Context, id int64) error
	UpdateNetworkServerName(ctx context.Context, nsID int64, name string) error
}

type gpObject map[uuid.UUID]string
type spObject struct {
	count int64
}
type dpObject struct {
	count int64
}
type gwObject map[lorawan.EUI64]gwd.Gateway
type nsObject struct {
	name  string
	gpMap gpObject
	spMap spObject
	dpMap dpObject
	gwMap map[string]gwObject
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

func (a *Server) inspectNetworkServerSettings(ctx context.Context) (map[int64]*nsObject, []string, error) {
	var nsList []nsd.NetworkServer
	var err error
	nsMap := make(map[int64]*nsObject)
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
		nsObj.gwMap = make(map[string]gwObject)
		nsMap[v.ID] = &nsObj

		result := fmt.Sprintf("network server: id=%d, name=%s \n", v.ID, v.Name)

		// get gateway profiles
		result += fmt.Sprintf("gateway profiles: \n")
		limitgp, err := a.st.GetGatewayProfileCountForNetworkServerID(ctx, v.ID)
		if err != nil {
			return nsMap, response, err
		}
		gpList, err := a.st.GetGatewayProfilesForNetworkServerID(ctx, v.ID, limitgp, 0)
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
		limitsp, err := a.st.GetServiceProfileCountForNetworkServerID(ctx, v.ID)
		if err != nil {
			return nsMap, response, err
		}
		result += fmt.Sprintf("    nsID=%d, nsName=%s, count=%d \n", v.ID, v.Name, limitsp)
		nsObj.spMap.count = int64(limitsp)

		// get device profiles
		result += fmt.Sprintf("device profiles: \n")
		limitdp, err := a.st.GetDeviceProfileCountForNetworkServerID(ctx, v.ID)
		if err != nil {
			return nsMap, response, err
		}
		result += fmt.Sprintf("    nsID=%d, nsName=%s, count=%d \n", v.ID, v.Name, limitdp)
		nsObj.dpMap.count = int64(limitdp)

		// get gateways
		result += fmt.Sprintf("gateways: \n")
		limitgw, err := a.st.GetGatewaysCountForNetworkServerID(ctx, v.ID)
		if err != nil {
			return nsMap, response, err
		}
		gwList, err := a.st.GetGatewaysForNetworkServerID(ctx, v.ID, int(limitgw), 0)
		if err != nil {
			return nsMap, response, err
		}
		for _, v := range gwList {
			if nsObj.gwMap[*v.GatewayProfileID] == nil {
				nsObj.gwMap[*v.GatewayProfileID] = make(map[lorawan.EUI64]gwd.Gateway)
			}
			nsObj.gwMap[*v.GatewayProfileID][v.MAC] = v
		}
		for gwpID, count := range nsObj.gwMap {
			result += fmt.Sprintf("    nsID=%d, gateway_profile_id=%s, count=%d \n", v.ID, gwpID, count)
		}

		response = append(response, result)
	}

	return nsMap, response, nil
}

// CorrectNetworkServerSettings
// - removes network server and all related settings from DB except for the given network server id
// - rename the only network server to default_network_server
// - ensure default_gateway_profile is set
func (a *Server) CorrectNetworkServerSettings(ctx context.Context, req *pb.CorrectNetworkServerSettingsRequest) (*pb.CorrectNetworkServerSettingsResponse, error) {
	nsMap, _, err := a.inspectNetworkServerSettings(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if len(nsMap) == 0 {
		return &pb.CorrectNetworkServerSettingsResponse{
			Report: "no network servers set yet",
		}, nil
	}

	// clean up redundant network server and all related settings
	report := ""

	if nsMap[req.NetworkServerId] == nil {
		return nil, status.Errorf(codes.InvalidArgument, "network server %d does not exist", req.NetworkServerId)
	}
	defaultGpID, err := a.ensureDefaultGatewayProfile(ctx, req.NetworkServerId, &report, nsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	defaultNsID := req.NetworkServerId

	for nsID, nsObj := range nsMap {
		if nsID == defaultNsID {
			// default network server or gateway profiles referring to default network server won't be processed
			continue
		}
		// reset network server id for device profile and service profile
		// reset device profiles' network server id
		dpChanged, spChanged, err := correctDeviceProfilesAndServiceProfiles(ctx, a.st, defaultNsID, nsID,
			nsObj.dpMap.count, nsObj.spMap.count)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
		}
		report += fmt.Sprintf("INFO: changed device profiles' network server id from %d to %d \n", nsID, defaultNsID)
		if dpChanged != nsObj.dpMap.count {
			report += fmt.Sprintf("WARNING: reset device profiles' network server id from %d to %d, expect done %d, actually done %d \n",
				nsID, defaultNsID, nsObj.dpMap.count, dpChanged)
		}
		report += fmt.Sprintf("INFO: changed service profiles' network server id from %d to %d \n", nsID, defaultNsID)
		if spChanged != nsObj.spMap.count {
			report += fmt.Sprintf("WARNING: reset service profiles' network server id from %d to %d, expect done %d, actually done %d \n",
				nsID, defaultNsID, nsObj.spMap.count, spChanged)
		}

		// update nsMap
		nsMap[defaultNsID].dpMap.count += dpChanged
		nsObj.dpMap.count -= dpChanged
		nsMap[defaultNsID].spMap.count += spChanged
		nsObj.spMap.count -= spChanged

		// reset gateways' network server id and gateway profile
		for gpIDStr, gwList := range nsObj.gwMap {
			count := len(gwList)
			if count == 0 {
				continue
			}
			gpID := uuid.UUID{}
			if err = gpID.UnmarshalText([]byte(gpIDStr)); err != nil {
				return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			if err = batchUpdateGateway(ctx, gwList, defaultNsID, defaultGpID, a.st, a.nsCli, &report,
				nsMap[defaultNsID].gwMap[defaultGpID.String()]); err != nil {
				return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			if len(gwList) != 0 {
				// not all gateways have been processed, log it
				report += fmt.Sprintf("WARNING: changed gateways number does not match for gpID=%s: expected %d, got %d \n",
					gpIDStr, count, count-len(gwList))
			}
		}

		// delete redundant gateway profiles
		for gpID, gpName := range nsObj.gpMap {
			// check whether gateway profile is still referred by gateways
			count, err := a.st.GetGatewaysCountForGatewayProfileID(ctx, gpID)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			if count != 0 {
				report += fmt.Sprintf("WARNING: cannot delete gateway profile %s(%s), it's still referred by %d gateways \n",
					gpID.String(), gpName, nsObj.gwMap[gpID.String()])
				continue
			}

			// remove gateway profile
			if err = gp.DeleteGatewayProfile(ctx, a.gpSt, gpID, a.nsCli); err != nil {
				return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
			}
			delete(nsMap[nsID].gpMap, gpID)
			report += fmt.Sprintf("INFO: delete gateway profile %s \n", gpID.String())
		}
	}

	err = a.ensureDefaultNetworkServerName(ctx, defaultNsID, &report, nsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s \n %v", report, err)
	}

	return &pb.CorrectNetworkServerSettingsResponse{Report: report}, nil
}

func batchUpdateGateway(ctx context.Context, gwMap map[lorawan.EUI64]gwd.Gateway, defaultNsID int64, defaultGpID uuid.UUID,
	st Store, nsCli *nscli.Client, report *string, gwMapDS map[lorawan.EUI64]gwd.Gateway) error {
	reportStr := *report
	for mac, gateway := range gwMap {
		nsID := gateway.NetworkServerID
		gpID := *gateway.GatewayProfileID
		updateReq := ns.UpdateGatewayRequest{
			Gateway: &ns.Gateway{
				Id:               mac[:],
				GatewayProfileId: defaultGpID.Bytes(),
				Location: &common.Location{
					Latitude:  gateway.Latitude,
					Longitude: gateway.Longitude,
					Altitude:  gateway.Altitude,
				},
			},
		}
		client, err := nsCli.GetNetworkServerServiceClient(defaultNsID)
		if err != nil {
			return err
		}
		_, err = client.UpdateGateway(ctx, &updateReq)
		if err != nil {
			return err
		}

		// update network id server and gateway profile id local gateway
		ra, err := st.UpdateNetworkServerIDAndGatewayProfileIDForGateway(ctx, defaultNsID, defaultGpID, gateway.MAC)
		if err != nil {
			return err
		}
		if ra != 0 {
			delete(gwMap, mac)
			gateway.NetworkServerID = defaultNsID
			*gateway.GatewayProfileID = defaultGpID.String()
			gwMapDS[mac] = gateway
			reportStr += fmt.Sprintf("INFO: updated network server id from %d to %d, gateway profile id from %s to "+
				"%s for gateway %s \n", nsID, defaultNsID, gpID, defaultGpID.String(), mac.String())
		}
	}
	*report = reportStr
	return nil
}

func correctDeviceProfilesAndServiceProfiles(ctx context.Context, st Store, nsIDKept, nsID int64,
	dpCount int64, spCount int64) (int64, int64, error) {
	if dpCount == 0 && spCount == 0 {
		return 0, 0, nil
	}
	dpChanged, spChanged, err := st.BatchSetNetworkServerIDForDeviceProfileAndServiceProfile(ctx, nsID, nsIDKept)
	if err != nil {
		return 0, 0, err
	}
	return dpChanged, spChanged, nil
}

func (a *Server) ensureDefaultGatewayProfile(ctx context.Context, defaultNsID int64, report *string, nsMap map[int64]*nsObject) (uuid.UUID, error) {
	var gatewayProfileID uuid.UUID
	reportStr := *report
	defer func() {
		*report = reportStr
	}()

	gpID, _, err := a.st.GetDefaultGatewayProfile(ctx)
	if err != nil && err != errHandler.ErrDoesNotExist {
		return gatewayProfileID, err
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
			return gatewayProfileID, err
		}
		// add default gateway profile to nsMap
		nsMap[defaultNsID].gpMap[*newGpID] = nsd.DefaultGatewayProfileName
		reportStr += fmt.Sprintf("INFO: created default_gateway_profile with network server id %d and gatway "+
			"profile id %s \n", defaultNsID, newGpID.String())
		return *newGpID, nil
	} else if err == nil {
		gatewayProfileID = *gpID
	}

	gatewayProfile, err := gp.GetGatewayProfile(ctx, gatewayProfileID, a.gpSt, a.nsCli)
	if err != nil {
		return gatewayProfileID, err
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
			return gatewayProfileID, err
		}
		// update nsMap
		delete(nsMap[gatewayProfile.NetworkServerID].gpMap, gatewayProfileID)
		nsMap[defaultNsID].gpMap[gatewayProfileID] = nsd.DefaultGatewayProfileName
		reportStr += fmt.Sprintf("INFO: updated default_gateway_profile %s, change network server id from %d "+
			"to %d \n", gatewayProfileID.String(), gatewayProfile.NetworkServerID, defaultNsID)
	}
	return gatewayProfileID, nil
}

func (a *Server) ensureDefaultNetworkServerName(ctx context.Context, defaultNsID int64, report *string, nsMap map[int64]*nsObject) error {
	reportStr := *report
	var defaultNsObj *nsObject
	defer func() {
		*report = reportStr
	}()

	for nsID, nsObj := range nsMap {
		if nsID == defaultNsID {
			defaultNsObj = nsObj
			continue
		}
		// is network server used by device profiles?
		count, err := a.st.GetDeviceProfileCountForNetworkServerID(ctx, nsID)
		if err != nil {
			return err
		}
		if count != 0 {
			reportStr += fmt.Sprintf("WARNING: cannot delete network server %d, "+
				"it's still referred by device profiles \n", nsID)
			continue
		}
		// is network server used by service profiles?
		count, err = a.st.GetDeviceProfileCountForNetworkServerID(ctx, nsID)
		if err != nil {
			return err
		}
		if count != 0 {
			reportStr += fmt.Sprintf("WARNING: cannot delete network server %d, "+
				"it's still referred by service profiles \n", nsID)
			continue
		}
		// is network server used by gateway profiles?
		count, err = a.st.GetGatewayProfileCountForNetworkServerID(ctx, nsID)
		if err != nil {
			return err
		}
		if count != 0 {
			reportStr += fmt.Sprintf("WARNING: cannot delete network server %d, "+
				"it's still referred by gateway profiles \n", nsID)
			continue
		}
		// is network server used by gateways?
		count64, err := a.st.GetGatewaysCountForNetworkServerID(ctx, nsID)
		if err != nil {
			return err
		}
		if count64 != 0 {
			reportStr += fmt.Sprintf("WARNING: cannot delete network server %d, "+
				"it's still referred by gateways \n", nsID)
			continue
		}
		// delete network server **LOCALLY ONLY**
		if err = a.st.DeleteNetworkServer(ctx, nsID); err != nil {
			return err
		}
		delete(nsMap, nsID)
		reportStr += fmt.Sprintf("INFO: delete network server %d \n", nsID)
	}

	if len(nsMap) != 1 {
		reportStr += fmt.Sprintf("WARNING: network server is not correted, %d left \n", len(nsMap))
		return nil
	}
	if defaultNsObj == nil {
		reportStr += fmt.Sprintf("WARNING: default network server is not assigned \n")
		return nil
	}

	if defaultNsObj.name != nsd.DefaultNetworkServerName {
		if err := a.st.UpdateNetworkServerName(ctx, defaultNsID, nsd.DefaultNetworkServerName); err != nil {
			return fmt.Errorf("update network server name to %s error: %v", nsd.DefaultNetworkServerName, err)
		}
		reportStr += fmt.Sprintf("INFO: update name for network server %d, from %s to %s \n",
			defaultNsID, defaultNsObj.name, nsd.DefaultNetworkServerName)
	}

	return nil
}
