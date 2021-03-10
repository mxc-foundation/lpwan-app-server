// Package app is responsible for starting all the routines required for
// running appserver.
package app

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations/code"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/bonus"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/dhx"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpapisrv"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/shopify"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

// App represents the running appserver
type App struct {
	// postgres server client
	pgstore *pgstore.PgStore
	// bonus distribution service
	bonus *bonus.Service
	// mxprotocol server client
	mxpCli *mxpcli.Client
	// mxprotocol server API
	mxpSrv *mxpapisrv.MXPAPIServer
	// shopify service
	shopify *shopify.Service
	// smtp service
	mailer *email.Mailer
}

// Start starts all the routines required for appserver and returns the App
// structure or error.
func Start(ctx context.Context, cfg config.Config) (*App, error) {
	app := &App{}
	if err := app.externalServices(ctx, cfg); err != nil {
		// we already have an error
		_ = app.Close()
		return nil, err
	}
	if err := app.systemManager(ctx, cfg); err != nil {
		// we already have an error
		_ = app.Close()
		return nil, err
	}
	if err := app.initInParallel(ctx, cfg); err != nil {
		// we already have an error
		_ = app.Close()
		return nil, err
	}
	if err := app.startAPIs(ctx, cfg); err != nil {
		// we already have an error
		_ = app.Close()
		return nil, err
	}

	return app, nil
}

// Close stops the appserver
func (app *App) Close() error {
	if app.mxpSrv != nil {
		app.mxpSrv.Stop()
	}
	if app.mxpCli != nil {
		if err := app.mxpCli.Close(); err != nil {
			logrus.Warnf("error shutting down MXP server connection: %v", err)
		}
	}
	if app.bonus != nil {
		app.bonus.Stop()
	}
	if app.shopify != nil {
		app.shopify.Stop()
	}
	return nil
}

// establish connections to external services
func (app *App) externalServices(ctx context.Context, cfg config.Config) error {
	var err error
	// postgres
	app.pgstore, err = pgstore.Setup(cfg.PostgreSQL)
	if err != nil {
		return err
	}
	// data migrations
	err = code.Setup(store.NewStore(), cfg.PostgreSQL.Automigrate)
	if err != nil {
		return err
	}
	// mxprotocol server client
	app.mxpCli, err = mxpcli.Connect(cfg.M2MServer)
	if err != nil {
		return err
	}
	mxpcli.Global = app.mxpCli
	// bonus distribution service
	app.bonus = bonus.Start(ctx, cfg.ApplicationServer.Airdrop,
		app.pgstore, app.mxpCli.GetDistributeBonusServiceClient())
	// shopify service
	app.shopify, err = shopify.Start(ctx, cfg.ShopifyConfig, app.pgstore, app.mxpCli.GetDistributeBonusServiceClient())
	if err != nil {
		return err
	}

	return nil
}

// this part should be gradually removed, no service should be managed by this
func (app *App) systemManager(ctx context.Context, cfg config.Config) error {
	// init config in all modules
	if err := mgr.SetupSystemSettings(cfg); err != nil {
		logrus.WithError(err).Error("set up configuration error")
		return err
	}

	// set up log level
	logrus.SetLevel(logrus.Level(uint8(serverinfo.GetSettings().LogLevel)))

	// print start message
	logrus.WithFields(logrus.Fields{
		"version": config.AppserverVersion,
		"docs":    "https://mxc.wiki/",
	}).Info("starting Lpwan Application Server")

	if err := mgr.SetupSystemModules(); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

// services that can be initialized in parallel, put here the ones that don't
// depend on other services
func (app *App) initInParallel(ctx context.Context, cfg config.Config) error {
	// register on dhx-center server
	if err := dhx.Register(cfg.General.ServerAddr, cfg.DHXCenter); err != nil {
		return fmt.Errorf("couldn't register on DHX server: %v", err)
	}
	// set up the email system
	var err error
	app.mailer, err = email.NewMailer(cfg.Operator, cfg.SMTP, email.ServerInfo{
		ServerAddr:      cfg.General.ServerAddr,
		DefaultLanguage: cfg.General.DefaultLanguage,
	})
	if err != nil {
		return err
	}
	// start bonus distribution service
	return nil
}

// the gRPC servers
func (app *App) startAPIs(ctx context.Context, cfg config.Config) error {
	var err error
	// API for mxprotocol server
	if app.mxpSrv, err = mxpapisrv.Start(app.pgstore, cfg.ApplicationServer.APIForM2M, app.mailer); err != nil {
		return err
	}
	// API for external clients
	if err = external.Start(store.NewStore(), external.RESTApiServer{
		S:                      cfg.ApplicationServer.ExternalAPI,
		ApplicationServerID:    cfg.ApplicationServer.ID,
		ServerAddr:             cfg.General.ServerAddr,
		Recaptcha:              cfg.Recaptcha,
		Enable2FA:              cfg.General.Enable2FALogin,
		ServerRegion:           cfg.General.ServerRegion,
		PasswordHashIterations: cfg.General.PasswordHashIterations,
		EnableSTC:              cfg.General.EnableSTC,
		ExternalAuth:           cfg.ExternalAuth,
		ShopifyConfig:          cfg.ShopifyConfig,
		OperatorLogo:           cfg.Operator.OperatorLogo,
		Mailer:                 app.mailer,
		MXPCli:                 app.mxpCli,
	}); err != nil {
		return err
	}

	return nil
}
