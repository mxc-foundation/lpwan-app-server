// Package app is responsible for starting all the routines required for
// running appserver.
package app

import (
	"context"
	"fmt"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/dhx"
)

// App represents the running appserver
type App struct {
}

// Start starts all the routines required for appserver and returns the App
// structure or error.
func Start(ctx context.Context, cfg config.Config) (*App, error) {
	app := &App{}
	if err := app.initInParallel(ctx, cfg); err != nil {
		// we already have an error
		_ = app.Close()
		return nil, err
	}

	return app, nil
}

// Close stops the appserver
func (app *App) Close() error {
	return nil
}

// services that can be initialized in parallel, put here the ones that don't
// depend on other services
func (app *App) initInParallel(ctx context.Context, cfg config.Config) error {
	if err := dhx.Register(cfg.General.ServerAddr, cfg.DHXCenter); err != nil {
		return fmt.Errorf("couldn't register on DHX server: %v", err)
	}
	return nil
}
