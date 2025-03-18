package bucketlist

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/M-kos/anti_brutforce/internal/db"
)

type BucketItem struct {
	Count     uint      `json:"count"`
	StartTime time.Time `json:"startTime"`
}

type DbBucket struct {
	db db.Db
}

func NewDbBucket(db db.Db) DbBucket {
	return DbBucket{db: db}
}

func (b DbBucket) Get(ctx context.Context, key string) (uint, time.Time, error) {
	valStr, err := b.db.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get value from db: %s", err.Error())
	}

	var val BucketItem

	err = json.Unmarshal([]byte(valStr), &val)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to unmarshal value from db: %s", err.Error())
	}

	return val.Count, val.StartTime, nil
}

func (b DbBucket) Set(ctx context.Context, key string, count uint, startTime time.Time) error {
	data, err := json.Marshal(&BucketItem{Count: count, StartTime: startTime})
	if err != nil {
		return fmt.Errorf("failed to marshal value for db: %s", err.Error())
	}

	_, err = b.db.Client.Set(ctx, key, string(data), time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to set value to db: %s", err.Error())
	}

	return nil
}

func (b DbBucket) Remove(ctx context.Context, key string) error {
	_, err := b.db.Client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to remove value from db: %s", err.Error())
	}

	return nil
}
