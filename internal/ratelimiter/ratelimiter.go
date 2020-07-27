// Package ratelimiter provides a way to limit the frequency with which certain
// calls can be made
package ratelimiter

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
)

var rateLimitCnt = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "rate_limit_event",
		Help: "the number of rate limit errors",
	},
	[]string{"call_name"},
)

func init() {
	prometheus.MustRegister(rateLimitCnt)
}

// RateLimiter is an object that provides a method to rate limit some calls
type RateLimiter struct {
	r *redis.Pool
}

// New returns a new rate limiter that uses the specified redis instance
func New(r *redis.Pool) *RateLimiter {
	return &RateLimiter{
		r: r,
	}
}

// Check checks if the call should be rate limited. It returns nil if the call
// should go ahead or an error if it should be rejected. name parameter
// identifies the call, id identifies the parameter based on which the call
// should be limited (e.g. IP address of the caller, or email). allowedPeriod
// is how frequently the call can be made. If user attempts to make the call
// more frequently they will be blocked for some time. The duration of the
// block doubles with every call, but it is limited to maxBlock interval.
func (rl *RateLimiter) Check(name, id string, allowedPeriod, maxBlock time.Duration) error {
	conn := rl.r.Get()
	defer conn.Close()

	key := fmt.Sprintf("ratelimiter:%s:%s", name, id)
	allowSeconds := int64(allowedPeriod.Seconds())
	maxSeconds := int64(maxBlock.Seconds())
	res, err := conn.Do("SET", key, 1, "NX", "EX", allowSeconds)
	if err != nil {
		return err
	}
	if res == nil {
		rateLimitCnt.WithLabelValues(name).Inc()
		// key already exists, we should increase the block time. We do it as follows:
		// - we double remaining block time
		// - add allowedRate to the result
		// - make sure it doesn't exceed the maxBlock time
		// - update the TTL
		ttl, err := redis.Int64(conn.Do("TTL", key))
		if err != nil {
			return err
		}
		ttl = ttl*2 + allowSeconds
		if ttl > maxSeconds {
			ttl = maxSeconds
		}
		_, err = conn.Do("EXPIRE", key, ttl)
		if err != nil {
			return err
		}
		return fmt.Errorf("operation has been rate limited")
	}
	return nil
}
