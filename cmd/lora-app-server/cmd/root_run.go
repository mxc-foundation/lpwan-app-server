package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/mxc-foundation/lpwan-app-server/internal/app"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
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
	log.SetLevel(log.Level(uint8(serverinfo.GetSettings().LogLevel)))
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

	a, err := app.Start(ctx, config.C)
	if err != nil {
		log.Errorf("failed to start: %v", err)
		return err
	}
	defer a.Close()

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
