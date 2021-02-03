package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	mxprotocolconn "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"
	"github.com/mxc-foundation/lpwan-app-server/internal/shopify"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"

	"github.com/mxc-foundation/lpwan-app-server/internal/dhx"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	servermod "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
)

func run(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init config in all modules
	if err := mgr.SetupSystemSettings(config.C); err != nil {
		log.WithError(err).Fatal("set up configuration error")
	}

	// set up log level
	log.SetLevel(log.Level(uint8(servermod.GetSettings().LogLevel)))
	// set up syslog
	if err = setSyslog(); err != nil {
		log.Fatal(err)
	}
	// print start message
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "https://mxc.wiki/",
	}).Info("starting Lpwan Application Server")

	if err := mgr.SetupSystemModules(); err != nil {
		log.Fatal(err)
	}

	if err := dhx.Setup(dhx.Config{
		Enable:      config.C.DHXCenter.Enable,
		SupernodeID: config.C.General.ServerAddr,
		DHXServer:   config.C.DHXCenter.DHXServer,
	}); err != nil {
		log.Fatal(err)
	}

	if err := shopify.Start(config.C.ShopifyConfig, pgstore.New(), mxprotocolconn.GetDistributeBonusServiceClient()); err != nil {
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
