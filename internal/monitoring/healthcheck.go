package monitoring

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func healthCheckHandlerFunc(w http.ResponseWriter, r *http.Request) {
	_, err := storage.RedisClient().Ping().Result()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(errors.Wrap(err, "redis ping error").Error()))
	}

	err = storage.DB().Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(errors.Wrap(err, "postgresql ping error").Error()))
	}

	w.WriteHeader(http.StatusOK)
}
