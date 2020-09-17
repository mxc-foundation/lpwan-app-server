// +build windows

package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/modules/serverinfo"
)

func setSyslog() error {
	if serverinfo.GetSettings().LogToSyslog {
		log.Fatal("syslog logging is not supported on Windows")
	}

	return nil
}
