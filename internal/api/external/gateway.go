package external

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	gwm "github.com/mxc-foundation/lpwan-app-server/internal/gateway-manager"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"
	"github.com/mxc-foundation/lpwan-server/api/common"
	"github.com/mxc-foundation/lpwan-server/api/ns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GatewayLocationsRedisKey defines the gateway location based on redis key
const GatewayLocationsRedisKey = "gateway_locations"

// GatewayAPI exports the Gateway related functions.
type GatewayAPI struct {
	validator auth.Validator
}

// NewGatewayAPI creates a new GatewayAPI.
func NewGatewayAPI(validator auth.Validator) *GatewayAPI {
	return &GatewayAPI{
		validator: validator,
	}
}

// BatchResetDefaultGatewatConfig reset gateways config to default config matching organization list
func (a *GatewayAPI) BatchResetDefaultGatewatConfig(ctx context.Context, req *pb.BatchResetDefaultGatewatConfigRequest) (*pb.BatchResetDefaultGatewatConfigResponse, error) {
	log.WithFields(log.Fields{
		"organization_list": req.OrganizationList,
	}).Info("BatchResetDefaultGatewatConfig is called")

	// check user permission, only global admin allowed
	isAdmin, err := a.validator.GetIsAdmin(ctx)
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
		count, err := storage.GetOrganizationCount(ctx, storage.DB(), "")
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		limit := 100
		for offset := 0; offset <= count/limit; offset ++ {
			list, err := storage.GetOrganizationIDList(storage.DB(), limit, offset, "")
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

		err := resetDefaultGatewayConfigByOrganizationID(ctx, int64(v))
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

func resetDefaultGatewayConfigByOrganizationID(ctx context.Context, orgID int64) error {
	count, err := storage.GetGatewayCountForOrganizationID(ctx, storage.DB(), orgID, "")
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(fmt.Sprintf("There is no gateway in organization : %d", orgID))
	}

	limit := 100
	for offset := 0; offset <= count/limit; offset ++ {
		gwList, err := storage.GetGatewaysForOrganizationID(ctx, storage.DB(), orgID, limit, offset, "")
		if err != nil {
			return err
		}

		for _, v := range gwList {
			err := gwm.GetDefaultGatewayConfig(ctx, &v)
			if err != nil {
				return err
			}

			err = storage.UpdateGatewayConfigByGwId(ctx, storage.DB(), v.Config, v.MAC)
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
	isAdmin, err := a.validator.GetIsAdmin(ctx)
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

	gw, err := storage.GetGateway(ctx, storage.DB(), mac, true)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	err = gwm.GetDefaultGatewayConfig(ctx, &gw)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	err = storage.UpdateGatewayConfigByGwId(ctx, storage.DB(), gw.Config, gw.MAC)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	
	return &pb.ResetDefaultGatewatConfigByIDResponse{}, status.Error(codes.OK, "")
}

// InsertNewDefaultGatewayConfig insert given new default gateway config
func (a *GatewayAPI) InsertNewDefaultGatewayConfig(ctx context.Context, req *pb.InsertNewDefaultGatewayConfigRequest) (*pb.InsertNewDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model": req.Model,
		"region": req.Region,
	}).Info("InsertNewDefaultGatewayConfig is called")

	// check user permission, only global admin allowed
	isAdmin, err := a.validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := storage.DefaultGatewayConfig{
		Model:         req.Model,
		Region:        req.Region,
		DefaultConfig: strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", storage.SupernodeAddr, -1),
	}

	err = storage.GetDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err == nil {
		// config already exist, no need to insert
		return nil, status.Errorf(codes.AlreadyExists, "model=%s, region=%s", req.Model, req.Region)
	} else if err != storage.ErrDoesNotExist {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = storage.AddNewDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.InsertNewDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// UpdateNewDefaultGatewayConfig update default gateway config matching model and region
func (a *GatewayAPI) UpdateDefaultGatewayConfig(ctx context.Context, req *pb.UpdateDefaultGatewayConfigRequest) (*pb.UpdateDefaultGatewayConfigResponse, error){
	log.WithFields(log.Fields{
		"model": req.Model,
		"region": req.Region,
	}).Info("UpdateDefaultGatewayConfig is called")

	// check user permission
	isAdmin, err := a.validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := storage.DefaultGatewayConfig{
		Model:         req.Model,
		Region:        req.Region,
	}

	err = storage.GetDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	defaultGatewayConfig.DefaultConfig = strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", storage.SupernodeAddr, -1)
	err = storage.UpdateDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.UpdateDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// GetDefaultGatewayConfig get content of default gateway config matching model and region
func (a *GatewayAPI) GetDefaultGatewayConfig(ctx context.Context, req *pb.GetDefaultGatewayConfigRequest) (*pb.GetDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model": req.Model,
		"region": req.Region,
	}).Info("GetDefaultGatewayConfig is called")

	// check user permission
	isAdmin, err := a.validator.GetIsAdmin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "not authorized for this operation")
	}

	defaultGatewayConfig := storage.DefaultGatewayConfig{
		Model:         req.Model,
		Region:        req.Region,
	}

	err = storage.GetDefaultGatewayConfig(storage.DB(), &defaultGatewayConfig)
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

	err := a.validator.Validate(ctx, auth.ValidateGatewaysAccess(auth.Create, req.Gateway.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// also validate that the network-server is accessible for the given organization
	err = a.validator.Validate(ctx, auth.ValidateOrganizationNetworkServerAccess(auth.Read, req.Gateway.OrganizationId, req.Gateway.NetworkServerId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err := storeGateway(ctx, req.Gateway, &storage.Gateway{}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func storeGateway(ctx context.Context, req *pb.Gateway, defaultGw *storage.Gateway) (err error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	createReq := ns.CreateGatewayRequest{
		Gateway: &ns.Gateway{
			Id:               mac[:],
			Location:         req.Location,
			RoutingProfileId: applicationServerID.Bytes(),
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
	err = gwm.GetDefaultGatewayConfig(ctx, defaultGw)
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		err = storage.CreateGateway(ctx, tx, &storage.Gateway{
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

		n, err := storage.GetNetworkServer(ctx, tx, req.NetworkServerId)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		_, err = nsClient.CreateGateway(ctx, &createReq)
		if err != nil && status.Code(err) != codes.AlreadyExists {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	redisConn := storage.RedisPool().Get()
	defer redisConn.Close()

	_, _ = redisConn.Do("DEL", GatewayLocationsRedisKey)

	return nil
}

// Get returns the gateway matching the given Mac.
func (a *GatewayAPI) Get(ctx context.Context, req *pb.GetGatewayRequest) (*pb.GetGatewayResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	gw, err := storage.GetGateway(ctx, storage.DB(), mac, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := storage.GetNetworkServer(ctx, storage.DB(), gw.NetworkServerID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	getResp, err := nsClient.GetGateway(ctx, &ns.GetGatewayRequest{
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

	return &resp, err
}

// List lists the gateways.
func (a *GatewayAPI) List(ctx context.Context, req *pb.ListGatewayRequest) (*pb.ListGatewayResponse, error) {
	err := a.validator.Validate(ctx, auth.ValidateGatewaysAccess(auth.List, req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var count int
	var gws []storage.Gateway

	if req.OrganizationId == 0 {
		isAdmin, err := a.validator.GetIsAdmin(ctx)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		if isAdmin {
			// in case of admin user list all gateways
			count, err = storage.GetGatewayCount(ctx, storage.DB(), req.Search)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}

			gws, err = storage.GetGateways(ctx, storage.DB(), int(req.Limit), int(req.Offset), req.Search)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		} else {
			// filter result based on user
			username, err := a.validator.GetUsername(ctx)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
			count, err = storage.GetGatewayCountForUser(ctx, storage.DB(), username, req.Search)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
			gws, err = storage.GetGatewaysForUser(ctx, storage.DB(), username, int(req.Limit), int(req.Offset), req.Search)
			if err != nil {
				return nil, helpers.ErrToRPCError(err)
			}
		}
	} else {
		count, err = storage.GetGatewayCountForOrganizationID(ctx, storage.DB(), req.OrganizationId, req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		gws, err = storage.GetGatewaysForOrganizationID(ctx, storage.DB(), req.OrganizationId, int(req.Limit), int(req.Offset), req.Search)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
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

	redisConn := storage.RedisPool().Get()
	defer redisConn.Close()

	resultJSON, err := redis.Bytes(redisConn.Do("GET", GatewayLocationsRedisKey))
	if err == nil {
		json.Unmarshal(resultJSON, &result)
	}

	if len(result) == 0 {
		gwsLoc, err := storage.GetGatewaysLoc(ctx, storage.DB(), viper.GetInt("application_server.gateways_locations_limit"))
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Update, mac))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		gw, err := storage.GetGateway(ctx, tx, mac, true)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		gw.Name = req.Gateway.Name
		gw.Description = req.Gateway.Description
		gw.Ping = req.Gateway.DiscoveryEnabled
		gw.Latitude = req.Gateway.Location.Latitude
		gw.Longitude = req.Gateway.Location.Longitude
		gw.Altitude = req.Gateway.Location.Altitude

		err = storage.UpdateGateway(ctx, tx, &gw)
		if err != nil {
			return helpers.ErrToRPCError(err)
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
				return status.Error(codes.InvalidArgument, err.Error())
			}
			updateReq.Gateway.GatewayProfileId = gpID.Bytes()
		}

		for _, board := range req.Gateway.Boards {
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

			updateReq.Gateway.Boards = append(updateReq.Gateway.Boards, &gwBoard)
		}

		n, err := storage.GetNetworkServer(ctx, tx, gw.NetworkServerID)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		_, err = nsClient.UpdateGateway(ctx, &updateReq)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	redisConn := storage.RedisPool().Get()
	defer redisConn.Close()

	redisConn.Do("DEL", GatewayLocationsRedisKey)

	return &empty.Empty{}, nil
}

// Delete deletes the gateway matching the given ID.
func (a *GatewayAPI) Delete(ctx context.Context, req *pb.DeleteGatewayRequest) (*empty.Empty, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Delete, mac))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		err = storage.DeleteGateway(ctx, tx, mac)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	redisConn := storage.RedisPool().Get()
	defer redisConn.Close()

	redisConn.Do("DEL", GatewayLocationsRedisKey)

	return &empty.Empty{}, nil
}

// GetStats gets the gateway statistics for the gateway with the given Mac.
func (a *GatewayAPI) GetStats(ctx context.Context, req *pb.GetGatewayStatsRequest) (*pb.GetGatewayStatsResponse, error) {
	var gatewayID lorawan.EUI64
	if err := gatewayID.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, gatewayID))
	if err != nil {
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

	metrics, err := storage.GetMetrics(ctx, storage.RedisPool(), storage.AggregationInterval(strings.ToUpper(req.Interval)), "gw:"+gatewayID.String(), start, end)
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	ping, pingRX, err := storage.GetLastGatewayPingAndRX(ctx, storage.DB(), mac)
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

	err := a.validator.Validate(srv.Context(), auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := storage.GetNetworkServerForGatewayMAC(srv.Context(), storage.DB(), mac)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	streamClient, err := nsClient.StreamFrameLogsForGateway(srv.Context(), &ns.StreamFrameLogsForGatewayRequest{
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

		up, down, err := convertUplinkAndDownlinkFrames(resp.GetUplinkFrameSet(), resp.GetDownlinkFrame(), false)
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	gwConfig, err := storage.GetGatewayConfigByGwId(ctx, storage.DB(), mac)
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	if err := storage.UpdateGatewayConfigByGwId(ctx, storage.DB(), req.Conf, mac); err != nil {
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
	err := a.validator.Validate(ctx, auth.ValidateGatewaysAccess(auth.Create, req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}

	// register gateway with current supernode on remote provisioning server
	provReq := api.RegisterGWRequest{
		Sn:            req.Sn,
		SuperNodeAddr: storage.SupernodeAddr,
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
	_, err = storage.GetGatewayFirmware(storage.DB(), resp.Model, true)
	if err == storage.ErrDoesNotExist {
		if _, err = storage.AddGatewayFirmware(storage.DB(), &storage.GatewayFirmware{
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

	nServers, err := storage.GetNetworkServerForGatewayProfileID(ctx, storage.DB(), gpID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to load network servers: %s", err.Error())
	}

	gateway.NetworkServerId = nServers.ID

	// create gateway
	if err := storeGateway(ctx, &gateway, &storage.Gateway{
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	provConf := config.C.ProvisionServer
	provClient, err := provisionserver.CreateClientWithCert(provConf.ProvisionServer, provConf.CACert,
		provConf.TLSCert, provConf.TLSKey)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to connect to provisioning server")
	}

	resp, err := provClient.GetRootPWD(context.Background(), &api.GetRootPWDRequest{
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

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %s", err)
	}

	if err := storage.SetAutoUpdateFirmware(ctx, storage.DB(), mac, req.AutoUpdate); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &pb.SetAutoUpdateFirmwareResponse{Message: "Auto update firmware set successfully"}, nil
}
