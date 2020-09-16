package redis

import (
	"errors"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

func SettingsSetup(s RedisStruct) error {
	ctrl = &controller{
		redis: s,
	}

	return nil
}

// SetupRedis :
func SetupRedis() error {
	log.Info("storage: setting up Redis client")
	if len(ctrl.redis.Servers) == 0 {
		return errors.New("at least one redis server must be configured")
	}

	if ctrl.redis.Cluster {
		redisClient.UniversalClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    ctrl.redis.Servers,
			PoolSize: ctrl.redis.PoolSize,
			Password: ctrl.redis.Password,
		})
	} else if ctrl.redis.MasterName != "" {
		redisClient.UniversalClient = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       ctrl.redis.MasterName,
			SentinelAddrs:    ctrl.redis.Servers,
			SentinelPassword: ctrl.redis.Password,
			DB:               ctrl.redis.Database,
			PoolSize:         ctrl.redis.PoolSize,
			Password:         ctrl.redis.Password,
		})
	} else {
		redisClient.UniversalClient = redis.NewClient(&redis.Options{
			Addr:     ctrl.redis.Servers[0],
			DB:       ctrl.redis.Database,
			Password: ctrl.redis.Password,
			PoolSize: ctrl.redis.PoolSize,
		})
	}

	return nil
}

type RedisStruct struct {
	URL        string   `mapstructure:"url"` // deprecated
	Servers    []string `mapstructure:"servers"`
	Cluster    bool     `mapstructure:"cluster"`
	MasterName string   `mapstructure:"master_name"`
	PoolSize   int      `mapstructure:"pool_size"`
	Password   string   `mapstructure:"password"`
	Database   int      `mapstructure:"database"`
}

type controller struct {
	redis RedisStruct
}

var ctrl *controller
