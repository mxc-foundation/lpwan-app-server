package external

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

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
	applicationPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/pgstore"
	devicePg "github.com/mxc-foundation/lpwan-app-server/internal/modules/device/pgstore"
	gwPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/pgstore"
	networkServerPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/networkserver/pgstore"
	organizationPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization/pgstore"
	userPg "github.com/mxc-foundation/lpwan-app-server/internal/modules/user/pgstore"
)

func SetupCusAPI(grpcServer *grpc.Server) error {
	jwtValidator := jwt.NewJWTValidator("HS256", config.C.ApplicationServer.ExternalAPI.JWTSecret)
	otpValidator, err := otp.NewValidator("lpwan-app-server", config.C.ApplicationServer.ExternalAPI.OTPSecret, pgstore.New(storage.DB().DB.DB))
	if err != nil {
		return err
	}

	tx, err := storage.DB().Beginx()
	if err != nil {
		return err
	}

	authcus.SetupCred(authPg.New(storage.DB().DB), jwtValidator, otpValidator)

	// device
	api.RegisterDeviceServiceServer(grpcServer, device.NewDeviceAPI(device.DeviceAPI{
		Validator:            device.NewValidator(otpValidator),
		Store:                devicePg.New(tx.Tx, storage.DB().DB),
		AppplicationServerID: applicationServerID,
	}))

	// user
	api.RegisterUserServiceServer(grpcServer, user.NewUserAPI(user.UserAPI{
		Validator: user.NewValidator(otpValidator),
		Store:     userPg.New(tx.Tx, storage.DB().DB),
	}))

	api.RegisterInternalServiceServer(grpcServer, user.NewInternalUserAPI(user.InternalUserAPI{
		Validator: user.NewValidator(otpValidator),
		Store:     userPg.New(tx.Tx, storage.DB().DB),
	}))

	// gateway
	api.RegisterGatewayServiceServer(grpcServer, gateway.NewGatewayAPI(gateway.GatewayAPI{
		Validator:           gateway.NewValidator(otpValidator),
		Store:               gwPg.New(tx.Tx, storage.DB().DB),
		ApplicationServerID: applicationServerID,
	}))

	api.RegisterServerInfoServiceServer(grpcServer, serverinfo.NewServerInfoAPI(serverinfo.ServerInfoAPI{
		Validator: serverinfo.NewValidator(otpValidator),
	}))

	api.RegisterSettingsServiceServer(grpcServer, serverinfo.NewSettingsServerAPI(serverinfo.SettingsServerAPI{
		Validator: serverinfo.NewValidator(otpValidator),
	}))

	api.RegisterStakingServiceServer(grpcServer, staking.NewStakingServerAPI(staking.StakingServerAPI{
		Validator: staking.NewValidator(otpValidator),
	}))

	api.RegisterTopUpServiceServer(grpcServer, topup.NewTopUpServerAPI(topup.TopUpServerAPI{
		Validator: topup.NewValidator(otpValidator),
	}))

	api.RegisterWalletServiceServer(grpcServer, wallet.NewWalletServerAPI(wallet.WalletServerAPI{
		Validator: wallet.NewValidator(otpValidator),
	}))

	api.RegisterWithdrawServiceServer(grpcServer, withdraw.NewWithdrawServerAPI(withdraw.WithdrawServerAPI{
		Validator: withdraw.NewValidator(otpValidator),
	}))

	api.RegisterNetworkServerServiceServer(grpcServer, networkserver.NewNetworkServerAPI(networkserver.NetworkServerAPI{
		Validator: networkserver.NewValidator(otpValidator),
		Store:     networkServerPg.New(tx.Tx, storage.DB().DB),
	}))

	api.RegisterApplicationServiceServer(grpcServer, application.NewApplicationAPI(application.ApplicationAPI{
		Validator: application.NewValidator(application.Validator{
			Store:       applicationPg.New(tx.Tx, storage.DB().DB),
			Credentials: authcus.NewCredentials(),
		}),
		Store: applicationPg.New(tx.Tx, storage.DB().DB),
	}))

	api.RegisterOrganizationServiceServer(grpcServer, organization.NewOrganizationAPI(organization.OrganizationAPI{
		Validator: organization.NewValidator(otpValidator),
		Store:     organizationPg.New(tx.Tx, storage.DB().DB),
	}))

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
