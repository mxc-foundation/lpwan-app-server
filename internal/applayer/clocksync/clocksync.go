package clocksync

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/applayer/clocksync"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// HandleClockSyncCommand handles an uplink clock synchronization command.
func HandleClockSyncCommand(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64,
	timeSinceGPSEpoch time.Duration, b []byte, nsCli *nscli.Client) error {
	var cmd clocksync.Command

	if err := cmd.UnmarshalBinary(true, b); err != nil {
		return errors.Wrap(err, "unmarshal command error")
	}

	switch cmd.CID {
	case clocksync.AppTimeReq:
		pl, ok := cmd.Payload.(*clocksync.AppTimeReqPayload)
		if !ok {
			return fmt.Errorf("expected *clocksync.AppTimeReqPayload, got: %T", cmd.Payload)
		}
		if err := handleAppTimeReq(ctx, handler, devEUI, timeSinceGPSEpoch, pl, nsCli); err != nil {
			return errors.Wrap(err, "handle AppTimeReq error")
		}
	default:
		return fmt.Errorf("CID not implemented: %s", cmd.CID)
	}

	return nil
}

func handleAppTimeReq(ctx context.Context, handler *store.Handler, devEUI lorawan.EUI64, timeSinceGPSEpoch time.Duration,
	pl *clocksync.AppTimeReqPayload, nsCli *nscli.Client) error {
	deviceGPSTime := int64(pl.DeviceTime)
	networkGPSTime := int64((timeSinceGPSEpoch / time.Second) % (1 << 32))

	log.WithFields(log.Fields{
		"dev_eui":      devEUI,
		"device_time":  pl.DeviceTime,
		"ans_required": pl.Param.AnsRequired,
		"token_req":    pl.Param.TokenReq,
	}).Info("AppTimeReq received")

	ans := clocksync.Command{
		CID: clocksync.AppTimeAns,
		Payload: &clocksync.AppTimeAnsPayload{
			TimeCorrection: int32(networkGPSTime - deviceGPSTime),
			Param: clocksync.AppTimeAnsPayloadParam{
				TokenAns: pl.Param.TokenReq,
			},
		},
	}
	b, err := ans.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "marshal command error")
	}

	_, err = device.EnqueueDownlinkPayload(ctx, handler, devEUI, false, uint8(clocksync.DefaultFPort), b, nsCli)
	if err != nil {
		return errors.Wrap(err, "enqueue downlink payload error")
	}

	log.WithFields(log.Fields{
		"dev_eui":         devEUI,
		"time_correction": int32(networkGPSTime - deviceGPSTime),
		"token_ans":       pl.Param.TokenReq,
	}).Info("AppTimeAns enqueued")

	return nil
}
