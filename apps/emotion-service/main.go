package main

import (
    "log"
	"os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"

    "github.com/syuhei0519/Runote/apps/emotion-service/handlers"
    "github.com/syuhei0519/Runote/apps/emotion-service/redis"
    "github.com/syuhei0519/Runote/apps/emotion-service/mysql"
)

func main() {
    env := os.Getenv("ENV")
    envFile := "dev.env"
    if env == "test" {
        envFile = "test.env"
    }

    if err := godotenv.Load(envFile); err != nil {
        log.Printf("%s が見つかりませんでした（スキップ）\n", envFile)
    }

    // Redis 初期化
	log.Println("🔁 Redis 接続開始")
    redis.InitRedis()
	log.Println("✅ Redis 接続完了")

    // MySQL初期化
    log.Println("🔁 MySQL 接続開始")
    mysql.InitMySQL()
	log.Println("✅ MySQL 接続完了")

    // Gin のルーティング設定
    r := gin.Default()
    r.GET("/emotions/:post_id/:user_id", handlers.GetEmotion)
    r.PUT("/emotions/:post_id/:user_id", handlers.UpdateEmotion)
    r.POST("/emotions", handlers.RegisterEmotion)
    r.DELETE("/emotions/:post_id/:user_id", handlers.DeleteEmotion)

    log.Println("Starting server on :8080")
    r.Run(":8080")
}