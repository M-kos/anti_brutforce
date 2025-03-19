package bucketlist

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/M-kos/anti_brutforce/internal/db"
)

type CredentialItem struct {
	Count     uint      `json:"count"`
	StartTime time.Time `json:"startTime"`
}

type List struct {
	Items []string `json:"item"`
}

type DbBucket struct {
	db db.Db
}

func NewDbBucketList(db db.Db) DbBucket {
	return DbBucket{db: db}
}

func (b DbBucket) GetItem(ctx context.Context, key string) (uint, time.Time, error) {
	valStr, err := b.db.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get value from db: %s", err.Error())
	}

	var val CredentialItem

	err = json.Unmarshal([]byte(valStr), &val)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to unmarshal value from db: %s", err.Error())
	}

	return val.Count, val.StartTime, nil
}

func (b DbBucket) AddItem(ctx context.Context, key string, count uint, startTime time.Time) error {
	data, err := json.Marshal(&CredentialItem{Count: count, StartTime: startTime})
	if err != nil {
		return fmt.Errorf("failed to marshal value for db: %s", err.Error())
	}

	_, err = b.db.Client.Set(ctx, key, string(data), time.Hour).Result()
	if err != nil {
		return fmt.Errorf("failed to set value to db: %s", err.Error())
	}

	return nil
}

func (b DbBucket) RemoveItem(ctx context.Context, key string) error {
	_, err := b.db.Client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to remove value from db: %s", err.Error())
	}

	return nil
}

func (b DbBucket) GetList(ctx context.Context, key string) ([]string, error) {
	valStr, err := b.db.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get list from db: %s", err.Error())
	}

	var val List

	err = json.Unmarshal([]byte(valStr), &val)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal list from db: %s", err.Error())
	}

	return val.Items, nil
}

func (b DbBucket) AddToList(ctx context.Context, key string, value string) error {
	list, err := b.GetList(ctx, key)
	if err != nil {
		return err
	}

	list = append(list, value)

	return b.setToList(ctx, key, list)
}

func (b DbBucket) RemoveFromList(ctx context.Context, key string, value string) error {
	list, err := b.GetList(ctx, key)
	if err != nil {
		return err
	}

	for i, v := range list {
		if v == value {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}

	return b.setToList(ctx, key, list)
}

func (b DbBucket) setToList(ctx context.Context, key string, list []string) error {
	data, err := json.Marshal(&List{Items: list})
	if err != nil {
		return fmt.Errorf("failed to marshal list to db: %s", err.Error())
	}

	_, err = b.db.Client.Set(ctx, key, string(data), time.Hour*72).Result()
	if err != nil {
		return fmt.Errorf("failed to set list to db: %s", err.Error())
	}

	return nil
}
