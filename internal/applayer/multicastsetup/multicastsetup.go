package multicastsetup

import (
	"context"
	"fmt"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/applayer/multicastsetup"
	"github.com/brocaar/lorawan/gps"

	. "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "multicastsetup"

type controller struct {
	name string
	s    MulticastStruct

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{
		name: moduleName,
		s:    conf.ApplicationServer.RemoteMulticastSetup,
	}
	return nil
}

// Setup configures the package.
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	go SyncRemoteMulticastSetupLoop()
	go SyncRemoteMulticastClassCSessionLoop()

	return nil
}

// SyncRemoteMulticastSetupLoop syncs the multicast setup with the devices.
func SyncRemoteMulticastSetupLoop() {
	for {
		ctxID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("new uuid error")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

		err = storage.Transaction(func(ctx context.Context, handler *store.Handler) error {
			return syncRemoteMulticastSetup(ctx, handler)
		})

		if err != nil {
			log.WithError(err).Error("sync remote multicast setup error")
		}
		time.Sleep(ctrl.s.SyncInterval)
	}
}

// SyncRemoteMulticastClassCSessionLoop syncs the multicast Class-C session
// with the devices.
func SyncRemoteMulticastClassCSessionLoop() {
	for {
		ctxID, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Error("new uuid error")
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

		err = storage.Transaction(func(ctx context.Context, handler *store.Handler) error {
			return syncRemoteMulticastClassCSession(ctx, handler)
		})

		if err != nil {
			log.WithError(err).Error("sync remote multicast class-c session error")
		}
		time.Sleep(ctrl.s.SyncInterval)
	}
}

// HandleRemoteMulticastSetupCommand handles an uplink remote multicast setup command.
func HandleRemoteMulticastSetupCommand(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, b []byte) error {
	var cmd multicastsetup.Command

	if err := cmd.UnmarshalBinary(true, b); err != nil {
		return errors.Wrap(err, "unmarshal command error")
	}

	switch cmd.CID {
	case multicastsetup.McGroupSetupAns:
		pl, ok := cmd.Payload.(*multicastsetup.McGroupSetupAnsPayload)
		if !ok {
			return fmt.Errorf("expected *multicastsetup.McGroupSetupAnsPayload, got: %T", cmd.Payload)
		}
		if err := handleMcGroupSetupAns(ctx, handler, devEUI, pl); err != nil {
			return errors.Wrap(err, "handle McGroupSetupAns error")
		}
	case multicastsetup.McGroupDeleteAns:
		pl, ok := cmd.Payload.(*multicastsetup.McGroupDeleteAnsPayload)
		if !ok {
			return fmt.Errorf("expected *multicastsetup.McGroupDeleteAnsPayload, got: %T", cmd.Payload)
		}
		if err := handleMcGroupDeleteAns(ctx, handler, devEUI, pl); err != nil {
			return errors.Wrap(err, "handle McGroupDeleteAns error")
		}
	case multicastsetup.McClassCSessionAns:
		pl, ok := cmd.Payload.(*multicastsetup.McClassCSessionAnsPayload)
		if !ok {
			return fmt.Errorf("expected *multicastsetup.McClassCSessionAnsPayload, got: %T", cmd.Payload)
		}
		if err := handleMcClassCSessionAns(ctx, handler, devEUI, pl); err != nil {
			return errors.Wrap(err, "handle McClassCSessionAns error")
		}
	default:
		return fmt.Errorf("CID not implemented: %s", cmd.CID)
	}

	return nil
}

func handleMcGroupSetupAns(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, pl *multicastsetup.McGroupSetupAnsPayload) error {
	log.WithFields(log.Fields{
		"dev_eui":     devEUI,
		"id_error":    pl.McGroupIDHeader.IDError,
		"mc_group_id": pl.McGroupIDHeader.McGroupID,
		"ctx_id":      ctx.Value(logging.ContextIDKey),
	}).Info("McGroupSetupAns received")

	if pl.McGroupIDHeader.IDError {
		return fmt.Errorf("IDError for McGroupID: %d", pl.McGroupIDHeader.McGroupID)
	}

	rms, err := storage.GetRemoteMulticastSetupByGroupID(ctx, handler, devEUI, int(pl.McGroupIDHeader.McGroupID), true)
	if err != nil {
		return errors.Wrap(err, "get remote multicast-setup by group id error")
	}

	rms.StateProvisioned = true
	if err := storage.UpdateRemoteMulticastSetup(ctx, handler, &rms); err != nil {
		return errors.Wrap(err, "update remote multicast-setup error")
	}

	return nil
}

