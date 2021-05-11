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
		// for now only mining-report and staking-report should be allowed to download
		// e.g. /download/tmp/mining-report/a1a1a1a1a1a1a1a1a1a1/file_name.csv
		// e.g. /download/tmp/mining-report/a1a1a1a1a1a1a1a1a1a1/file_name.pdf
		if len(pathList) < 6 {
			w.WriteHeader(400)
			return
		}
		if pathList[1] != "download" || pathList[2] != "tmp" {
			w.WriteHeader(400)
			return
		}
		switch pathList[3] {
		case "mining-report":
			file := filepath.Join("/", filepath.Join(pathList[2:]...))
			http.ServeFile(w, r, file)
		default:
			w.WriteHeader(400)
			return
		}
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
