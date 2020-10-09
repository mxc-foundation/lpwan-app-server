package external

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	. "github.com/mxc-foundation/lpwan-app-server/internal/api/external/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	authPg "github.com/mxc-foundation/lpwan-app-server/internal/authentication/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	user "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	otpst "github.com/mxc-foundation/lpwan-app-server/internal/otp/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "external"

type controller struct {
	name                string
	s                   ExternalAPIStruct
	applicationServerID uuid.UUID
	serverAddr          string
	recaptcha           user.RecaptchaStruct
	enable2FA           bool
	serverRegion        string
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) (err error) {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name:         moduleName,
		s:            conf.ApplicationServer.ExternalAPI,
		serverAddr:   conf.General.ServerAddr,
		recaptcha:    conf.Recaptcha,
		enable2FA:    conf.General.Enable2FALogin,
		serverRegion: conf.General.ServerRegion,
	}
	ctrl.applicationServerID, err = uuid.FromString(conf.ApplicationServer.ID)
	if err != nil {
		return errors.Wrap(err, "application-server id to uuid error")
	}

	return nil
}
func GetApplicationServerID() uuid.UUID {
	return ctrl.applicationServerID
}
func GetJWTSecret() string {
	return ctrl.s.JWTSecret
}
func GetOTPSecret() string {
	return ctrl.s.OTPSecret
}

// Setup configures the API endpoints.
func Setup(name string, h *store.Handler) (err error) {
	if name != moduleName {
		return errors.New(fmt.Sprintf("intend to call Setup function for %s, actually calling %s", name, moduleName))
	}

	// Bind external api port to listen to requests to all services
	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)

	if err := SetupCusAPI(h, grpcServer, GetApplicationServerID()); err != nil {
		return err
	}

	// setup the client http interface variable
	// we need to start the gRPC service first, as it is used by the
	// grpc-gateway
	var clientHTTPHandler http.Handler

	// switch between gRPC and "plain" http handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if clientHTTPHandler == nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			if ctrl.s.CORSAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", ctrl.s.CORSAllowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Grpc-Metadata-Authorization")

				if r.Method == "OPTIONS" {
					return
				}
			}

			clientHTTPHandler.ServeHTTP(w, r)
		}
	})

	// start the API server
	go func() {
		log.WithFields(log.Fields{
			"bind":     ctrl.s.Bind,
			"tls-cert": ctrl.s.TLSCert,
			"tls-key":  ctrl.s.TLSKey,
		}).Info("api/external: starting api server")

		if ctrl.s.TLSCert == "" || ctrl.s.TLSKey == "" {
			log.Fatal(http.ListenAndServe(ctrl.s.Bind, h2c.NewHandler(handler, &http2.Server{})))
		} else {
			log.Fatal(http.ListenAndServeTLS(
				ctrl.s.Bind,
				ctrl.s.TLSCert,
				ctrl.s.TLSKey,
				h2c.NewHandler(handler, &http2.Server{}),
			))
		}
	}()

	// give the http server some time to start
	time.Sleep(time.Millisecond * 100)

	// setup the HTTP handler
	clientHTTPHandler, err = setupHTTPAPI()
	if err != nil {
		return err
	}

	return nil
}

