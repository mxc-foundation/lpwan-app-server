package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"

	"github.com/jmoiron/sqlx"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/pgstore"
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

	appmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	asmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/as"
	devmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	gwmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	gpmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	miningmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	orgmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	servermod "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	usermod "github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := []func() error{
		setLogLevel,
		setSyslog,
		printStartMessage,
		startPProf,
		setupStorage,
		setupNetworkserver,
		setupClient,
		migrateGatewayStats,
		migrateToClusterKeys,
		setupIntegration,
		setupSMTP,
		setupCodec,
		handleDataDownPayloads,
		startGatewayPing,
		setupMulticastSetup,
		setupFragmentation,
		setupFUOTA,

		setupModules,
		setupUpdateFirmwareFromPs,
		setupAPI,
		setupMonitoring,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
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
	return pprof.Setup(config.C.PProf)
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
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
	if err := storage.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup storage error")
	}

	return nil
}

func setupSMTP() error {
	if err := email.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup SMTP error")
	}

	return nil
}

func setupIntegration() error {
	if err := integration.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup integration error")
	}

	return nil
}

func setupCodec() error {
	if err := jscodec.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup codec error")
	}

	return nil
}

func setupNetworkserver() error {
	if err := networkserver.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup networkserver pool error")
	}

	return nil
}

func setupClient() error {
	if err := nscli.Setup(); err != nil {
		return errors.Wrap(err, "setup networkserver connection error")
	}

	if err := m2mcli.Setup(config.C.M2MServer); err != nil {
		return errors.Wrap(err, "setup m2m-server connection error")
	}

	if err := pscli.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup provisioning server connection error")
	}

	return nil
}

func setupUpdateFirmwareFromPs() error {
	if err := gwmod.Service.UpdateFirmwareFromProvisioningServer(context.TODO()); err != nil {
		return errors.Wrap(err, "setup update firmware error")
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
	return code.Migrate("migrate_to_cluster_keys", func(db sqlx.Ext) error {
		return code.MigrateToClusterKeys(config.C)
	})
}

func handleDataDownPayloads() error {
	downChan := integration.ForApplicationID(0).DataDownChan()
	go downlink.HandleDataDownPayloads(downChan)
	return nil
}

func startGatewayPing() error {
	go gwping.SendPingLoop()

	return nil
}

func setupMulticastSetup() error {
	if err := multicastsetup.Setup(config.C); err != nil {
		return errors.Wrap(err, "multicastsetup setup error")
	}
	return nil
}

func setupFragmentation() error {
	if err := fragmentation.Setup(config.C); err != nil {
		return errors.Wrap(err, "fragmentation setup error")
	}
	return nil
}

func setupFUOTA() error {
	if err := fuota.Setup(config.C); err != nil {
		return errors.Wrap(err, "fuota setup error")
	}
	return nil
}

func setupModules() (err error) {

	if err = gwmod.Setup(config.C, pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = devmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = appmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = gpmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = miningmod.Setup(config.C.ApplicationServer.MiningSetUp, pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = nsmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = orgmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = usermod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = servermod.Setup(config.C, pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	if err = asmod.Setup(pgstore.New(storage.DB().DB)); err != nil {
		return err
	}

	return nil
}

func setupAPI() error {
	if err := api.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup api error")
	}
	return nil
}

func setupMonitoring() error {
	if err := monitoring.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup monitoring error")
	}
	return nil
}
