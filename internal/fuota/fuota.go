package fuota

import (
	"context"
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/applayer/fragmentation"
	"github.com/brocaar/lorawan/applayer/multicastsetup"

	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/device"
	fts "github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation/data"
	mcss "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	fds "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment/data"
	mgs "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"

	. "github.com/mxc-foundation/lpwan-app-server/internal/fuota/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

type Config struct {
	RemoteMulticastSetupRetries       int
	RemoteFragmentationSessionRetries int
	ApplicationServerID               string
}

type controller struct {
	s                FuotaStruct
	interval         time.Duration
	batchSize        int
	config           Config
	routingProfileID uuid.UUID
	nsCli            *nscli.Client
}

// Setup configures the package.
func Setup(name string, h *store.Handler, conf Config, s FuotaStruct, nsCli *nscli.Client) error {
	var err error
	ctrl := &controller{
		s:         s,
		config:    conf,
		interval:  time.Second,
		batchSize: 1,
		nsCli:     nsCli,
	}

	ctrl.routingProfileID, err = uuid.FromString(conf.ApplicationServerID)
	if err != nil {
		return errors.Wrap(err, "application-server id to uuid error")
	}

	go ctrl.fuotaDeploymentLoop(h)

	return nil
}

func (c *controller) fuotaDeploymentLoop(h *store.Handler) {
	for {
		ctxID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("new uuid error")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

		err = h.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
			return c.fuotaDeployments(ctx, handler)
		})
		if err != nil {
			log.WithError(err).Error("fuota deployment error")
		}
		time.Sleep(c.interval)
	}
}

func (c *controller) fuotaDeployments(ctx context.Context, handler *store.Handler) error {
	items, err := handler.GetPendingFUOTADeployments(ctx, c.batchSize)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := c.fuotaDeployment(ctx, handler, item); err != nil {
			return errors.Wrap(err, "fuota deployment error")
		}
	}

	return nil
}

func (c *controller) fuotaDeployment(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	switch item.State {
	case fds.FUOTADeploymentMulticastCreate:
		return c.stepMulticastCreate(ctx, handler, item)
	case fds.FUOTADeploymentMulticastSetup:
		return c.stepMulticastSetup(ctx, handler, item)
	case fds.FUOTADeploymentFragmentationSessSetup:
		return c.stepFragmentationSessSetup(ctx, handler, item)
	case fds.FUOTADeploymentMulticastSessCSetup:
		return c.stepMulticastSessCSetup(ctx, handler, item)
	case fds.FUOTADeploymentEnqueue:
		return c.stepEnqueue(ctx, handler, item)
	case fds.FUOTADeploymentStatusRequest:
		return c.stepStatusRequest(ctx, handler, item)
	case fds.FUOTADeploymentSetDeviceStatus:
		return c.stepSetDeviceStatus(ctx, handler, item)
	case fds.FUOTADeploymentCleanup:
		return c.stepCleanup(ctx, handler, item)
	default:
		return fmt.Errorf("unexpected state: %s", item.State)
	}
}

