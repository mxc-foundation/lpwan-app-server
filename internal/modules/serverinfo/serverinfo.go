package serverinfo

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/api/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/oidc"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	joinserver "github.com/mxc-foundation/lpwan-app-server/internal/api/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation"
	"github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup"
	m2mcli "github.com/mxc-foundation/lpwan-app-server/internal/clients/mxprotocol-server"
	"github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	jscodec "github.com/mxc-foundation/lpwan-app-server/internal/codec/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/email"
	"github.com/mxc-foundation/lpwan-app-server/internal/fuota"
	"github.com/mxc-foundation/lpwan-app-server/internal/integration"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/store"
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	"github.com/mxc-foundation/lpwan-app-server/internal/monitoring"
	"github.com/mxc-foundation/lpwan-app-server/internal/pprof"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

type controller struct {
	st      *store.Handler
	general config.ServerSettingsStruct
}

var ctrl *controller

func Setup(h *store.Handler) error {
	ctrl.st = h
	return nil
}

func GetSettings() config.ServerSettingsStruct {
	return ctrl.general
}

// SettingsSetup init settings extracted values from toml file then assign each modules
func SettingsSetup(conf config.Config) error {
	ctrl = &controller{
		general: conf.General,
	}

	if err := storage.SettingsSetup(storage.SettingStruct{
		Db:                  conf.PostgreSQL,
		Metrics:             conf.Metrics,
		JWTSecret:           conf.ApplicationServer.ExternalAPI.JWTSecret,
		ApplicationServerID: conf.ApplicationServer.ID,
	}); err != nil {
		return err
	}

	if err := redis.SettingsSetup(conf.Redis); err != nil {
		return err
	}

	if err := email.SettingsSetup(conf.SMTP, conf.Operator, email.ServerInfoStruct{
		ServerAddr:      conf.General.ServerAddr,
		DefaultLanguage: conf.General.DefaultLanguage,
	}); err != nil {
		return err
	}

	if err := m2mcli.SettingsSetup(conf.M2MServer); err != nil {
		return err
	}

	if err := gws.SettingsSetup(conf.ApplicationServer.APIForGateway); err != nil {
		return err
	}

	if err := psconn.SettingsSetup(conf.ProvisionServer); err != nil {
		return err
	}

	if err := user.SettingsSetup(user.Config{
		Recaptcha:      conf.Recaptcha,
		Enable2FALogin: conf.General.Enable2FALogin,
	}); err != nil {
		return err
	}

	if err := external.SettingsSetup(conf.ApplicationServer.ExternalAPI, conf.ApplicationServer.ID); err != nil {
		return err
	}

	if err := oidc.SettingsSetup(conf.ApplicationServer.UserAuthentication, conf.ApplicationServer.ExternalAPI.JWTSecret); err != nil {
		return err
	}

	if err := jscodec.SettingsSetup(conf.ApplicationServer.Codec); err != nil {
		return err
	}

	if err := integration.SettingsSetup(conf.ApplicationServer.Integration); err != nil {
		return err
	}

	if err := as.SettingsSetup(conf.ApplicationServer.API); err != nil {
		return err
	}

	if err := m2m.SettingsSetup(conf.ApplicationServer.APIForM2M); err != nil {
		return err
	}

	if err := multicastsetup.SettingsSetup(conf.ApplicationServer.RemoteMulticastSetup); err != nil {
		return err
	}

	if err := fragmentation.SettingsSetup(conf.ApplicationServer.FragmentationSession); err != nil {
		return err
	}

	if err := fuota.SettingsSetup(conf.ApplicationServer.FUOTADeployment, fuota.Config{
		RemoteMulticastSetupRetries:       conf.ApplicationServer.RemoteMulticastSetup.SyncRetries,
		RemoteFragmentationSessionRetries: conf.ApplicationServer.FragmentationSession.SyncRetries,
		ApplicationServerID:               conf.ApplicationServer.ID,
	}); err != nil {
		return err
	}

	if err := mining.SettingsSetup(conf.ApplicationServer.MiningSetUp); err != nil {
		return err
	}

	if err := joinserver.SettingsSetup(conf.JoinServer); err != nil {
		return err
	}

	if err := monitoring.SettingsSerup(conf.Monitoring); err != nil {
		return err
	}

	if err := pprof.SettingsSetup(conf.PProf); err != nil {
		return err
	}

	if err := gateway.SettingsSetup(gateway.Config{
		ServerAddr: conf.General.ServerAddr,
	}); err != nil {
		return err
	}

	if err := device.SettingsSetup(device.Config{
		ApplicationServerID: conf.ApplicationServer.ID,
	}); err != nil {
		return err
	}

	return nil
}
