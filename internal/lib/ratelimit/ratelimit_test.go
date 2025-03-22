package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/M-kos/anti_brutforce/internal/models"
	"github.com/stretchr/testify/require"
)

type StubBucket struct {
	store map[string]models.Bucket
}

func (b StubBucket) AddItem(ctx context.Context, bucket models.Bucket, expiration time.Duration) error {
	b.store[bucket.Key] = bucket

	return nil
}

func (b StubBucket) GetItem(ctx context.Context, key string) (*models.Bucket, error) {
	if v, ok := b.store[key]; ok {
		return &v, nil
	}

	return &models.Bucket{}, nil
}

func (b StubBucket) RemoveItem(ctx context.Context, key string) error {
	delete(b.store, key)
	return nil
}

func TestRatelimit(t *testing.T) {
	tests := []struct {
		title       string
		limit       uint
		interval    time.Duration
		bucket      BucketProvider
		expectedErr error
		expectedOk  bool
	}{
		{
			title:    "failed test",
			limit:    0,
			interval: 60 * time.Second,
			bucket: StubBucket{
				store: map[string]models.Bucket{
					"test": {
						Key:       "test",
						Count:     0,
						StartTime: time.Now(),
					},
				},
			},
			expectedErr: nil,
			expectedOk:  false,
		},
		{
			title:    "successful test",
			limit:    10,
			interval: 60 * time.Second,
			bucket: StubBucket{
				store: map[string]models.Bucket{
					"test": {
						Key:       "test",
						Count:     5,
						StartTime: time.Now(),
					},
				},
			},
			expectedErr: nil,
			expectedOk:  true,
		},
		{
			title:    "exceeded the number of requests test",
			limit:    5,
			interval: 60 * time.Second,
			bucket: StubBucket{
				store: map[string]models.Bucket{
					"test": {
						Key:       "test",
						Count:     5,
						StartTime: time.Now(),
					},
				},
			},
			expectedErr: nil,
			expectedOk:  false,
		},
		{
			title:    "successful test after interval",
			limit:    5,
			interval: 60 * time.Second,
			bucket: StubBucket{
				store: map[string]models.Bucket{
					"test": {
						Key:       "test",
						Count:     0,
						StartTime: time.Now().Add(-61 * time.Second),
					},
				},
			},
			expectedErr: nil,
			expectedOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()

			rl := NewRateLimiter(tt.limit, tt.interval, tt.bucket)

			ok, err := rl.Check(context.Background(), "test")

			require.NoError(t, err)
			require.Equal(t, tt.expectedOk, ok)
		})
	}
}
