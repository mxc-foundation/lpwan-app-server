package external

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

	api "github.com/mxc-foundation/lpwan-app-server/api/extapi"
	. "github.com/mxc-foundation/lpwan-app-server/internal/api/external/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dfi"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dhx"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/mqttauth"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/report"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/staking"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpcauth"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// ExtAPIServer represents gRPC server serving external api
type ExtAPIServer struct {
	gs *grpc.Server
}

// Config defines all attributes for ext api service
type Config struct {
	S                      ExternalAPIStruct
	ApplicationServerID    string
	ServerAddr             string
	Recaptcha              user.RecaptchaConfig
	Enable2FA              bool
	ServerRegion           string
	PasswordHashIterations int
	EnableSTC              bool
	ExternalAuth           user.ExternalAuthentication
	ShopifyConfig          user.Shopify
	OperatorLogo           string
	Mailer                 *email.Mailer
	MXPCli                 *mxpcli.Client
}

// Stop gracefully stops gRPC server
func (srv *ExtAPIServer) Stop() {
	srv.gs.GracefulStop()
}

// Start configures the API endpoints.
func Start(h *store.Handler, conf Config) (*ExtAPIServer, error) {
	var err error
	// Bind external api port to listen to requests to all services
	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)

	srv := &ExtAPIServer{
		gs: grpcServer,
	}
	if err := srv.SetupCusAPI(h, conf); err != nil {
		return nil, err
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

			if conf.S.CORSAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", conf.S.CORSAllowOrigin)
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
			"bind":     conf.S.Bind,
			"tls-cert": conf.S.TLSCert,
			"tls-key":  conf.S.TLSKey,
		}).Info("api/external: starting api server")

		if conf.S.TLSCert == "" || conf.S.TLSKey == "" {
			log.Fatal(http.ListenAndServe(conf.S.Bind, h2c.NewHandler(handler, &http2.Server{})))
		} else {
			log.Fatal(http.ListenAndServeTLS(
				conf.S.Bind,
				conf.S.TLSCert,
				conf.S.TLSKey,
				h2c.NewHandler(handler, &http2.Server{}),
			))
		}
	}()

	// give the http server some time to start
	time.Sleep(time.Millisecond * 100)

	// setup the HTTP handler
	clientHTTPHandler, err = srv.setupHTTPAPI(conf)
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// SetupCusAPI registers all ext api services
func (srv *ExtAPIServer) SetupCusAPI(h *store.Handler, conf Config) error {
	jwtSecret := conf.S.JWTSecret
	if jwtSecret == "" {
		return errors.New("jwt_secret must be set")
	}
	jwtTTL := conf.S.JWTDefaultTTL
	pgs := pgstore.New()

	jwtValidator := jwt.NewValidator("HS256", []byte(jwtSecret), jwtTTL)
	otpValidator, err := otp.NewValidator("lpwan-app-server", conf.S.OTPSecret, pgs)
	if err != nil {
		return err
	}
	grpcAuth := grpcauth.New(pgs, jwtValidator, otpValidator)
	authcus.SetupCred(pgs, jwtValidator, otpValidator)

	rpID, err := uuid.FromString(conf.ApplicationServerID)
	if err != nil {
		return fmt.Errorf("failed to convert application server id from string to uuid: %v", err)
	}

	pb.RegisterFUOTADeploymentServiceServer(srv.gs, NewFUOTADeploymentAPI(h))
	pb.RegisterDeviceQueueServiceServer(srv.gs, NewDeviceQueueAPI(h))
	pb.RegisterMulticastGroupServiceServer(srv.gs, NewMulticastGroupAPI(rpID, h))
	pb.RegisterServiceProfileServiceServer(srv.gs, NewServiceProfileServiceAPI(h))
	pb.RegisterDeviceProfileServiceServer(srv.gs, NewDeviceProfileServiceAPI(h))
	// device
	api.RegisterDeviceServiceServer(srv.gs, NewDeviceAPI(rpID, h))
	// gateway
	psCli, err := pscli.GetPServerClient()
	if err != nil {
		return err
	}
	api.RegisterGatewayServiceServer(srv.gs, NewGatewayAPI(
		h.PgStore,
		grpcAuth,
		GwConfig{
			ApplicationServerID: rpID,
			ServerAddr:          conf.ServerAddr,
			EnableSTC:           conf.EnableSTC,
		},
		psCli,
	))

	// gateway profile
	api.RegisterGatewayProfileServiceServer(srv.gs, NewGatewayProfileAPI(h))
	// application
	api.RegisterApplicationServiceServer(srv.gs, NewApplicationAPI(h))
	// network server
	api.RegisterNetworkServerServiceServer(srv.gs, NewNetworkServerAPI(h))
	// orgnization
	api.RegisterOrganizationServiceServer(srv.gs, NewOrganizationAPI(h))
	// user
	pwhasher, err := pwhash.New(16, conf.PasswordHashIterations)
	if err != nil {
		return err
	}
	userSrv := user.NewServer(
		pgs,
		conf.Mailer,
		grpcAuth,
		jwtValidator,
		otpValidator,
		pwhasher,
		user.Config{
			Recaptcha:        conf.Recaptcha,
			Enable2FALogin:   conf.Enable2FA,
			OperatorLogoPath: conf.OperatorLogo,
			WeChatLogin:      conf.ExternalAuth.WechatAuth,
			DebugWeChatLogin: conf.ExternalAuth.DebugWechatAuth,
			ShopifyConfig:    conf.ShopifyConfig,
		},
	)
	api.RegisterUserServiceServer(srv.gs, userSrv)
	api.RegisterInternalServiceServer(srv.gs, userSrv)
	api.RegisterExternalUserServiceServer(srv.gs, userSrv)

	api.RegisterServerInfoServiceServer(srv.gs, NewServerInfoAPI(conf.ServerRegion))
	api.RegisterSettingsServiceServer(srv.gs, NewSettingsServerAPI())
	api.RegisterTopUpServiceServer(srv.gs, NewTopUpServerAPI(grpcAuth))

	api.RegisterWalletServiceServer(srv.gs, NewWalletServerAPI(
		h,
		grpcAuth,
		conf.EnableSTC,
	))

	api.RegisterWithdrawServiceServer(srv.gs, NewWithdrawServerAPI(grpcAuth))

	api.RegisterStakingServiceServer(srv.gs, staking.NewServer(
		mxpcli.Global.GetStakingServiceClient(),
		grpcAuth,
	))

	api.RegisterDHXServcieServer(srv.gs, dhx.NewServer(
		mxpcli.Global.GetDHXServiceClient(),
		grpcAuth,
		pgs,
	))

	api.RegisterShopifyIntegrationServer(srv.gs, user.NewShopifyServiceServer(
		grpcAuth,
		pgs,
	))

	api.RegisterDFIServiceServer(srv.gs, dfi.NewServer(
		pgs,
	))

	api.RegisterReportServiceServer(srv.gs, report.NewServer(
		conf.MXPCli.GetFianceReportClient(),
		grpcAuth,
		conf.ServerAddr,
	))

	eventTopicRegexp, err := regexp.Compile(mqttauth.EventTopicTemplate)
	if err != nil {
		return fmt.Errorf("compile regexp error: %v", err)
	}
	commandTopicRegexp, err := regexp.Compile(mqttauth.CommandTopicTemplate)
	if err != nil {
		return fmt.Errorf("compile regexp error: %v", err)
	}
	allEventsTopicRegexp, err := regexp.Compile(mqttauth.AllEventsTopicTemplate)
	if err != nil {
		return fmt.Errorf("compile regexp error: %v", err)
	}
	api.RegisterMosquittoAuthServiceServer(srv.gs, mqttauth.NewServer(
		pgs,
		grpcAuth,
		jwtValidator,
		eventTopicRegexp,
		commandTopicRegexp,
		allEventsTopicRegexp,
	))
	return nil
}

