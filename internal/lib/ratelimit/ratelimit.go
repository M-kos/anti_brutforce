package ratelimit

import (
	"context"
	"log"
	"time"

	"github.com/M-kos/anti_brutforce/internal/models"
)

type BucketProvider interface {
	GetItem(ctx context.Context, key string) (*models.Bucket, error)
	AddItem(ctx context.Context, bucket models.Bucket, expiration time.Duration) error
	RemoveItem(ctx context.Context, key string) error
}

type RateLimitter struct {
	limit          uint
	interval       time.Duration
	bucketProvider BucketProvider
}

func NewRateLimiter(limit uint, interval time.Duration, bucketProvider BucketProvider) *RateLimitter {
	return &RateLimitter{
		limit:          limit,
		interval:       interval,
		bucketProvider: bucketProvider,
	}
}

func (rl *RateLimitter) Check(ctx context.Context, key string) (bool, error) {
	bucket, err := rl.bucketProvider.GetItem(ctx, key)
	if err != nil {
		log.Println(err.Error()) // TODO: обработать ошибку??? возвращаем true, т.к. не можем проверить сколько было запросов???
	}

	if time.Since(bucket.StartTime) > rl.interval {
		return true, rl.bucketProvider.AddItem(ctx, models.Bucket{Count: 1, Key: key, StartTime: time.Now()}, rl.interval)
	}

	if bucket.Count >= rl.limit {
		return false, nil
	}

	return true, rl.bucketProvider.AddItem(ctx, models.Bucket{Count: bucket.Count + 1, Key: key, StartTime: bucket.StartTime}, rl.interval)
}

func (rl *RateLimitter) Remove(ctx context.Context, key string) error {
	return rl.bucketProvider.RemoveItem(ctx, key)
}
