package setdefault

import (
	"context"

	"github.com/brocaar/lorawan/band"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	gwp "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	gws "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile/data"
	nsmod "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "set_default"

type controller struct {
	st                          *store.Handler
	applicationServerID         uuid.UUID
	applicationServerPublicHost string

	moduleUp bool
}

var ctrl *controller

// SettingsSetup :
func SettingsSetup(name string, conf config.Config) error {

	appServerID, err := uuid.FromString(conf.ApplicationServer.ID)
	if err != nil {
		return errors.Wrap(err, "failed to convert applicationserver id from string to uuid")
	}

	ctrl = &controller{
		applicationServerID:         appServerID,
		applicationServerPublicHost: conf.ApplicationServer.API.PublicHost,
	}

	return nil
}

// Setup :
func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}
	defer func() {
		ctrl.moduleUp = true
	}()

	ctrl.st = h
	ctx := context.Background()

	has, err := hasActiveNetworkServer(ctx)
	if err != nil {
		return errors.Wrap(err, "check active network servers error")
	}

	if !has {
		if err := createDefaultNetworkServer(ctx); err != nil {
			return errors.Wrap(err, "create default network server error")
		}
	}

	if err := createDefaultGatewayProfileForNetworkserver(ctx); err != nil {
		return errors.Wrap(err, "create default gateway profile for network server error")
	}

	return nil
}

func createDefaultGatewayProfileForNetworkserver(ctx context.Context) error {
	count, err := ctrl.st.GetNetworkServerCount(ctx, nsd.NetworkServerFilters{
		OrganizationID: 0,
	})
	if err != nil {
		return errors.Wrap(err, "get network server count error")
	}

	if count == 0 {
		return errors.New("cannot create default gateway profile, no network server set")
	}

	limit := 10
	for offset := 0; offset <= count/limit; offset += limit {
		nsList, err := ctrl.st.GetNetworkServers(ctx, nsd.NetworkServerFilters{
			OrganizationID: 0,
			Limit:          limit,
			Offset:         offset,
		})
		if err != nil {
			return errors.Wrap(err, "get network servers error")
		}

		for _, v := range nsList {
			has, err := hasActiveGatewayProfileForNetworkServer(ctx, &v)
			if err != nil {
				return errors.Wrap(err, "check active gateway profile for network server error")
			}

			if !has {
				// create gateway profile based on network server region
				var gatewayProfile ns.GatewayProfile

				if v.Region == string(band.AS923) {
					gatewayProfile = ns.GatewayProfile{
						Channels:      []uint32{0, 1},
						ExtraChannels: []*ns.GatewayProfileExtraChannel{},
					}
				} else {
					// all other region
					gatewayProfile = ns.GatewayProfile{
						Channels:      []uint32{0, 1, 2},
						ExtraChannels: []*ns.GatewayProfileExtraChannel{},
					}
				}

				if err := gwp.CreateGatewayProfile(ctx, &gws.GatewayProfile{
					NetworkServerID: v.ID,
					Name:            "default_gateway_profile",
					GatewayProfile:  gatewayProfile,
				}); err != nil {
					return errors.Wrap(err, "create gateway profile error")
				}
			}

		}
	}

	return nil

}

func hasActiveGatewayProfileForNetworkServer(ctx context.Context, nServer *nsd.NetworkServer) (bool, error) {
	count, err := ctrl.st.GetGatewayProfileCountForNetworkServerID(ctx, nServer.ID)
	if err != nil {
		return false, errors.Wrap(err, "get gateway profile count for network server id error")
	}

	if count == 0 {
		return false, nil
	}

	limit := 10
	for offset := 0; offset <= count/limit; offset += limit {
		gwpList, err := ctrl.st.GetGatewayProfilesForNetworkServerID(ctx, nServer.ID, limit, offset)
		if err != nil {
			return false, errors.Wrap(err, "get gateway profile for network server id error")
		}

		for _, v := range gwpList {
			client, err := (&nsmod.NSStruct{
				Server:  nServer.Server,
				CACert:  nServer.CACert,
				TLSCert: nServer.TLSCert,
				TLSKey:  nServer.TLSKey,
			}).GetNetworkServiceClient()
			if err != nil {
				return false, errors.Wrap(err, "get network server client error")
			}

			resp, err := client.GetGatewayProfile(ctx, &ns.GetGatewayProfileRequest{
				Id: v.GatewayProfileID.Bytes(),
			})
			if err == nil {
				if resp.GatewayProfile != nil {
					continue
				}
			}

			var gatewayProfile ns.GatewayProfile
			if nServer.Region == string(band.AS923) {
				gatewayProfile = ns.GatewayProfile{
					Id:            v.GatewayProfileID.Bytes(),
					Channels:      []uint32{0, 1},
					ExtraChannels: []*ns.GatewayProfileExtraChannel{},
				}
			} else {
				// all other region
				gatewayProfile = ns.GatewayProfile{
					Id:            v.GatewayProfileID.Bytes(),
					Channels:      []uint32{0, 1, 2},
					ExtraChannels: []*ns.GatewayProfileExtraChannel{},
				}
			}
			// create gateway profile in network server
			_, err = client.CreateGatewayProfile(ctx, &ns.CreateGatewayProfileRequest{
				GatewayProfile: &gatewayProfile,
			})
			if err != nil {
				return false, errors.Wrap(err, "create gateway profile in network server error")
			}

		}

	}

	return true, nil
}

