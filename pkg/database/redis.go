package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-api-test/config"
	"log"
)

type Redis struct {
	cfg  *config.Config
	conn *redis.Client
}

func NewRedis(cfg *config.Config) *Redis {
	return &Redis{cfg: cfg}
}

func (r *Redis) Init() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.cfg.RedisHost, r.cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	return rdb, nil
}
