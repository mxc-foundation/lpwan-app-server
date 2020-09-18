package devprofile

import (
	"database/sql"
	"time"

	"google.golang.org/grpc/status"

	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lib/pq/hstore"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/ns"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
)

// DeviceProfileServiceAPI exports the ServiceProfile related functions.
type DeviceProfileServiceAPI struct {
	st *store.Handler
}

// NewDeviceProfileServiceAPI creates a new DeviceProfileServiceAPI.
func NewDeviceProfileServiceAPI() *DeviceProfileServiceAPI {
	return &DeviceProfileServiceAPI{
		st: ctrl.st,
	}
}

// Create creates the given device-profile.
func (a *DeviceProfileServiceAPI) Create(ctx context.Context, req *pb.CreateDeviceProfileRequest) (*pb.CreateDeviceProfileResponse, error) {
	if req.DeviceProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "deviceProfile expected")
	}

	if valid, err := NewValidator().ValidateDeviceProfilesAccess(ctx, authcus.Create, req.DeviceProfile.OrganizationId, 0); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var err error
	var uplinkInterval time.Duration
	if req.DeviceProfile.UplinkInterval != nil {
		uplinkInterval, err = ptypes.Duration(req.DeviceProfile.UplinkInterval)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
	}

	dp := store.DeviceProfile{
		OrganizationID:       req.DeviceProfile.OrganizationId,
		NetworkServerID:      req.DeviceProfile.NetworkServerId,
		Name:                 req.DeviceProfile.Name,
		PayloadCodec:         req.DeviceProfile.PayloadCodec,
		PayloadEncoderScript: req.DeviceProfile.PayloadEncoderScript,
		PayloadDecoderScript: req.DeviceProfile.PayloadDecoderScript,
		Tags: hstore.Hstore{
			Map: make(map[string]sql.NullString),
		},
		UplinkInterval: uplinkInterval,
		DeviceProfile: ns.DeviceProfile{
			SupportsClassB:     req.DeviceProfile.SupportsClassB,
			ClassBTimeout:      req.DeviceProfile.ClassBTimeout,
			PingSlotPeriod:     req.DeviceProfile.PingSlotPeriod,
			PingSlotDr:         req.DeviceProfile.PingSlotDr,
			PingSlotFreq:       req.DeviceProfile.PingSlotFreq,
			SupportsClassC:     req.DeviceProfile.SupportsClassC,
			ClassCTimeout:      req.DeviceProfile.ClassCTimeout,
			MacVersion:         req.DeviceProfile.MacVersion,
			RegParamsRevision:  req.DeviceProfile.RegParamsRevision,
			RxDelay_1:          req.DeviceProfile.RxDelay_1,
			RxDrOffset_1:       req.DeviceProfile.RxDrOffset_1,
			RxDatarate_2:       req.DeviceProfile.RxDatarate_2,
			RxFreq_2:           req.DeviceProfile.RxFreq_2,
			MaxEirp:            req.DeviceProfile.MaxEirp,
			MaxDutyCycle:       req.DeviceProfile.MaxDutyCycle,
			SupportsJoin:       req.DeviceProfile.SupportsJoin,
			RfRegion:           req.DeviceProfile.RfRegion,
			Supports_32BitFCnt: req.DeviceProfile.Supports_32BitFCnt,
			FactoryPresetFreqs: req.DeviceProfile.FactoryPresetFreqs,
		},
	}

	for k, v := range req.DeviceProfile.Tags {
		dp.Tags.Map[k] = sql.NullString{Valid: true, String: v}
	}

	// as this also performs a remote call to create the device-profile
	// on the network-server, wrap it in a transaction
	err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		return handler.CreateDeviceProfile(ctx, &dp)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	dpID, err := uuid.FromBytes(dp.DeviceProfile.Id)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &pb.CreateDeviceProfileResponse{
		Id: dpID.String(),
	}, nil
}