func hasActiveNetworkServer(ctx context.Context) (bool, error) {
	count, err := ctrl.st.GetNetworkServerCount(ctx, nsd.NetworkServerFilters{
		OrganizationID: 0,
	})
	if err != nil {
		return false, errors.Wrap(err, "get network server count error")
	}

	if count == 0 {
		return false, nil
	}

	limit := 10
	for offset := 0; offset <= count/limit; offset += limit {
		nsList, err := ctrl.st.GetNetworkServers(ctx, nsd.NetworkServerFilters{
			OrganizationID: 0,
			Limit:          limit,
			Offset:         offset,
		})
		if err != nil {
			return false, errors.Wrap(err, "get network servers error")
		}

		for _, v := range nsList {
			// check if networkserver profile is created in network server
			client, err := (&nsmod.NSStruct{
				Server:  v.Server,
				CACert:  v.CACert,
				TLSCert: v.TLSCert,
				TLSKey:  v.TLSKey,
			}).GetNetworkServiceClient()
			if err != nil {
				return false, errors.Wrap(err, "get network server client error")
			}

			// make sure every local network server profile has routing profile created in cooresponding network server
			if _, err := client.GetRoutingProfile(ctx, &ns.GetRoutingProfileRequest{
				Id: ctrl.applicationServerID.Bytes(),
			}); err != nil {
				_, err = client.CreateRoutingProfile(ctx, &ns.CreateRoutingProfileRequest{
					RoutingProfile: &ns.RoutingProfile{
						Id:      ctrl.applicationServerID.Bytes(),
						AsId:    ctrl.applicationServerPublicHost,
						CaCert:  v.RoutingProfileCACert,
						TlsCert: v.RoutingProfileTLSCert,
						TlsKey:  v.RoutingProfileTLSKey,
					},
				})
				if err != nil {
					return false, errors.Wrap(err, "create routing-profile error")
				}
			}

		}
	}

	return true, nil
}

func createDefaultNetworkServer(ctx context.Context) error {
	nsServer := "network-server:8000"
	client, err := (&nsmod.NSStruct{
		Server:  nsServer,
		CACert:  "",
		TLSCert: "",
		TLSKey:  "",
	}).GetNetworkServiceClient()
	if err != nil {
		return errors.Wrap(err, "get network server client error")
	}

	res, err := client.GetVersion(ctx, &empty.Empty{})
	if err != nil {
		return errors.Wrap(err, "get version for default network server (network-server:8000) error")
	}

	region := res.Region.String()
	version := res.Version

	if err := nsmod.CreateNetworkServer(ctx, &nsd.NetworkServer{
		Name:                        "default_network_server",
		Server:                      nsServer,
		CACert:                      "",
		TLSCert:                     "",
		TLSKey:                      "",
		RoutingProfileCACert:        "",
		RoutingProfileTLSCert:       "",
		RoutingProfileTLSKey:        "",
		GatewayDiscoveryEnabled:     false,
		GatewayDiscoveryInterval:    0,
		GatewayDiscoveryTXFrequency: 0,
		GatewayDiscoveryDR:          0,
		Region:                      region,
		Version:                     version,
	}, ctrl.st); err != nil {
		return errors.Wrap(err, "create network server error")
	}

	return nil
}
