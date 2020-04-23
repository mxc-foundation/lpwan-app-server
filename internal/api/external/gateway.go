package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	api "github.com/mxc-foundation/lpwan-app-server/api/ps_serves_appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/provisionserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"google.golang.org/grpc/status"
	"net/textproto"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/brocaar/lorawan"
	pb "github.com/mxc-foundation/lpwan-app-server/api/appserver_serves_ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/auth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-server/api/common"
	"github.com/mxc-foundation/lpwan-server/api/ns"
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

// Create creates the given gateway.
func (a *GatewayAPI) Create(ctx context.Context, req *pb.CreateGatewayRequest) (*empty.Empty, error) {
	if req.Gateway == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "gateway must not be nil")
	}

	if req.Gateway.Location == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "gateway.location must not be nil")
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewaysAccess(auth.Create, req.Gateway.OrganizationId))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// also validate that the network-server is accessible for the given organization
	err = a.validator.Validate(ctx, auth.ValidateOrganizationNetworkServerAccess(auth.Read, req.Gateway.OrganizationId, req.Gateway.NetworkServerId))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Gateway.Id)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	createReq := ns.CreateGatewayRequest{
		Gateway: &ns.Gateway{
			Id:               mac[:],
			Location:         req.Gateway.Location,
			RoutingProfileId: applicationServerID.Bytes(),
		},
	}

	if req.Gateway.GatewayProfileId != "" {
		gpID, err := uuid.FromString(req.Gateway.GatewayProfileId)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		createReq.Gateway.GatewayProfileId = gpID.Bytes()
	}

	for _, board := range req.Gateway.Boards {
		var gwBoard ns.GatewayBoard

		if board.FpgaId != "" {
			var fpgaID lorawan.EUI64
			if err := fpgaID.UnmarshalText([]byte(board.FpgaId)); err != nil {
				return nil, grpc.Errorf(codes.InvalidArgument, "fpga_id: %s", err)
			}
			gwBoard.FpgaId = fpgaID[:]
		}

		if board.FineTimestampKey != "" {
			var key lorawan.AES128Key
			if err := key.UnmarshalText([]byte(board.FineTimestampKey)); err != nil {
				return nil, grpc.Errorf(codes.InvalidArgument, "fine_timestamp_key: %s", err)
			}
			gwBoard.FineTimestampKey = key[:]
		}

		createReq.Gateway.Boards = append(createReq.Gateway.Boards, &gwBoard)
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		err = storage.CreateGateway(ctx, tx, &storage.Gateway{
			MAC:             mac,
			Name:            req.Gateway.Name,
			Description:     req.Gateway.Description,
			OrganizationID:  req.Gateway.OrganizationId,
			Ping:            req.Gateway.DiscoveryEnabled,
			NetworkServerID: req.Gateway.NetworkServerId,
			Latitude:        req.Gateway.Location.Latitude,
			Longitude:       req.Gateway.Location.Longitude,
			Altitude:        req.Gateway.Location.Altitude,
			FirstHeartbeat: 0,
			LastHeartbeat: 0,
		})
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		n, err := storage.GetNetworkServer(ctx, tx, req.Gateway.NetworkServerId)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		_, err = nsClient.CreateGateway(ctx, &createReq)
		if err != nil && grpc.Code(err) != codes.AlreadyExists {
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

// Get returns the gateway matching the given Mac.
func (a *GatewayAPI) Get(ctx context.Context, req *pb.GetGatewayRequest) (*pb.GetGatewayResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return nil, grpc.Errorf(codes.InvalidArgument, "gateway must not be nil")
	}

	if req.Gateway.Location == nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "gateway.location must not be nil")
	}

	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Gateway.Id)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Update, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
				return grpc.Errorf(codes.InvalidArgument, err.Error())
			}
			updateReq.Gateway.GatewayProfileId = gpID.Bytes()
		}

		for _, board := range req.Gateway.Boards {
			var gwBoard ns.GatewayBoard

			if board.FpgaId != "" {
				var fpgaID lorawan.EUI64
				if err := fpgaID.UnmarshalText([]byte(board.FpgaId)); err != nil {
					return grpc.Errorf(codes.InvalidArgument, "fpga_id: %s", err)
				}
				gwBoard.FpgaId = fpgaID[:]
			}

			if board.FineTimestampKey != "" {
				var key lorawan.AES128Key
				if err := key.UnmarshalText([]byte(board.FineTimestampKey)); err != nil {
					return grpc.Errorf(codes.InvalidArgument, "fine_timestamp_key: %s", err)
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
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Delete, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		//check if the gateway is in the board table. If not, the gateway is not belong to the MatchX, delete it.
		bd, err := storage.GetBoard(tx, mac)
		if err != nil {
			if err == storage.ErrDoesNotExist {

				err = storage.DeleteGateway(ctx, tx, mac)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}
				return nil
			} else {
				return helpers.ErrToRPCError(err)
			}
		}

		// if gateway is in the board table, send the request to provision server
		provReq := api.UnregisterGwRequest{
			Mac: mac.String(),
		}

		provConf := config.C.ProvisionServer

		provClient, err := provisionserver.GetPool().Get(provConf.ProvisionServer, []byte(provConf.CACert),
			[]byte(provConf.TLSCert), []byte(provConf.TLSKey))
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		resp, err := provClient.UnregisterGw(ctx, &provReq)
		if err != nil && grpc.Code(err) != codes.AlreadyExists {
			return err
		}

		// if the response is true, check the gateway connection
		// 0 = RequestStatus_SUCCESSFUL
		if resp.Status == 0 {
			// if cannot connect to the gateway (gateway offline), delete the data
			// ToDo: need to set the timeout
			_, err := mxConfDGet(ctx, bd.VpnAddr, "STAT", 250)
			if err != nil {
				err = storage.DeleteBoardByMac(ctx, tx, &bd)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}

				err = storage.DeleteGateway(ctx, tx, mac)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}
				return nil
			}

			// or to ping the gateway
			/*out, _ := exec.Command("ping", "192.168.0.124", "-c 5", "-w 1").Output()
			if strings.Contains(string(out), "Destination Host Unreachable") {
				fmt.Println("TANGO DOWN")
			} else {
				fmt.Println("IT'S ALIVEEE")
			}*/

			// if gateway is still online, user should first disconnect the gateway and delete again
			err = errors.New("The gateway is still online")

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

// GetStats gets the gateway statistics for the gateway with the given Mac.
func (a *GatewayAPI) GetStats(ctx context.Context, req *pb.GetGatewayStatsRequest) (*pb.GetGatewayStatsResponse, error) {
	var gatewayID lorawan.EUI64
	if err := gatewayID.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, gatewayID))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	start, err := ptypes.Timestamp(req.StartTimestamp)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	end, err := ptypes.Timestamp(req.EndTimestamp)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	_, ok := ns.AggregationInterval_value[strings.ToUpper(req.Interval)]
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad interval: %s", req.Interval)
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
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return grpc.Errorf(codes.InvalidArgument, "mac: %s", err)
	}

	err := a.validator.Validate(srv.Context(), auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
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
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	openVPNaddr, err := storage.GetOpenVPNByMac(ctx, storage.DB(), mac)
	if err != nil {
		log.WithError(err).Error("cannot get openVPN address from DB")
	}

	message, err := mxConfDGet(ctx, openVPNaddr, "GGLC", 250)
	if err != nil {
		log.WithError(err).Error("cannot connect to gw")
	}

	comments := []string{"/* radio_1 provides clock to concentrator */", "/* dBm */", "/* 8 channels maximum */",
		"/* dB */", "/* antenna gain, in dBi */", "/* [126..250] KHz */", "/* Lora MAC channel, 125kHz, all SF, 868.1 MHz */",
		"/* Lora MAC channel, 125kHz, all SF, 868.3 MHz */", "/* Lora MAC channel, 125kHz, all SF, 868.5 MHz */",
		"/* Lora MAC channel, 125kHz, all SF, 868.8 MHz */", "/* Lora MAC channel, 125kHz, all SF, 864.7 MHz */",
		"/* Lora MAC channel, 125kHz, all SF, 864.9 MHz */", "/* Lora MAC channel, 125kHz, all SF, 865.1 MHz */",
		"/* Lora MAC channel, 125kHz, all SF, 865.3 MHz */", "/* Lora MAC channel, 250kHz, SF7, 868.3 MHz */",
		"/* FSK 50kbps channel, 868.8 MHz */", "/* TX gain table, index 0 */", "/* TX gain table, index 1 */", "/* TX gain table, index 2 */",
		"/* TX gain table, index 3 */", "/* TX gain table, index 4 */", "/* TX gain table, index 5 */", "/* TX gain table, index 6 */",
		"/* TX gain table, index 7 */", "/* TX gain table, index 8 */", "/* TX gain table, index 9 */", "/* TX gain table, index 10 */",
		"/* TX gain table, index 11 */", "/* TX gain table, index 12 */", "/* TX gain table, index 13 */", "/* TX gain table, index 14 */",
		"/* TX gain table, index 15 */", "/* change with default server address/ports, or overwrite in local_conf.json */",
		"/* adjust the following parameters for your network */", "/* forward only valid packets */", "/* GPS configuration */",
		"/* GPS reference coordinates */"}

	for _, v := range comments {
		message = strings.Replace(message, v, "", -1)
	}

	return &pb.GetGwConfigResponse{
		Conf: message,
	}, nil
}

