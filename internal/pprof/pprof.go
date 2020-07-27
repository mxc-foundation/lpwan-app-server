package pprof

import (
	"net/http"
	"net/http/pprof"
)

type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Bind    string `mapstructure:"bind"`
}

func Setup(c Config) error {
	if !c.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := http.Server{
		Addr:    c.Bind,
		Handler: mux,
	}
	go srv.ListenAndServe()

	return nil
}
