package redis

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisStore interface {
	Subscribe(channels ...string) RedisPubSub
	Publish(channel string, message interface{}) RedisIntCmd
	FlushAll() RedisStatusCmd
	Set(key string, value interface{}, expiration time.Duration) RedisStatusCmd
	Del(keys ...string) RedisIntCmd
	Get(key string) RedisStringCmd
	LRange(key string, start, stop int64) RedisStringSliceCmd
	TxPipeline() redis.Pipeliner
	Pipeline() redis.Pipeliner
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	HGetAll(key string) *redis.StringStringMapCmd
	Keys(pattern string) RedisStringSliceCmd
	Ping() RedisStatusCmd
}

type RedisPubSub interface {
	Receive() (interface{}, error)
	Channel() <-chan *redis.Message
	Close() error
}
type RedisIntCmd interface {
	Val() int64
	Result() (int64, error)
	Uint64() (uint64, error)
	String() string
	Err() error
}
type RedisStatusCmd interface {
	Val() string
	Result() (string, error)
	String() string
	Err() error
}
type RedisStringCmd interface {
	Val() string
	Result() (string, error)
	Bytes() ([]byte, error)
	Int() (int, error)
	Int64() (int64, error)
	Uint64() (uint64, error)
	Float32() (float32, error)
	Float64() (float64, error)
	Time() (time.Time, error)
	Scan(val interface{}) error
	String() string
}
type RedisStringSliceCmd interface {
	Val() []string
	Result() ([]string, error)
	String() string
	ScanSlice(container interface{}) error
}
