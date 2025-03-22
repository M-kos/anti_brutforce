package db

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/M-kos/anti_brutforce/internal/config"
	"github.com/M-kos/anti_brutforce/internal/lib/logger"
	"github.com/M-kos/anti_brutforce/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	Client *redis.Client
	log    logger.LoggerProvider
}

func NewRedisDb(conf *config.Config, log logger.LoggerProvider) *RedisDb {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.DbHost, conf.DbPort),
		Username: "admin", // TODO change it
		Password: "admin",
		DB:       0,
	})

	return &RedisDb{
		Client: rdb,
		log:    log,
	}
}

func (rdb *RedisDb) GetItem(ctx context.Context, key string) (*models.Bucket, error) {
	valStr, err := rdb.Client.Get(ctx, key).Result()
	if err != nil {
		rdb.log.Error(err.Error())
		return nil, fmt.Errorf("Redis: GetItem: failed to get value from db: %w", err)
	}

	var val models.Bucket

	err = json.Unmarshal([]byte(valStr), &val)
	if err != nil {
		rdb.log.Error(err.Error())
		return nil, fmt.Errorf("Redis: GetItem: failed to unmarshal value from db: %w", err)
	}

	return &val, nil
}

func (rdb *RedisDb) AddItem(ctx context.Context, bucket models.Bucket, expiration time.Duration) error {
	data, err := json.Marshal(&bucket)
	if err != nil {
		rdb.log.Error(err.Error())
		return fmt.Errorf("Redis: AddItem: failed to marshal value for db: %w", err)
	}

	_, err = rdb.Client.Set(ctx, bucket.Key, string(data), time.Hour).Result()
	if err != nil {
		rdb.log.Error(err.Error())
		return fmt.Errorf("Redis: AddItem: failed to set value to db: %w", err)
	}

	return nil
}

func (rdb *RedisDb) RemoveItem(ctx context.Context, key string) error {
	_, err := rdb.Client.Del(ctx, key).Result()
	if err != nil {
		rdb.log.Error(err.Error())
		return fmt.Errorf("Redis: RemoveItem: failed to remove value from db: %w", err)
	}

	return nil
}

func (rdb *RedisDb) GetList(ctx context.Context, key string) (*models.LabelList, error) {
	valStr, err := rdb.Client.Get(ctx, key).Result()
	if err != nil {
		rdb.log.Error(err.Error())
		return nil, fmt.Errorf("Redis: GetList: failed to get list from db: %w", err)
	}

	var val models.LabelList

	err = json.Unmarshal([]byte(valStr), &val)
	if err != nil {
		rdb.log.Error(err.Error())
		return nil, fmt.Errorf("Redis: GetList: failed to unmarshal list from db: %w", err)
	}

	return &val, nil
}

func (rdb *RedisDb) AddToList(ctx context.Context, key string, value string) error {
	list, err := rdb.GetList(ctx, key)
	if err != nil {
		return err
	}

	list.Values = append(list.Values, value)

	return rdb.setList(ctx, key, list)
}

func (rdb *RedisDb) RemoveFromList(ctx context.Context, key string, value string) error {
	list, err := rdb.GetList(ctx, key)
	if err != nil {
		return err
	}

	for i, v := range list.Values {
		if v == value {
			list.Values = slices.Delete(list.Values, i, i+1)
			break
		}
	}

	return rdb.setList(ctx, key, list)
}

func (rdb *RedisDb) setList(ctx context.Context, key string, list *models.LabelList) error {
	data, err := json.Marshal(list)
	if err != nil {
		rdb.log.Error(err.Error())
		return fmt.Errorf("Redis: setList: failed to marshal list to db: %w", err)
	}

	_, err = rdb.Client.Set(ctx, key, string(data), 72*time.Hour).Result()
	if err != nil {
		rdb.log.Error(err.Error())
		return fmt.Errorf("Redis: setList: failed to set list to db: %w", err)
	}

	return nil
}
