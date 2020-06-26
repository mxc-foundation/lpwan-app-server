package pgstore

import (
	"github.com/pkg/errors"

	gwmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway"
)

func (h *GWHandler) AddNewDefaultGatewayConfig(defaultConfig *gwmod.DefaultGatewayConfig) error {
	_, err := h.db.Exec(`
		insert into default_gateway_config (
		    model, region, created_at, updated_at, default_config
		) values (
		    $1, $2, now(), now(), $3
		)`,
		defaultConfig.Model,
		defaultConfig.Region,
		defaultConfig.DefaultConfig)

	if err != nil {
		return errors.Wrap(err, "AddNewDefaultGatewayConfig")
	}

	return errors.Wrap(err, "AddNewDefaultGatewayConfig")
}

func (h *GWHandler) UpdateDefaultGatewayConfig(defaultConfig *gwmod.DefaultGatewayConfig) error {
	_, err := h.db.Exec(`
		update 
		    default_gateway_config 
		set
		    updated_at = now(), 
		    default_config = $1
		where 
		    id = $2 `,
		defaultConfig.DefaultConfig,
		defaultConfig.ID)

	if err != nil {
		return errors.Wrap(err, "UpdateDefaultGatewayConfig")
	}

	return errors.Wrap(err, "UpdateDefaultGatewayConfig")
}

func (h *GWHandler) GetDefaultGatewayConfig(defaultConfig *gwmod.DefaultGatewayConfig) error {
	err := h.db.QueryRowx(`
		select 
		    id, model, region, created_at, updated_at, default_config 
		from 
		    default_gateway_config
		where 
		    model = $1 and region = $2 `,
		defaultConfig.Model,
		defaultConfig.Region).Scan(
		&defaultConfig.ID,
		&defaultConfig.Model,
		&defaultConfig.Region,
		&defaultConfig.CreatedAt,
		&defaultConfig.UpdatedAt,
		&defaultConfig.DefaultConfig)

	if err != nil {
		return handlePSQLError(Select, err, "select error")
	}

	return errors.Wrap(err, "GetDefaultGatewayConfig")
}
