package pgstore

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	errHandler "github.com/mxc-foundation/lpwan-app-server/internal/errors"
	. "github.com/mxc-foundation/lpwan-app-server/internal/modules/gateway/data"
)

func (ps *PgStore) AddNewDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	_, err := ps.db.ExecContext(ctx, `
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

func (ps *PgStore) UpdateDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	_, err := ps.db.ExecContext(ctx, `
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

func (ps *PgStore) GetDefaultGatewayConfig(ctx context.Context, defaultConfig *DefaultGatewayConfig) error {
	err := ps.db.QueryRowxContext(ctx, `
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
		if err == sql.ErrNoRows {
			return errHandler.ErrDoesNotExist
		}
		return errors.Wrap(err, "select error")
	}

	return errors.Wrap(err, "GetDefaultGatewayConfig")
}
