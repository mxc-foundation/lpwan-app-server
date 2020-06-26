package externalcus

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
	"github.com/mxc-foundation/lpwan-app-server/internal/otp/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	devicePg "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	gwPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/staking"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/topup"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	userPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/wallet"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/withdraw"
)

func SetupAPI(grpcServer *grpc.Server) error {
	otpStore := pgstore.New(storage.DB().DB.DB)
	otpValidator, err := otp.NewValidator("lpwan-app-server", config.C.ApplicationServer.ExternalAPI.OTPSecret, otpStore)
	if err != nil {
		return err
	}

	validator := authcus.NewJWTValidator(storage.DB(), "HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret, otpValidator)

	api.RegisterDeviceServiceServer(grpcServer, device.NewDeviceAPI(validator, devicePg.New(storage.DB().DB)))
	api.RegisterUserServiceServer(grpcServer, user.NewUserAPI(validator, userPg.New(storage.DB().DB)))
	api.RegisterInternalServiceServer(grpcServer, user.NewInternalUserAPI(validator, otpValidator, userPg.New(storage.DB().DB)))
	api.RegisterGatewayServiceServer(grpcServer, gateway.NewGatewayAPI(validator, gwPg.New(storage.DB().DB)))
	api.RegisterServerInfoServiceServer(grpcServer, serverinfo.NewServerInfoAPI())
	api.RegisterDSDeviceServiceServer(grpcServer, device.NewM2MDeviceAPI(validator))
	api.RegisterGSGatewayServiceServer(grpcServer, gateway.NewM2MGatewayAPI(validator))
	api.RegisterSettingsServiceServer(grpcServer, serverinfo.NewSettingsServerAPI(validator))
	api.RegisterStakingServiceServer(grpcServer, staking.NewStakingServerAPI(validator))
	api.RegisterTopUpServiceServer(grpcServer, topup.NewTopUpServerAPI(validator))
	api.RegisterWalletServiceServer(grpcServer, wallet.NewWalletServerAPI(validator))
	api.RegisterWithdrawServiceServer(grpcServer, withdraw.NewWithdrawServerAPI(validator))
	api.RegisterM2MServerInfoServiceServer(grpcServer, serverinfo.NewM2MServerAPI(validator))

	return nil
}

func CusGetJSONGateway(ctx context.Context, mux *runtime.ServeMux, apiEndpoint string, grpcDialOpts []grpc.DialOption) error {

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
	if err := api.RegisterServerInfoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register server info handler error")
	}
	if err := api.RegisterDSDeviceServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}
	if err := api.RegisterGSGatewayServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
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
	if err := api.RegisterM2MServerInfoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return errors.Wrap(err, "register proxy request handler error")
	}

	return nil
}
