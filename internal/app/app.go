// Package app is responsible for starting all the routines required for
// running appserver.
package app

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/bonus"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/devprovision"
	"github.com/mxc-foundation/lpwan-app-server/internal/downlink"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations/code"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpapisrv"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/pscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/shopify"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
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
	// network server client
	nsCli *nscli.Client
	// provisioning server client
	psCli *pscli.Client
	// mxprotocol server API
	mxpSrv *mxpapisrv.MXPAPIServer
	// network server API
	nsSrv *as.NetworkServerAPIServer
	// external API server
	extAPISrv *external.ExtAPIServer
	// gateway server for new gateway model
	newGwSrv *gateway.Server
	// gateway server for old gateway model
	oldGwSrv *gateway.Server
	// shopify service
	shopify *shopify.Service
	// integration handlers
	integrations []models.IntegrationHandler
	// smtp service
	mailer *email.Mailer
	// uuid of appserver used by other internal servers to identify appserver
	applicationServerID uuid.UUID
	// device provisioning session list
	devSessionList *devprovision.DeviceSessionList
}

// Start starts all the routines required for appserver and returns the App
// structure or error.
func Start(ctx context.Context, cfg config.Config) (*App, error) {
	app := &App{}
	var err error
	app.applicationServerID, err = uuid.FromString(cfg.ApplicationServer.ID)
	if err != nil {
		return nil, err
	}
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

	logrus.Info("Successfully start all services in server")
	return app, nil
}

// Close stops the appserver
func (app *App) Close() error {
	if app.mxpSrv != nil {
		app.mxpSrv.Stop()
	}
	if app.nsSrv != nil {
		app.nsSrv.Stop()
	}
	if app.extAPISrv != nil {
		app.extAPISrv.Stop()
	}
	if app.newGwSrv != nil {
		app.newGwSrv.Stop()
	}
	if app.oldGwSrv != nil {
		app.oldGwSrv.Stop()
	}
	if app.mxpCli != nil {
		if err := app.mxpCli.Close(); err != nil {
			logrus.Warnf("error shutting down MXP server connection: %v", err)
		}
	}
	if app.psCli != nil {
		if err := app.psCli.Close(); err != nil {
			logrus.Warnf("error shutting down Provisioning server connection: %v", err)
		}
	}
	if app.nsCli != nil {
		if err := app.nsCli.Close(); err != nil {
			logrus.Warnf("error shutting down network server connection: %v", err)
		}
	}
	if app.bonus != nil {
		app.bonus.Stop()
	}
	if app.shopify != nil {
		app.shopify.Stop()
	}
	for _, v := range app.integrations {
		if err := v.Close(); err != nil {
			logrus.Warnf("error shutting down integrations: %v", err)
		}
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
	// network server client (also used by code migration)
	if err = app.networkServer(ctx, cfg); err != nil {
		return err
	}
	// provisioning server client
	app.psCli, err = pscli.Connect(cfg.ProvisionServer)
	if err != nil {
		return err
	}
	// bonus distribution service
	app.bonus = bonus.Start(ctx, cfg.ApplicationServer.Airdrop,
		app.pgstore, app.mxpCli.GetDistributeBonusServiceClient())
	// shopify service
	app.shopify, err = shopify.Start(ctx, cfg.ShopifyConfig, app.pgstore, app.mxpCli.GetDistributeBonusServiceClient())
	if err != nil {
		return err
	}

	// integration service clients
	app.integrations, err = integration.SetupGlobalIntegrations(cfg.ApplicationServer.Integration)
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
	// set up the email system
	var err error
	app.mailer, err = email.NewMailer(cfg.Operator, cfg.SMTP, email.ServerInfo{
		ServerAddr:      cfg.General.ServerAddr,
		DefaultLanguage: cfg.General.DefaultLanguage,
	})
	if err != nil {
		return err
	}
	return nil
}

// the gRPC servers
func (app *App) startAPIs(ctx context.Context, cfg config.Config) error {
	var err error
	app.devSessionList, err = devprovision.Start(app.psCli, app.nsCli)
	if err != nil {
		return err
	}
	// API for network-server
	if app.nsSrv, err = as.Start(store.NewStore(), cfg.ApplicationServer.API, app.integrations,
		app.psCli, app.nsCli, app.devSessionList); err != nil {
		return err
	}
	// API for mxprotocol server
	if app.mxpSrv, err = mxpapisrv.Start(app.pgstore, cfg.ApplicationServer.APIForM2M, app.mailer); err != nil {
		return err
	}
	// API for external clients
	if app.extAPISrv, err = external.Start(store.NewStore(), external.ExtAPIConfig{
		S:                      cfg.ApplicationServer.ExternalAPI,
		ApplicationServerID:    app.applicationServerID,
		ServerAddr:             cfg.General.ServerAddr,
		BindNewGateway:         cfg.ApplicationServer.APIForGateway.NewGateway.Bind,
		BindOldGateway:         cfg.ApplicationServer.APIForGateway.OldGateway.Bind,
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
		PSCli:                  app.psCli,
		NSCli:                  app.nsCli,
	}); err != nil {
		return err
	}

	app.newGwSrv, app.oldGwSrv, err = gateway.Start(app.pgstore, cfg.General.ServerAddr,
		app.psCli, cfg.ApplicationServer.APIForGateway, cfg.ProvisionServer.UpdateSchedule)
	if err != nil {
		return err
	}

	downlink.Start(store.NewStore(), app.integrations)
	return nil
}
