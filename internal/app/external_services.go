package app

import (
	"context"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	setdefault "github.com/mxc-foundation/lpwan-app-server/internal/modules/set_default"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func (app *App) networkServer(ctx context.Context, cfg config.Config) error {
	// get network server list (normally there should be only one network server saved in db)
	nsList, err := app.pgstore.GetNetworkServers(ctx, nsd.NetworkServerFilters{Limit: 999, Offset: 0})
	if err != nil {
		return err
	}
	var nscfg []nscli.NetworkServerConfig
	for _, v := range nsList {
		nscfg = append(nscfg, nscli.NetworkServerConfig{
			NetworkServerID: v.ID,
			ConnOptions: grpccli.ConnectionOpts{
				Server:  v.Server,
				CACert:  v.CACert,
				TLSCert: v.TLSCert,
				TLSKey:  v.TLSKey,
			},
		})
	}
	app.nsCli = &nscli.Client{}
	if err := app.nsCli.Connect(nscfg); err != nil {
		return err
	}
	// set default networkserver, gateway profile
	err = setdefault.Setup(ctx, store.NewStore(), app.applicationServerID,
		cfg.ApplicationServer.API.PublicHost, app.nsCli)
	if err != nil {
		return err
	}
	return nil
}