func handleMcGroupDeleteAns(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, pl *multicastsetup.McGroupDeleteAnsPayload) error {
	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"mc_group_id":        pl.McGroupIDHeader.McGroupID,
		"mc_group_undefined": pl.McGroupIDHeader.McGroupUndefined,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("McGroupDeleteAns received")

	if pl.McGroupIDHeader.McGroupUndefined {
		return fmt.Errorf("McGroupUndefined for McGroupID: %d", pl.McGroupIDHeader.McGroupID)
	}

	rms, err := storage.GetRemoteMulticastSetupByGroupID(ctx, handler, devEUI, int(pl.McGroupIDHeader.McGroupID), true)
	if err != nil {
		return errors.Wrap(err, "get remote multicast-setup by group id error")
	}

	rms.StateProvisioned = true
	if err := storage.UpdateRemoteMulticastSetup(ctx, handler, &rms); err != nil {
		return errors.Wrap(err, "update remote multicast-setup error")
	}

	if err := multicast.RemoveDeviceFromMulticastGroup(ctx, rms.MulticastGroupID, devEUI); err != nil {
		if err == errHandler.ErrDoesNotExist {
			log.WithFields(log.Fields{
				"dev_eui":            devEUI,
				"multicast_group_id": rms.MulticastGroupID,
				"ctx_id":             ctx.Value(logging.ContextIDKey),
			}).Info("applayer/multicastsetup: removing device from multicast group, but device does not exist")
		} else {
			return errors.Wrap(err, "remove device from multicast group error")
		}
	}

	return nil
}

func handleMcClassCSessionAns(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, pl *multicastsetup.McClassCSessionAnsPayload) error {
	log.WithFields(log.Fields{
		"dev_eui":            devEUI,
		"time_to_start":      pl.TimeToStart,
		"mc_group_undefined": pl.StatusAndMcGroupID.McGroupUndefined,
		"freq_error":         pl.StatusAndMcGroupID.FreqError,
		"dr_error":           pl.StatusAndMcGroupID.DRError,
		"mc_group_id":        pl.StatusAndMcGroupID.McGroupID,
		"ctx_id":             ctx.Value(logging.ContextIDKey),
	}).Info("McClassCSessionAns received")

	if pl.StatusAndMcGroupID.DRError || pl.StatusAndMcGroupID.FreqError || pl.StatusAndMcGroupID.McGroupUndefined {
		return fmt.Errorf("DRError: %t, FreqError: %t, McGroupUndefined: %t for McGroupID: %d", pl.StatusAndMcGroupID.DRError, pl.StatusAndMcGroupID.FreqError, pl.StatusAndMcGroupID.McGroupUndefined, pl.StatusAndMcGroupID.McGroupID)
	}

	sess, err := storage.GetRemoteMulticastClassCSessionByGroupID(ctx, handler, devEUI, int(pl.StatusAndMcGroupID.McGroupID), true)
	if err != nil {
		return errors.Wrap(err, "get remote multicast class-c session error")
	}

	sess.StateProvisioned = true
	if err := storage.UpdateRemoteMulticastClassCSession(ctx, handler, &sess); err != nil {
		return errors.Wrap(err, "update remote multicast class-c session error")
	}

	if err := multicast.AddDeviceToMulticastGroup(ctx, sess.MulticastGroupID, devEUI); err != nil {
		if err == errHandler.ErrAlreadyExists {
			log.WithFields(log.Fields{
				"dev_eui":            devEUI,
				"multicast_group_id": sess.MulticastGroupID,
				"ctx_id":             ctx.Value(logging.ContextIDKey),
			}).Warning("applayer/multicastsetup: adding device to multicast group, but device was already added")
		} else {
			return errors.Wrap(err, "add device to multicast group error")
		}
	}

	return nil
}

func syncRemoteMulticastSetup(ctx context.Context, handler *store.Handler) error {
	items, err := storage.GetPendingRemoteMulticastSetupItems(ctx, handler, ctrl.s.SyncBatchSize, ctrl.s.SyncRetries)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := syncRemoteMulticastSetupItem(ctx, handler, item); err != nil {
			return errors.Wrap(err, "sync remote multicast-setup error")
		}
	}

	return nil
}

