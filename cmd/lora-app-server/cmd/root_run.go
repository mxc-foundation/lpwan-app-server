package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mxc-foundation/lpwan-app-server/internal/api"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/m2m_client"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/codec"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/downlink"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/fuota"
	gw "github.com/mxc-foundation/lpwan-app-server/internal/gateway-manager"
	"github.com/mxc-foundation/lpwan-app-server/internal/gwping"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/multi"
	"github.com/mxc-foundation/lpwan-app-server/internal/metrics"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations/code"
	"github.com/mxc-foundation/lpwan-app-server/internal/mining"
	"github.com/mxc-foundation/lpwan-app-server/internal/pprof"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := []func() error{
		setLogLevel,
		printStartMessage,
		startPProf,
		setupStorage,
		setupClient,
		setupUpdateFirmwareFromPs,
		setupDefaultEnv,
		setupLoadGatewayTemplates,
		migrateGatewayStats,
		setupIntegration,
		setupSMTP,
		setupCodec,
		handleDataDownPayloads,
		startGatewayPing,
		setupMulticastSetup,
		setupFragmentation,
		setupFUOTA,
		setupMetrics,
		setupMining,
		setupAPI,
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
		"docs":    "https://www.loraserver.io/",
	}).Info("starting LPWAN App Server")
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
	var confs []interface{}

	for _, name := range config.C.ApplicationServer.Integration.Enabled {
		switch name {
		case "aws_sns":
			confs = append(confs, config.C.ApplicationServer.Integration.AWSSNS)
		case "azure_service_bus":
			confs = append(confs, config.C.ApplicationServer.Integration.AzureServiceBus)
		case "mqtt":
			confs = append(confs, config.C.ApplicationServer.Integration.MQTT)
		case "gcp_pub_sub":
			confs = append(confs, config.C.ApplicationServer.Integration.GCPPubSub)
		case "postgresql":
			confs = append(confs, config.C.ApplicationServer.Integration.PostgreSQL)
		default:
			return fmt.Errorf("unknown integration type: %s", name)
		}
	}

	mi, err := multi.New(confs)
	if err != nil {
		return errors.Wrap(err, "setup integrations error")
	}
	mi.Add(application.New())
	integration.SetIntegration(mi)

	return nil
}

func setupCodec() error {
	if err := codec.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup codec error")
	}
	return nil
}

func setupClient() error {
	if err := setupNetworkServer(); err != nil {
		return err
	}

	if err := setupM2MServer(); err != nil {
		return err
	}

	return nil
}

func setupM2MServer() error {
	if err := m2m_client.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup m2m-server error")
	}
	return nil
}

func setupNetworkServer() error {
	if err := networkserver.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup networkserver error")
	}
	return nil
}

func setupUpdateFirmwareFromPs() error {
	if err := storage.UpdateFirmwareFromProvisioningServer(config.C); err != nil {
		return errors.Wrap(err, "setup update firmware error")
	}
	return nil
}

func setupDefaultEnv() error {
	if err := storage.SetupDefault(); err != nil {
		return errors.Wrap(err, "setup default error")
	}
	return nil
}

func setupLoadGatewayTemplates() error {
	if err := gw.LoadTemplates(); err != nil {
		return errors.Wrap(err, "load gateway config template error")
	}
	return nil
}

func migrateGatewayStats() error {
	if err := code.Migrate("migrate_gw_stats", code.MigrateGatewayStats); err != nil {
		return errors.Wrap(err, "migration error")
	}

	return nil
}

func handleDataDownPayloads() error {
	go downlink.HandleDataDownPayloads()
	return nil
}

func setupAPI() error {
	if err := api.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup api error")
	}
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

func setupMetrics() error {
	if err := metrics.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup metrics error")
	}
	return nil
}

func setupMining() error {
	if config.C.ApplicationServer.MiningSetUp.Mining == false {
		log.Info("Stop mining function")
		return nil
	}

	if err := mining.Setup(config.C); err != nil {
		return errors.Wrap(err, "setup service mining error")
	}
	return nil
}
