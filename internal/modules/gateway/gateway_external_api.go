package gateway

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	m2mServer "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	nsClient "github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
)

// GatewayAPI exports the Gateway related functions.
type GatewayAPI struct {
	st                  GatewayStore
	txSt                store.Store
	ApplicationServerID uuid.UUID
}

// NewGatewayAPI creates a new GatewayAPI.
func NewGatewayAPI(applicationID uuid.UUID) *GatewayAPI {
	st := store.New(storage.DB().DB)
	return &GatewayAPI{
		st:                  st,
		txSt:                st,
		ApplicationServerID: applicationID,
	}
}

func (a *GatewayAPI) GetGatewayMiningList(ctx context.Context, time, limit int64) ([]lorawan.EUI64, error) {
	return a.st.GetGatewayMiningList(ctx, time, limit)
}

func (a *GatewayAPI) UpdateFirstHeartbeatToZero(ctx context.Context, mac lorawan.EUI64) error {
	return a.st.UpdateFirstHeartbeatToZero(ctx, mac)
}

func (a *GatewayAPI) UpdateFirmwareFromProvisioningServer(conf config.Config) error {
	log.WithFields(log.Fields{
		"provisioning-server": conf.ProvisionServer.ProvisionServer,
		"caCert":              conf.ProvisionServer.CACert,
		"tlsCert":             conf.ProvisionServer.TLSCert,
		"tlsKey":              conf.ProvisionServer.TLSKey,
		"schedule":            conf.ProvisionServer.UpdateSchedule,
	}).Info("Start schedule to update gateway firmware...")
	Service.SupernodeAddr = os.Getenv("APPSERVER")
	if strings.HasPrefix(Service.SupernodeAddr, "https://") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "https://", "", -1)
	}
	if strings.HasPrefix(Service.SupernodeAddr, "http://") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "http://", "", -1)
	}
	if strings.HasSuffix(Service.SupernodeAddr, ":8080") {
		Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, ":8080", "", -1)
	}
	Service.SupernodeAddr = strings.Replace(Service.SupernodeAddr, "/", "", -1)

	var bindPortOldGateway string
	var bindPortNewGateway string

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.OldGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for OldGateway: %s", conf.ApplicationServer.APIForGateway.OldGateway.Bind))
	} else {
		bindPortOldGateway = strArray[1]
	}

	if strArray := strings.Split(conf.ApplicationServer.APIForGateway.NewGateway.Bind, ":"); len(strArray) != 2 {
		return errors.New(fmt.Sprintf("Invalid API Bind settings for NewGateway: %s", conf.ApplicationServer.APIForGateway.NewGateway.Bind))
	} else {
		bindPortNewGateway = strArray[1]
	}

	c := cron.New()
	err := c.AddFunc(conf.ProvisionServer.UpdateSchedule, func() {
		log.Info("Check firmware update...")
		gwFwList, err := a.st.GetGatewayFirmwareList()
		if err != nil {
			log.WithError(err).Errorf("Failed to get gateway firmware list.")
			return
		}

		// send update
		psClient, err := provisionserver.CreateClientWithCert(conf.ProvisionServer.ProvisionServer,
			conf.ProvisionServer.CACert,
			conf.ProvisionServer.TLSCert,
			conf.ProvisionServer.TLSKey)
		if err != nil {
			log.WithError(err).Errorf("Create Provisioning server client error")
			return
		}

		for _, v := range gwFwList {
			res, err := psClient.GetUpdate(context.Background(), &psPb.GetUpdateRequest{
				Model:          v.Model,
				SuperNodeAddr:  Service.SupernodeAddr,
				PortOldGateway: bindPortOldGateway,
				PortNewGateway: bindPortNewGateway,
			})
			if err != nil {
				log.WithError(err).Errorf("Failed to get update for gateway model: %s", v.Model)
				continue
			}

			var md5sum types.MD5SUM
			if err := md5sum.UnmarshalText([]byte(res.FirmwareHash)); err != nil {
				log.WithError(err).Errorf("Failed to unmarshal firmware hash: %s", res.FirmwareHash)
				continue
			}

			gatewayFw := GatewayFirmware{
				Model:        v.Model,
				ResourceLink: res.ResourceLink,
				FirmwareHash: md5sum,
			}

			model, _ := a.st.UpdateGatewayFirmware(&gatewayFw)
			if model == "" {
				log.Warnf("No row updated for gateway_firmware at model=%s", v.Model)
			}

		}
	})
	if err != nil {
		log.Fatalf("Failed to set update schedule when set up provisioning server config: %s", err.Error())
	}

	go c.Start()

	return nil
}