func SetupCusAPI(h *store.Handler, grpcServer *grpc.Server, rpID uuid.UUID) error {
	jwtSecret := GetJWTSecret()
	if jwtSecret == "" {
		return errors.New("jwt_secret must be set")
	}

	jwtValidator := jwt.NewJWTValidator("HS256", []byte(jwtSecret))
	otpValidator, err := otp.NewValidator("lpwan-app-server", GetOTPSecret(), otpst.NewStore(pgstore.New()))
	if err != nil {
		return err
	}
	authcus.SetupCred(authPg.NewStore(pgstore.New()), jwtValidator, otpValidator)

	pb.RegisterFUOTADeploymentServiceServer(grpcServer, NewFUOTADeploymentAPI(h))
	pb.RegisterDeviceQueueServiceServer(grpcServer, NewDeviceQueueAPI(h))
	pb.RegisterMulticastGroupServiceServer(grpcServer, NewMulticastGroupAPI(rpID, h))
	pb.RegisterServiceProfileServiceServer(grpcServer, NewServiceProfileServiceAPI(h))
	pb.RegisterDeviceProfileServiceServer(grpcServer, NewDeviceProfileServiceAPI(h))
	// device
	api.RegisterDeviceServiceServer(grpcServer, NewDeviceAPI(rpID, h))
	// gateway
	api.RegisterGatewayServiceServer(grpcServer, NewGatewayAPI(rpID, h, ctrl.serverAddr))
	// gateway profile
	api.RegisterGatewayProfileServiceServer(grpcServer, NewGatewayProfileAPI(h))
	// application
	api.RegisterApplicationServiceServer(grpcServer, NewApplicationAPI(h))
	// network server
	api.RegisterNetworkServerServiceServer(grpcServer, NewNetworkServerAPI(h))
	// orgnization
	api.RegisterOrganizationServiceServer(grpcServer, NewOrganizationAPI(h))
	// user
	api.RegisterUserServiceServer(grpcServer, NewUserAPI(h))
	api.RegisterInternalServiceServer(grpcServer, NewInternalUserAPI(h, user.Config{
		Recaptcha:      ctrl.recaptcha,
		Enable2FALogin: ctrl.enable2FA,
	}))

	api.RegisterServerInfoServiceServer(grpcServer, NewServerInfoAPI(ctrl.serverRegion))
	api.RegisterSettingsServiceServer(grpcServer, NewSettingsServerAPI())
	api.RegisterStakingServiceServer(grpcServer, NewStakingServerAPI())
	api.RegisterTopUpServiceServer(grpcServer, NewTopUpServerAPI())
	api.RegisterWalletServiceServer(grpcServer, NewWalletServerAPI())
	api.RegisterWithdrawServiceServer(grpcServer, NewWithdrawServerAPI())

	return nil
}

func setupHTTPAPI() (http.Handler, error) {
	r := mux.NewRouter()

	// setup json api handler
	jsonHandler, err := getJSONGateway(context.Background())
	if err != nil {
		return nil, err
	}

	log.WithField("path", "/api").Info("api/external: registering rest api handler and documentation endpoint")
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.Asset("swagger/index.html")
		if err != nil {
			log.WithError(err).Error("get swagger template error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	}).Methods("get")
	r.PathPrefix("/api").Handler(jsonHandler)

	if err := oidc.Setup(r); err != nil {
		return nil, errors.Wrap(err, "setup openid connect error")
	}

	// setup static file server
	r.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "",
	}))

	return wsproxy.WebsocketProxy(r), nil
}

func getJSONGateway(ctx context.Context) (http.Handler, error) {
	// dial options for the grpc-gateway
	var grpcDialOpts []grpc.DialOption

	if ctrl.s.TLSCert == "" || ctrl.s.TLSKey == "" {
		grpcDialOpts = append(grpcDialOpts, grpc.WithInsecure())
	} else {
		b, err := ioutil.ReadFile(ctrl.s.TLSCert)
		if err != nil {
			return nil, errors.Wrap(err, "read external api tls cert error")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, errors.Wrap(err, "failed to append certificate")
		}
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			// given the grpc-gateway is always connecting to localhost, does
			// InsecureSkipVerify=true cause any security issues?
			InsecureSkipVerify: true,
			RootCAs:            cp,
		})))
	}

	bindParts := strings.SplitN(ctrl.s.Bind, ":", 2)
	if len(bindParts) != 2 {
		log.Fatal("get port from bind failed")
	}
	apiEndpoint := fmt.Sprintf("localhost:%s", bindParts[1])

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			EnumsAsInts:  false,
			EmitDefaults: true,
		},
	))

	if err := pb.RegisterDeviceQueueServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register downlink queue handler error")
	}
	if err := pb.RegisterServiceProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register service-profile handler error")
	}
	if err := pb.RegisterDeviceProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register device-profile handler error")
	}
	if err := pb.RegisterMulticastGroupServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register multicast-group handler error")
	}
	if err := pb.RegisterFUOTADeploymentServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register fuota deployment handler error")
	}

	if err := api.RegisterServerInfoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register server info handler error")
	}
	if err := api.RegisterStakingServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterTopUpServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterWalletServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterWithdrawServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterSettingsServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register proxy request handler error")
	}

	if err := api.RegisterApplicationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register application handler error")
	}
	if err := api.RegisterDeviceServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register node handler error")
	}
	if err := api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register user handler error")
	}
	if err := api.RegisterInternalServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register internal handler error")
	}
	if err := api.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register gateway handler error")
	}
	if err := api.RegisterGatewayProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register gateway-profile handler error")
	}
	if err := api.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register organization handler error")
	}
	if err := api.RegisterNetworkServerServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register network-server handler error")
	}

	return mux, nil
}
