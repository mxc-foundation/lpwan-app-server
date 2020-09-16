package monitoring

import (
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/storage"
)

func healthCheckHandlerFunc(w http.ResponseWriter, r *http.Request) {
	_, err := rs.RedisClient().Ping().Result()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(errors.Wrap(err, "redis ping error").Error()))
	}

	err = storage.DBTest().Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(errors.Wrap(err, "postgresql ping error").Error()))
	}

	w.WriteHeader(http.StatusOK)
}
