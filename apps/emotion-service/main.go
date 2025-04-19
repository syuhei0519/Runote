package main

import (
    "log"
	"os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"

    "github.com/syuhei0519/Runote/apps/emotion-service/handlers"
    myredis "github.com/syuhei0519/Runote/apps/emotion-service/redis"
    "github.com/syuhei0519/Runote/apps/emotion-service/mysql"

    "gorm.io/gorm"

    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "github.com/syuhei0519/Runote/apps/emotion-service/docs"
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

    var db *gorm.DB
	var redisClient *redis.Client

    // Redis 初期化
	log.Println("🔁 Redis 接続開始")
    redisClient = myredis.InitRedis()
	log.Println("✅ Redis 接続完了")

    // MySQL初期化
    log.Println("🔁 MySQL 接続開始")
    db = mysql.InitMySQL()
	log.Println("✅ MySQL 接続完了")

    r := gin.Default()

    if os.Getenv("NODE_ENV") == "test" {
        r.POST("/test/cleanup", handlers.TestCleanupHandler(db, redisClient))
    }

    // Swagger UI を追加
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Gin のルーティング設定
    // Emotion マスタ
    r.POST("/emotions", handlers.RegisterEmotion(db))
    r.GET("/emotions", handlers.GetEmotionList(db))
    r.GET("/emotions/unused", handlers.GetUnusedEmotions(db))
    r.PUT("/emotions/:id", handlers.UpdateEmotionName(db))
    r.DELETE("/emotions/:id", handlers.DeleteEmotionByID(db))

    // 投稿に紐づく感情
    r.GET("/post-emotions/:post_id/:user_id", handlers.GetEmotion(db, redisClient))
    r.PUT("/post-emotions/:post_id/:user_id", handlers.UpdateEmotion(db, redisClient))
    r.DELETE("/post-emotions/:post_id/:user_id", handlers.DeleteEmotion(db, redisClient))


    log.Println("Starting server on :8080")
    r.Run(":8080")
}