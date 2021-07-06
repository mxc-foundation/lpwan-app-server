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
)

func run(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	a, err := app.Start(ctx, config.C)
	if err != nil {
		log.Errorf("failed to start: %v", err)
		return err
	}
	defer a.Close()

	sigChan := make(chan os.Signal, 1)
	exitChan := make(chan struct{}, 1)
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