// BatchResetDefaultGatewatConfig reset gateways config to default config matching organization list
func (a *GatewayAPI) BatchResetDefaultGatewatConfig(ctx context.Context, req *pb.BatchResetDefaultGatewatConfigRequest) (*pb.BatchResetDefaultGatewatConfigResponse, error) {
	log.WithFields(log.Fields{
		"organization_list": req.OrganizationList,
	}).Info("BatchResetDefaultGatewatConfig is called")

	// check user permission, only global admin allowed
	isAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	// if process for any organizaiton failed, return this error message to user for retry
	var failedList []string
	var succeededList []string
	var organizationList []int

	if req.OrganizationList == "all" {
		// reset for all organizations
		count, err := organization.Service.St.GetOrganizationCount(ctx, organization.OrganizationFilters{})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		limit := 100
		for offset := 0; offset <= count/limit; offset++ {
			list, err := organization.Service.St.GetOrganizationIDList(limit, offset, "")
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			organizationList = append(organizationList, list...)
		}

	} else {
		// parse organization list
		strOrgList := strings.Split(req.OrganizationList, ",")
		for _, v := range strOrgList {
			orgID, err := strconv.Atoi(v)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid organization list format: %s (correct example: '2, 3, 4, 5' or '1' or 'all')", req.OrganizationList)
			}

			organizationList = append(organizationList, orgID)
		}
	}

	// proceed when organizationList is complete
	for _, v := range organizationList {
		if v == 0 {
			log.Warn("0 is in organization list")
			continue
		}

		err := a.resetDefaultGatewayConfigByOrganizationID(ctx, int64(v))
		if err != nil {
			log.WithError(err).Errorf("failed to reset default gateway config for organization %d", v)
			failedList = append(failedList, strconv.Itoa(v))
			continue
		}

		succeededList = append(succeededList, strconv.Itoa(v))
	}

	return &pb.BatchResetDefaultGatewatConfigResponse{
		Status: fmt.Sprintf("following organization failed: %s \n following organization succeeded: %s",
			strings.Join(failedList, ","), strings.Join(succeededList, ",")),
	}, status.Error(codes.OK, "")
}

