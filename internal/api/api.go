package api

import (
	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/as"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/external"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/gws"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/js"
	"github.com/mxc-foundation/lpwan-app-server/internal/api/m2m"
	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// Setup configures the API endpoints.
func Setup(conf config.Config) error {
	if err := as.Setup(conf); err != nil {
		return errors.Wrap(err, "setup application-server api error")
	}

	if err := external.Setup(conf); err != nil {
		return errors.Wrap(err, "setup external api error")
	}

	if err := js.Setup(conf); err != nil {
		return errors.Wrap(err, "setup join-server api error")
	}

	if err := gws.Setup(); err != nil {
		return errors.Wrap(err, "setup gateway api error")
	}

	if err := m2m.Setup(conf); err != nil {
		return errors.Wrap(err, "setup m2m api error")
	}

	return nil
}
