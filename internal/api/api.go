package api

import (
	"github.com/gofrs/uuid"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/brocaar/chirpstack-api/go/v3/as/external/api"

	api "github.com/mxc-foundation/lpwan-app-server/api/appserver-serves-ui"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/helpers"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/m2m"
	authcus "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	authPg "github.com/mxc-foundation/lpwan-app-server/internal/authentication/pgstore"
	"github.com/mxc-foundation/lpwan-app-server/internal/jwt"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	devprofile "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile"
	fuotamod "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment"
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

// Setup configures the API endpoints.
func Setup(h *store.Handler) error {
	if err := as.Setup(h); err != nil {
		return errors.Wrap(err, "setup application-server api error")
	}

	if err := js.Setup(); err != nil {
		return errors.Wrap(err, "setup join-server api error")
	}

	if err := gws.Setup(); err != nil {
		return errors.Wrap(err, "setup gateway api error")
	}

	if err := m2m.Setup(); err != nil {
		return errors.Wrap(err, "setup m2m api error")
	}

	// Bind external api port to listen to requests to all services
	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)

	if err := SetupCusAPI(grpcServer, external.GetApplicationServerID()); err != nil {
		return err
	}

	if err := external.Setup(grpcServer); err != nil {
		return errors.Wrap(err, "setup external api error")
	}

	return nil
}

func SetupCusAPI(grpcServer *grpc.Server, rpID uuid.UUID) error {
	jwtSecret := external.GetJWTSecret()
	if jwtSecret == "" {
		return errors.New("jwt_secret must be set")
	}

	jwtValidator := jwt.NewJWTValidator("HS256", []byte(jwtSecret))
	otpValidator, err := otp.NewValidator("lpwan-app-server", external.GetOTPSecret(), pgstore.New(storage.DBTest().DB.DB))
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
