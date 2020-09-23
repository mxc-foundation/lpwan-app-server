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

type client struct {
	rc redis.UniversalClient
}

func (r client) Subscribe(channels ...string) RedisPubSub {
	return r.rc.Subscribe(channels...)
}

func (r client) Publish(channel string, message interface{}) RedisIntCmd {
	return r.rc.Publish(channel, message)
}

func (r client) FlushAll() RedisStatusCmd {
	return r.rc.FlushAll()
}

func (r client) Set(key string, value interface{}, expiration time.Duration) RedisStatusCmd {
	return r.rc.Set(key, value, expiration)
}

func (r client) Del(keys ...string) RedisIntCmd {
	return r.rc.Del(keys...)
}

func (r client) Get(key string) RedisStringCmd {
	return r.rc.Get(key)
}

func (r client) LRange(key string, start, stop int64) RedisStringSliceCmd {
	return r.rc.LRange(key, start, stop)
}

func (r client) TxPipeline() RedisPipeliner {
	return r.rc.TxPipeline()
}

func (r client) Pipeline() RedisPipeliner {
	return r.rc.Pipeline()
}

func (r client) SetNX(key string, value interface{}, expiration time.Duration) RedisBoolCmd {
	return r.rc.SetNX(key, value, expiration)
}

func (r client) HGetAll(key string) RedisStringStringMapCmd {
	return r.rc.HGetAll(key)
}

func (r client) Keys(pattern string) RedisStringSliceCmd {
	return r.rc.Keys(pattern)
}

func (r client) Ping() RedisStatusCmd {
	return r.rc.Ping()
}
