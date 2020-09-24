package gwping

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/brocaar/chirpstack-api/go/v3/as"
	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver/mock"
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/test"
)

type testStore struct {
	store.Store
	inTx           bool
	expectRollback bool
	Organization   map[int64]*store.Organization
	Gateway        map[lorawan.EUI64]*store.Gateway
	GatewayPing    map[int64]*store.GatewayPing
	NetworkServer  map[int64]*store.NetworkServer
	GatewayPingRX  map[int64]*store.GatewayPingRX
}

func (ts *testStore) TxBegin(ctx context.Context) (store.Store, error) {
	ts.inTx = true
	return ts, nil
}

func (ts *testStore) TxCommit(ctx context.Context) error {
	ts.inTx = false
	if ts.expectRollback {
		ts.expectRollback = false
		panic("expected rollback")
	}
	return nil
}

func (ts *testStore) TxRollback(ctx context.Context) error {
	if ts.inTx {
		if ts.expectRollback {
			ts.expectRollback = false
			ts.inTx = false
		} else {
			panic("can't really rollback")
		}
	}
	return nil
}

func (ts *testStore) IsErrorRepeat(err error) bool {
	return false
}

func (ts *testStore) CreateOrganization(ctx context.Context, org *store.Organization) error {
	ts.Organization[org.ID] = org
	return nil
}
func (ts *testStore) CreateNetworkServer(ctx context.Context, n *store.NetworkServer) error {
	ts.NetworkServer[n.ID] = n
	return nil
}
func (ts *testStore) CreateGateway(ctx context.Context, gw *store.Gateway) error {
	ts.Gateway[gw.MAC] = gw
	return nil
}
func (ts *testStore) UpdateNetworkServer(ctx context.Context, n *store.NetworkServer) error {
	ts.NetworkServer[n.ID] = n
	return nil
}
func (ts *testStore) GetGateway(ctx context.Context, mac lorawan.EUI64, forUpdate bool) (store.Gateway, error) {
	return *ts.Gateway[mac], nil
}
func (ts *testStore) GetGatewayPing(ctx context.Context, id int64) (store.GatewayPing, error) {
	return *ts.GatewayPing[id], nil
}
func (ts *testStore) GetLastGatewayPingAndRX(ctx context.Context, mac lorawan.EUI64) (store.GatewayPing, []store.GatewayPingRX, error) {
	var pingID int64
	var pingRXList []store.GatewayPingRX

	if v, ok := ts.Gateway[mac]; !ok {
		panic("no gateway found with mac " + mac.String())
	} else {
		pingID = *v.LastPingID
	}

	for _, v := range ts.GatewayPingRX {
		if v.PingID == pingID {
			pingRXList = append(pingRXList, *v)
		}
	}

	return *ts.GatewayPing[pingID], pingRXList, nil
}
func (ts *testStore) GetGatewayForPing(ctx context.Context) (*store.Gateway, error) {
	for _, v := range ts.NetworkServer {
		if v.GatewayDiscoveryEnabled == false {
			continue
		}

		for _, item := range ts.Gateway {
			if item.NetworkServerID == v.ID && item.Ping == true {
				if item.LastPingSentAt == nil || item.LastPingSentAt.Second() <= (time.Now().Second()-24/v.GatewayDiscoveryInterval) {
					return item, nil
				}
			}
		}
	}

	return nil, nil
}
func (ts *testStore) GetNetworkServer(ctx context.Context, id int64) (store.NetworkServer, error) {
	return *ts.NetworkServer[id], nil
}
func (ts *testStore) CreateGatewayPing(ctx context.Context, ping *store.GatewayPing) error {
	ts.GatewayPing[ping.ID] = ping
	return nil
}
func (ts *testStore) UpdateGateway(ctx context.Context, gw *store.Gateway) error {
	ts.Gateway[gw.MAC] = gw
	return nil
}
func (ts *testStore) CreateGatewayPingRX(ctx context.Context, rx *store.GatewayPingRX) error {
	ts.GatewayPingRX[rx.ID] = rx
	return nil
}

type testEnv struct {
	ctx context.Context
	h   *store.Handler
	t   *testing.T
}

func newTestEnv(t *testing.T) *testEnv {
	te := &testEnv{
		ctx: context.Background(),
		t:   t,
	}

	te.h, _ = store.New(&testStore{
		Organization:  make(map[int64]*store.Organization),
		Gateway:       make(map[lorawan.EUI64]*store.Gateway),
		GatewayPing:   make(map[int64]*store.GatewayPing),
		NetworkServer: make(map[int64]*store.NetworkServer),
		GatewayPingRX: make(map[int64]*store.GatewayPingRX),
	})
	return te
}