// Get returns the device-profile matching the given id.
func (a *DeviceProfileServiceAPI) Get(ctx context.Context, req *pb.GetDeviceProfileRequest) (*pb.GetDeviceProfileResponse, error) {
	dpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceProfileAccess(ctx, authcus.Read, dpID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	dp, err := a.st.GetDeviceProfile(ctx, dpID, false)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	n, err := a.st.GetNetworkServer(ctx, dp.NetworkServerID)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	nstruct := nscli.NSStruct{
		Server:  n.Server,
		CACert:  n.CACert,
		TLSCert: n.TLSCert,
		TLSKey:  n.TLSKey,
	}

	nsClient, err := nstruct.GetNetworkServiceClient()
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	res, err := nsClient.GetDeviceProfile(ctx, &ns.GetDeviceProfileRequest{
		Id: dpID.Bytes(),
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	if res.DeviceProfile == nil {
		return nil, helpers.ErrToRPCError(err)
	}

	dp.DeviceProfile = *res.DeviceProfile

	resp := pb.GetDeviceProfileResponse{
		DeviceProfile: &pb.DeviceProfile{
			Id:                   dpID.String(),
			Name:                 dp.Name,
			OrganizationId:       dp.OrganizationID,
			NetworkServerId:      dp.NetworkServerID,
			PayloadCodec:         string(dp.PayloadCodec),
			PayloadEncoderScript: dp.PayloadEncoderScript,
			PayloadDecoderScript: dp.PayloadDecoderScript,
			SupportsClassB:       dp.DeviceProfile.SupportsClassB,
			ClassBTimeout:        dp.DeviceProfile.ClassBTimeout,
			PingSlotPeriod:       dp.DeviceProfile.PingSlotPeriod,
			PingSlotDr:           dp.DeviceProfile.PingSlotDr,
			PingSlotFreq:         dp.DeviceProfile.PingSlotFreq,
			SupportsClassC:       dp.DeviceProfile.SupportsClassC,
			ClassCTimeout:        dp.DeviceProfile.ClassCTimeout,
			MacVersion:           dp.DeviceProfile.MacVersion,
			RegParamsRevision:    dp.DeviceProfile.RegParamsRevision,
			RxDelay_1:            dp.DeviceProfile.RxDelay_1,
			RxDrOffset_1:         dp.DeviceProfile.RxDrOffset_1,
			RxDatarate_2:         dp.DeviceProfile.RxDatarate_2,
			RxFreq_2:             dp.DeviceProfile.RxFreq_2,
			MaxEirp:              dp.DeviceProfile.MaxEirp,
			MaxDutyCycle:         dp.DeviceProfile.MaxDutyCycle,
			SupportsJoin:         dp.DeviceProfile.SupportsJoin,
			RfRegion:             dp.DeviceProfile.RfRegion,
			Supports_32BitFCnt:   dp.DeviceProfile.Supports_32BitFCnt,
			FactoryPresetFreqs:   dp.DeviceProfile.FactoryPresetFreqs,
			Tags:                 make(map[string]string),
			UplinkInterval:       ptypes.DurationProto(dp.UplinkInterval),
		},
	}

	resp.CreatedAt, err = ptypes.TimestampProto(dp.CreatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}
	resp.UpdatedAt, err = ptypes.TimestampProto(dp.UpdatedAt)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	for k, v := range dp.Tags.Map {
		resp.DeviceProfile.Tags[k] = v.String
	}

	return &resp, nil
}

// Update updates the given device-profile.
func (a *DeviceProfileServiceAPI) Update(ctx context.Context, req *pb.UpdateDeviceProfileRequest) (*empty.Empty, error) {
	if req.DeviceProfile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "deviceProfile expected")
	}

	dpID, err := uuid.FromString(req.DeviceProfile.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceProfileAccess(ctx, authcus.Update, dpID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// As this also performs a remote call to update the device-profile
	// on the network-server, wrap it in a transaction.
	// This also locks the local device-profile record in the database.
	err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		dp, err := handler.GetDeviceProfile(ctx, dpID, true)
		if err != nil {
			return err
		}

		n, err := a.st.GetNetworkServer(ctx, dp.NetworkServerID)
		if err != nil {
			return err
		}

		nstruct := nscli.NSStruct{
			Server:  n.Server,
			CACert:  n.CACert,
			TLSCert: n.TLSCert,
			TLSKey:  n.TLSKey,
		}

		nsClient, err := nstruct.GetNetworkServiceClient()
		if err != nil {
			return err
		}

		res, err := nsClient.GetDeviceProfile(ctx, &ns.GetDeviceProfileRequest{
			Id: dpID.Bytes(),
		})
		if err != nil {
			return err
		}
		if res.DeviceProfile == nil {
			return err
		}

		dp.DeviceProfile = *res.DeviceProfile

		var uplinkInterval time.Duration
		if req.DeviceProfile.UplinkInterval != nil {
			uplinkInterval, err = ptypes.Duration(req.DeviceProfile.UplinkInterval)
			if err != nil {
				return err
			}
		}

		dp.Name = req.DeviceProfile.Name
		dp.PayloadCodec = req.DeviceProfile.PayloadCodec
		dp.PayloadEncoderScript = req.DeviceProfile.PayloadEncoderScript
		dp.PayloadDecoderScript = req.DeviceProfile.PayloadDecoderScript
		dp.Tags = hstore.Hstore{
			Map: make(map[string]sql.NullString),
		}
		dp.UplinkInterval = uplinkInterval
		dp.DeviceProfile = ns.DeviceProfile{
			Id:                 dpID.Bytes(),
			SupportsClassB:     req.DeviceProfile.SupportsClassB,
			ClassBTimeout:      req.DeviceProfile.ClassBTimeout,
			PingSlotPeriod:     req.DeviceProfile.PingSlotPeriod,
			PingSlotDr:         req.DeviceProfile.PingSlotDr,
			PingSlotFreq:       req.DeviceProfile.PingSlotFreq,
			SupportsClassC:     req.DeviceProfile.SupportsClassC,
			ClassCTimeout:      req.DeviceProfile.ClassCTimeout,
			MacVersion:         req.DeviceProfile.MacVersion,
			RegParamsRevision:  req.DeviceProfile.RegParamsRevision,
			RxDelay_1:          req.DeviceProfile.RxDelay_1,
			RxDrOffset_1:       req.DeviceProfile.RxDrOffset_1,
			RxDatarate_2:       req.DeviceProfile.RxDatarate_2,
			RxFreq_2:           req.DeviceProfile.RxFreq_2,
			MaxEirp:            req.DeviceProfile.MaxEirp,
			MaxDutyCycle:       req.DeviceProfile.MaxDutyCycle,
			SupportsJoin:       req.DeviceProfile.SupportsJoin,
			RfRegion:           req.DeviceProfile.RfRegion,
			Supports_32BitFCnt: req.DeviceProfile.Supports_32BitFCnt,
			FactoryPresetFreqs: req.DeviceProfile.FactoryPresetFreqs,
		}

		for k, v := range req.DeviceProfile.Tags {
			dp.Tags.Map[k] = sql.NullString{Valid: true, String: v}
		}

		return handler.UpdateDeviceProfile(ctx, &dp)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// Delete deletes the device-profile matching the given id.
func (a *DeviceProfileServiceAPI) Delete(ctx context.Context, req *pb.DeleteDeviceProfileRequest) (*empty.Empty, error) {
	dpID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "uuid error: %s", err)
	}

	if valid, err := NewValidator().ValidateDeviceProfileAccess(ctx, authcus.Delete, dpID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}
	// as this also performs a remote call to delete the device-profile
	// on the network-server, wrap it in a transaction
	err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		return handler.DeleteDeviceProfile(ctx, dpID)
	})
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	return &empty.Empty{}, nil
}

