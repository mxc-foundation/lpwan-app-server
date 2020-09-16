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

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
)

type ExternalAPIStruct struct {
	Bind            string
	TLSCert         string `mapstructure:"tls_cert"`
	TLSKey          string `mapstructure:"tls_key"`
	JWTSecret       string `mapstructure:"jwt_secret"`
	OTPSecret       string `mapstructure:"otp_secret"`
	CORSAllowOrigin string `mapstructure:"cors_allow_origin"`
}

type controller struct {
	s                   ExternalAPIStruct
	applicationServerID string
}

var ctrl *controller

func SettingsSetup(apiStruct ExternalAPIStruct, applicationServerID string) error {
	ctrl = &controller{
		s:                   apiStruct,
		applicationServerID: applicationServerID,
	}

	return nil
}
func GetApplicationServerID() string {
	return ctrl.applicationServerID
}
func GetJWTSecret() string {
	return ctrl.s.JWTSecret
}

// Setup configures the API package.
func Setup() error {
	if ctrl.s.JWTSecret == "" {
		return errors.New("jwt_secret must be set")
	}

	return setupAPI()
}

func setupAPI() (err error) {
	/*	validator := auth.NewJWTValidator(storage.DB(), "HS256", jwtSecret)*/
	rpID, err := uuid.FromString(ctrl.applicationServerID)
	if err != nil {
		return errors.Wrap(err, "application-server id to uuid error")
	}

	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)

	if err := SetupCusAPI(grpcServer, rpID); err != nil {
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
		w.Write(data)
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

	if err := CusGetJSONGateway(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, err
	}

	return mux, nil
}
