package db

import (
	"github.com/redis/go-redis/v9"
)

type Db struct {
	Client *redis.Client
}

func NewDb() Db {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "admin",
		Password: "admin",
		DB:       0,
	})

	return Db{
		Client: rdb,
	}
}
