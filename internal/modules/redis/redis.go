package redis

import (
	"errors"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

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
	redis   RedisStruct
	handler *RedisHandler
}

var ctrl *controller

func SettingsSetup(s RedisStruct) error {
	ctrl = &controller{
		redis:   s,
		handler: &RedisHandler{},
	}

	return nil
}

// Setup :
func Setup() (err error) {
	log.Info("storage: setting up Redis client")
	if len(ctrl.redis.Servers) == 0 {
		return errors.New("at least one redis server must be configured")
	}

	var st RedisStore
	if ctrl.redis.Cluster {
		st = NewRedisStore(&client{
			rc: redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    ctrl.redis.Servers,
				PoolSize: ctrl.redis.PoolSize,
				Password: ctrl.redis.Password,
			}),
		})
	} else if ctrl.redis.MasterName != "" {
		st = NewRedisStore(&client{
			rc: redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:       ctrl.redis.MasterName,
				SentinelAddrs:    ctrl.redis.Servers,
				SentinelPassword: ctrl.redis.Password,
				DB:               ctrl.redis.Database,
				PoolSize:         ctrl.redis.PoolSize,
				Password:         ctrl.redis.Password,
			}),
		})
	} else {
		st = NewRedisStore(&client{
			rc: redis.NewClient(&redis.Options{
				Addr:     ctrl.redis.Servers[0],
				DB:       ctrl.redis.Database,
				Password: ctrl.redis.Password,
				PoolSize: ctrl.redis.PoolSize,
			}),
		})
	}

	SetupRedisHandler(st)

	return nil
}

func SetupRedisHandler(store RedisStore) {
	ctrl.handler.S = store
}

type RedisHandler struct {
	S RedisStore
}

// RedisClient returns the RedisClient.
func RedisClient() RedisStore {
	return ctrl.handler.S
}
