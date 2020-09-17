package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mxc-foundation/lpwan-app-server/internal/api"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	nscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/networkserver"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	jscodec "github.com/mxc-foundation/lpwan-app-server/internal/codec/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/downlink"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/fuota"
	"github.com/mxc-foundation/lpwan-app-server/internal/gwping"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations/code"
	"github.com/mxc-foundation/lpwan-app-server/internal/monitoring"
	"github.com/mxc-foundation/lpwan-app-server/internal/pprof"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	appmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	asmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/as"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	devprofilemod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile"
	fuotamod "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment"
	gwmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	gpmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	miningmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"
	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	orgmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/pgstore"
	servermod "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	serviceprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	usermod "github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pwh, err := pwhash.New(16, servermod.GetSettings().PasswordHashIterations)
	if err != nil {
		return err
	}

	handler, _ := store.New(pgstore.New(storage.DBTest().DB, pgstore.Settings{
		ApplicationServerID:         external.GetApplicationServerID(),
		JWTSecret:                   external.GetJWTSecret(),
		ApplicationServerPublicHost: as.GetSettings().PublicHost,
		PWH:                         pwh,
	}))

	if err := setLogLevel(); err != nil {
		log.Fatal(err)
	}
	if err := setSyslog(); err != nil {
		log.Fatal(err)
	}
	if err := printStartMessage(); err != nil {
		log.Fatal(err)
	}
	if err := startPProf(); err != nil {
		log.Fatal(err)
	}
	if err := setupStorage(); err != nil {
		log.Fatal(err)
	}
	if err := setupNetworkserver(); err != nil {
		log.Fatal(err)
	}
	if err := setupClient(); err != nil {
		log.Fatal(err)
	}
	if err := migrateGatewayStats(); err != nil {
		log.Fatal(err)
	}
	if err := migrateToClusterKeys(); err != nil {
		log.Fatal(err)
	}
	if err := setupIntegration(); err != nil {
		log.Fatal(err)
	}
	if err := setupSMTP(); err != nil {
		log.Fatal(err)
	}
	if err := setupCodec(); err != nil {
		log.Fatal(err)
	}
	if err := handleDataDownPayloads(); err != nil {
		log.Fatal(err)
	}
	if err := startGatewayPing(handler); err != nil {
		log.Fatal(err)
	}
	if err := setupMulticastSetup(); err != nil {
		log.Fatal(err)
	}
	if err := setupFragmentation(); err != nil {
		log.Fatal(err)
	}
	if err := setupFUOTA(); err != nil {
		log.Fatal(err)
	}

	if err := setupModules(handler); err != nil {
		log.Fatal(err)
	}
	if err := setupAPI(); err != nil {
		log.Fatal(err)
	}
	if err := setupMonitoring(); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping lora-app-server")
		// todo: handle graceful shutdown?
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func startPProf() error {
	return pprof.Setup()
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(servermod.GetSettings().LogLevel)))
	return nil
}

func printStartMessage() error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://www.chirpstack.io/",
	}).Info("starting ChirpStack Application Server")
	return nil
}

func setupStorage() error {
	if err := storage.Setup(); err != nil {
		return errors.Wrap(err, "setup storage error")
	}

	return nil
}

func setupSMTP() error {
	if err := email.Setup(); err != nil {
		return errors.Wrap(err, "setup SMTP error")
	}

	return nil
}

func setupIntegration() error {
	if err := integration.Setup(); err != nil {
		return errors.Wrap(err, "setup integration error")
	}

	return nil
}

func setupCodec() error {
	if err := jscodec.Setup(); err != nil {
		return errors.Wrap(err, "setup codec error")
	}

	return nil
}

func setupNetworkserver() error {
	if err := networkserver.Setup(); err != nil {
		return errors.Wrap(err, "setup networkserver pool error")
	}

	return nil
}

func setupClient() error {
	if err := nscli.Setup(); err != nil {
		return errors.Wrap(err, "setup networkserver connection error")
	}

	if err := m2mcli.Setup(); err != nil {
		return errors.Wrap(err, "setup m2m-server connection error")
	}

	if err := pscli.Setup(); err != nil {
		return errors.Wrap(err, "setup provisioning server connection error")
	}

	return nil
}

func migrateGatewayStats() error {
	if err := code.Migrate("migrate_gw_stats", code.MigrateGatewayStats); err != nil {
		return errors.Wrap(err, "migration error")
	}

	return nil
}

func migrateToClusterKeys() error {
	return code.Migrate("migrate_to_cluster_keys", func(handler *store.Handler) error {
		return code.MigrateToClusterKeys(config.C)
	})
}

func handleDataDownPayloads() error {
	downChan := integration.ForApplicationID(0).DataDownChan()
	go downlink.HandleDataDownPayloads(downChan)
	return nil
}

func startGatewayPing(h *store.Handler) error {
	go gwping.SendPingLoop(h)

	return nil
}

func setupMulticastSetup() error {
	if err := multicastsetup.Setup(); err != nil {
		return errors.Wrap(err, "multicastsetup setup error")
	}
	return nil
}

func setupFragmentation() error {
	if err := fragmentation.Setup(); err != nil {
		return errors.Wrap(err, "fragmentation setup error")
	}
	return nil
}

func setupFUOTA() error {
	if err := fuota.Setup(); err != nil {
		return errors.Wrap(err, "fuota setup error")
	}
	return nil
}

func setupModules(h *store.Handler) (err error) {
	if err = gwmod.Setup(h); err != nil {
		return err
	}

	if err = devmod.Setup(h); err != nil {
		return err
	}

	if err = appmod.Setup(h); err != nil {
		return err
	}

	if err = gpmod.Setup(h); err != nil {
		return err
	}

	if err = miningmod.Setup(h); err != nil {
		return err
	}

	if err = nsmod.Setup(h); err != nil {
		return err
	}

	if err = orgmod.Setup(h); err != nil {
		return err
	}

	if err = usermod.Setup(h); err != nil {
		return err
	}

	if err = servermod.Setup(h); err != nil {
		return err
	}

	if err = asmod.Setup(h); err != nil {
		return err
	}

	if err = devprofilemod.Setup(h); err != nil {
		return err
	}

	if err = serviceprofile.Setup(h); err != nil {
		return err
	}

	if err = multicast.Setup(h); err != nil {
		return err
	}

	if err = fuotamod.Setup(h); err != nil {
		return err
	}

	return nil
}

func setupAPI() error {
	if err := api.Setup(); err != nil {
		return errors.Wrap(err, "setup api error")
	}
	return nil
}

func setupMonitoring() error {
	if err := monitoring.Setup(); err != nil {
		return errors.Wrap(err, "setup monitoring error")
	}
	return nil
}
