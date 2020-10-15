package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"

	"github.com/mxc-foundation/lpwan-app-server/cmd/lora-app-server/cmd"

	// execute init() for following packages
	_ "github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/applayer/fragmentation"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/applayer/multicastsetup"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/clients/psconn"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/codec/js"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/downlink"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/email"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/fuota"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/gwping"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/integration"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/js"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/migrations/code"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/application"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/as"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/device"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/mining"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/monitoring"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/mxp_portal"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/networkserver_portal"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/pprof"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/storage"

	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/device-profile"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/fuota-deployment"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway-profile"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/multicast-group"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/organization"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/set_default"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/modules/user"
	_ "github.com/mxc-foundation/lpwan-app-server/internal/monitoring"
)

// grpcLogger implements a wrapper around the logrus Logger to make it
// compatible with the grpc LoggerV2. It seems that V is not (always)
// called, therefore the Info* methods are overridden as we want to
// log these as debug info.
type grpcLogger struct {
	*log.Logger
}

func (gl *grpcLogger) V(l int) bool {
	level, ok := map[log.Level]int{
		log.DebugLevel: 0,
		log.InfoLevel:  1,
		log.WarnLevel:  2,
		log.ErrorLevel: 3,
		log.FatalLevel: 4,
	}[log.GetLevel()]
	if !ok {
		return false
	}

	return l >= level
}

func (gl *grpcLogger) Info(args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(args...)
	}
}

func (gl *grpcLogger) Infoln(args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debug(args...)
	}
}

func (gl *grpcLogger) Infof(format string, args ...interface{}) {
	if log.GetLevel() == log.DebugLevel {
		log.Debugf(format, args...)
	}
}

func init() {
	grpclog.SetLoggerV2(&grpcLogger{log.StandardLogger()})

	// the default is passthrough, see:
	// https://github.com/grpc/grpc-go/issues/1783
	resolver.SetDefaultScheme("dns")
}

var version string // set by the compiler

func main() {
	cmd.Execute(version)
}
