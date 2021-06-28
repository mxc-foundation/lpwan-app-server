package external

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lib/pq/hstore"

	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	pb "github.com/mxc-foundation/lpwan-app-server/api/m2m-serves-appserver"
	psPb "github.com/mxc-foundation/lpwan-app-server/api/ps-serves-appserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/mannr"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/auth"
	metricsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/types"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"

	gw "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
	orgs "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
)

// GatewayAPI exports the Gateway related functions.
type GatewayAPI struct {
	st     *pgstore.PgStore
	config GwConfig
	auth   auth.Authenticator
	pscli  psPb.ProvisionClient
	mxpCli pb.GSGatewayServiceClient

	serverAddr     string
	bindOldGateway string
	bindNewGateway string
}

type GwConfig struct {
	ApplicationServerID uuid.UUID
	ServerAddr          string
	EnableSTC           bool
}

// NewGatewayAPI creates a new GatewayAPI.
func NewGatewayAPI(st *pgstore.PgStore, auth auth.Authenticator, config GwConfig, pscli *pscli.Client,
	mxpCli *mxpcli.Client, serverAddr, bindOldGateway, bindNewGateway string) *GatewayAPI {
	return &GatewayAPI{
		st:             st,
		auth:           auth,
		config:         config,
		pscli:          pscli.GetPServerClient(),
		mxpCli:         mxpCli.GetM2MGatewayServiceClient(),
		serverAddr:     serverAddr,
		bindNewGateway: bindNewGateway,
		bindOldGateway: bindOldGateway,
	}
}

