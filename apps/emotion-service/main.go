package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    "github.com/syuhei0519/Runote/apps/emotion-service/handlers"
    "github.com/syuhei0519/Runote/apps/emotion-service/redis"
)

func main() {
    // dev.env を読み込み（必要に応じて切り替え）
    err := godotenv.Load("dev.env")
    if err != nil {
        log.Println("dev.env が見つかりませんでした（スキップ）")
    }

    // Redis 初期化
    redis.InitRedis()

    // Gin のルーティング設定
    r := gin.Default()
    r.POST("/emotions", handlers.RegisterEmotion)

    log.Println("Starting server on :8080")
    r.Run(":8080")
}