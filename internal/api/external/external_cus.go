package external

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/staking"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/topup"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/wallet"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/withdraw"

	authPg "github.com/mxc-foundation/lpwan-app-server/internal/authentication/pgstore"
	gatewayprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
)

func SetupCusAPI(grpcServer *grpc.Server) error {
	jwtValidator := jwt.NewJWTValidator("HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret)
	otpValidator, err := otp.NewValidator("lpwan-app-server", config.C.ApplicationServer.ExternalAPI.OTPSecret, pgstore.New(storage.DB().DB.DB))
	if err != nil {
		return err
	}
	authcus.SetupCred(authPg.New(storage.DB().DB), jwtValidator, otpValidator)

	// device
	api.RegisterDeviceServiceServer(grpcServer, device.NewDeviceAPI(applicationServerID))
	// gateway
	api.RegisterGatewayServiceServer(grpcServer, gateway.NewGatewayAPI(applicationServerID))
	// gateway profile
	api.RegisterGatewayProfileServiceServer(grpcServer, gatewayprofile.NewGatewayProfileAPI())
	// application
	api.RegisterApplicationServiceServer(grpcServer, application.NewApplicationAPI())
	// network server
	api.RegisterNetworkServerServiceServer(grpcServer, networkserver.NewNetworkServerAPI())
	// orgnization
	api.RegisterOrganizationServiceServer(grpcServer, organization.NewOrganizationAPI())
	// user
	api.RegisterUserServiceServer(grpcServer, user.NewUserAPI())
	api.RegisterInternalServiceServer(grpcServer, user.NewInternalUserAPI())

	api.RegisterServerInfoServiceServer(grpcServer, serverinfo.NewServerInfoAPI())
	api.RegisterSettingsServiceServer(grpcServer, serverinfo.NewSettingsServerAPI())
	api.RegisterStakingServiceServer(grpcServer, staking.NewStakingServerAPI())
	api.RegisterTopUpServiceServer(grpcServer, topup.NewTopUpServerAPI())
	api.RegisterWalletServiceServer(grpcServer, wallet.NewWalletServerAPI())
	api.RegisterWithdrawServiceServer(grpcServer, withdraw.NewWithdrawServerAPI())

	return nil
}

func CusGetJSONGateway(ctx context.Context, mux *runtime.ServeMux, apiEndpoint string, grpcDialOpts []grpc.DialOption) error {

	if err := api.RegisterServerInfoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register server info handler error")
	}
	if err := api.RegisterStakingServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterTopUpServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterWalletServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterWithdrawServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterSettingsServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}

	return nil
}