// List lists the available device-profiles.
func (a *DeviceProfileServiceAPI) List(ctx context.Context, req *pb.ListDeviceProfileRequest) (*pb.ListDeviceProfileResponse, error) {
	if req.ApplicationId != 0 {
		if valid, err := NewValidator().ValidateDeviceProfilesAccess(ctx, authcus.List, 0, req.ApplicationId); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	} else {
		if valid, err := NewValidator().ValidateDeviceProfilesAccess(ctx, authcus.List, req.OrganizationId, 0); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	filters := store.DeviceProfileFilters{
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
		OrganizationID: req.OrganizationId,
		ApplicationID:  req.ApplicationId,
	}

	user, err := NewValidator().GetUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	// Filter on user ID when org and app ID are not set and user is not
	// global admin.
	if !user.IsGlobalAdmin && filters.OrganizationID == 0 && filters.ApplicationID == 0 {
		filters.UserID = user.ID
	}

	count, err := a.st.GetDeviceProfileCount(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	dps, err := a.st.GetDeviceProfiles(ctx, filters)
	if err != nil {
		return nil, helpers.ErrToRPCError(err)
	}

	resp := pb.ListDeviceProfileResponse{
		TotalCount: int64(count),
	}

	for _, dp := range dps {
		row := pb.DeviceProfileListItem{
			Id:                dp.DeviceProfileID.String(),
			Name:              dp.Name,
			OrganizationId:    dp.OrganizationID,
			NetworkServerId:   dp.NetworkServerID,
			NetworkServerName: dp.NetworkServerName,
		}

		row.CreatedAt, err = ptypes.TimestampProto(dp.CreatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}
		row.UpdatedAt, err = ptypes.TimestampProto(dp.UpdatedAt)
		if err != nil {
			return nil, helpers.ErrToRPCError(err)
		}

		resp.Result = append(resp.Result, &row)
	}

	return &resp, nil
}