func TestGatewayPing(t *testing.T) {
	te := newTestEnv(t)
	Setup(te.h)
	_ = rs.SettingsSetup(rs.RedisStruct{})

	rc := &test.TestRedisClient{
		Data: make(map[string]test.TestRedisResult),
	}
	rs.SetupRedisHandler(rs.NewRedisStore(rc))

	Convey("Given a clean database and a gateway", t, func() {
		nsClient := mock.NewClient()
		networkserver.SetPool(mock.NewPool(nsClient))

		org := store.Organization{
			ID:   1,
			Name: "test-org",
		}
		So(te.h.CreateOrganization(te.ctx, &org), ShouldBeNil)

		n := store.NetworkServer{
			Name:                        "test-ns",
			Server:                      "test-ns:1234",
			GatewayDiscoveryEnabled:     true,
			GatewayDiscoveryDR:          5,
			GatewayDiscoveryTXFrequency: 868100000,
			GatewayDiscoveryInterval:    1,
		}
		So(te.h.CreateNetworkServer(te.ctx, &n), ShouldBeNil)

		gw1 := store.Gateway{
			MAC:             lorawan.EUI64{1, 2, 3, 4, 5, 6, 7, 8},
			Name:            "test-gw",
			Description:     "test gateway",
			OrganizationID:  org.ID,
			Ping:            true,
			NetworkServerID: n.ID,
		}
		So(te.h.CreateGateway(te.ctx, &gw1), ShouldBeNil)

		Convey("When gateway discovery is disabled on the network-server", func() {
			n.GatewayDiscoveryEnabled = false
			So(te.h.UpdateNetworkServer(te.ctx, &n), ShouldBeNil)

			Convey("When calling sendGatewayPing", func() {
				So(sendGatewayPing(te.ctx, te.h), ShouldBeNil)
			})

			Convey("Then no ping was sent", func() {
				gwGet, err := te.h.GetGateway(te.ctx, gw1.MAC, false)
				So(err, ShouldBeNil)
				So(gwGet.LastPingID, ShouldBeNil)
				So(gwGet.LastPingSentAt, ShouldBeNil)
			})
		})

		Convey("When calling sendGatewayPing", func() {
			So(sendGatewayPing(te.ctx, te.h), ShouldBeNil)

			Convey("Then the gateway ping fields have been set", func() {
				gwGet, err := te.h.GetGateway(te.ctx, gw1.MAC, false)
				So(err, ShouldBeNil)
				So(gwGet.LastPingID, ShouldNotBeNil)
				So(gwGet.LastPingSentAt, ShouldNotBeNil)

				Convey("Then a gateway ping records has been created", func() {
					gwPing, err := te.h.GetGatewayPing(te.ctx, *gwGet.LastPingID)
					So(err, ShouldBeNil)
					So(gwPing.GatewayMAC, ShouldEqual, gwGet.MAC)
					So(gwPing.DR, ShouldEqual, n.GatewayDiscoveryDR)
					So(gwPing.Frequency, ShouldEqual, n.GatewayDiscoveryTXFrequency)
				})

				Convey("Then the expected ping has been sent to the network-server", func() {
					So(nsClient.SendProprietaryPayloadChan, ShouldHaveLength, 1)
					req := <-nsClient.SendProprietaryPayloadChan
					So(req.Dr, ShouldEqual, uint32(n.GatewayDiscoveryDR))
					So(req.Frequency, ShouldEqual, uint32(n.GatewayDiscoveryTXFrequency))
					So(req.GatewayMacs, ShouldResemble, [][]byte{{1, 2, 3, 4, 5, 6, 7, 8}})
					So(req.PolarizationInversion, ShouldBeFalse)

					var mic lorawan.MIC
					copy(mic[:], req.Mic)
					So(mic, ShouldNotEqual, lorawan.MIC{})

					Convey("Then a ping lookup has been created", func() {
						id, err := GetPingLookup(mic)
						So(err, ShouldBeNil)
						So(id, ShouldEqual, *gwGet.LastPingID)
					})

					Convey("When calling HandleReceivedPing", func() {
						gw2 := store.Gateway{
							MAC:             lorawan.EUI64{8, 7, 6, 5, 4, 3, 2, 1},
							Name:            "test-gw-2",
							Description:     "test gateway 2",
							OrganizationID:  org.ID,
							NetworkServerID: n.ID,
						}
						So(te.h.CreateGateway(te.ctx, &gw2), ShouldBeNil)

						now := time.Now().UTC().Truncate(time.Millisecond)

						pong := as.HandleProprietaryUplinkRequest{
							Mic: mic[:],
							RxInfo: []*gw.UplinkRXInfo{
								{
									GatewayId: gw2.MAC[:],
									Rssi:      -10,
									LoraSnr:   5.5,
									Location: &common.Location{
										Latitude:  1.12345,
										Longitude: 1.23456,
										Altitude:  10,
									},
								},
							},
						}
						pong.RxInfo[0].Time, _ = ptypes.TimestampProto(now)
						So(HandleReceivedPing(te.ctx, &pong), ShouldBeNil)

						Convey("Then the ping lookup has been deleted", func() {
							_, err := GetPingLookup(mic)
							So(err, ShouldNotBeNil)
						})

						Convey("Then the received ping has been stored to the database", func() {
							ping, rx, err := te.h.GetLastGatewayPingAndRX(te.ctx, gw1.MAC)
							So(err, ShouldBeNil)

							So(ping.ID, ShouldEqual, *gwGet.LastPingID)
							So(rx, ShouldHaveLength, 1)
							So(rx[0].GatewayMAC, ShouldEqual, gw2.MAC)
							So(rx[0].ReceivedAt.Equal(now), ShouldBeTrue)
							So(rx[0].RSSI, ShouldEqual, -10)
							So(rx[0].LoRaSNR, ShouldEqual, 5.5)
							So(rx[0].Location, ShouldResemble, store.GPSPoint{
								Latitude:  1.12345,
								Longitude: 1.23456,
							})
							So(rx[0].Altitude, ShouldEqual, 10)
						})
					})
				})
			})
		})
	})
}
