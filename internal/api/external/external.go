package external

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/download"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dfi"

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
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dhx"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/staking"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/grpcauth"
	"github.com/mxc-foundation/lpwan-app-server/internal/mxpcli"
	"github.com/mxc-foundation/lpwan-app-server/internal/oidc"

	pscli "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	"github.com/mxc-foundation/lpwan-app-server/internal/pwhash"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// RESTApiServer defines all attributes for REST api service
type RESTApiServer struct {
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

// Start configures the API endpoints.
func Start(h *store.Handler, srv RESTApiServer) (err error) {
	// Bind external api port to listen to requests to all services
	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)

	if err := srv.SetupCusAPI(h, grpcServer); err != nil {
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

			if srv.S.CORSAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", srv.S.CORSAllowOrigin)
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
			"bind":     srv.S.Bind,
			"tls-cert": srv.S.TLSCert,
			"tls-key":  srv.S.TLSKey,
		}).Info("api/external: starting api server")

		if srv.S.TLSCert == "" || srv.S.TLSKey == "" {
			log.Fatal(http.ListenAndServe(srv.S.Bind, h2c.NewHandler(handler, &http2.Server{})))
		} else {
			log.Fatal(http.ListenAndServeTLS(
				srv.S.Bind,
				srv.S.TLSCert,
				srv.S.TLSKey,
				h2c.NewHandler(handler, &http2.Server{}),
			))
		}
	}()

	// give the http server some time to start
	time.Sleep(time.Millisecond * 100)

	// setup the HTTP handler
	clientHTTPHandler, err = srv.setupHTTPAPI()
	if err != nil {
		return err
	}

	return nil
}

// SetupCusAPI registers all REST api services
func (srv *RESTApiServer) SetupCusAPI(h *store.Handler, grpcServer *grpc.Server) error {
	jwtSecret := srv.S.JWTSecret
	if jwtSecret == "" {
		return errors.New("jwt_secret must be set")
	}
	jwtTTL := srv.S.JWTDefaultTTL
	pgs := pgstore.New()

	jwtValidator := jwt.NewValidator("HS256", []byte(jwtSecret), jwtTTL)
	otpValidator, err := otp.NewValidator("lpwan-app-server", srv.S.OTPSecret, pgs)
	if err != nil {
		return err
	}
	grpcAuth := grpcauth.New(pgs, jwtValidator, otpValidator)
	authcus.SetupCred(pgs, jwtValidator, otpValidator)

	rpID, err := uuid.FromString(srv.ApplicationServerID)
	if err != nil {
		return fmt.Errorf("failed to convert application server id from string to uuid: %v", err)
	}

	pb.RegisterFUOTADeploymentServiceServer(grpcServer, NewFUOTADeploymentAPI(h))
	pb.RegisterDeviceQueueServiceServer(grpcServer, NewDeviceQueueAPI(h))
	pb.RegisterMulticastGroupServiceServer(grpcServer, NewMulticastGroupAPI(rpID, h))
	pb.RegisterServiceProfileServiceServer(grpcServer, NewServiceProfileServiceAPI(h))
	pb.RegisterDeviceProfileServiceServer(grpcServer, NewDeviceProfileServiceAPI(h))
	// device
	api.RegisterDeviceServiceServer(grpcServer, NewDeviceAPI(rpID, h))
	// gateway
	psCli, err := pscli.GetPServerClient()
	if err != nil {
		return err
	}
	api.RegisterGatewayServiceServer(grpcServer, NewGatewayAPI(
		h,
		grpcAuth,
		GwConfig{
			ApplicationServerID: rpID,
			ServerAddr:          srv.ServerAddr,
			EnableSTC:           srv.EnableSTC,
		},
		psCli,
	))

	// gateway profile
	api.RegisterGatewayProfileServiceServer(grpcServer, NewGatewayProfileAPI(h))
	// application
	api.RegisterApplicationServiceServer(grpcServer, NewApplicationAPI(h))
	// network server
	api.RegisterNetworkServerServiceServer(grpcServer, NewNetworkServerAPI(h))
	// orgnization
	api.RegisterOrganizationServiceServer(grpcServer, NewOrganizationAPI(h))
	// user
	pwhasher, err := pwhash.New(16, srv.PasswordHashIterations)
	if err != nil {
		return err
	}
	userSrv := user.NewServer(
		pgs,
		srv.Mailer,
		grpcAuth,
		jwtValidator,
		otpValidator,
		pwhasher,
		user.Config{
			Recaptcha:        srv.Recaptcha,
			Enable2FALogin:   srv.Enable2FA,
			OperatorLogoPath: srv.OperatorLogo,
			WeChatLogin:      srv.ExternalAuth.WechatAuth,
			DebugWeChatLogin: srv.ExternalAuth.DebugWechatAuth,
			ShopifyConfig:    srv.ShopifyConfig,
		},
	)
	api.RegisterUserServiceServer(grpcServer, userSrv)
	api.RegisterInternalServiceServer(grpcServer, userSrv)
	api.RegisterExternalUserServiceServer(grpcServer, userSrv)

	api.RegisterServerInfoServiceServer(grpcServer, NewServerInfoAPI(srv.ServerRegion))
	api.RegisterSettingsServiceServer(grpcServer, NewSettingsServerAPI())
	api.RegisterTopUpServiceServer(grpcServer, NewTopUpServerAPI(grpcAuth))

	api.RegisterWalletServiceServer(grpcServer, NewWalletServerAPI(
		h,
		grpcAuth,
		srv.EnableSTC,
	))

	api.RegisterWithdrawServiceServer(grpcServer, NewWithdrawServerAPI(grpcAuth))

	api.RegisterStakingServiceServer(grpcServer, staking.NewServer(
		mxpcli.Global.GetStakingServiceClient(),
		grpcAuth,
	))

	api.RegisterDHXServcieServer(grpcServer, dhx.NewServer(
		mxpcli.Global.GetDHXServiceClient(),
		grpcAuth,
		pgs,
	))

	api.RegisterShopifyIntegrationServer(grpcServer, user.NewShopifyServiceServer(
		grpcAuth,
		pgs,
	))

	api.RegisterDFIServiceServer(grpcServer, dfi.NewServer(
		pgs,
		srv.MXPCli,
	))

	api.RegisterDownloadServiceServer(grpcServer, download.NewServer(
		srv.MXPCli,
		grpcAuth,
		srv.ServerAddr,
	))
	return nil
}

func (srv *RESTApiServer) setupHTTPAPI() (http.Handler, error) {
	r := mux.NewRouter()

	// setup json api handler
	jsonHandler, err := srv.getJSONGateway(context.Background())
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

func (srv *RESTApiServer) getJSONGateway(ctx context.Context) (http.Handler, error) {
	// dial options for the grpc-gateway
	var grpcDialOpts []grpc.DialOption

	if srv.S.TLSCert == "" || srv.S.TLSKey == "" {
		grpcDialOpts = append(grpcDialOpts, grpc.WithInsecure())
	} else {
		b, err := ioutil.ReadFile(srv.S.TLSCert)
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

	bindParts := strings.SplitN(srv.S.Bind, ":", 2)
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

	err = api.RegisterDownloadServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts)
	log.Infof("register download service handler: %v", err)

	return mux, nil
}
