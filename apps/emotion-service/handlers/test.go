package handlers

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func TestCleanupHandler(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("NODE_ENV") != "test" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		if db == nil {
			log.Println("❌ DB is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB is not initialized"})
			return
		}

		if redisClient == nil {
			log.Println("❌ Redis is nil")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis is not initialized"})
			return
		}

		// DBデータ削除
		if err := db.Exec("DELETE FROM emotions").Error; err != nil {
			log.Println("❌ Failed to delete emotions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete emotions"})
			return
		}

		if err := db.Exec("DELETE FROM post_emotions").Error; err != nil {
			log.Println("❌ Failed to delete post_emotions:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post_emotions"})
			return
		}

		// Redis削除
		if err := redisClient.FlushDB(context.Background()).Err(); err != nil {
			log.Println("❌ Failed to flush Redis:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to flush Redis"})
			return
		}

		// 204 No Content には body を返さない
		c.Status(http.StatusNoContent)
	}
}