// UpdateGwConfig gateway configuration file
func (a *GatewayAPI) UpdateGwConfig(ctx context.Context, req *pb.UpdateGwConfigRequest) (*pb.UpdateGwConfigResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	openVPNaddr, err := storage.GetOpenVPNByMac(ctx, storage.DB(), mac)
	if err != nil {
		return &pb.UpdateGwConfigResponse{Status: "Update config failed."},
			grpc.Errorf(codes.Unauthenticated, "cannot get gateway address from db: %s", err)
	}

	if err := mxConfUpdate(openVPNaddr, req.Conf); err != nil {
		log.WithError(err).Error("Update conf to gw failed")
		return &pb.UpdateGwConfigResponse{Status: "Update config failed, please check your gateway connection."},
			grpc.Errorf(codes.Unauthenticated, "cannot update gateway config: %s", err)
	}

	return &pb.UpdateGwConfigResponse{Status: "successful"}, nil
}

// Register will first try to get the gateway from provision server
func (a *GatewayAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Sn == "" {
		return &pb.RegisterResponse{Status: "gateway sn number must not be nil"}, status.Errorf(codes.InvalidArgument, "")
	}
	err := a.validator.Validate(ctx, auth.ValidateGatewaysAccess(auth.Create, req.OrganizationId))
	if err != nil {
		return &pb.RegisterResponse{Status: "authentication failed"}, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// also validate that the network-server is accessible for the given organization
	/*err = a.validator.Validate(ctx, auth.ValidateOrganizationNetworkServerAccess(auth.Read, req.Gateway.OrganizationId, req.Gateway.NetworkServerId))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}*/

	snAddr := os.Getenv("REMOTE_SERVER_NAME")

	// send the req to provision server
	provReq := api.RegisterGWRequest{
		Sn:            req.Sn,
		SuperNodeAddr: snAddr,
		OrgId:         req.OrganizationId,
	}

	provConf := config.C.ProvisionServer

	provClient, err := provisionserver.GetPool().Get(provConf.ProvisionServer, []byte(provConf.CACert),
		[]byte(provConf.TLSCert), []byte(provConf.TLSKey))
	if err != nil {
		return &pb.RegisterResponse{Status: "cannot connect to provisioning server"}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := provClient.RegisterGW(ctx, &provReq)
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		return &pb.RegisterResponse{Status: "cannot get the response from provisioning server"}, status.Errorf(codes.Unavailable, err.Error())
	}

	switch resp.Status {
	case 2:
		return &pb.RegisterResponse{Status: "please turn on your gateway"}, nil
	case 1:
		return &pb.RegisterResponse{Status: "please delete the gateway from previous supernode"}, nil
	case 0:
		err = storage.Transaction(func(tx sqlx.Ext) error {
			// get mac from provision server (resp.mac)
			var mac lorawan.EUI64
			if err := mac.UnmarshalText([]byte(resp.Mac)); err != nil {
				return status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
			}

			//check if gw already in db. If yes, update board table. No, create gw and insert board table.
			gw, err := storage.GetGateway(ctx, tx, mac, false)
			if err != nil {
				if err != storage.ErrDoesNotExist {
					return helpers.ErrToRPCError(err)
				}
			}

			if gw.Name == "" {
				var defLocation common.Location

				createReq := ns.CreateGatewayRequest{
					Gateway: &ns.Gateway{
						Id:               mac[:],
						Location:         &defLocation,
						RoutingProfileId: applicationServerID.Bytes(),
					},
				}

				// gateway profileID and gateway boards has been deleted in this func

				NetworkServers, err := storage.GetNetworkServers(ctx, tx, 1, 0)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}

				defNetworkServerID := NetworkServers[0].ID

				err = storage.Transaction(func(tx sqlx.Ext) error {
					err = storage.CreateGateway(ctx, tx, &storage.Gateway{
						MAC:             mac,
						Name:            "Default_Gateway",
						Description:     "MXC_Gateway",
						OrganizationID:  req.OrganizationId,
						Ping:            false,
						NetworkServerID: defNetworkServerID,
						Latitude:        0,
						Longitude:       0,
						Altitude:        0,
					})
					if err != nil {
						return helpers.ErrToRPCError(err)
					}

					n, err := storage.GetNetworkServer(ctx, tx, defNetworkServerID)
					if err != nil {
						return helpers.ErrToRPCError(err)
					}

					nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
					if err != nil {
						return helpers.ErrToRPCError(err)
					}
					_, err = nsClient.CreateGateway(ctx, &createReq)
					if err != nil && grpc.Code(err) != codes.AlreadyExists {
						return err
					}

					// store data to the board table
					var bd storage.Board

					bd.MAC = mac
					bd.SN = &resp.Sn
					bd.Model = resp.Model
					bd.VpnAddr = resp.VpnAddr
					bd.OsVersion = &resp.OsVersion

					err = storage.CreateBoard(tx, &bd)
					if err != nil {
						return helpers.ErrToRPCError(err)
					}
					return nil
				})
				if err != nil {
					return err
				}

				redisConn := storage.RedisPool().Get()
				defer redisConn.Close()
				redisConn.Do("DEL", GatewayLocationsRedisKey)

			} else {
				bd, err := storage.GetBoard(tx, mac)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}

				bd.SN = &req.Sn
				bd.VpnAddr = resp.VpnAddr

				err = storage.UpdateVPNAddr(ctx, tx, &bd)
				if err != nil {
					return helpers.ErrToRPCError(err)
				}
			}
			return nil
		})
		if err != nil {
			return &pb.RegisterResponse{Status: "storage transaction error"}, status.Errorf(codes.Unavailable, err.Error())
		}
		return &pb.RegisterResponse{Status: resp.Status.String()}, nil
	}

	return &pb.RegisterResponse{Status: resp.Status.String()}, nil
}

