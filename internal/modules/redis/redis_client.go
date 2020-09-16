package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type RSGetResult interface {
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
type redisStore interface {
	RSSet(key string, value interface{}) error
	RSGet(key string) RSGetResult
	RSDelete(key string) error
}

type RedisClientType struct {
	redis.UniversalClient
	redisStore
}

var redisClient RedisClientType

const (
	micLookupExpire = time.Second * 10
	micLookupTempl  = "lora:as:gwping:%s"
)

func (r RedisClientType) RSSet(keyWord string, value interface{}) error {
	key := fmt.Sprintf(micLookupTempl, keyWord)
	return r.Set(key, value, micLookupExpire).Err()

}
func (r RedisClientType) RSGet(key string) RSGetResult {
	return r.Get(key)
}
func (r RedisClientType) RSDelete(key string) error {
	return r.Del(key).Err()
}

// RedisClient returns the RedisClient.
func RedisClient() RedisClientType {
	return redisClient
}