func (c *controller) stepMulticastCreate(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	var devAddr lorawan.DevAddr
	if _, err := rand.Read(devAddr[:]); err != nil {
		return errors.Wrap(err, "read random bytes error")
	}

	var mcKey lorawan.AES128Key
	if _, err := rand.Read(mcKey[:]); err != nil {
		return errors.Wrap(err, "read random bytes error")
	}

	mcAppSKey, err := multicastsetup.GetMcAppSKey(mcKey, devAddr)
	if err != nil {
		return errors.Wrap(err, "get McAppSKey error")
	}

	mcNetSKey, err := multicastsetup.GetMcNetSKey(mcKey, devAddr)
	if err != nil {
		return errors.Wrap(err, "get McNetSKey error")
	}

	spID, err := handler.GetServiceProfileIDForFUOTADeployment(ctx, item.ID)
	if err != nil {
		return errors.Wrap(err, "get service-profile for fuota deployment error")
	}

	mg := mgs.MulticastGroup{
		Name:             fmt.Sprintf("fuota-%s", item.ID),
		MCAppSKey:        mcAppSKey,
		MCKey:            mcKey,
		ServiceProfileID: spID,
		MulticastGroup: ns.MulticastGroup{
			McAddr:           devAddr[:],
			McNwkSKey:        mcNetSKey[:],
			FCnt:             0,
			Dr:               uint32(item.DR),
			Frequency:        uint32(item.Frequency),
			PingSlotPeriod:   uint32(item.PingSlotPeriod),
			ServiceProfileId: spID.Bytes(),
			RoutingProfileId: c.routingProfileID.Bytes(),
		},
	}

	switch item.GroupType {
	case fds.FUOTADeploymentGroupTypeB:
		mg.MulticastGroup.GroupType = ns.MulticastGroupType_CLASS_B
	case fds.FUOTADeploymentGroupTypeC:
		mg.MulticastGroup.GroupType = ns.MulticastGroupType_CLASS_C
	default:
		return fmt.Errorf("unknown group-type: %s", item.GroupType)
	}

	err = multicast.CreateMulticastGroup(ctx, &mg, c.nsCli)
	if err != nil {
		return errors.Wrap(err, "create multicast-group error")
	}

	var mgID uuid.UUID
	copy(mgID[:], mg.MulticastGroup.Id)

	item.MulticastGroupID = &mgID
	item.State = fds.FUOTADeploymentMulticastSetup
	item.NextStepAfter = time.Now()

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepMulticastSetup(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	mcg, err := multicast.GetMulticastGroup(ctx, *item.MulticastGroupID, false, false, c.nsCli)
	if err != nil {
		return errors.Wrap(err, "get multicast group error")
	}

	// query all device-keys that relate to this FUOTA deployment
	deviceKeys, err := handler.GetDeviceKeysFromFuotaDevelopmentDevice(ctx, item.ID)
	if err != nil {
		return errors.Wrap(err, "sql select error")
	}

	for _, dk := range deviceKeys {
		var nullKey lorawan.AES128Key

		// get the encrypted McKey.
		var mcKeyEncrypted, mcRootKey lorawan.AES128Key
		if dk.AppKey != nullKey {
			mcRootKey, err = multicastsetup.GetMcRootKeyForAppKey(dk.AppKey)
			if err != nil {
				return errors.Wrap(err, "get McRootKey for AppKey error")
			}
		} else {
			mcRootKey, err = multicastsetup.GetMcRootKeyForGenAppKey(dk.GenAppKey)
			if err != nil {
				return errors.Wrap(err, "get McRootKey for GenAppKey error")
			}
		}

		mcKEKey, err := multicastsetup.GetMcKEKey(mcRootKey)
		if err != nil {
			return errors.Wrap(err, "get McKEKey error")
		}

		block, err := aes.NewCipher(mcKEKey[:])
		if err != nil {
			return errors.Wrap(err, "new cipher error")
		}
		block.Decrypt(mcKeyEncrypted[:], mcg.MCKey[:])

		// create remote multicast setup record for device
		rms := mcss.RemoteMulticastSetup{
			DevEUI:           dk.DevEUI,
			MulticastGroupID: *item.MulticastGroupID,
			McGroupID:        c.s.McGroupID,
			McKeyEncrypted:   mcKeyEncrypted,
			MinMcFCnt:        0,
			MaxMcFCnt:        (1 << 32) - 1,
			State:            mcss.RemoteMulticastSetupSetup,
			RetryInterval:    item.UnicastTimeout,
		}
		copy(rms.McAddr[:], mcg.MulticastGroup.McAddr)

		err = handler.CreateRemoteMulticastSetup(ctx, &rms)
		if err != nil {
			return errors.Wrap(err, "create remote multicast setup error")
		}
	}

	item.State = fds.FUOTADeploymentFragmentationSessSetup
	item.NextStepAfter = time.Now().Add(time.Duration(c.config.RemoteMulticastSetupRetries) * item.UnicastTimeout)

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepFragmentationSessSetup(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	if item.FragSize == 0 {
		return errors.New("FragSize must not be 0")
	}

	devEUIs, err := handler.GetDevEUIsWithMulticastSetup(ctx, item.MulticastGroupID)
	if err != nil {
		return errors.Wrap(err, "get devices with multicast setup error")
	}

	padding := (item.FragSize - (len(item.Payload) % item.FragSize)) % item.FragSize
	nbFrag := (len(item.Payload) + padding) / item.FragSize

	for _, devEUI := range devEUIs {
		// delete existing fragmentation session if it exist
		err = handler.DeleteRemoteFragmentationSession(ctx, devEUI, c.s.FragIndex)
		if err != nil && err != errHandler.ErrDoesNotExist {
			return errors.Wrap(err, "delete remote fragmentation session error")
		}

		fs := fts.RemoteFragmentationSession{
			DevEUI:              devEUI,
			FragIndex:           c.s.FragIndex,
			MCGroupIDs:          []int{c.s.McGroupID},
			NbFrag:              nbFrag,
			FragSize:            item.FragSize,
			FragmentationMatrix: item.FragmentationMatrix,
			BlockAckDelay:       item.BlockAckDelay,
			Padding:             padding,
			Descriptor:          item.Descriptor,
			State:               mcss.RemoteMulticastSetupSetup,
			RetryInterval:       item.UnicastTimeout,
		}
		err = handler.CreateRemoteFragmentationSession(ctx, &fs)
		if err != nil {
			return errors.Wrap(err, "create remote fragmentation session error")
		}
	}

	item.State = fds.FUOTADeploymentMulticastSessCSetup
	item.NextStepAfter = time.Now().Add(time.Duration(c.config.RemoteFragmentationSessionRetries) * item.UnicastTimeout)

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepMulticastSessCSetup(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	mcg, err := multicast.GetMulticastGroup(ctx, *item.MulticastGroupID, false, false, c.nsCli)
	if err != nil {
		return errors.Wrap(err, "get multicast group error")
	}

	devEUIs, err := handler.GetDevEUIsWithFragmentationSessionSetup(ctx, item.MulticastGroupID, c.s.FragIndex)
	if err != nil {
		return errors.Wrap(err, "get devices with fragmentation session setup error")
	}

	for _, devEUI := range devEUIs {
		rmccs := mcss.RemoteMulticastClassCSession{
			DevEUI:           devEUI,
			MulticastGroupID: *item.MulticastGroupID,
			McGroupID:        c.s.McGroupID,
			DLFrequency:      int(mcg.MulticastGroup.Frequency),
			DR:               int(mcg.MulticastGroup.Dr),
			SessionTime:      time.Now().Add(time.Duration(c.config.RemoteMulticastSetupRetries) * item.UnicastTimeout),
			SessionTimeOut:   item.MulticastTimeout,
			RetryInterval:    item.UnicastTimeout,
		}
		err = handler.CreateRemoteMulticastClassCSession(ctx, &rmccs)
		if err != nil {
			return errors.Wrap(err, "create remote multicast class-c session error")
		}
	}

	item.State = fds.FUOTADeploymentEnqueue
	item.NextStepAfter = time.Now().Add(time.Duration(c.config.RemoteMulticastSetupRetries) * item.UnicastTimeout)

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepEnqueue(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	// fragment the payload
	padding := (item.FragSize - (len(item.Payload) % item.FragSize)) % item.FragSize
	fragments, err := fragmentation.Encode(append(item.Payload, make([]byte, padding)...), item.FragSize, item.Redundancy)
	if err != nil {
		return errors.Wrap(err, "fragment payload error")
	}

	// wrap the payloads into data-fragment payloads
	var payloads [][]byte
	for i := range fragments {
		cmd := fragmentation.Command{
			CID: fragmentation.DataFragment,
			Payload: &fragmentation.DataFragmentPayload{
				IndexAndN: fragmentation.DataFragmentPayloadIndexAndN{
					FragIndex: uint8(c.s.FragIndex),
					N:         uint16(i + 1),
				},
				Payload: fragments[i],
			},
		}
		b, err := cmd.MarshalBinary()
		if err != nil {
			return errors.Wrap(err, "marshal binary error")
		}

		payloads = append(payloads, b)
	}

	// enqueue the payloads
	_, err = multicast.EnqueueMultiple(ctx, handler, *item.MulticastGroupID, fragmentation.DefaultFPort, payloads, c.nsCli)
	if err != nil {
		return errors.Wrap(err, "enqueue multiple error")
	}

	item.State = fds.FUOTADeploymentStatusRequest

	switch item.GroupType {
	case fds.FUOTADeploymentGroupTypeC:
		item.NextStepAfter = time.Now().Add(time.Second * time.Duration(1<<uint(item.MulticastTimeout)))
	default:
		return fmt.Errorf("group-type not implemented: %s", item.GroupType)
	}

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepStatusRequest(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	devEUIs, err := handler.GetDevEUIsWithFragmentationSessionSetup(ctx, item.MulticastGroupID, c.s.FragIndex)
	if err != nil {
		return errors.Wrap(err, "get devices with fragmentation session setup error")
	}

	for _, devEUI := range devEUIs {
		cmd := fragmentation.Command{
			CID: fragmentation.FragSessionStatusReq,
			Payload: &fragmentation.FragSessionStatusReqPayload{
				FragStatusReqParam: fragmentation.FragSessionStatusReqPayloadFragStatusReqParam{
					FragIndex:    uint8(c.s.FragIndex),
					Participants: true,
				},
			},
		}
		b, err := cmd.MarshalBinary()
		if err != nil {
			return errors.Wrap(err, "marshal binary error")
		}

		_, err = device.EnqueueDownlinkPayload(ctx, handler, devEUI, false, fragmentation.DefaultFPort, b, c.nsCli)
		if err != nil {
			return errors.Wrap(err, "enqueue downlink payload error")
		}
	}

	item.State = fds.FUOTADeploymentSetDeviceStatus
	item.NextStepAfter = time.Now().Add(item.UnicastTimeout)

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepSetDeviceStatus(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID == nil {
		return errors.New("MulticastGroupID must not be nil")
	}

	err := handler.SetFromRemoteMulticastSetup(ctx, item.ID, *item.MulticastGroupID)
	if err != nil {
		return errors.Wrap(err, "set remote multicast setup error error")
	}

	err = handler.SetFromRemoteFragmentationSession(ctx, item.ID, c.s.FragIndex)
	if err != nil {
		return errors.Wrap(err, "set fragmentation session setup error error")
	}

	err = handler.SetIncompleteFuotaDevelopment(ctx, item.ID)
	if err != nil {
		return errors.Wrap(err, "set incomplete fuota deployment error")
	}

	item.State = fds.FUOTADeploymentCleanup
	item.NextStepAfter = time.Now()

	err = handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}

func (c *controller) stepCleanup(ctx context.Context, handler *store.Handler, item fds.FUOTADeployment) error {
	if item.MulticastGroupID != nil {
		if err := multicast.DeleteMulticastGroup(ctx, *item.MulticastGroupID, c.nsCli); err != nil {
			return errors.Wrap(err, "delete multicast group error")
		}
	}

	item.MulticastGroupID = nil
	item.State = fds.FUOTADeploymentDone

	err := handler.UpdateFUOTADeployment(ctx, &item)
	if err != nil {
		return errors.Wrap(err, "update fuota deployment error")
	}

	return nil
}
