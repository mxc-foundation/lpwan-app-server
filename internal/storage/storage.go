package storage

import (
	"context"
	"fmt"
	"strings"
	"text/template"
	"time"

	uuid "github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
	"github.com/mxc-foundation/lpwan-app-server/internal/migrations"
	"github.com/mxc-foundation/lpwan-app-server/internal/static"
	"github.com/mxc-foundation/lpwan-server/api/ns"
)

// HashIterations configures the Hash iteration
var (
	jwtsecret           []byte
	HashIterations      = 100000
	DemoUser            = ""
	applicationServerID uuid.UUID
)

// Setup configures the storage package.
func Setup(c config.Config) error {
	log.Info("storage: setting up storage package")

	jwtsecret = []byte(c.ApplicationServer.ExternalAPI.JWTSecret)
	HashIterations = c.General.PasswordHashIterations
	DemoUser = c.General.DemoUser

	if err := applicationServerID.UnmarshalText([]byte(c.ApplicationServer.ID)); err != nil {
		return errors.Wrap(err, "decode application_server.id error")
	}

	log.Info("storage: setup metrics")
	// setup aggregation intervals
	var intervals []AggregationInterval
	for _, agg := range c.Metrics.Redis.AggregationIntervals {
		intervals = append(intervals, AggregationInterval(strings.ToUpper(agg)))
	}
	if err := SetAggregationIntervals(intervals); err != nil {
		return errors.Wrap(err, "set aggregation intervals error")
	}

	// setup timezone
	if err := SetTimeLocation(c.Metrics.Timezone); err != nil {
		return errors.Wrap(err, "set time location error")
	}

	// setup storage TTL
	SetMetricsTTL(
		c.Metrics.Redis.MinuteAggregationTTL,
		c.Metrics.Redis.HourAggregationTTL,
		c.Metrics.Redis.DayAggregationTTL,
		c.Metrics.Redis.MonthAggregationTTL,
	)

	log.Info("storage: setting up Redis pool")
	redisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: c.Redis.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(c.Redis.URL,
				redis.DialReadTimeout(redisDialReadTimeout),
				redis.DialWriteTimeout(redisDialWriteTimeout),
			)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Now().Sub(t) < onBorrowPingInterval {
				return nil
			}

			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}

	log.Info("storage: connecting to PostgreSQL database")
	d, err := sqlx.Open("postgres", c.PostgreSQL.DSN)
	if err != nil {
		return errors.Wrap(err, "storage: PostgreSQL connection error")
	}
	d.SetMaxOpenConns(c.PostgreSQL.MaxOpenConnections)
	d.SetMaxIdleConns(c.PostgreSQL.MaxIdleConnections)
	for {
		if err := d.Ping(); err != nil {
			log.WithError(err).Warning("storage: ping PostgreSQL database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	db = &DBLogger{d}

	if c.PostgreSQL.Automigrate {
		log.Info("storage: applying PostgreSQL data migrations")
		m := &migrate.AssetMigrationSource{
			Asset:    migrations.Asset,
			AssetDir: migrations.AssetDir,
			Dir:      "",
		}
		n, err := migrate.Exec(db.DB.DB, "postgres", m, migrate.Up)
		if err != nil {
			return errors.Wrap(err, "storage: applying PostgreSQL data migrations error")
		}
		log.WithField("count", n).Info("storage: PostgreSQL data migrations applied")
	}

	return nil
}

func SetupDefault() error {
	ctx := context.Background()
	count, err := GetGatewayProfileCount(ctx, DB())
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Failed to load gateway profiles")
	}

	if count != 0 {
		// check if default gateway profile already exists
		gpList, err := GetGatewayProfiles(ctx, DB(), count, 0)
		if err != nil {
			return errors.Wrap(err, "Failed to load gateway profiles")
		}

		for _, v := range gpList {
			if v.Name == "default_gateway_profile" {
				return nil
			}
		}
	}

	// none default_gateway_profile exists, add one
	var networkServer NetworkServer
	n, err := GetNetworkServers(ctx, DB(), 1, 0)
	if err != nil && err != ErrDoesNotExist {
		return errors.Wrap(err, "Load network server internal error")
	}

	if len(n) >= 1 {
		networkServer = n[0]
	} else {
		// insert default one
		err := Transaction(func(tx sqlx.Ext) error {
			return CreateNetworkServer(ctx, DB(), &NetworkServer{
				Name:                    "default_network_server",
				Server:                  "network-server:8000",
				GatewayDiscoveryEnabled: false,
			})
		})
		if err != nil {
			return nil
		}

		// get network-server id
		networkServer, err = GetDefaultNetworkServer(ctx, DB())
		if err != nil {
			return err
		}
	}

	gp := GatewayProfile{
		NetworkServerID: networkServer.ID,
		Name:            "default_gateway_profile",
		GatewayProfile: ns.GatewayProfile{
			Channels:      []uint32{0, 1, 2},
			ExtraChannels: []*ns.GatewayProfileExtraChannel{},
		},
	}

	err = Transaction(func(tx sqlx.Ext) error {
		return CreateGatewayProfile(ctx, tx, &gp)
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create default gateway profile")
	}

	return nil
}

func LoadTemplates() error {
	// load gateway config templates
	GatewayConfigTemplate = template.Must(template.New("gateway-config/global_conf.json.sx1250.CN490").Parse(
			string(static.MustAsset("gateway-config/global_conf.json.sx1250.CN490"))))

	return nil
}
