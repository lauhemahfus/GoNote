package database

import (
    "fmt"
    "gonote/internal/config"
    "github.com/go-redis/redis/v8"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
        Password: cfg.RedisPassword,
        DB:       0,
    })
}