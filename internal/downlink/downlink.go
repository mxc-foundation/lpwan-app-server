package downlink

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	"github.com/mxc-foundation/lpwan-app-server/internal/codec"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/logging"
	apps "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	ds "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

type controller struct {
	h             *store.Handler
	gIntegrations []models.IntegrationHandler
}

// Server represents downlink service
type Server struct {
	downLinkChan chan models.DataDownPayload
}

// Start starts service which handles received downlink payloads to be emitted to the devices.
func Start(h *store.Handler, gIntegrations []models.IntegrationHandler) *Server {
	ctrl := &controller{
		h:             h,
		gIntegrations: gIntegrations,
	}
	downChan := make(chan models.DataDownPayload)

	go func() {
		downChan = integration.ForApplicationID(0, gIntegrations).DataDownChan()
		for pl := range downChan {
			go func(pl models.DataDownPayload) {
				ctxID, err := uuid.NewV4()
				if err != nil {
					log.WithError(err).Error("new uuid error")
					return
				}

				ctx := context.Background()
				ctx = context.WithValue(ctx, logging.ContextIDKey, ctxID)

				if err := ctrl.handleDataDownPayload(ctx, pl); err != nil {
					log.WithFields(log.Fields{
						"dev_eui":        pl.DevEUI,
						"application_id": pl.ApplicationID,
					}).Errorf("handle data-down payload error: %s", err)
				}
			}(pl)
		}
	}()

	return &Server{downLinkChan: downChan}
}

// Stop closes down link channel
func (s *Server) Stop() {
	select {
	case <-s.downLinkChan:
		return
	default:
	}
	// close a closed channel will cause panic
	close(s.downLinkChan)
}

func (c *controller) handleDataDownPayload(ctx context.Context, pl models.DataDownPayload) error {
	return c.h.Tx(ctx, func(ctx context.Context, handler *store.Handler) error {
		// lock the device so that a concurrent Enqueue action will block
		// until this transaction has been completed
		d, err := handler.GetDevice(ctx, pl.DevEUI, true)
		if err != nil {
			return fmt.Errorf("get device error: %s", err)
		}

		// Validate that the ApplicationID matches the actual DevEUI.
		// This is needed as authorisation might be performed on MQTT topic level
		// where it is unknown if the given ApplicationID matches the given
		// DevEUI.
		if d.ApplicationID != pl.ApplicationID {
			return errors.New("enqueue downlink payload: device does not exist for given application")
		}

		// if Object is set, try to encode it to bytes using the application codec
		//if pl.Object != nil && string(pl.Object) != "null" {
		if pl.Object != nil && string(pl.Object) != "null" {
			app, err := handler.GetApplication(ctx, d.ApplicationID)
			if err != nil {
				return errors.Wrap(err, "get application error")
			}

			dp, err := handler.GetDeviceProfile(ctx, d.DeviceProfileID, false)
			if err != nil {
				return errors.Wrap(err, "get device-profile error")
			}

			// TODO: in the next major release, remove this and always use the
			// device-profile codec fields.
			payloadCodec := app.PayloadCodec
			payloadEncoderScript := app.PayloadEncoderScript

			if dp.PayloadCodec != "" {
				payloadCodec = string(dp.PayloadCodec)
				payloadEncoderScript = dp.PayloadEncoderScript
			}

			pl.Data, err = codec.JSONToBinary(codec.Type(payloadCodec), pl.FPort, d.Variables, payloadEncoderScript, pl.Object)
			if err != nil {
				c.logCodecError(ctx, app, d, err)
				return errors.Wrap(err, "encode object error")
			}
		}

		if _, err := devmod.EnqueueDownlinkPayload(ctx, handler, pl.DevEUI, pl.Confirmed, pl.FPort, pl.Data); err != nil {
			return errors.Wrap(err, "enqueue downlink device-queue item error")
		}

		return nil
	})
}

func (c *controller) logCodecError(ctx context.Context, a apps.Application, d ds.Device, err error) {
	errEvent := pb.ErrorEvent{
		ApplicationId:   uint64(a.ID),
		ApplicationName: a.Name,
		DeviceName:      d.Name,
		DevEui:          d.DevEUI[:],
		Type:            pb.ErrorType_DOWNLINK_CODEC,
		Error:           err.Error(),
		Tags:            make(map[string]string),
	}

	for k, v := range d.Tags.Map {
		if v.Valid {
			errEvent.Tags[k] = v.String
		}
	}

	vars := make(map[string]string)
	for k, v := range d.Variables.Map {
		if v.Valid {
			vars[k] = v.String
		}
	}

	if err := integration.ForApplicationID(a.ID, c.gIntegrations).HandleErrorEvent(ctx, vars, errEvent); err != nil {
		log.WithError(err).WithField("ctx_id", ctx.Value(logging.ContextIDKey)).Error("send error event to integration error")
	}
}
