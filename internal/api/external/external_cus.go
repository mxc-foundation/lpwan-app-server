package external

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	authPg "github.com/mxc-foundation/lpwan-app-server/internal/authentication/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	devprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile"
	fuotamod "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota_deployment"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	gatewayprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	serviceprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/staking"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/topup"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/wallet"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/withdraw"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func SetupCusAPI(grpcServer *grpc.Server, rpID uuid.UUID) error {
	jwtValidator := jwt.NewJWTValidator("HS256", []byte(config.C.ApplicationServer.ExternalAPI.JWTSecret))
	otpValidator, err := otp.NewValidator("lpwan-app-server", config.C.ApplicationServer.ExternalAPI.OTPSecret, pgstore.New(storage.DBTest().DB.DB))
	if err != nil {
		return err
	}
	authcus.SetupCred(authPg.New(storage.DBTest().DB), jwtValidator, otpValidator)

	pb.RegisterFUOTADeploymentServiceServer(grpcServer, fuotamod.NewFUOTADeploymentAPI())
	pb.RegisterDeviceQueueServiceServer(grpcServer, device.NewDeviceQueueAPI())
	pb.RegisterMulticastGroupServiceServer(grpcServer, multicast.NewMulticastGroupAPI(rpID))
	pb.RegisterServiceProfileServiceServer(grpcServer, serviceprofile.NewServiceProfileServiceAPI())
	pb.RegisterDeviceProfileServiceServer(grpcServer, devprofile.NewDeviceProfileServiceAPI())
	// device
	api.RegisterDeviceServiceServer(grpcServer, device.NewDeviceAPI(rpID))
	// gateway
	api.RegisterGatewayServiceServer(grpcServer, gateway.NewGatewayAPI(rpID))
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

	if err := api.RegisterApplicationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register application handler error")
	}
	if err := api.RegisterDeviceServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register node handler error")
	}
	if err := api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register user handler error")
	}
	if err := api.RegisterInternalServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register internal handler error")
	}
	if err := api.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register gateway handler error")
	}
	if err := api.RegisterGatewayProfileServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register gateway-profile handler error")
	}
	if err := api.RegisterOrganizationServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register organization handler error")
	}
	if err := api.RegisterNetworkServerServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register network-server handler error")
	}

	return nil
}
