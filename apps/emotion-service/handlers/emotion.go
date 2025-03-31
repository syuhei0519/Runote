package handlers

import (
	"net/http"
	"context"
    "log"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/syuhei0519/Runote/apps/emotion-service/redis"
	goredis "github.com/redis/go-redis/v9"
)

type Emotion struct {
    PostID  string `json:"post_id"`
    UserID  string `json:"user_id"`
    Emotion string `json:"emotion"`
}

// Emotion リクエスト用構造体
type EmotionRequest struct {
	PostID   string `json:"post_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	Emotion  string `json:"emotion" binding:"required"` // 例: "嬉しい", "疲れた"
}

// RegisterEmotion 感情ログ登録ハンドラー
func RegisterEmotion(c *gin.Context) {
	var emotion Emotion
	if err := c.ShouldBindJSON(&emotion); err != nil {
		log.Printf("❌ JSON bind エラー: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	ctx := context.Background()
	key := "emotion:" + emotion.PostID + ":" + emotion.UserID

	if err := redis.Client.Set(ctx, key, emotion.Emotion, 0).Err(); err != nil {
		log.Printf("❌ Redis 保存エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis への保存に失敗しました"})
		return
	}

	log.Printf("✅ 感情ログ登録成功: %s", key)
	c.JSON(http.StatusOK, gin.H{
		"message": "Emotion logged",
		"data":    emotion,
	})
}

func GetEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")

	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Emotion not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id": postID,
		"user_id": userID,
		"emotion": val,
	})
}

func UpdateEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	var req struct {
		Emotion string `json:"emotion"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Emotion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// すでに存在するか確認（存在しないなら 404）
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err == goredis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emotion not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
		return
	}

	log.Printf("🔄 旧: %s → 新: %s", val, req.Emotion)

	// 更新
	err = redis.Client.Set(redis.Ctx, key, req.Emotion, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update emotion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Emotion updated",
		"data": gin.H{
			"post_id": postID,
			"user_id": userID,
			"emotion": req.Emotion,
		},
	})
}

func DeleteEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	// 削除前に存在確認
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err == goredis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Emotion not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
		return
	}

	log.Printf("🗑️ 削除対象の感情: %s", val)

	if err := redis.Client.Del(redis.Ctx, key).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete emotion"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
