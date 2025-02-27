package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Stub8ucket struct {
	store map[string]struct {
		count     uint
		startTime time.Time
	}
}

func (b *Stub8ucket) Set(key string, count uint, startTime time.Time) error {
	b.store[key] = struct {
		count     uint
		startTime time.Time
	}{
		count:     count,
		startTime: startTime,
	}

	return nil
}

func (b *Stub8ucket) Get(key string) (uint, time.Time, error) {
	if v, ok := b.store[key]; ok {
		return v.count, v.startTime, nil
	}

	return 0, time.Time{}, nil
}

func TestRatelimit(t *testing.T) {
	tests := []struct {
		title       string
		limit       uint
		interval    time.Duration
		bucket      Bucket
		expectedErr error
		expectedOk  bool
	}{
		{
			title:    "failed test",
			limit:    0,
			interval: 60 * time.Second,
			bucket: &Stub8ucket{
				store: map[string]struct {
					count     uint
					startTime time.Time
				}{
					"test": {
						count:     0,
						startTime: time.Now(),
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
			bucket: &Stub8ucket{
				store: map[string]struct {
					count     uint
					startTime time.Time
				}{
					"test": {
						count:     5,
						startTime: time.Now(),
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
			bucket: &Stub8ucket{
				store: map[string]struct {
					count     uint
					startTime time.Time
				}{
					"test": {
						count:     5,
						startTime: time.Now(),
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
			bucket: &Stub8ucket{
				store: map[string]struct {
					count     uint
					startTime time.Time
				}{
					"test": {
						count:     0,
						startTime: time.Now().Add(-61 * time.Second),
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

			ok, err := rl.Check("test")

			require.NoError(t, err)
			require.Equal(t, tt.expectedOk, ok)
		})
	}
}
