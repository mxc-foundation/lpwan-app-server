package download

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func download() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURL := r.URL.Path
		pathList := strings.Split(requestURL, "/")
		if pathList[1] != "download" {
			w.WriteHeader(400)
		}
		file := filepath.Join("/", filepath.Join(pathList[2:]...))
		http.ServeFile(w, r, file)
	})
}

// Config defines configuration for download service
type Config struct {
	Bind string `mapstructure:"bind"`
}

// Setup sets up dowanload service API
func Setup(config Config) {
	handler := mux.NewRouter()
	handler.PathPrefix("/download").Handler(download())

	// start the download service
	go func() {
		logrus.WithFields(logrus.Fields{
			"bind": config.Bind,
		}).Info("starting download server")

		logrus.Fatal(http.ListenAndServe(config.Bind, h2c.NewHandler(handler, &http2.Server{})))
	}()
}
