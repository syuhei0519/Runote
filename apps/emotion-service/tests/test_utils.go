package tests

import (
    "path/filepath"
    "runtime"
    "os"
    "testing"
    "github.com/syuhei0519/Runote/apps/emotion-service/redis"
    "github.com/joho/godotenv"
)

func Setup(t *testing.T) {
    _, b, _, _ := runtime.Caller(0)
    basePath := filepath.Dir(b)
    envPath := filepath.Join(basePath, "..", "test.env")

    os.Setenv("ENV", "test")
    if err := godotenv.Load(envPath); err != nil {
        t.Fatalf("❌ test.env が読み込めません: %v", err)
    }

    redis.InitRedis()
    redis.FlushAll()
}
