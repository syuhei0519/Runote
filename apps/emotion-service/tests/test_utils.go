package tests

import (
    "os"
    "testing"
    "github.com/syuhei0519/Runote/apps/emotion-service/redis"
    "github.com/joho/godotenv"
)

func Setup(t *testing.T) {
    os.Setenv("ENV", "test")
    if err := godotenv.Load("test.env"); err != nil {
        t.Fatal("❌ test.env が読み込めません")
    }
    redis.InitRedis()
    redis.FlushAll()
}