func (a *GatewayAPI) resetDefaultGatewayConfigByOrganizationID(ctx context.Context, orgID int64) error {
	count, err := a.st.GetGatewayCountForOrganizationID(ctx, orgID, "")
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(fmt.Sprintf("There is no gateway in organization : %d", orgID))
	}

	limit := 100
	for offset := 0; offset <= count/limit; offset++ {
		gwList, err := a.st.GetGatewaysForOrganizationID(ctx, orgID, limit, offset, "")
		if err != nil {
			return err
		}

		for _, v := range gwList {
			err := a.getDefaultGatewayConfig(ctx, &v)
			if err != nil {
				return err
			}

			err = a.st.UpdateGatewayConfigByGwId(ctx, v.Config, v.MAC)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResetDefaultGatewatConfigByID reste gateway config to default config matching gateway id
func (a *GatewayAPI) ResetDefaultGatewatConfigByID(ctx context.Context, req *pb.ResetDefaultGatewatConfigByIDRequest) (*pb.ResetDefaultGatewatConfigByIDResponse, error) {
	log.WithFields(log.Fields{
		"gateway_id": req.Id,
	}).Info("ResetDefaultGatewatConfigByID is called")

	// check user permission, only global admin allowed
	isAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	gw, err := a.st.GetGateway(ctx, mac, true)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	err = a.getDefaultGatewayConfig(ctx, &gw)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	err = a.st.UpdateGatewayConfigByGwId(ctx, gw.Config, gw.MAC)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.ResetDefaultGatewatConfigByIDResponse{}, status.Error(codes.OK, "")
}

// InsertNewDefaultGatewayConfig insert given new default gateway config
func (a *GatewayAPI) InsertNewDefaultGatewayConfig(ctx context.Context, req *pb.InsertNewDefaultGatewayConfigRequest) (*pb.InsertNewDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("InsertNewDefaultGatewayConfig is called")

	// check user permission, only global admin allowed
	isAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:         req.Model,
		Region:        req.Region,
		DefaultConfig: strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", Service.SupernodeAddr, -1),
	}

	err = a.st.GetDefaultGatewayConfig(&defaultGatewayConfig)
	if err == nil {
		// config already exist, no need to insert
		return nil, status.Errorf(codes.AlreadyExists, "model=%s, region=%s", req.Model, req.Region)
	} else if err != storage.ErrDoesNotExist {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = a.st.AddNewDefaultGatewayConfig(&defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.InsertNewDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// UpdateNewDefaultGatewayConfig update default gateway config matching model and region
func (a *GatewayAPI) UpdateDefaultGatewayConfig(ctx context.Context, req *pb.UpdateDefaultGatewayConfigRequest) (*pb.UpdateDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("UpdateDefaultGatewayConfig is called")

	// check user permission
	isAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  req.Model,
		Region: req.Region,
	}

	err = a.st.GetDefaultGatewayConfig(&defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	defaultGatewayConfig.DefaultConfig = strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", Service.SupernodeAddr, -1)
	err = a.st.UpdateDefaultGatewayConfig(&defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.UpdateDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// GetDefaultGatewayConfig get content of default gateway config matching model and region
func (a *GatewayAPI) GetDefaultGatewayConfig(ctx context.Context, req *pb.GetDefaultGatewayConfigRequest) (*pb.GetDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("GetDefaultGatewayConfig is called")

	// check user permission
	isAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  req.Model,
		Region: req.Region,
	}

	err = a.st.GetDefaultGatewayConfig(&defaultGatewayConfig)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &pb.GetDefaultGatewayConfigResponse{DefaultConfig: defaultGatewayConfig.DefaultConfig}, status.Error(codes.OK, "")
}

// Create creates the given gateway.
func (a *GatewayAPI) Create(ctx context.Context, req *pb.CreateGatewayRequest) (*empty.Empty, error) {
	if req.Gateway == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway must not be nil")
	}

	if req.Gateway.Location == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway.location must not be nil")
	}

	if valid, err := NewValidator().ValidateGlobalGatewaysAccess(ctx, authcus.Create, req.Gateway.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// also validate that the network-server is accessible for the given organization
	if valid, err := NewValidator().ValidateOrganizationNetworkServerAccess(ctx, authcus.Read,
		req.Gateway.OrganizationId, req.Gateway.NetworkServerId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := a.storeGateway(ctx, req.Gateway, &Gateway{}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (a *GatewayAPI) getDefaultGatewayConfig(ctx context.Context, gw *Gateway) error {
	if !strings.HasPrefix(gw.Model, "MX19") {
		return nil
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get network server %d", gw.NetworkServerID)
		return errors.Wrapf(err, "GetDefaultGatewayConfig")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  gw.Model,
		Region: n.Region,
	}

	err = a.st.GetDefaultGatewayConfig(&defaultGatewayConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to get default gateway config for model= %s, region= %s", defaultGatewayConfig.Model, defaultGatewayConfig.Region)
	}

	gw.Config = strings.Replace(defaultGatewayConfig.DefaultConfig, "{{ .GatewayID }}", gw.MAC.String(), -1)
	return nil
}

func (a *GatewayAPI) storeGateway(ctx context.Context, req *pb.Gateway, defaultGw *Gateway) (err error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	createReq := ns.CreateGatewayRequest{
		Gateway: &ns.Gateway{
			Id:               mac[:],
			Location:         req.Location,
			RoutingProfileId: a.ApplicationServerID.Bytes(),
		},
	}

	if req.GatewayProfileId != "" {
		gpID, err := uuid.FromString(req.GatewayProfileId)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		createReq.Gateway.GatewayProfileId = gpID.Bytes()
	}

	for _, board := range req.Boards {
		var gwBoard ns.GatewayBoard

		if board.FpgaId != "" {
			var fpgaID lorawan.EUI64
			if err := fpgaID.UnmarshalText([]byte(board.FpgaId)); err != nil {
				return status.Errorf(codes.InvalidArgument, "fpga_id: %s", err)
			}
			gwBoard.FpgaId = fpgaID[:]
		}

		if board.FineTimestampKey != "" {
			var key lorawan.AES128Key
			if err := key.UnmarshalText([]byte(board.FineTimestampKey)); err != nil {
				return status.Errorf(codes.InvalidArgument, "fine_timestamp_key: %s", err)
			}
			gwBoard.FineTimestampKey = key[:]
		}

		createReq.Gateway.Boards = append(createReq.Gateway.Boards, &gwBoard)
	}

	defaultGw.MAC = mac
	defaultGw.NetworkServerID = req.NetworkServerId
	err = a.getDefaultGatewayConfig(ctx, defaultGw)
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	// TODO: this part needs UI modification
	/*
		tags := hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}

			for k, v := range req.Tags {
			tags.Map[k] = sql.NullString{Valid: true, String: v}
		}*/

	// A transaction is needed as:
	//  * A remote gRPC call is performed and in case of error, we want to
	//    rollback the transaction.
	//  * We want to lock the organization so that we can validate the
	//    max gateway count.

	tx, err := a.txSt.TxBegin(ctx)
	if err != nil {
		return err
	}
	defer tx.TxRollback(ctx)

	org, err := organization.Service.St.GetOrganization(ctx, req.OrganizationId, true)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	// Validate max. gateway count when != 0.
	if org.MaxGatewayCount != 0 {
		count, err := tx.GetGatewayCount(ctx, "")
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		if count >= org.MaxGatewayCount {
			return helpers.ErrToRPCError(storage.ErrOrganizationMaxGatewayCount)
		}
	}

	err = tx.CreateGateway(ctx, &Gateway{
		MAC:             mac,
		Name:            req.Name,
		Description:     req.Description,
		OrganizationID:  req.OrganizationId,
		Ping:            req.DiscoveryEnabled,
		NetworkServerID: req.NetworkServerId,
		Latitude:        req.Location.Latitude,
		Longitude:       req.Location.Longitude,
		Altitude:        req.Location.Altitude,
		Model:           defaultGw.Model,
		FirstHeartbeat:  0,
		LastHeartbeat:   0,
		Config:          defaultGw.Config,
		OsVersion:       defaultGw.OsVersion,
		Statistics:      defaultGw.Statistics,
		SerialNumber:    defaultGw.SerialNumber,
		FirmwareHash:    types.MD5SUM{},
	})
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServer(ctx, req.NetworkServerId)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	_, err = client.CreateGateway(ctx, &createReq)
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return err
	}

	if err := tx.TxCommit(ctx); err != nil {
		return err
	}

	return nil
}

// Get returns the gateway matching the given Mac.
func (a *GatewayAPI) Get(ctx context.Context, req *pb.GetGatewayRequest) (*pb.GetGatewayResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gw, err := a.st.GetGateway(ctx, mac, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := networkserver.GetNetworkServerAPI().Store.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	getResp, err := client.GetGateway(ctx, &ns.GetGatewayRequest{
		Id: mac[:],
	})
	if err != nil {
		return nil, err
	}

	resp := pb.GetGatewayResponse{
		Gateway: &pb.Gateway{
			Id:               mac.String(),
			Name:             gw.Name,
			Description:      gw.Description,
			OrganizationId:   gw.OrganizationID,
			DiscoveryEnabled: gw.Ping,
			Location: &common.Location{
				Latitude:  gw.Latitude,
				Longitude: gw.Longitude,
				Altitude:  gw.Altitude,
			},
			NetworkServerId: gw.NetworkServerID,
			// TODO: UI
			/*			Tags:            make(map[string]string),
						Metadata:        make(map[string]string),*/
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(gw.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(gw.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	if gw.FirstSeenAt != nil {
		resp.FirstSeenAt, err = ptypes.TimestampProto(*gw.FirstSeenAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	if gw.LastSeenAt != nil {
		resp.LastSeenAt, err = ptypes.TimestampProto(*gw.LastSeenAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	if len(getResp.Gateway.GatewayProfileId) != 0 {
		gpID, err := uuid.FromBytes(getResp.Gateway.GatewayProfileId)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		resp.Gateway.GatewayProfileId = gpID.String()
	}

	for i := range getResp.Gateway.Boards {
		var gwBoard pb.GatewayBoard

		if len(getResp.Gateway.Boards[i].FpgaId) != 0 {
			var fpgaID lorawan.EUI64
			copy(fpgaID[:], getResp.Gateway.Boards[i].FpgaId)
			gwBoard.FpgaId = fpgaID.String()
		}

		if len(getResp.Gateway.Boards[i].FineTimestampKey) != 0 {
			var key lorawan.AES128Key
			copy(key[:], getResp.Gateway.Boards[i].FineTimestampKey)
			gwBoard.FineTimestampKey = key.String()
		}

		resp.Gateway.Boards = append(resp.Gateway.Boards, &gwBoard)
	}

	// TODO: UI
	/*	for k, v := range gw.Tags.Map {
			resp.Gateway.Tags[k] = v.String
		}
		for k, v := range gw.Metadata.Map {
			resp.Gateway.Metadata[k] = v.String
		}*/

	return &resp, err
}

// List lists the gateways.
func (a *GatewayAPI) List(ctx context.Context, req *pb.ListGatewayRequest) (*pb.ListGatewayResponse, error) {
	if valid, err := NewValidator().ValidateGlobalGatewaysAccess(ctx, authcus.List, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	filters := GatewayFilters{
		Search:         req.Search,
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
		OrganizationID: req.OrganizationId,
	}

	u, err := NewValidator().GetUser(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	// Filter on username when OrganizationID is not set and the user is
	// not a global admin.
	if !u.IsGlobalAdmin && filters.OrganizationID == 0 {
		filters.UserID = u.ID
	}

	count, err := a.st.GetGatewayCount(ctx, "")
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	gws, err := a.st.GetGateways(ctx, req.Limit, req.Offset, "")
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListGatewayResponse{
		TotalCount: int64(count),
	}

	for _, gw := range gws {
		row := pb.GatewayListItem{
			Id:              gw.MAC.String(),
			Name:            gw.Name,
			Description:     gw.Description,
			OrganizationId:  gw.OrganizationID,
			NetworkServerId: gw.NetworkServerID,
			Location: &common.Location{
				Latitude:  gw.Latitude,
				Longitude: gw.Longitude,
				Altitude:  gw.Altitude,
			},
		}

		row.CreatedAt, err = ptypes.TimestampProto(gw.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(gw.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		if gw.FirstSeenAt != nil {
			row.FirstSeenAt, err = ptypes.TimestampProto(*gw.FirstSeenAt)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}
		if gw.LastSeenAt != nil {
			row.LastSeenAt, err = ptypes.TimestampProto(*gw.LastSeenAt)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}

// ListLocations lists the gateway locations.
func (a *GatewayAPI) ListLocations(ctx context.Context, req *pb.ListGatewayLocationsRequest) (*pb.ListGatewayLocationsResponse, error) {
	var result []*pb.GatewayLocationListItem
	/*
		redisConn := storage.RedisPool().Get()
		defer redisConn.Close()

		resultJSON, err := redis.Bytes(redisConn.Do("GET", GatewayLocationsRedisKey))
		if err == nil {
			json.Unmarshal(resultJSON, &result)
		}

		if len(result) == 0 {
			gwsLoc, err := a.st.GetGatewaysLoc(ctx, storage.DB(), viper.GetInt("application_server.gateways_locations_limit"))
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}

			for _, loc := range gwsLoc {
				result = append(result, &pb.GatewayLocationListItem{
					Location: &pb.GatewayLocation{
						Latitude:  loc.Latitude,
						Longitude: loc.Longitude,
						Altitude:  loc.Altitude,
					},
				})
			}

			bytes, err := json.Marshal(&result)
			if err == nil {
				redisConn.Do("SET", GatewayLocationsRedisKey, bytes)
			}
		}
	*/
	resp := pb.ListGatewayLocationsResponse{
		Result: result,
	}

	return &resp, nil
}

// Update updates the given gateway.
func (a *GatewayAPI) Update(ctx context.Context, req *pb.UpdateGatewayRequest) (*empty.Empty, error) {
	if req.Gateway == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway must not be nil")
	}

	if req.Gateway.Location == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway.location must not be nil")
	}

	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Gateway.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Update, mac); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// TODO: UI
	/*
		tags := hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}
			for k, v := range req.Gateway.Tags {
				tags.Map[k] = sql.NullString{Valid: true, String: v}
			}
	*/
	tx, err := a.txSt.TxBegin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't start transaction: %v", err)
	}
	defer tx.TxRollback(ctx)

	gw, err := tx.GetGateway(ctx, mac, true)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	gw.Name = req.Gateway.Name
	gw.Description = req.Gateway.Description
	gw.Ping = req.Gateway.DiscoveryEnabled
	gw.Latitude = req.Gateway.Location.Latitude
	gw.Longitude = req.Gateway.Location.Longitude
	gw.Altitude = req.Gateway.Location.Altitude
	//gw.Tags = tags

	err = tx.UpdateGateway(ctx, &gw)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	updateReq := ns.UpdateGatewayRequest{
		Gateway: &ns.Gateway{
			Id:       mac[:],
			Location: req.Gateway.Location,
		},
	}

	if req.Gateway.GatewayProfileId != "" {
		gpID, err := uuid.FromString(req.Gateway.GatewayProfileId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		updateReq.Gateway.GatewayProfileId = gpID.Bytes()
	}

	for _, board := range req.Gateway.Boards {
		var gwBoard ns.GatewayBoard

		if board.FpgaId != "" {
			var fpgaID lorawan.EUI64
			if err := fpgaID.UnmarshalText([]byte(board.FpgaId)); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "fpga_id: %s", err)
			}
			gwBoard.FpgaId = fpgaID[:]
		}

		if board.FineTimestampKey != "" {
			var key lorawan.AES128Key
			if err := key.UnmarshalText([]byte(board.FineTimestampKey)); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "fine_timestamp_key: %s", err)
			}
			gwBoard.FineTimestampKey = key[:]
		}

		updateReq.Gateway.Boards = append(updateReq.Gateway.Boards, &gwBoard)
	}

	n, err := networkserver.Service.St.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	_, err = client.UpdateGateway(ctx, &updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%v", err)
	}

	if err := tx.TxCommit(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't commit transaction: %v", err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the gateway matching the given ID.
func (a *GatewayAPI) Delete(ctx context.Context, req *pb.DeleteGatewayRequest) (*empty.Empty, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Delete, mac); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err := a.st.DeleteGateway(ctx, mac)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// GetStats gets the gateway statistics for the gateway with the given Mac.
func (a *GatewayAPI) GetStats(ctx context.Context, req *pb.GetGatewayStatsRequest) (*pb.GetGatewayStatsResponse, error) {
	var gatewayID lorawan.EUI64
	if err := gatewayID.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, gatewayID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	start, err := ptypes.Timestamp(req.StartTimestamp)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	end, err := ptypes.Timestamp(req.EndTimestamp)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, ok := ns.AggregationInterval_value[strings.ToUpper(req.Interval)]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "bad interval: %s", req.Interval)
	}

	metrics, err := storage.GetMetrics(ctx, storage.AggregationInterval(strings.ToUpper(req.Interval)), "gw:"+gatewayID.String(), start, end)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	result := make([]*pb.GatewayStats, len(metrics))
	for i, m := range metrics {
		result[i] = &pb.GatewayStats{
			RxPacketsReceived:   int32(m.Metrics["rx_count"]),
			RxPacketsReceivedOk: int32(m.Metrics["rx_ok_count"]),
			TxPacketsReceived:   int32(m.Metrics["tx_count"]),
			TxPacketsEmitted:    int32(m.Metrics["tx_ok_count"]),
		}

		result[i].Timestamp, err = ptypes.TimestampProto(m.Time)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	return &pb.GetGatewayStatsResponse{
		Result: result,
	}, nil
}

// GetLastPing returns the last emitted ping and gateways receiving this ping.
func (a *GatewayAPI) GetLastPing(ctx context.Context, req *pb.GetLastPingRequest) (*pb.GetLastPingResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	ping, pingRX, err := a.st.GetLastGatewayPingAndRX(ctx, mac)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.GetLastPingResponse{
		Frequency: uint32(ping.Frequency),
		Dr:        uint32(ping.DR),
	}

	resp.CreatedAt, err = ptypes.TimestampProto(ping.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	for _, rx := range pingRX {
		resp.PingRx = append(resp.PingRx, &pb.PingRX{
			GatewayId: rx.GatewayMAC.String(),
			Rssi:      int32(rx.RSSI),
			LoraSnr:   rx.LoRaSNR,
			Latitude:  rx.Location.Latitude,
			Longitude: rx.Location.Longitude,
			Altitude:  rx.Altitude,
		})
	}

	return &resp, nil
}

// StreamFrameLogs streams the uplink and downlink frame-logs for the given mac.
// Note: these are the raw LoRaWAN frames and this endpoint is intended for debugging.
func (a *GatewayAPI) StreamFrameLogs(req *pb.StreamGatewayFrameLogsRequest, srv pb.GatewayService_StreamFrameLogsServer) error {
	var mac lorawan.EUI64

	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return status.Errorf(codes.InvalidArgument, "mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(srv.Context(), authcus.Read, mac); !valid || err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := networkserver.Service.St.GetNetworkServerForGatewayMAC(srv.Context(), mac)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	client, err := nsClient.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	streamClient, err := client.StreamFrameLogsForGateway(srv.Context(), &ns.StreamFrameLogsForGatewayRequest{
		GatewayId: mac[:],
	})
	if err != nil {
		return err
	}

	for {
		resp, err := streamClient.Recv()
		if err != nil {
			return err
		}

		up, down, err := device.ConvertUplinkAndDownlinkFrames(resp.GetUplinkFrameSet(), resp.GetDownlinkFrame(), false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		var frameResp pb.StreamGatewayFrameLogsResponse
		if up != nil {
			frameResp.Frame = &pb.StreamGatewayFrameLogsResponse_UplinkFrame{
				UplinkFrame: up,
			}
		}

		if down != nil {
			frameResp.Frame = &pb.StreamGatewayFrameLogsResponse_DownlinkFrame{
				DownlinkFrame: down,
			}
		}

		err = srv.Send(&frameResp)
		if err != nil {
			return err
		}
	}
}

// GetGwConfig gets the gateway config file
func (a *GatewayAPI) GetGwConfig(ctx context.Context, req *pb.GetGwConfigRequest) (*pb.GetGwConfigResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	gwConfig, err := a.st.GetGatewayConfigByGwId(ctx, mac)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "GetGwConfig/unable to get gateway config from DB %s", err)
	}

	return &pb.GetGwConfigResponse{Conf: gwConfig}, nil
}

// UpdateGwConfig gateway configuration file
func (a *GatewayAPI) UpdateGwConfig(ctx context.Context, req *pb.UpdateGwConfigRequest) (*pb.UpdateGwConfigResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	if err := a.st.UpdateGatewayConfigByGwId(ctx, req.Conf, mac); err != nil {
		log.WithError(err).Error("Update conf to gw failed")
		return &pb.UpdateGwConfigResponse{Status: "Update config failed, please check your gateway connection."},
			status.Errorf(codes.Internal, "cannot update gateway config: %s", err)
	}

	return &pb.UpdateGwConfigResponse{
		Status: "Update gateway config file successful",
	}, nil
}

// Register will first try to get the gateway from provision server
func (a *GatewayAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.WithFields(log.Fields{
		"Sn":             req.Sn,
		"OrganizationID": req.OrganizationId,
	}).Info("API Register is called")

	if req.Sn == "" {
		return nil, status.Error(codes.InvalidArgument, "gateway sn number must not be empty string")
	}

	if valid, err := NewValidator().ValidateGlobalGatewaysAccess(ctx, authcus.Create, req.OrganizationId); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	// register gateway with current supernode on remote provisioning server
	provReq := psPb.RegisterGWRequest{
		Sn:            req.Sn,
		SuperNodeAddr: Service.SupernodeAddr,
		OrgId:         req.OrganizationId,
	}

	provConf := config.C.ProvisionServer

	provClient, err := provisionserver.CreateClientWithCert(provConf.ProvisionServer, provConf.CACert,
		provConf.TLSCert, provConf.TLSKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	resp, err := provClient.RegisterGW(ctx, &provReq)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	// add new firmware if new model is registered
	_, err = a.st.GetGatewayFirmware(resp.Model, true)
	if err == storage.ErrDoesNotExist {
		if _, err = a.st.AddGatewayFirmware(&GatewayFirmware{
			Model: resp.Model,
		}); err != nil {
			log.WithError(err).Errorf("Failed to add new firmware: %s", resp.Model)
		}
	}

	gateway := pb.Gateway{
		Id:          resp.Mac,
		Name:        fmt.Sprintf("Gateway_%s", resp.Sn),
		Description: fmt.Sprintf("Gateway Model: %s\nGateway OsVersion: %s\n", resp.Model, resp.OsVersion),
		Location: &common.Location{
			Latitude:  52.520008,
			Longitude: 13.404954,
			Altitude:  0,
			Source:    0,
			Accuracy:  0,
		},
		OrganizationId:   req.OrganizationId,
		DiscoveryEnabled: true,
		NetworkServerId:  0,
		GatewayProfileId: "",
		Boards:           []*pb.GatewayBoard{},
	}

	// get gateway profile id, always use the default one
	count, err := storage.GetGatewayProfileCount(ctx, storage.DB())
	if err != nil && err != storage.ErrDoesNotExist || count == 0 {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	gpList, err := storage.GetGatewayProfiles(ctx, storage.DB(), count, 0)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	for _, v := range gpList {
		if v.Name != "default_gateway_profile" {
			continue
		}

		gateway.GatewayProfileId = v.GatewayProfileID.String()
	}

	if gateway.GatewayProfileId == "" {
		return nil, status.Error(codes.NotFound, "Default gateway profile does not exist")
	}

	// get network server from gateway profile
	gpID, err := uuid.FromString(gateway.GatewayProfileId)
	if err != nil {
		log.WithError(err).Error("Gateway profile ID invalid")
		return nil, status.Error(codes.DataLoss, "Gateway profile ID invalid")
	}

	nServers, err := networkserver.Service.St.GetNetworkServerForGatewayProfileID(ctx, gpID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to load network servers: %s", err.Error())
	}

	gateway.NetworkServerId = nServers.ID

	// create gateway
	if err := a.storeGateway(ctx, &gateway, &Gateway{
		Model:        resp.Model,
		OsVersion:    resp.OsVersion,
		Statistics:   "",
		SerialNumber: resp.Sn,
	}); err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Status: "Successful",
	}, nil
}

func (a *GatewayAPI) GetGwPwd(ctx context.Context, req *pb.GetGwPwdRequest) (*pb.GetGwPwdResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	provConf := config.C.ProvisionServer
	provClient, err := provisionserver.CreateClientWithCert(provConf.ProvisionServer, provConf.CACert,
		provConf.TLSCert, provConf.TLSKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to connect to provisioning server")
	}

	resp, err := provClient.GetRootPWD(context.Background(), &psPb.GetRootPWDRequest{
		Sn:  req.Sn,
		Mac: req.GatewayId,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.GetGwPwdResponse{Password: resp.RootPassword}, nil
}

func (a *GatewayAPI) SetAutoUpdateFirmware(ctx context.Context, req *pb.SetAutoUpdateFirmwareRequest) (*pb.SetAutoUpdateFirmwareResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if valid, err := NewValidator().ValidateGatewayAccess(ctx, authcus.Read, mac); !valid || err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	if err := a.st.SetAutoUpdateFirmware(ctx, mac, req.AutoUpdate); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.SetAutoUpdateFirmwareResponse{Message: "Auto update firmware set successfully"}, nil
}

// GetGatewayList defines the get Gateway list request and response
func (a *GatewayAPI) GetGatewayList(ctx context.Context, req *api.GetGatewayListRequest) (*api.GetGatewayListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayList org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetGatewayListResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayList(ctx, &m2mServer.GetGatewayListRequest{
		OrgId:  req.OrgId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayListResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	var gatewayProfileList []*api.GSGatewayProfile
	for _, item := range resp.GwProfile {
		gatewayProfile := &api.GSGatewayProfile{
			Id:          item.Id,
			Mac:         item.Mac,
			FkGwNs:      item.FkGwNs,
			FkWallet:    item.FkWallet,
			Mode:        api.GatewayMode(item.Mode),
			CreateAt:    item.CreateAt,
			LastSeenAt:  item.LastSeenAt,
			OrgId:       item.OrgId,
			Description: item.Description,
			Name:        item.Name,
		}

		gatewayProfileList = append(gatewayProfileList, gatewayProfile)
	}

	return &api.GetGatewayListResponse{
		GwProfile: gatewayProfileList,
		Count:     resp.Count,
	}, status.Error(codes.OK, "")
}

// GetGatewayProfile defines the get Gateway Profile request and response
func (a *GatewayAPI) GetGatewayProfile(ctx context.Context, req *api.GetGSGatewayProfileRequest) (*api.GetGSGatewayProfileResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayProfile org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayProfile(ctx, &m2mServer.GetGSGatewayProfileRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGSGatewayProfileResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGSGatewayProfileResponse{
		GwProfile: &api.GSGatewayProfile{
			Id:          resp.GwProfile.Id,
			Mac:         resp.GwProfile.Mac,
			FkGwNs:      resp.GwProfile.FkGwNs,
			FkWallet:    resp.GwProfile.FkWallet,
			Mode:        api.GatewayMode(resp.GwProfile.Mode),
			CreateAt:    resp.GwProfile.CreateAt,
			LastSeenAt:  resp.GwProfile.LastSeenAt,
			OrgId:       resp.GwProfile.OrgId,
			Description: resp.GwProfile.Description,
			Name:        resp.GwProfile.Name,
		},
	}, status.Error(codes.OK, "")
}

// GetGatewayHistory defines the get Gateway History request and response
func (a *GatewayAPI) GetGatewayHistory(ctx context.Context, req *api.GetGatewayHistoryRequest) (*api.GetGatewayHistoryResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayHistory org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.GetGatewayHistory(ctx, &m2mServer.GetGatewayHistoryRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.GetGatewayHistoryResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.GetGatewayHistoryResponse{
		GwHistory: resp.GwHistory,
	}, status.Error(codes.OK, "")
}

// SetGatewayMode defines the set Gateway mode request and response
func (a *GatewayAPI) SetGatewayMode(ctx context.Context, req *api.SetGatewayModeRequest) (*api.SetGatewayModeResponse, error) {
	logInfo := "api/appserver_serves_ui/SetGatewayMode org=" + strconv.FormatInt(req.OrgId, 10)

	// verify if user is global admin
	userIsAdmin, err := NewValidator().IsGlobalAdmin(ctx)
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Internal, "unable to verify user: %s", err.Error())
	}
	// is user is not global admin, user must have accesss to this organization
	if !userIsAdmin {
		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, authcus.Read, req.OrgId); !valid || err != nil {
			return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
		}
	}

	m2mClient, err := m2m_client.GetPool().Get(config.C.M2MServer.M2MServer, []byte(config.C.M2MServer.CACert),
		[]byte(config.C.M2MServer.TLSCert), []byte(config.C.M2MServer.TLSKey))
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	gwClient := m2mServer.NewGSGatewayServiceClient(m2mClient)

	resp, err := gwClient.SetGatewayMode(ctx, &m2mServer.SetGatewayModeRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		GwMode: m2mServer.GatewayMode(req.GwMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetGatewayModeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