func (a *GatewayAPI) GetGwPwd(ctx context.Context, req *pb.GetGwPwdRequest) (*pb.GetGwPwdResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	err := a.validator.Validate(ctx, auth.ValidateGatewayAccess(auth.Read, mac))
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// send the req to provision server
	provReq := api.GetRootPWDRequest{
		Sn: req.Sn,
	}

	provConf := config.C.ProvisionServer

	provClient, err := provisionserver.GetPool().Get(provConf.ProvisionServer, []byte(provConf.CACert),
		[]byte(provConf.TLSCert), []byte(provConf.TLSKey))
	if err != nil {
		return &pb.GetGwPwdResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	resp, err := provClient.GetRootPWD(ctx, &provReq)
	if err != nil && grpc.Code(err) != codes.AlreadyExists {
		return &pb.GetGwPwdResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &pb.GetGwPwdResponse{
		Password: resp.RootPWD,
	}, nil
}

// mxConfDGet connects to gateway though openVPN and get config file
func mxConfDGet(ctx context.Context, ip string, cmd string, rCode int) (string, error) {
	rpConn, err := textproto.Dial("tcp", fmt.Sprintf("%s:75", ip))
	if err != nil {
		return "", err
	}
	defer rpConn.Close()
	if _, _, err = rpConn.ReadResponse(220); err != nil {
		return "", err
	}
	if _, err = rpConn.Cmd(cmd); err != nil {
		return "", err
	}
	_, message, err := rpConn.ReadResponse(rCode)
	if err != nil {
		return "", err
	}
	if _, err = rpConn.Cmd("QUIT"); err != nil {
		return "", err
	}
	if _, _, err = rpConn.ReadResponse(221); err != nil {
		return "", err
	}
	return message, nil
}

// mxConfUpdate connects to gateway though openVPN and update config file
func mxConfUpdate(ip string, conf string) error {
	rpConn, err := textproto.Dial("tcp", fmt.Sprintf("%s:75", ip))
	if err != nil {
		return err
	}
	defer rpConn.Close()

	if _, _, err = rpConn.ReadResponse(220); err != nil {
		return err
	}
	if _, err = rpConn.Cmd("WGLC"); err != nil {
		return err
	}
	if _, _, err = rpConn.ReadResponse(354); err != nil {
		return err
	}

	w := rpConn.DotWriter()
	if _, err = w.Write([]byte(conf)); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	if _, _, err = rpConn.ReadResponse(250); err != nil {
		return err
	}

	if _, err = rpConn.Cmd("QUIT"); err != nil {
		return err
	}

	return nil
}
