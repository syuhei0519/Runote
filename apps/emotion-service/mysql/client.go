package mysql

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/syuhei0519/Runote/apps/emotion-service/models"
)

func InitMySQL() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	maxRetries := 10
	var db *gorm.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		log.Printf("🔁 MySQL 接続試行中 (%d/%d)...", i, maxRetries)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("❌ MySQL 接続失敗: %v", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ 最終的に MySQL に接続できませんでした: %v", err)
	}

	// マイグレーション
	err = db.AutoMigrate(&models.Emotion{}, &models.Post{}, &models.PostEmotion{})
	if err != nil {
		log.Fatalf("❌ マイグレーション失敗: %v", err)
	}

	log.Println("✅ MySQL 接続 & マイグレーション完了")
	return db
}