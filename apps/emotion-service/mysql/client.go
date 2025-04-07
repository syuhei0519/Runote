package mysql

import (
	"log"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/syuhei0519/Runote/apps/emotion-service/models"
)

var DB *gorm.DB

func InitMySQL() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ DB接続失敗: %v", err)
	}

	// 必要なモデルをマイグレーション
	err = db.AutoMigrate(&models.Emotion{}, &models.Post{}, &models.PostEmotion{})
	if err != nil {
		log.Fatalf("❌ マイグレーション失敗: %v", err)
	}

	DB = db
	log.Println("✅ DB接続完了 & マイグレーション完了")
}