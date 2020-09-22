package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	MicLookupExpire = time.Second * 10
	MicLookupTempl  = "lora:as:gwping:%s"
)

func NewRedisStore(client RedisStore) RedisStore {
	return client
}

type redisClient struct {
	rc redis.UniversalClient
}

func (r redisClient) Subscribe(channels ...string) *redis.PubSub {
	return r.Subscribe(channels...)
}

func (r redisClient) Publish(channel string, message interface{}) *redis.IntCmd {
	return r.Publish(channel, message)
}

func (r redisClient) FlushAll() *redis.StatusCmd {
	return r.FlushAll()
}

func (r redisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Set(key, value, expiration)
}

func (r redisClient) Del(keys ...string) *redis.IntCmd {
	return r.Del(keys...)
}

func (r redisClient) Get(key string) *redis.StringCmd {
	return r.Get(key)
}

func (r redisClient) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return r.LRange(key, stop, stop)
}

func (r redisClient) TxPipeline() redis.Pipeliner {
	return r.TxPipeline()
}

func (r redisClient) Pipeline() redis.Pipeliner {
	return r.Pipeline()
}

func (r redisClient) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.SetNX(key, value, expiration)
}

func (r redisClient) HGetAll(key string) *redis.StringStringMapCmd {
	return r.HGetAll(key)
}

func (r redisClient) Keys(pattern string) *redis.StringSliceCmd {
	return r.Keys(pattern)
}

func (r redisClient) Ping() *redis.StatusCmd {
	return r.Ping()
}
