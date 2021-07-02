package external

import (
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"
	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	serviceprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// MulticastGroupAPI implements the multicast-group api.
type MulticastGroupAPI struct {
	st               *store.Handler
	routingProfileID uuid.UUID
	nsCli            *nscli.Client
}

// NewMulticastGroupAPI creates a new multicast-group API.
func NewMulticastGroupAPI(routingProfileID uuid.UUID, h *store.Handler, nsCli *nscli.Client) *MulticastGroupAPI {
	return &MulticastGroupAPI{
		st:               h,
		routingProfileID: routingProfileID,
		nsCli:            nsCli,
	}
}

// Create creates the given multicast-group.
func (a *MulticastGroupAPI) Create(ctx context.Context, req *pb.CreateMulticastGroupRequest) (*pb.CreateMulticastGroupResponse, error) {
	if req.MulticastGroup == nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group must not be nil")
	}

	spID, err := uuid.FromString(req.MulticastGroup.ServiceProfileId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	sp, err := serviceprofile.GetServiceProfile(ctx, a.st, spID, a.nsCli, true) // local-only, as we only want to fetch the org. id
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupsAccess(ctx, auth.Create, sp.OrganizationID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	var mcAddr lorawan.DevAddr
	if err = mcAddr.UnmarshalText([]byte(req.MulticastGroup.McAddr)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_app_s_key: %s", err)
	}

	var mcNwkSKey lorawan.AES128Key
	if err = mcNwkSKey.UnmarshalText([]byte(req.MulticastGroup.McNwkSKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_net_s_key: %s", err)
	}

	mg := MulticastGroup{
		Name:             req.MulticastGroup.Name,
		ServiceProfileID: spID,
		MulticastGroup: ns.MulticastGroup{
			McAddr:           mcAddr[:],
			McNwkSKey:        mcNwkSKey[:],
			GroupType:        ns.MulticastGroupType(req.MulticastGroup.GroupType),
			Dr:               req.MulticastGroup.Dr,
			Frequency:        req.MulticastGroup.Frequency,
			PingSlotPeriod:   req.MulticastGroup.PingSlotPeriod,
			ServiceProfileId: spID.Bytes(),
			RoutingProfileId: a.routingProfileID.Bytes(),
			FCnt:             req.MulticastGroup.FCnt,
		},
	}

	if err = mg.MCAppSKey.UnmarshalText([]byte(req.MulticastGroup.McAppSKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_app_s_key: %s", err)
	}

	if err := multicast.CreateMulticastGroup(ctx, &mg); err != nil {
		return nil, err
	}

	var mgID uuid.UUID
	copy(mgID[:], mg.MulticastGroup.Id)

	return &pb.CreateMulticastGroupResponse{
		Id: mgID.String(),
	}, nil
}

// Get returns a multicast-group given an ID.
func (a *MulticastGroupAPI) Get(ctx context.Context, req *pb.GetMulticastGroupRequest) (*pb.GetMulticastGroupResponse, error) {
	mgID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupAccess(ctx, auth.Read, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	mg, err := multicast.GetMulticastGroup(ctx, mgID, false, false)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	var mcAddr lorawan.DevAddr
	var mcNwkSKey lorawan.AES128Key
	copy(mcAddr[:], mg.MulticastGroup.McAddr)
	copy(mcNwkSKey[:], mg.MulticastGroup.McNwkSKey)

	out := pb.GetMulticastGroupResponse{
		MulticastGroup: &pb.MulticastGroup{
			Id:               mgID.String(),
			Name:             mg.Name,
			McAddr:           mcAddr.String(),
			McNwkSKey:        mcNwkSKey.String(),
			McAppSKey:        mg.MCAppSKey.String(),
			FCnt:             mg.MulticastGroup.FCnt,
			GroupType:        pb.MulticastGroupType(mg.MulticastGroup.GroupType),
			Dr:               mg.MulticastGroup.Dr,
			Frequency:        mg.MulticastGroup.Frequency,
			PingSlotPeriod:   mg.MulticastGroup.PingSlotPeriod,
			ServiceProfileId: mg.ServiceProfileID.String(),
		},
	}

	out.CreatedAt = timestamppb.New(mg.CreatedAt)
	out.UpdatedAt = timestamppb.New(mg.UpdatedAt)

	return &out, nil
}

// Update updates the given multicast-group.
func (a *MulticastGroupAPI) Update(ctx context.Context, req *pb.UpdateMulticastGroupRequest) (*empty.Empty, error) {
	if req.MulticastGroup == nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group must not be nil")
	}

	mgID, err := uuid.FromString(req.MulticastGroup.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupAccess(ctx, auth.Update, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	mg, err := multicast.GetMulticastGroup(ctx, mgID, false, false)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	var mcAddr lorawan.DevAddr
	if err = mcAddr.UnmarshalText([]byte(req.MulticastGroup.McAddr)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_app_s_key: %s", err)
	}

	var mcNwkSKey lorawan.AES128Key
	if err = mcNwkSKey.UnmarshalText([]byte(req.MulticastGroup.McNwkSKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_net_s_key: %s", err)
	}

	mg.Name = req.MulticastGroup.Name
	mg.MulticastGroup = ns.MulticastGroup{
		Id:               mg.MulticastGroup.Id,
		McAddr:           mcAddr[:],
		McNwkSKey:        mcNwkSKey[:],
		GroupType:        ns.MulticastGroupType(req.MulticastGroup.GroupType),
		Dr:               req.MulticastGroup.Dr,
		Frequency:        req.MulticastGroup.Frequency,
		PingSlotPeriod:   req.MulticastGroup.PingSlotPeriod,
		ServiceProfileId: mg.MulticastGroup.ServiceProfileId,
		RoutingProfileId: mg.MulticastGroup.RoutingProfileId,
		FCnt:             req.MulticastGroup.FCnt,
	}

	if err = mg.MCAppSKey.UnmarshalText([]byte(req.MulticastGroup.McAppSKey)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "mc_app_s_key: %s", err)
	}

	if err = multicast.UpdateMulticastGroup(ctx, &mg); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Delete deletes a multicast-group given an ID.
func (a *MulticastGroupAPI) Delete(ctx context.Context, req *pb.DeleteMulticastGroupRequest) (*empty.Empty, error) {
	mgID, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	}

	if err = multicast.DeleteMulticastGroup(ctx, mgID); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// List lists the available multicast-groups.
func (a *MulticastGroupAPI) List(ctx context.Context, req *pb.ListMulticastGroupRequest) (*pb.ListMulticastGroupResponse, error) {
	var err error
	var idFilter bool

	filters := MulticastGroupFilters{
		OrganizationID: req.OrganizationId,
		Search:         req.Search,
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
	}

	// if org. filter has been set, validate the client has access to the given org
	if filters.OrganizationID != 0 {
		idFilter = true

		if valid, err := organization.NewValidator().ValidateOrganizationAccess(ctx, auth.Read, req.OrganizationId); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
		}
	}

	// if sp filter has been set, validate the client has access to the given sp
	if req.ServiceProfileId != "" {
		idFilter = true

		filters.ServiceProfileID, err = uuid.FromString(req.ServiceProfileId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "service_profile_id: %s", err)
		}

		if valid, err := serviceprofile.NewValidator(a.st).ValidateServiceProfileAccess(ctx, auth.Read, filters.ServiceProfileID); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication error: %s", err)
		}
	}

	// if devEUI has been set, validate the client has access to the given device
	if req.DevEui != "" {
		idFilter = true

		if err = filters.DevEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "dev_eui: %s", err)
		}

		if valid, err := multicast.NewValidator().ValidateNodeAccess(ctx, auth.Read, filters.DevEUI); !valid || err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "authentication error: %s", err)
		}
	}

	// listing all stored objects is for global admin only
	if !idFilter {
		user, err := multicast.NewValidator().GetUser(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "%s", err)
		}

		if !user.IsGlobalAdmin {
			return nil, status.Errorf(codes.Unauthenticated, "client must be global admin for unfiltered request")
		}
	}

	count, err := a.st.GetMulticastGroupCount(ctx, filters)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	items, err := a.st.GetMulticastGroups(ctx, filters)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	out := pb.ListMulticastGroupResponse{
		TotalCount: int64(count),
	}

	for _, item := range items {
		out.Result = append(out.Result, &pb.MulticastGroupListItem{
			Id:                 item.ID.String(),
			Name:               item.Name,
			ServiceProfileId:   item.ServiceProfileID.String(),
			ServiceProfileName: item.ServiceProfileName,
		})
	}

	return &out, nil
}