func syncRemoteMulticastSetupItem(ctx context.Context, handler *store.Handler, item storage.RemoteMulticastSetup) error {
	var cmd multicastsetup.Command

	switch item.State {
	case storage.RemoteMulticastSetupSetup:
		cmd = multicastsetup.Command{
			CID: multicastsetup.McGroupSetupReq,
			Payload: &multicastsetup.McGroupSetupReqPayload{
				McGroupIDHeader: multicastsetup.McGroupSetupReqPayloadMcGroupIDHeader{
					McGroupID: uint8(item.McGroupID),
				},
				McAddr:         item.McAddr,
				McKeyEncrypted: item.McKeyEncrypted,
				MinMcFCnt:      item.MinMcFCnt,
				MaxMcFCnt:      item.MaxMcFCnt,
			},
		}
	case storage.RemoteMulticastSetupDelete:
		cmd = multicastsetup.Command{
			CID: multicastsetup.McGroupDeleteReq,
			Payload: &multicastsetup.McGroupDeleteReqPayload{
				McGroupIDHeader: multicastsetup.McGroupDeleteReqPayloadMcGroupIDHeader{
					McGroupID: uint8(item.McGroupID),
				},
			},
		}
	default:
		return fmt.Errorf("invalid state: %s", item.State)
	}

	b, err := cmd.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "marshal binary error")
	}

	_, err = storage.EnqueueDownlinkPayload(ctx, handler, item.DevEUI, false, multicastsetup.DefaultFPort, b)
	if err != nil {
		return errors.Wrap(err, "enqueue downlink payload error")
	}

	log.WithFields(log.Fields{
		"dev_eui":     item.DevEUI,
		"mc_group_id": item.McGroupID,
		"ctx_id":      ctx.Value(logging.ContextIDKey),
	}).Infof("%s enqueued", cmd.CID)

	item.RetryCount++
	item.RetryAfter = time.Now().Add(item.RetryInterval)

	err = storage.UpdateRemoteMulticastSetup(ctx, handler, &item)
	if err != nil {
		return errors.Wrap(err, "update remote multicast-setup error")
	}

	return nil
}

func syncRemoteMulticastClassCSession(ctx context.Context, handler *store.Handler) error {
	items, err := storage.GetPendingRemoteMulticastClassCSessions(ctx, handler, ctrl.s.SyncBatchSize, ctrl.s.SyncRetries)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := syncRemoteMulticastClassCSessionItem(ctx, handler, item); err != nil {
			return errors.Wrap(err, "sync remote multicast class-c session error")
		}
	}

	return nil
}

func syncRemoteMulticastClassCSessionItem(ctx context.Context, handler *store.Handler, item storage.RemoteMulticastClassCSession) error {
	cmd := multicastsetup.Command{
		CID: multicastsetup.McClassCSessionReq,
		Payload: &multicastsetup.McClassCSessionReqPayload{
			McGroupIDHeader: multicastsetup.McClassCSessionReqPayloadMcGroupIDHeader{
				McGroupID: uint8(item.McGroupID),
			},
			SessionTime: uint32((gps.Time(item.SessionTime).TimeSinceGPSEpoch() / time.Second) % (1 << 32)),
			SessionTimeOut: multicastsetup.McClassCSessionReqPayloadSessionTimeOut{
				TimeOut: uint8(item.SessionTimeOut),
			},
			DLFrequency: uint32(item.DLFrequency),
			DR:          uint8(item.DR),
		},
	}

	b, err := cmd.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "marshal binary error")
	}

	_, err = storage.EnqueueDownlinkPayload(ctx, handler, item.DevEUI, false, multicastsetup.DefaultFPort, b)
	if err != nil {
		return errors.Wrap(err, "enqueue downlink payload error")
	}

	log.WithFields(log.Fields{
		"dev_eui":     item.DevEUI,
		"mc_group_id": item.McGroupID,
	}).Infof("%s enqueued", cmd.CID)

	item.RetryCount++
	item.RetryAfter = time.Now().Add(item.RetryInterval)

	err = storage.UpdateRemoteMulticastClassCSession(ctx, handler, &item)
	if err != nil {
		return errors.Wrap(err, "update remote multicast class-c session error")
	}

	return nil
}
