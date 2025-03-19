package ratelimit

import (
	"context"
	"log"
	"time"
)

type Bucket interface {
	GetItem(ctx context.Context, key string) (uint, time.Time, error)
	AddItem(ctx context.Context, key string, count uint, startTime time.Time) error
	RemoveItem(ctx context.Context, key string) error
}

type RateLimitter struct {
	limit    uint
	interval time.Duration
	bucket   Bucket
}

func (rl *RateLimitter) Check(ctx context.Context, key string) (bool, error) {
	count, startTime, err := rl.bucket.GetItem(ctx, key)
	if err != nil {
		log.Println(err.Error()) // TODO: обработать ошибку??? возвращаем true, т.к. не можем проверить сколько было запросов???
	}

	if time.Since(startTime) > rl.interval {
		return true, rl.bucket.AddItem(ctx, key, 1, time.Now())
	}

	if count >= rl.limit {
		return false, nil
	}

	return true, rl.bucket.AddItem(ctx, key, count+1, startTime)
}

func (rl *RateLimitter) Remove(ctx context.Context, key string) error {
	return rl.bucket.RemoveItem(ctx, key)
}

func NewRateLimiter(limit uint, interval time.Duration, bucket Bucket) *RateLimitter {
	return &RateLimitter{
		limit:    limit,
		interval: interval,
		bucket:   bucket,
	}
}