func (a *GatewayAPI) RegisterReseller(ctx context.Context, req *api.RegisterResellerRequest) (*api.RegisterResellerResponse, error) {
	// organization admin user or organization user should both be able to register organization as reseller for the gateway
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// check whether register reseller feature is enabled
	if !a.config.EnableSTC {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied: stc not enabled")
	}

	// check whether the gateway already has a reseller
	stcOrgID, err := a.st.GetSTCOrgIDForGateway(ctx, req.ManufacturerNr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if stcOrgID != 0 {
		return nil, status.Errorf(codes.PermissionDenied, "a reseller has already been registered for this gateway")
	}

	// check whether the gateway is a new and unused gateway
	resp, err := a.pscli.VerifyGWStatus(ctx, &psPb.VerifyGWStatusRequest{
		ManufacturerNr: req.ManufacturerNr,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	if resp.RegisteredSupernode != "" || resp.TotalOnlineSeconds != 0 {
		return nil, status.Error(codes.PermissionDenied, "not brand new gateway")
	}

	err = a.st.AddGatewayReseller(ctx, req.ManufacturerNr, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &api.RegisterResellerResponse{}, nil
}

func (a *GatewayAPI) ManualTriggerUpdateFirmware(ctx context.Context, req *api.ManualTriggerUpdateFirmwareRequest) (*api.ManualTriggerUpdateFirmwareResponse, error) {
	log.Info("ManualTriggerUpdateFirmware is called")

	if err := gw.UpdateFirmware(ctx, a.st, a.bindOldGateway, a.bindNewGateway, a.serverAddr, a.pscli); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.ManualTriggerUpdateFirmwareResponse{}, nil
}

// BatchResetDefaultGatewatConfig reset gateways config to default config matching organization list
func (a *GatewayAPI) BatchResetDefaultGatewatConfig(ctx context.Context, req *api.BatchResetDefaultGatewatConfigRequest) (*api.BatchResetDefaultGatewatConfigResponse, error) {
	log.WithFields(log.Fields{
		"organization_list": req.OrganizationList,
	}).Info("BatchResetDefaultGatewatConfig is called")

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// if process for any organizaiton failed, return this error message to user for retry
	var failedList []string
	var succeededList []string
	var organizationList []int

	if req.OrganizationList == "all" {
		// reset for all organizations
		count, err := organization.GetOrganizationCount(ctx, orgs.OrganizationFilters{})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		limit := 100
		for offset := 0; offset <= count/limit; offset += limit {
			list, err := organization.GetOrganizationIDList(ctx, limit, offset, "")
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

	return &api.BatchResetDefaultGatewatConfigResponse{
		Status: fmt.Sprintf("following organization failed: %s \n following organization succeeded: %s",
			strings.Join(failedList, ","), strings.Join(succeededList, ",")),
	}, status.Error(codes.OK, "")
}

func (a *GatewayAPI) resetDefaultGatewayConfigByOrganizationID(ctx context.Context, orgID int64) error {
	filters := GatewayFilters{OrganizationID: orgID}
	count, err := a.st.GetGatewayCount(ctx, filters)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New(fmt.Sprintf("There is no gateway in organization : %d", orgID))
	}

	limit := 100
	for offset := 0; offset <= count/limit; offset += limit {
		filters.Limit = limit
		filters.Offset = offset
		gwList, err := a.st.GetGateways(ctx, filters)
		if err != nil {
			return err
		}

		for _, v := range gwList {
			err := a.getDefaultGatewayConfig(ctx, &Gateway{
				MAC:             v.MAC,
				CreatedAt:       v.CreatedAt,
				UpdatedAt:       v.UpdatedAt,
				FirstSeenAt:     v.FirstSeenAt,
				LastSeenAt:      v.LastSeenAt,
				Name:            v.Name,
				Description:     v.Description,
				OrganizationID:  v.OrganizationID,
				NetworkServerID: v.NetworkServerID,
				Model:           v.Model,
				Config:          v.Config,
			})
			if err != nil {
				return err
			}

			err = a.st.UpdateGatewayConfigByGwID(ctx, v.Config, v.MAC)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ResetDefaultGatewatConfigByID reste gateway config to default config matching gateway id
func (a *GatewayAPI) ResetDefaultGatewatConfigByID(ctx context.Context, req *api.ResetDefaultGatewatConfigByIDRequest) (*api.ResetDefaultGatewatConfigByIDResponse, error) {
	log.WithFields(log.Fields{
		"gateway_id": req.Id,
	}).Info("ResetDefaultGatewatConfigByID is called")

	// check user permission, only global admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
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

	err = a.st.UpdateGatewayConfigByGwID(ctx, gw.Config, gw.MAC)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.ResetDefaultGatewatConfigByIDResponse{}, status.Error(codes.OK, "")
}

// InsertNewDefaultGatewayConfig insert given new default gateway config
func (a *GatewayAPI) InsertNewDefaultGatewayConfig(ctx context.Context, req *api.InsertNewDefaultGatewayConfigRequest) (*api.InsertNewDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("InsertNewDefaultGatewayConfig is called")

	// check user permission, only global admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:         req.Model,
		Region:        req.Region,
		DefaultConfig: strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", a.config.ServerAddr, -1),
	}

	err = a.st.GetDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err == nil {
		// config already exist, no need to insert
		return nil, status.Errorf(codes.AlreadyExists, "model=%s, region=%s", req.Model, req.Region)
	} else if err != errHandler.ErrDoesNotExist {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = a.st.AddNewDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.InsertNewDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// UpdateNewDefaultGatewayConfig update default gateway config matching model and region
func (a *GatewayAPI) UpdateDefaultGatewayConfig(ctx context.Context, req *api.UpdateDefaultGatewayConfigRequest) (*api.UpdateDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("UpdateDefaultGatewayConfig is called")

	// check user permission, only global admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  req.Model,
		Region: req.Region,
	}

	err = a.st.GetDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	defaultGatewayConfig.DefaultConfig = strings.Replace(req.DefaultConfig, "{{ .ServerAddr }}", a.config.ServerAddr, -1)
	err = a.st.UpdateDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.UpdateDefaultGatewayConfigResponse{}, status.Error(codes.OK, "")
}

// GetDefaultGatewayConfig get content of default gateway config matching model and region
func (a *GatewayAPI) GetDefaultGatewayConfig(ctx context.Context, req *api.GetDefaultGatewayConfigRequest) (*api.GetDefaultGatewayConfigResponse, error) {
	log.WithFields(log.Fields{
		"model":  req.Model,
		"region": req.Region,
	}).Info("GetDefaultGatewayConfig is called")

	// check user permission, only global admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  req.Model,
		Region: req.Region,
	}

	err = a.st.GetDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	return &api.GetDefaultGatewayConfigResponse{DefaultConfig: defaultGatewayConfig.DefaultConfig}, status.Error(codes.OK, "")
}

// Create creates the given gateway.
func (a *GatewayAPI) Create(ctx context.Context, req *api.CreateGatewayRequest) (*empty.Empty, error) {
	var response empty.Empty
	if req.Gateway == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway must not be nil")
	}

	if req.Gateway.Location == nil {
		return nil, status.Error(codes.InvalidArgument, "gateway.location must not be nil")
	}

	// check user permission, only gateway admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.Gateway.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	nsAllowed, err := a.st.OrganizationCanAccessNetworkServer(ctx, req.Gateway.OrganizationId, req.Gateway.NetworkServerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if !cred.IsGatewayAdmin || !nsAllowed {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if err := a.storeGateway(ctx, req.Gateway, &Gateway{}); err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *GatewayAPI) getDefaultGatewayConfig(ctx context.Context, gw *Gateway) error {
	if !strings.HasPrefix(gw.Model, "MX19") {
		return nil
	}

	n, err := a.st.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get network server %d", gw.NetworkServerID)
		return errors.Wrapf(err, "GetDefaultGatewayConfig")
	}

	defaultGatewayConfig := DefaultGatewayConfig{
		Model:  gw.Model,
		Region: n.Region,
	}

	err = a.st.GetDefaultGatewayConfig(ctx, &defaultGatewayConfig)
	if err != nil {
		return errors.Wrapf(err, "Failed to get default gateway config for model= %s, region= %s", defaultGatewayConfig.Model, defaultGatewayConfig.Region)
	}

	gw.Config = strings.Replace(defaultGatewayConfig.DefaultConfig, "{{ .GatewayID }}", gw.MAC.String(), -1)
	return nil
}

func (a *GatewayAPI) storeGateway(ctx context.Context, req *api.Gateway, defaultGw *Gateway) (err error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	createReq := ns.CreateGatewayRequest{
		Gateway: &ns.Gateway{
			Id:               mac[:],
			Location:         req.Location,
			RoutingProfileId: a.config.ApplicationServerID.Bytes(),
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

	tags := hstore.Hstore{
		Map: make(map[string]sql.NullString),
	}
	for k, v := range req.Tags {
		tags.Map[k] = sql.NullString{Valid: true, String: v}
	}

	metadata := hstore.Hstore{
		Map: make(map[string]sql.NullString),
	}
	for k, v := range req.Metadata {
		tags.Map[k] = sql.NullString{Valid: true, String: v}
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, store *pgstore.PgStore) error {
		if err := gw.AddGateway(ctx, store, &Gateway{
			MAC:                mac,
			Name:               req.Name,
			Description:        req.Description,
			OrganizationID:     req.OrganizationId,
			Ping:               req.DiscoveryEnabled,
			NetworkServerID:    req.NetworkServerId,
			GatewayProfileID:   &req.GatewayProfileId,
			Latitude:           req.Location.Latitude,
			Longitude:          req.Location.Longitude,
			Altitude:           req.Location.Altitude,
			Tags:               tags,
			Metadata:           metadata,
			Model:              defaultGw.Model,
			FirstHeartbeat:     0,
			LastHeartbeat:      0,
			Config:             defaultGw.Config,
			OsVersion:          defaultGw.OsVersion,
			Statistics:         defaultGw.Statistics,
			SerialNumber:       defaultGw.SerialNumber,
			FirmwareHash:       types.MD5SUM{},
			AutoUpdateFirmware: true,
		}, createReq, a.mxpCli); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Get returns the gateway matching the given Mac.
func (a *GatewayAPI) Get(ctx context.Context, req *api.GetGatewayRequest) (*api.GetGatewayResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}
	// first check whether user is an authorized user
	_, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	item, err := a.st.GetGateway(ctx, mac, false)
	if err != nil {
		if err != errHandler.ErrDoesNotExist {
			return nil, status.Errorf(codes.NotFound, "gateway with mac %s does not exist", req.Id)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// check whether user has access to the gateway information
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(item.OrganizationID))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		// user is neither global admin nor organization user, check whether user is reseller of the gateway
		if a.config.EnableSTC && item.STCOrgID != nil && *item.STCOrgID != 0 {
			stcOrgID := *item.STCOrgID
			_, err = a.st.GetOrganizationUser(ctx, stcOrgID, cred.UserID)
			if err != nil {
				if err == errHandler.ErrDoesNotExist {
					return nil, status.Errorf(codes.PermissionDenied, "permission denied")
				}
				return nil, status.Errorf(codes.Internal, err.Error())
			}

		} else {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}
	}

	gw, err := a.st.GetGateway(ctx, mac, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := a.st.GetNetworkServer(ctx, gw.NetworkServerID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nStruct := &nsmod.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nStruct.GetNetworkServiceClient()
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	getResp, err := client.GetGateway(ctx, &ns.GetGatewayRequest{
		Id: mac[:],
	})
	if err != nil {
		return nil, err
	}

	resp := api.GetGatewayResponse{
		Gateway: &api.Gateway{
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
		var gwBoard api.GatewayBoard

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

func (a *GatewayAPI) listGateways(ctx context.Context, filters GatewayFilters) (int64, []*api.GatewayListItem, error) {
	var gwList []*api.GatewayListItem
	count, err := a.st.GetGatewayCount(ctx, filters)

	if err != nil {
		return 0, nil, helpers.ErrToRPCError(err)
	}

	gws, err := a.st.GetGateways(ctx, filters)
	if err != nil {
		return 0, nil, helpers.ErrToRPCError(err)
	}

	for _, gw := range gws {
		row := api.GatewayListItem{
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

		// add prefix 'stc_' to gateway name if gateway's stc_org_id equals to req.OrganizationID
		// if stc is not enabled, gw.STCOrgID should either be nil or *gw.STCOrgID does not equal to req.OrganizationID
		if gw.STCOrgID != nil && *gw.STCOrgID == filters.OrganizationID {
			row.Name = "STC_" + row.Name
		}

		row.CreatedAt, err = ptypes.TimestampProto(gw.CreatedAt)
		if err != nil {
			return 0, nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(gw.UpdatedAt)
		if err != nil {
			return 0, nil, helpers.ErrToRPCError(err)
		}

		if gw.FirstSeenAt != nil {
			row.FirstSeenAt, err = ptypes.TimestampProto(*gw.FirstSeenAt)
			if err != nil {
				return 0, nil, helpers.ErrToRPCError(err)
			}
		}
		if gw.LastSeenAt != nil {
			row.LastSeenAt, err = ptypes.TimestampProto(*gw.LastSeenAt)
			if err != nil {
				return 0, nil, helpers.ErrToRPCError(err)
			}
		}

		gwList = append(gwList, &row)
	}

	return int64(count), gwList, nil
}

// List lists the gateways.
func (a *GatewayAPI) List(ctx context.Context, req *api.ListGatewayRequest) (*api.ListGatewayResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// to make sure the pagination of listing gateways work properly, need to take stc settings into consideration
	//   when calling GetGatewayCount and GetGateways
	filters := GatewayFilters{
		Search:         req.Search,
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
		OrganizationID: req.OrganizationId,
		EnabledSTC:     a.config.EnableSTC,
	}

	// Filter on username when OrganizationID is not set and the user is
	// not a global admin.
	if !cred.IsGlobalAdmin && filters.OrganizationID == 0 {
		filters.UserID = cred.UserID
	}

	resp := api.ListGatewayResponse{}
	resp.TotalCount, resp.Result, err = a.listGateways(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListNewGateways lists only new model gateways.
func (a *GatewayAPI) ListNewGateways(ctx context.Context, req *api.ListGatewayRequest) (*api.ListGatewayResponse, error) {
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// to make sure the pagination of listing gateways work properly, need to take stc settings into consideration
	//   when calling GetGatewayCount and GetGateways
	filters := GatewayFilters{
		Search:          req.Search,
		Limit:           int(req.Limit),
		Offset:          int(req.Offset),
		OrganizationID:  req.OrganizationId,
		EnabledSTC:      a.config.EnableSTC,
		NewGatewayModel: true,
	}

	// Filter on username when OrganizationID is not set and the user is
	// not a global admin.
	if !cred.IsGlobalAdmin && filters.OrganizationID == 0 {
		filters.UserID = cred.UserID
	}

	resp := api.ListGatewayResponse{}
	resp.TotalCount, resp.Result, err = a.listGateways(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// we don't want to reveal precise locations of the gateways, so we round them
// to reveal the location just within the radius of 500m
func roundCoordinates(lat, lon float64) (float64, float64) {
	// a degree of latitude is approximately 110km
	newLat := math.Round(lat*110) / 110
	// the length of a degree of longitude depends on latitude
	// (we just want a ballpark figure here)
	lonLength := (math.Pi / 180) * 6378.137 * math.Cos(math.Pi*lat/180)
	newLon := math.Round(lon*lonLength) / lonLength
	// if the approximated location is too close to the original we move it
	// 500m up or down
	if math.Abs(lat-newLat)*110 < 0.03 && math.Abs(lon-newLon)*lonLength < 0.03 {
		if newLat < lat {
			newLat = math.Round(lat+0.5/110) / 110
		} else {
			newLat = math.Round(lat-0.5/110) / 110
		}
		if newLat > 90 {
			newLat = 90
		}
		if newLat < -90 {
			newLat = -90
		}
	}
	return newLat, newLon
}

// ListLocations lists the gateway locations.
func (a *GatewayAPI) ListLocations(ctx context.Context, req *api.ListGatewayLocationsRequest) (*api.ListGatewayLocationsResponse, error) {
	var result []*api.GatewayLocationListItem

	gwsLoc, err := a.st.GetGatewaysLoc(ctx)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	for _, loc := range gwsLoc {
		latitude, longitude := roundCoordinates(loc.Latitude, loc.Longitude)
		result = append(result, &api.GatewayLocationListItem{
			Location: &api.GatewayLocation{
				Latitude:  latitude,
				Longitude: longitude,
				Altitude:  loc.Altitude,
			},
		})
	}

	resp := api.ListGatewayLocationsResponse{
		Result: result,
	}

	return &resp, nil
}

func (a *GatewayAPI) ensureGatewayAdmin(ctx context.Context, mac lorawan.EUI64) error {
	gw, err := a.st.GetGateway(ctx, mac, false)
	if err != nil {
		if err == errHandler.ErrDoesNotExist {
			return status.Errorf(codes.NotFound, "Gateway does not exist: %s", mac.String())
		}
		return status.Errorf(codes.Internal, "%v", err)
	}

	// check user permission, only gateway admin allowed
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(gw.OrganizationID))
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	if !cred.IsGatewayAdmin {
		return status.Errorf(codes.PermissionDenied, "permission denied")
	}
	return nil
}

// Update updates the given gateway.
func (a *GatewayAPI) Update(ctx context.Context, req *api.UpdateGatewayRequest) (*empty.Empty, error) {
	var response empty.Empty
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

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, handler *pgstore.PgStore) error {
		gw, err := handler.GetGateway(ctx, mac, true)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		gw.Name = req.Gateway.Name
		gw.Description = req.Gateway.Description
		gw.Ping = req.Gateway.DiscoveryEnabled
		gw.Latitude = req.Gateway.Location.Latitude
		gw.Longitude = req.Gateway.Location.Longitude
		gw.Altitude = req.Gateway.Location.Altitude
		//gw.Tags = tags

		err = handler.UpdateGateway(ctx, &gw)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
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

		n, err := handler.GetNetworkServer(ctx, gw.NetworkServerID)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		nStruct := &nsmod.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}
		client, err := nStruct.GetNetworkServiceClient()
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		_, err = client.UpdateGateway(ctx, &updateReq)
		if err != nil {
			return status.Errorf(codes.Unknown, "%v", err)
		}

		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	return &response, nil
}

// Delete deletes the gateway matching the given ID.
func (a *GatewayAPI) Delete(ctx context.Context, req *api.DeleteGatewayRequest) (*empty.Empty, error) {
	var response empty.Empty
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.Id)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	if err := a.st.Tx(ctx, func(ctx context.Context, store *pgstore.PgStore) error {
		if err := gw.DeleteGateway(ctx, mac, store, a.pscli); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &response, nil
}

// GetStats gets the gateway statistics for the gateway with the given Mac.
func (a *GatewayAPI) GetStats(ctx context.Context, req *api.GetGatewayStatsRequest) (*api.GetGatewayStatsResponse, error) {
	var gatewayID lorawan.EUI64
	if err := gatewayID.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}
	if err := a.ensureGatewayAdmin(ctx, gatewayID); err != nil {
		return nil, err
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

	metrics, err := metricsmod.GetMetrics(ctx, metricsmod.AggregationInterval(strings.ToUpper(req.Interval)), "gw:"+gatewayID.String(), start, end)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	result := make([]*api.GatewayStats, len(metrics))
	for i, m := range metrics {
		result[i] = &api.GatewayStats{
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

	return &api.GetGatewayStatsResponse{
		Result: result,
	}, nil
}

// GetLastPing returns the last emitted ping and gateways receiving this ping.
func (a *GatewayAPI) GetLastPing(ctx context.Context, req *api.GetLastPingRequest) (*api.GetLastPingResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	ping, pingRX, err := a.st.GetLastGatewayPingAndRX(ctx, mac)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := api.GetLastPingResponse{
		Frequency: uint32(ping.Frequency),
		Dr:        uint32(ping.DR),
	}

	resp.CreatedAt, err = ptypes.TimestampProto(ping.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	for _, rx := range pingRX {
		resp.PingRx = append(resp.PingRx, &api.PingRX{
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
func (a *GatewayAPI) StreamFrameLogs(req *api.StreamGatewayFrameLogsRequest, srv api.GatewayService_StreamFrameLogsServer) error {
	var mac lorawan.EUI64

	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return status.Errorf(codes.InvalidArgument, "mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(srv.Context(), mac); err != nil {
		return err
	}

	n, err := a.st.GetNetworkServerForGatewayMAC(srv.Context(), mac)
	if err != nil {
		return helpers.ErrToRPCError(err)
	}

	nStruct := &nsmod.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}
	client, err := nStruct.GetNetworkServiceClient()
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

		up, down, err := ConvertUplinkAndDownlinkFrames(resp.GetUplinkFrameSet(), resp.GetDownlinkFrame(), false)
		if err != nil {
			return helpers.ErrToRPCError(err)
		}

		var frameResp api.StreamGatewayFrameLogsResponse
		if up != nil {
			frameResp.Frame = &api.StreamGatewayFrameLogsResponse_UplinkFrame{
				UplinkFrame: up,
			}
		}

		if down != nil {
			frameResp.Frame = &api.StreamGatewayFrameLogsResponse_DownlinkFrame{
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
func (a *GatewayAPI) GetGwConfig(ctx context.Context, req *api.GetGwConfigRequest) (*api.GetGwConfigResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	gwConfig, err := a.st.GetGatewayConfigByGwID(ctx, mac)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "GetGwConfig/unable to get gateway config from DB %s", err)
	}

	return &api.GetGwConfigResponse{Conf: gwConfig}, nil
}

// UpdateGwConfig gateway configuration file
func (a *GatewayAPI) UpdateGwConfig(ctx context.Context, req *api.UpdateGwConfigRequest) (*api.UpdateGwConfigResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	if err := a.st.UpdateGatewayConfigByGwID(ctx, req.Conf, mac); err != nil {
		log.WithError(err).Error("Update conf to gw failed")
		return &api.UpdateGwConfigResponse{Status: "Update config failed, please check your gateway connection."},
			status.Errorf(codes.Internal, "cannot update gateway config: %s", err)
	}

	return &api.UpdateGwConfigResponse{
		Status: "Update gateway config file successful",
	}, nil
}

// Register will first try to get the gateway from provision server
func (a *GatewayAPI) Register(ctx context.Context, req *api.RegisterRequest) (*api.RegisterResponse, error) {
	log.WithFields(log.Fields{
		"Sn":             req.Sn,
		"OrganizationID": req.OrganizationId,
	}).Info("API Register is called")

	if req.Sn == "" {
		return nil, status.Error(codes.InvalidArgument, "gateway sn number must not be empty string")
	}

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if !cred.IsGlobalAdmin && !cred.IsOrgAdmin && !cred.IsGatewayAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	canHaveGateways, err := a.st.OrganizationCanHaveGateways(ctx, req.OrganizationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if !canHaveGateways {
		return nil, status.Error(codes.PermissionDenied, "organization cannot have gateways")
	}

	manufacturerNr := mannr.Serial2Manufacturer(req.Sn)
	stcOrgID, err := a.st.GetSTCOrgIDForGateway(ctx, manufacturerNr)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if stcOrgID == req.OrganizationId {
		stcOrgID = 0

		err := a.st.RemoveGatewayReseller(ctx, manufacturerNr)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	// check whether gateway is allowed to be registered on current supernode
	res, err := a.pscli.PreRegisterGW(ctx, &psPb.PreRegisterGWRequest{Sn: req.Sn})
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "get gateway info from ps by sn failed: %v", err)
	}
	if err = a.st.GatewayModelIsSupported(ctx, res.Model); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "gateway model %s is not supported on supernode: %s", res.Model, err.Error())
	}

	// register gateway with current supernode on remote provisioning server
	provReq := psPb.RegisterGWRequest{
		Sn:            req.Sn,
		SuperNodeAddr: a.config.ServerAddr,
		OrgId:         req.OrganizationId,
	}

	resp, err := a.pscli.RegisterGW(ctx, &provReq)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	// add new firmware if new model is registered
	_, err = a.st.GetGatewayFirmware(ctx, resp.Model, true)
	if err == errHandler.ErrDoesNotExist {
		if _, err = a.st.AddGatewayFirmware(ctx, &GatewayFirmware{
			Model: resp.Model,
		}); err != nil {
			log.WithError(err).Errorf("Failed to add new firmware: %s", resp.Model)
		}
	}

	gateway := api.Gateway{
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
		Boards:           []*api.GatewayBoard{},
	}

	// get gateway profile id, always use the default one
	count, err := a.st.GetGatewayProfileCount(ctx)
	if err != nil && err != errHandler.ErrDoesNotExist {
		return nil, status.Error(codes.Internal, err.Error())
	} else if err == errHandler.ErrDoesNotExist {
		return nil, status.Error(codes.NotFound, "")
	}

	gpList, err := a.st.GetGatewayProfiles(ctx, count, 0)
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

	nServers, err := a.st.GetNetworkServerForGatewayProfileID(ctx, gpID)
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

	if stcOrgID != 0 && a.config.EnableSTC {
		err := a.st.BindResellerToGateway(ctx, stcOrgID, resp.Sn)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	return &api.RegisterResponse{
		Status: "Successful",
	}, nil
}

func (a *GatewayAPI) GetGwPwd(ctx context.Context, req *api.GetGwPwdRequest) (*api.GetGwPwdResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	resp, err := a.pscli.GetRootPWD(context.Background(), &psPb.GetRootPWDRequest{
		Sn:  req.Sn,
		Mac: req.GatewayId,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.GetGwPwdResponse{Password: resp.RootPassword}, nil
}

func (a *GatewayAPI) SetAutoUpdateFirmware(ctx context.Context, req *api.SetAutoUpdateFirmwareRequest) (*api.SetAutoUpdateFirmwareResponse, error) {
	var mac lorawan.EUI64
	if err := mac.UnmarshalText([]byte(req.GatewayId)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad gateway mac: %s", err)
	}

	if err := a.ensureGatewayAdmin(ctx, mac); err != nil {
		return nil, err
	}

	if err := a.st.SetAutoUpdateFirmware(ctx, mac, req.AutoUpdate); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &api.SetAutoUpdateFirmwareResponse{Message: "Auto update firmware set successfully"}, nil
}

// GetGatewayList defines the get Gateway list request and response, data is obtained from mxpeotocol-server
func (a *GatewayAPI) GetGatewayList(ctx context.Context, req *api.GetGatewayListRequest) (*api.GetGatewayListResponse, error) {
	logInfo := "api/appserver_serves_ui/GetGatewayList org=" + strconv.FormatInt(req.OrgId, 10)

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err.Error())
	}
	if !cred.IsGlobalAdmin && !cred.IsOrgUser {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	resp, err := gwClient.GetGatewayList(ctx, &pb.GetGatewayListRequest{
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

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	resp, err := gwClient.GetGatewayProfile(ctx, &pb.GetGSGatewayProfileRequest{
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

	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions().WithOrgID(req.OrgId))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsOrgAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	resp, err := gwClient.GetGatewayHistory(ctx, &pb.GetGatewayHistoryRequest{
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
	cred, err := a.auth.GetCredentials(ctx, auth.NewOptions())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}
	if !cred.IsGlobalAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	gwClient := mxpcli.Global.GetM2MGatewayServiceClient()

	resp, err := gwClient.SetGatewayMode(ctx, &pb.SetGatewayModeRequest{
		OrgId:  req.OrgId,
		GwId:   req.GwId,
		GwMode: pb.GatewayMode(req.GwMode),
	})
	if err != nil {
		log.WithError(err).Error(logInfo)
		return &api.SetGatewayModeResponse{}, status.Errorf(codes.Unavailable, err.Error())
	}

	return &api.SetGatewayModeResponse{
		Status: resp.Status,
	}, status.Error(codes.OK, "")
}
