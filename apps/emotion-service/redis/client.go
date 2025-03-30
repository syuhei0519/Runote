package redis

import (
    "context"
    "github.com/redis/go-redis/v9"
    "log"
    "os"
)

var Client *redis.Client
var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() {
    Client = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
        Password: "",
        DB:       0,
    })

    _, err := Client.Ping(Ctx).Result()
    if err != nil {
        log.Fatalf("Redis connection failed: %v", err)
    }
}