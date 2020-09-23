package test

import (
	"fmt"
	"strconv"
	"time"

	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
)

type TestRedisResult struct {
	val interface{}
	TestRedisStatusCmd
	TestRedisIntCmd
	TestStringCmd
}

type TestRedisStatusCmd struct {
	val string
	err error
}

func (r TestRedisStatusCmd) Val() string {
	return r.val
}
func (r TestRedisStatusCmd) Result() (string, error) {
	return r.Val(), nil
}
func (r TestRedisStatusCmd) String() string {
	return r.val
}
func (r TestRedisStatusCmd) Err() error {
	return nil
}

type TestRedisIntCmd struct {
	val int64
	err error
}

func (r TestRedisIntCmd) Val() int64 {
	return r.val
}
func (r TestRedisIntCmd) Result() (int64, error) {
	return r.Val(), nil
}
func (r TestRedisIntCmd) Uint64() (uint64, error) {
	return uint64(r.Val()), nil
}
func (r TestRedisIntCmd) String() string {
	return fmt.Sprintf("%d", r.Val())
}
func (r TestRedisIntCmd) Err() error {
	return nil
}

type TestStringCmd struct {
	val string
}

func (r TestStringCmd) Val() string {
	return r.val
}
func (r TestStringCmd) Result() (string, error) {
	return r.Val(), nil
}
func (r TestStringCmd) Bytes() ([]byte, error) {
	return []byte(r.Val()), nil
}
func (r TestStringCmd) Int() (int, error) {
	return strconv.Atoi(r.Val())
}
func (r TestStringCmd) Int64() (int64, error) {
	return strconv.ParseInt(r.Val(), 10, 64)
}
func (r TestStringCmd) Uint64() (uint64, error) {
	return strconv.ParseUint(r.Val(), 10, 64)
}
func (r TestStringCmd) Float32() (float32, error) {
	f, err := strconv.ParseFloat(r.Val(), 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}
func (r TestStringCmd) Float64() (float64, error) {
	return strconv.ParseFloat(r.Val(), 64)
}
func (r TestStringCmd) Time() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, r.Val())
}
func (r TestStringCmd) Scan(val interface{}) error {
	return nil
}
func (r TestStringCmd) String() string {
	return r.Val()
}

type TestRedisClient struct {
	Data map[string]TestRedisResult
}

func (t TestRedisClient) Subscribe(channels ...string) rs.RedisPubSub {
	panic("implement me")
}

func (t TestRedisClient) Publish(channel string, message interface{}) rs.RedisIntCmd {
	panic("implement me")
}

func (t TestRedisClient) FlushAll() rs.RedisStatusCmd {
	panic("implement me")
}

func (t TestRedisClient) Set(key string, value interface{}, expiration time.Duration) rs.RedisStatusCmd {
	result := TestRedisResult{
		val: value,
	}
	t.Data[key] = result

	return result.TestRedisStatusCmd
}

func (t TestRedisClient) Del(keys ...string) rs.RedisIntCmd {
	for _, v := range keys {
		delete(t.Data, v)
	}

	return TestRedisResult{}.TestRedisIntCmd
}

func (t TestRedisClient) Get(key string) rs.RedisStringCmd {
	result := t.Data[key]
	if res, ok := result.val.(string); ok {
		result.TestStringCmd.val = res
	} else if res, ok := result.val.(int64); ok {
		result.TestStringCmd.val = fmt.Sprintf("%d", res)
	}

	return result.TestStringCmd
}

func (t TestRedisClient) LRange(key string, start, stop int64) rs.RedisStringSliceCmd {
	panic("implement me")
}

func (t TestRedisClient) TxPipeline() rs.RedisPipeliner {
	panic("implement me")
}

func (t TestRedisClient) Pipeline() rs.RedisPipeliner {
	panic("implement me")
}

func (t TestRedisClient) SetNX(key string, value interface{}, expiration time.Duration) rs.RedisBoolCmd {
	panic("implement me")
}

func (t TestRedisClient) HGetAll(key string) rs.RedisStringStringMapCmd {
	panic("implement me")
}

func (t TestRedisClient) Keys(pattern string) rs.RedisStringSliceCmd {
	panic("implement me")
}

func (t TestRedisClient) Ping() rs.RedisStatusCmd {
	panic("implement me")
}
