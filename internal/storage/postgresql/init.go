package postgresql

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-server/api/ns"
)

func SetupDefault() error {
	ctx := context.Background()
	count, err := GetGatewayProfileCount(ctx, DB())
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Failed to load gateway profiles")
	}

	if count != 0 {
		// check if default gateway profile already exists
		gpList, err := GetGatewayProfiles(ctx, DB(), count, 0)
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
	var networkServer NetworkServer
	n, err := GetNetworkServers(ctx, DB(), 1, 0)
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Load network server internal error")
	}

	if len(n) >= 1 {
		networkServer = n[0]
	} else {
		// insert default one
		err := Transaction(func(tx sqlx.Ext) error {
			return CreateNetworkServer(ctx, DB(), &NetworkServer{
				Name:                    "default_network_server",
				Server:                  "network-server:8000",
				GatewayDiscoveryEnabled: false,
			})
		})
		if err != nil {
			return nil
		}

		// get network-server id
		networkServer, err = GetDefaultNetworkServer(ctx, DB())
		if err != nil {
			return err
		}
	}

	gp := GatewayProfile{
		NetworkServerID: networkServer.ID,
		Name:            "default_gateway_profile",
		GatewayProfile: ns.GatewayProfile{
			Channels:      []uint32{0, 1, 2},
			ExtraChannels: []*ns.GatewayProfileExtraChannel{},
		},
	}

	err = Transaction(func(tx sqlx.Ext) error {
		return CreateGatewayProfile(ctx, tx, &gp)
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create default gateway profile")
	}

	return nil
}
