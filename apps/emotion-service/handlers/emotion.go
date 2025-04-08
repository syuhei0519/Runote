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

type SuccessResponse struct {
	Message string `json:"message" example:"操作が成功しました"`
}

// ErrorResponse 汎用のエラーレスポンス構造体
type ErrorResponse struct {
	Error string `json:"error" example:"不正な入力です"`
}

// RegisterEmotion godoc
// @Summary 感情を新規登録する
// @Description ユーザー定義の感情を登録する（重複不可）
// @Tags emotions
// @Accept  json
// @Produce  json
// @Param  emotion body EmotionCreateRequest true "感情登録リクエスト"
// @Success 201 {object} models.Emotion
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 409 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /emotions [post]
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

// GetEmotionList godoc
// @Summary 感情一覧を取得
// @Description すべての登録感情を取得
// @Tags emotions
// @Produce  json
// @Success 200 {array} models.Emotion
// @Failure 500 {object} handlers.ErrorResponse
// @Router /emotions [get]
func GetEmotionList(c *gin.Context) {
	var emotions []models.Emotion
	if err := mysql.DB.Order("created_at ASC").Find(&emotions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取得に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, emotions)
}

// GetUnusedEmotions godoc
// @Summary 未使用感情一覧を取得
// @Description 投稿に紐づいていない感情のみを返す
// @Tags emotions
// @Produce  json
// @Success 200 {array} models.Emotion
// @Failure 500 {object} handlers.ErrorResponse
// @Router /emotions/unused [get]
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

// GetEmotion godoc
// @Summary 投稿に紐づく感情を取得
// @Description 投稿ID・ユーザーIDに紐づく感情をRedisまたはDBから取得
// @Tags post-emotions
// @Produce  json
// @Param  post_id path string true "投稿ID"
// @Param  user_id path string true "ユーザーID"
// @Success 200 {object} models.PostEmotion
// @Failure 404 {object} handlers.ErrorResponse
// @Router /post-emotions/{post_id}/{user_id} [get]
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

// UpdateEmotion godoc
// @Summary 投稿に紐づく感情を更新
// @Description 投稿ID・ユーザーIDを指定して感情を更新、キャッシュも更新
// @Tags post-emotions
// @Accept  json
// @Produce  json
// @Param  post_id path string true "投稿ID"
// @Param  user_id path string true "ユーザーID"
// @Param  emotion body PostEmotionUpdateRequest true "感情更新内容"
// @Success 200 {object} models.PostEmotion
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /post-emotions/{post_id}/{user_id} [put]
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

// DeleteEmotion godoc
// @Summary 投稿に紐づく感情を削除
// @Description 指定された投稿ID・ユーザーIDに紐づく感情を削除、Redisキャッシュも削除
// @Tags post-emotions
// @Produce  json
// @Param  post_id path string true "投稿ID"
// @Param  user_id path string true "ユーザーID"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /post-emotions/{post_id}/{user_id} [delete]
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

// UpdateEmotionName godoc
// @Summary 感情名の更新
// @Description 感情の名前を変更する
// @Tags emotions
// @Accept  json
// @Produce  json
// @Param  id path string true "感情ID"
// @Param  emotion body EmotionNameUpdateRequest true "感情名更新内容"
// @Success 200 {object} models.Emotion
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /emotions/{id} [put]
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

// DeleteEmotionByID godoc
// @Summary 感情をIDで削除
// @Description 投稿で使用されていない感情を削除する
// @Tags emotions
// @Produce  json
// @Param  id path string true "感情ID"
// @Success 200 {object} handlers.SuccessResponse
// @Failure 409 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /emotions/{id} [delete]
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

