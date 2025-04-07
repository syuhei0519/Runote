package handlers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syuhei0519/Runote/apps/emotion-service/redis"
	"github.com/syuhei0519/Runote/apps/emotion-service/models"
	"github.com/syuhei0519/Runote/apps/emotion-service/mysql"
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

type EmotionCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type EmotionNameUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}

type PostEmotionUpdateRequest struct {
	EmotionID uint `json:"emotion_id" binding:"required"`
	Intensity int  `json:"intensity" binding:"required"`
}

// RegisterEmotion 感情ログ登録ハンドラー
func RegisterEmotion(c *gin.Context) {
	var req EmotionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	// 重複チェック
	var existing models.Emotion
	if err := mysql.DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "すでに存在する感情です"})
		return
	}

	emotion := models.Emotion{
		Name:     req.Name,
		IsPreset: false,
	}

	if err := mysql.DB.Create(&emotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "感情の登録に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, emotion)
}

func GetEmotionList(c *gin.Context) {
	var emotions []models.Emotion
	if err := mysql.DB.Order("created_at ASC").Find(&emotions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, emotions)
}

func GetUnusedEmotions(c *gin.Context) {
	var emotions []models.Emotion

	err := mysql.DB.
		Raw(`
			SELECT * FROM emotions 
			WHERE id NOT IN (
				SELECT DISTINCT emotion_id FROM post_emotions
			)
		`).Scan(&emotions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "未使用感情の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, emotions)
}

func GetEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	// Redisヒットチェック
	val, err := redis.Client.Get(redis.Ctx, key).Result()
	if err == nil {
		var cached models.PostEmotion
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// DB取得
	var emotion models.PostEmotion
	err = mysql.DB.
		Preload("Emotion").
		Where("post_id = ? AND user_id = ?", postID, userID).
		First(&emotion).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "感情が見つかりません"})
		return
	}

	// Redisに保存
	jsonBytes, _ := json.Marshal(emotion)
	redis.Client.Set(redis.Ctx, key, jsonBytes, time.Hour)

	c.JSON(http.StatusOK, emotion)
}

func UpdateEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	var req PostEmotionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正な入力"})
		return
	}

	// レコード更新
	var emotion models.PostEmotion
	err := mysql.DB.
		Model(&emotion).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Updates(map[string]interface{}{
			"emotion_id": req.EmotionID,
			"intensity":  req.Intensity,
		}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗"})
		return
	}

	// 再取得してキャッシュ更新
	mysql.DB.Preload("Emotion").
		Where("post_id = ? AND user_id = ?", postID, userID).
		First(&emotion)

	jsonBytes, _ := json.Marshal(emotion)
	redis.Client.Set(redis.Ctx, key, jsonBytes, time.Hour)

	c.JSON(http.StatusOK, emotion)
}

func DeleteEmotion(c *gin.Context) {
	postID := c.Param("post_id")
	userID := c.Param("user_id")
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)

	err := mysql.DB.
		Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&models.PostEmotion{}).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "削除に失敗しました"})
		return
	}

	// Redisキャッシュ削除
	redis.Client.Del(redis.Ctx, key)

	c.JSON(http.StatusOK, gin.H{"message": "削除完了"})
}

func UpdateEmotionName(c *gin.Context) {
	id := c.Param("id")

	var req EmotionNameUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正な入力です"})
		return
	}

	var emotion models.Emotion
	if err := mysql.DB.First(&emotion, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "感情が見つかりません"})
		return
	}

	// 名前更新
	emotion.Name = req.Name
	if err := mysql.DB.Save(&emotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "感情の更新に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, emotion)
}


func DeleteEmotionByID(c *gin.Context) {
	id := c.Param("id")

	// 使用中チェック
	var count int64
	mysql.DB.Model(&models.PostEmotion{}).
		Where("emotion_id = ?", id).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "この感情は使用中のため削除できません"})
		return
	}

	if err := mysql.DB.Delete(&models.Emotion{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "感情の削除に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "感情を削除しました"})
}

