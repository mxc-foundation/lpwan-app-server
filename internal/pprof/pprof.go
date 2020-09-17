package pprof

import (
	"net/http"
	"net/http/pprof"
)

type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Bind    string `mapstructure:"bind"`
}
type controller struct {
	s Config
}

var ctrl *controller

func SettingsSetup(s Config) error {
	ctrl = &controller{
		s: s,
	}
	return nil
}

func Setup() error {
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
	go srv.ListenAndServe()

	return nil
}
