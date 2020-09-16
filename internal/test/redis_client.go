package test

import (
	rs "github.com/mxc-foundation/lpwan-app-server/internal/modules/redis"
)

type testRedisClient rs.RedisClientType

func (trs *testRedisClient) RSSet(key string, value interface{}) error {

}
func (trs *testRedisClient) RSGet(key string) rs.RSGetResult {

}
func (trs *testRedisClient) RSDelete(key string) error {

}
