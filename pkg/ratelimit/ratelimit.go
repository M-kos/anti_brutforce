package ratelimit

import (
	"fmt"
	"time"
)

type Bucket interface {
  Get(key string) (uint, time.Time, error)
  Set(key string, count uint, startTime time.Time) error
}

type rateLimiter struct {
  limit uint
	interval time.Duration
	bucket Bucket
}

func (rl *rateLimiter) Check(key string) (bool, error) {
	count, startTime, err := rl.bucket.Get(key)

	fmt.Println("count", count, "startTime", startTime, "err", err)

	if err != nil {
		return true, err // TODO: обработать ошибку??? возвращаем true, т.к. не можем проверить сколько было запросов???
	}

	if time.Since(startTime) > rl.interval {
		return true, rl.bucket.Set(key, 1, time.Now())
	}

	if count >= rl.limit {
		return false, nil
	}

	return true, rl.bucket.Set(key, count + 1, startTime)
}

func NewRateLimiter(limit uint, interval time.Duration, bucket Bucket) *rateLimiter {
  return &rateLimiter{
    limit: limit,
    interval: interval,
    bucket: bucket,
  }
}

