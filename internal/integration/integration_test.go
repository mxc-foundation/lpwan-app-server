package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/integration"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver/mock"
	httpint "github.com/mxc-foundation/lpwan-app-server/internal/integration/http"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/marshaler"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration/models"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
	"github.com/mxc-foundation/lpwan-app-server/internal/test"
)

type testHTTPHandler struct {
	requests chan *http.Request
}

func (h *testHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(b))
	h.requests <- r
	w.WriteHeader(http.StatusOK)

}

type IntegrationTestSuite struct {
	suite.Suite

	httpServer   *httptest.Server
	httpRequests chan *http.Request
	integration  models.Integration
}

func (ts *IntegrationTestSuite) SetupSuite() {
	assert := require.New(ts.T())
	marshalType = marshaler.Protobuf

	// setup storage
	conf := test.GetConfig()
	assert.NoError(storage.Setup(conf))
	test.MustResetDB(storage.DB().DB)
	rs.RedisClient().FlushAll()

	// http request channel
	ts.httpRequests = make(chan *http.Request, 100)
	ts.httpServer = httptest.NewServer(&testHTTPHandler{
		requests: ts.httpRequests,
	})

	// mock ns client
	networkserver.SetPool(mock.NewPool(mock.NewClient()))

	// setup application with http integration
	ns := storage.NetworkServer{
		Name:   "test-ns",
		Server: "test:1234",
	}
	assert.NoError(storage.CreateNetworkServer(context.Background(), storage.DB(), &ns))

	org := storage.Organization{
		Name: "test-org",
	}
	assert.NoError(storage.CreateOrganization(context.Background(), storage.DB(), &org))

	sp := storage.ServiceProfile{
		Name:            "test-sp",
		OrganizationID:  org.ID,
		NetworkServerID: ns.ID,
	}
	assert.NoError(storage.CreateServiceProfile(context.Background(), storage.DB(), &sp))
	spID, err := uuid.FromBytes(sp.ServiceProfile.Id)
	assert.NoError(err)

	app := storage.Application{
		OrganizationID:   org.ID,
		Name:             "test-app",
		ServiceProfileID: spID,
	}
	assert.NoError(storage.CreateApplication(context.Background(), storage.DB(), &app))

	httpConfig := httpint.Config{
		DataUpURL:               ts.httpServer.URL + "/rx",
		JoinNotificationURL:     ts.httpServer.URL + "/join",
		ACKNotificationURL:      ts.httpServer.URL + "/ack",
		ErrorNotificationURL:    ts.httpServer.URL + "/error",
		StatusNotificationURL:   ts.httpServer.URL + "/status",
		LocationNotificationURL: ts.httpServer.URL + "/location",
		TxAckNotificationURL:    ts.httpServer.URL + "/txack",
	}
	configJSON, err := json.Marshal(httpConfig)
	assert.NoError(err)

	assert.NoError(storage.CreateIntegration(context.Background(), storage.DB(), &storage.Integration{
		ApplicationID: app.ID,
		Kind:          HTTP,
		Settings:      configJSON,
	}))

	ts.integration = ForApplicationID(app.ID)
}

func (ts *IntegrationTestSuite) TearDownSuite() {
	ts.httpServer.Close()
}

// TestHandleUplinkEvent tests that the http integration TestHandleUplinkEvent
// method was called. There is no need to test all methods, as this is already
// done within the multi testsuite.
func (ts *IntegrationTestSuite) TestHandleUplinkEvent() {
	assert := require.New(ts.T())
	assert.NoError(ts.integration.HandleUplinkEvent(context.Background(), nil, pb.UplinkEvent{
		ApplicationId: 1,
		DevEui:        []byte{1, 2, 3, 4, 5, 6, 7, 8},
	}))

	req := <-ts.httpRequests
	assert.Equal("/rx", req.URL.Path)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
