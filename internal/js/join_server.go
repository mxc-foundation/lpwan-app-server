package js

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"

	"github.com/brocaar/lorawan"
	"github.com/brocaar/lorawan/backend/joinserver"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	. "github.com/mxc-foundation/lpwan-app-server/internal/js/data"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "join_server"

type controller struct {
	name string
	s    JoinServerStruct

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	ctrl = &controller{
		name: moduleName,
		s:    conf.JoinServer,
	}
	return nil
}
func GetSettings() JoinServerStruct {
	return ctrl.s
}

// Setup configures the package.
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	log.WithFields(log.Fields{
		"bind":     ctrl.s.Bind,
		"ca_cert":  ctrl.s.CACert,
		"tls_cert": ctrl.s.TLSCert,
		"tls_key":  ctrl.s.TLSKey,
	}).Info("api/js: starting join-server api")

	handler, err := getHandler(h)
	if err != nil {
		return errors.Wrap(err, "get join-server handler error")
	}

	server := http.Server{
		Handler:   handler,
		Addr:      ctrl.s.Bind,
		TLSConfig: &tls.Config{},
	}

	if ctrl.s.CACert == "" && ctrl.s.TLSCert == "" && ctrl.s.TLSKey == "" {
		go func() {
			err := server.ListenAndServe()
			log.WithError(err).Fatal("join-server api error")
		}()

		return nil
	}

	if ctrl.s.CACert != "" {
		caCert, err := ioutil.ReadFile(ctrl.s.CACert)
		if err != nil {
			return errors.Wrap(err, "read ca certificate error")
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return errors.New("append ca certificate error")
		}

		server.TLSConfig.ClientCAs = caCertPool
		server.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert

		log.WithFields(log.Fields{
			"ca_cert": ctrl.s.CACert,
		}).Info("api/js: join-server is configured with client-certificate authentication")
	}

	go func() {
		err := server.ListenAndServeTLS(ctrl.s.TLSCert, ctrl.s.TLSKey)
		log.WithError(err).Fatal("api/js: join-server api error")
	}()

	return nil
}
func getHandler(st *store.Handler) (http.Handler, error) {
	jsConf := joinserver.HandlerConfig{
		Logger: log.StandardLogger(),
		GetDeviceKeysByDevEUIFunc: func(devEUI lorawan.EUI64) (joinserver.DeviceKeys, error) {
			dk, err := st.GetDeviceKeys(context.TODO(), devEUI)
			if err != nil {
				return joinserver.DeviceKeys{}, errors.Wrap(err, "get device-keys error")
			}

			if dk.JoinNonce == (1<<24)-1 {
				return joinserver.DeviceKeys{}, errors.New("join-nonce overflow")
			}
			dk.JoinNonce++
			if err := st.UpdateDeviceKeys(context.TODO(), &dk); err != nil {
				return joinserver.DeviceKeys{}, errors.Wrap(err, "update device-keys error")
			}

			return joinserver.DeviceKeys{
				DevEUI:    dk.DevEUI,
				NwkKey:    dk.NwkKey,
				AppKey:    dk.AppKey,
				JoinNonce: dk.JoinNonce,
			}, nil
		},
		GetKEKByLabelFunc: func(label string) ([]byte, error) {
			for _, kek := range ctrl.s.KEK.Set {
				if label == kek.Label {
					b, err := hex.DecodeString(kek.KEK)
					if err != nil {
						return nil, errors.Wrap(err, "decode hex encoded kek error")
					}

					return b, nil
				}
			}

			return nil, nil
		},
		GetASKEKLabelByDevEUIFunc: func(devEUI lorawan.EUI64) (string, error) {
			return ctrl.s.KEK.ASKEKLabel, nil
		},
		GetHomeNetIDByDevEUIFunc: func(devEUI lorawan.EUI64) (lorawan.NetID, error) {
			d, err := st.GetDevice(context.TODO(), devEUI, false)
			if err != nil {
				if errors.Cause(err) == errHandler.ErrDoesNotExist {
					return lorawan.NetID{}, joinserver.ErrDevEUINotFound
				}

				return lorawan.NetID{}, errors.Wrap(err, "get device error")
			}

			var netID lorawan.NetID

			netIDStr, ok := d.Variables.Map["home_netid"]
			if !ok {
				return netID, nil
			}

			if err := netID.UnmarshalText([]byte(netIDStr.String)); err != nil {
				return lorawan.NetID{}, errors.Wrap(err, "unmarshal netid error")
			}

			return netID, nil
		},
	}

	handler, err := joinserver.NewHandler(jsConf)
	if err != nil {
		return nil, errors.Wrap(err, "new join-server handler error")
	}

	return &prometheusMiddleware{
		handler:         handler,
		timingHistogram: metrics.GetMetricsSettings().Prometheus.APITimingHistogram,
	}, nil
}
