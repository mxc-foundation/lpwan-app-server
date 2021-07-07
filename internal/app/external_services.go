package app

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpccli"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func (app *App) networkServer(ctx context.Context, cfg config.Config) error {
	// get network server list (normally there should be only one network server saved in db)
	nsList, err := app.pgstore.GetNetworkServers(ctx, ns.NetworkServerFilters{Limit: 999, Offset: 0})
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

	// sync region and version from network server
	for _, v := range nsList {
		nsClient, err := app.nsCli.GetNetworkServerServiceClient(v.ID)
		if err != nil {
			return err
		}
		res, err := nsClient.GetVersion(ctx, &empty.Empty{})
		if err != nil {
			return err
		}
		if res.Region.String() != v.Region || res.Version != v.Version {
			err = app.pgstore.UpdateNetworkServerRegionAndVersion(ctx, v.ID, res.Region.String(), res.Version)
			if err != nil {
				return err
			}
		}
	}

	// if no network server stored, set default networkserver, gateway profile
	if 0 == app.nsCli.GetNumberOfNetworkServerClients() {
		// create default network server
		if err := ns.CreateNetworkServer(ctx, &ns.NetworkServer{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      ns.DefaultNetworkServerName,
			Server:    ns.DefaultNetworkServerAddress,
		}, store.NewStore(), app.pgstore, app.nsCli, app.applicationServerID, cfg.ApplicationServer.API.PublicHost); err != nil {
			return err
		}
	}
	return nil
}