func (srv *ExtAPIServer) setupHTTPAPI(conf Config) (http.Handler, error) {
	r := mux.NewRouter()

	// setup json api handler
	jsonHandler, err := srv.getJSONGateway(context.Background(), conf)
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

func (srv *ExtAPIServer) getJSONGateway(ctx context.Context, conf Config) (http.Handler, error) {
	// dial options for the grpc-gateway
	var grpcDialOpts []grpc.DialOption

	if conf.S.TLSCert == "" || conf.S.TLSKey == "" {
		grpcDialOpts = append(grpcDialOpts, grpc.WithInsecure())
	} else {
		b, err := ioutil.ReadFile(conf.S.TLSCert)
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

	bindParts := strings.SplitN(conf.S.Bind, ":", 2)
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

	err := pb.RegisterDeviceQueueServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register downlink queue handler: %v", err)

	err = pb.RegisterServiceProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register service-profile handler: %v", err)

	err = pb.RegisterDeviceProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register device-profile handler: %v", err)

	err = pb.RegisterMulticastGroupServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register multicast-group handler: %v", err)

	err = pb.RegisterFUOTADeploymentServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register fuota deployment handler: %v", err)

	err = api.RegisterServerInfoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register server info handler: %v", err)

	err = api.RegisterStakingServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register staking service handler: %v", err)

	err = api.RegisterTopUpServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register top up service handler: %v", err)

	err = api.RegisterWalletServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register wallet service  handler: %v", err)

	err = api.RegisterWithdrawServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register withdraw service  handler: %v", err)

	err = api.RegisterSettingsServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register settings service handler: %v", err)

	err = api.RegisterApplicationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register application service handler: %v", err)

	err = api.RegisterDeviceServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register device service handler: %v", err)

	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register user service handler: %v", err)

	err = api.RegisterInternalServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register internal service handler: %v", err)

	err = api.RegisterExternalUserServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register external user service handler: %v", err)

	err = api.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register gateway service handler: %v", err)

	err = api.RegisterGatewayProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register gateway profile service handler: %v", err)

	err = api.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register organization service handler: %v", err)

	err = api.RegisterNetworkServerServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register network server service handler: %v", err)

	err = api.RegisterDHXServcieHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register dhx service handler: %v", err)

	err = api.RegisterShopifyIntegrationHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register shopify integration service handler: %v", err)

	err = api.RegisterDFIServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register dfi service handler: %v", err)

	err = api.RegisterReportServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register download service handler: %v", err)

	err = api.RegisterMosquittoAuthServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register mosquitto auth service handler: %v", err)

	return mux, nil
}
