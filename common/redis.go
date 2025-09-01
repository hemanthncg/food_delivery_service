package common

import (
    "context"
    "github.com/go-redis/redis/v8"
    "sync"
    "os"
)

var (
    redisClient *redis.Client
    once        sync.Once
)

func GetRedisClient() *redis.Client {
    once.Do(func() {
        addr := os.Getenv("REDIS_ADDR")
        if addr == "" {
            addr = "redis:6379"
        }
        redisClient = redis.NewClient(
            &redis.Options{Addr: addr},
        )
    })
    return redisClient
}

var Ctx = context.Background()