// AddDevice adds the given device to the multicast-group.
func (a *MulticastGroupAPI) AddDevice(ctx context.Context, req *pb.AddDeviceToMulticastGroupRequest) (*empty.Empty, error) {
	mgID, err := uuid.FromString(req.MulticastGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group_id: %s", err)
	}

	var devEUI lorawan.EUI64
	if err = devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "dev_eui: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupAccess(ctx, auth.Update, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	// validate that the device is under the same service-profile as the multicast-group
	dev, err := a.st.GetDevice(ctx, devEUI, false)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	app, err := a.st.GetApplication(ctx, dev.ApplicationID)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	mg, err := multicast.GetMulticastGroup(ctx, mgID, false, true)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	if app.ServiceProfileID != mg.ServiceProfileID {
		return nil, status.Errorf(codes.FailedPrecondition, "service-profile of device != service-profile of multicast-group")
	}

	if err = multicast.AddDeviceToMulticastGroup(ctx, mgID, devEUI); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// RemoveDevice removes the given device from the multicast-group.
func (a *MulticastGroupAPI) RemoveDevice(ctx context.Context, req *pb.RemoveDeviceFromMulticastGroupRequest) (*empty.Empty, error) {
	mgID, err := uuid.FromString(req.MulticastGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group_id: %s", err)
	}

	var devEUI lorawan.EUI64
	if err = devEUI.UnmarshalText([]byte(req.DevEui)); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "dev_eui: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupAccess(ctx, auth.Update, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err = multicast.RemoveDeviceFromMulticastGroup(ctx, mgID, devEUI); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// Enqueue adds the given item to the multicast-queue.
func (a *MulticastGroupAPI) Enqueue(ctx context.Context, req *pb.EnqueueMulticastQueueItemRequest) (*pb.EnqueueMulticastQueueItemResponse, error) {
	var fCnt uint32

	if req.MulticastQueueItem == nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_queue_item must not be nil")
	}

	if req.MulticastQueueItem.FPort == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "f_port must be > 0")
	}

	mgID, err := uuid.FromString(req.MulticastQueueItem.MulticastGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group_id: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupQueueAccess(ctx, auth.Create, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	if err = a.st.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		var err error
		fCnt, err = multicast.Enqueue(ctx, handler, mgID, uint8(req.MulticastQueueItem.FPort), req.MulticastQueueItem.Data)
		if err != nil {
			return status.Errorf(codes.Internal, "enqueue multicast-group queue-item error: %s", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.EnqueueMulticastQueueItemResponse{
		FCnt: fCnt,
	}, nil
}

// FlushQueue flushes the multicast-group queue.
func (a *MulticastGroupAPI) FlushQueue(ctx context.Context, req *pb.FlushMulticastGroupQueueItemsRequest) (*empty.Empty, error) {
	mgID, err := uuid.FromString(req.MulticastGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group_id: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupQueueAccess(ctx, auth.Delete, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	n, err := a.st.GetNetworkServerForMulticastGroupID(ctx, mgID)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	_, err = nsClient.FlushMulticastQueueForMulticastGroup(ctx, &ns.FlushMulticastQueueForMulticastGroupRequest{
		MulticastGroupId: mgID.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ListQueue lists the items in the multicast-group queue.
func (a *MulticastGroupAPI) ListQueue(ctx context.Context, req *pb.ListMulticastGroupQueueItemsRequest) (*pb.ListMulticastGroupQueueItemsResponse, error) {
	mgID, err := uuid.FromString(req.MulticastGroupId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "multicast_group_id: %s", err)
	}

	if valid, err := multicast.NewValidator().ValidateMulticastGroupQueueAccess(ctx, auth.Read, mgID); !valid || err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %s", err)
	}

	queueItems, err := multicast.ListQueue(ctx, a.st, mgID)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "%s", err)
	}

	var resp pb.ListMulticastGroupQueueItemsResponse
	for i := range queueItems {
		resp.MulticastQueueItems = append(resp.MulticastQueueItems, &queueItems[i])
	}

	return &resp, nil
}
