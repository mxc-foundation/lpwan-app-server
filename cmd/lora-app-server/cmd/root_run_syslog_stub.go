// +build windows

package cmd

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
	log "github.com/sirupsen/logrus"
)

func setSyslog() error {
	if serverinfo.GetSettings().LogToSyslog {
		log.Fatal("syslog logging is not supported on Windows")
	}

	return nil
}
