package code

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/mxc-foundation/lpwan-app-server/internal/config"
)

// MigrateToClusterKeys migrates the keys to Redis Cluster compatible keys.
func MigrateToClusterKeys(conf config.Config) error {
	keys, err := rs.RedisClient().Keys("lora:as:metrics:*").Result()
	if err != nil {
		return errors.Wrap(err, "get keys error")
	}

	for i, key := range keys {
		if err := migrateKey(conf, key); err != nil {
			log.WithError(err).Error("migrations/code: migrate metrics key error")
		}

		if i > 0 && i%1000 == 0 {
			log.WithFields(log.Fields{
				"migrated":    i,
				"total_count": len(keys),
			}).Info("migrations/code: migrating metrics keys")
		}
	}

	return nil
}

func migrateKey(conf config.Config, key string) error {
	keyParts := strings.Split(key, ":")
	if len(keyParts) < 6 {
		return fmt.Errorf("key %s is invalid", key)
	}

	metricsStruct := storage.GetMetricsSettings()
	ttlMap := map[string]time.Duration{
		"MINUTE": metricsStruct.Redis.MinuteAggregationTTL,
		"HOUR":   metricsStruct.Redis.HourAggregationTTL,
		"DAY":    metricsStruct.Redis.DayAggregationTTL,
		"MONTH":  metricsStruct.Redis.MonthAggregationTTL,
	}

	ttl, ok := ttlMap[keyParts[len(keyParts)-2]]
	if !ok {
		return fmt.Errorf("key %s is invalid", key)
	}

	newKey := fmt.Sprintf("lora:as:metrics:{%s}:%s", strings.Join(keyParts[3:len(keyParts)-2], ":"), strings.Join(keyParts[len(keyParts)-2:], ":"))

	val, err := rs.RedisClient().HGetAll(key).Result()
	if err != nil {
		return errors.Wrap(err, "hgetall error")
	}

	pipe := rs.RedisClient().TxPipeline()
	for k, v := range val {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return errors.Wrap(err, "parse float error")
		}
		pipe.HIncrByFloat(newKey, k, f)
	}
	pipe.PExpire(key, ttl)

	if _, err := pipe.Exec(); err != nil {
		return errors.Wrap(err, "exec error")
	}

	return nil
}
