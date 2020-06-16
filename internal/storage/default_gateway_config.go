package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DefaultGatewayConfig struct {
	ID            int64      `db:"id"`
	Model         string     `db:"model"`
	Region        string     `db:"region"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DefaultConfig string     `db:"default_config"`
}

func AddNewDefaultGatewayConfig(db sqlx.Execer, defaultConfig *DefaultGatewayConfig) error {
	_, err := db.Exec(`
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

func UpdateDefaultGatewayConfig(db sqlx.Execer, defaultConfig *DefaultGatewayConfig) error {
	_, err := db.Exec(`
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


func GetDefaultGatewayConfig(db sqlx.Queryer, defaultConfig *DefaultGatewayConfig) error {
	err := db.QueryRowx(`
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
