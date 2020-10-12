package pprof

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	mgr "github.com/mxc-foundation/lpwan-app-server/internal/system_manager"

	. "github.com/mxc-foundation/lpwan-app-server/internal/pprof/data"
)

func init() {
	mgr.RegisterSettingsSetup(moduleName, SettingsSetup)
	mgr.RegisterModuleSetup(moduleName, Setup)
}

const moduleName = "pprof"

type controller struct {
	name string
	s    Config

	moduleUp bool
}

var ctrl *controller

// SettingsSetup initialize module settings on start
func SettingsSetup(name string, conf config.Config) error {
	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	ctrl = &controller{
		name: moduleName,
		s:    conf.PProf,
	}
	return nil
}

func Setup(name string, h *store.Handler) error {
	if ctrl.moduleUp == true {
		return nil
	}

	defer func() {
		ctrl.moduleUp = true
	}()

	if name != moduleName {
		return errors.New(fmt.Sprintf("Calling SettingsSetup for %s, but %s is called", name, moduleName))
	}

	if !ctrl.s.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := http.Server{
		Addr:    ctrl.s.Bind,
		Handler: mux,
	}
	go func() {
		_ = srv.ListenAndServe()
	}()

	return nil
}
