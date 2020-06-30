package postgresql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-server/api/ns"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func SetupDefault() error {
	ctx := context.Background()
	count, err := storage.GetGatewayProfileCount(ctx, storage.DB())
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Failed to load gateway profiles")
	}

	if count != 0 {
		// check if default gateway profile already exists
		gpList, err := storage.GetGatewayProfiles(ctx, storage.DB(), count, 0)
		if err != nil {
			return errors.Wrap(err, "Failed to load gateway profiles")
		}

		for _, v := range gpList {
			if v.Name == "default_gateway_profile" {
				return nil
			}
		}
	}

	// none default_gateway_profile exists, add one
	var networkServer storage.NetworkServer
	n, err := storage.GetNetworkServers(ctx, storage.DB(), 1, 0)
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Load network server internal error")
	}

	if len(n) >= 1 {
		networkServer = n[0]
	} else {
		// insert default one
		err := storage.Transaction(func(tx sqlx.Ext) error {
			return storage.CreateNetworkServer(ctx, storage.DB(), &storage.NetworkServer{
				Name:                    "default_network_server",
				Server:                  "network-server:8000",
				GatewayDiscoveryEnabled: false,
			})
		})
		if err != nil {
			return nil
		}

		// get network-server id

		networkServer, err = storage.GetDefaultNetworkServer(ctx, storage.DB())
		if err != nil {
			return err
		}
	}

	gp := storage.GatewayProfile{
		NetworkServerID: networkServer.ID,
		Name:            "default_gateway_profile",
		GatewayProfile: ns.GatewayProfile{
			Channels:      []uint32{0, 1, 2},
			ExtraChannels: []*ns.GatewayProfileExtraChannel{},
		},
	}

	err = storage.Transaction(func(tx sqlx.Ext) error {
		return storage.CreateGatewayProfile(ctx, tx, &gp)
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create default gateway profile")
	}

	return nil
}
