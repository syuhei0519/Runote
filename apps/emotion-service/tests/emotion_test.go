package tests

import (
	"context"
	"os"
	"testing"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/syuhei0519/Runote/apps/emotion-service/handlers"
	"github.com/syuhei0519/Runote/apps/emotion-service/redis"
	"net/http"
	"net/http/httptest"
	"strings"
)

// テスト用のRedisキー生成ヘルパー
func redisKey(postID, userID string) string {
	return "emotion:" + postID + ":" + userID
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/emotions/:post_id/:user_id", handlers.GetEmotion)
	r.POST("/emotions", handlers.RegisterEmotion)
	return r
}

func TestRegisterEmotion(t *testing.T) {
	_ = os.Setenv("REDIS_URL", "redis://redis:6379/1")
	redis.InitRedis() // Client に接続

	r := setupRouter()

	postID := "testPost123"
	userID := "testUser456"
	emotion := "ワクワク"

	payload := `{"post_id": "` + postID + `", "user_id": "` + userID + `", "emotion": "` + emotion + `"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/emotions", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Redis に保存されているか確認
	val, err := redis.Client.Get(context.Background(), redisKey(postID, userID)).Result()
	assert.NoError(t, err)
	assert.Equal(t, emotion, val)
}

// 今後の GET / PUT / DELETE に備えて、helper も追加しておく
func clearEmotionKey(postID, userID string) {
	_ = redis.Client.Del(context.Background(), redisKey(postID, userID)).Err()
}

func TestGetEmotion(t *testing.T) {
	Setup(t) // test.env読み込みとRedis初期化

	// テスト用データ
	postID := "testPost123"
	userID := "testUser456"
	emotion := "やった！"

	// Redis に事前に保存
	key := fmt.Sprintf("emotion:%s:%s", postID, userID)
	err := redis.Client.Set(redis.Ctx, key, emotion, 0).Err()
	assert.NoError(t, err, "Redisへの保存でエラー")

	// Gin のルーターセットアップ
	r := gin.Default()
	r.GET("/emotions/:post_id/:user_id", handlers.GetEmotion)

	// リクエスト送信
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/emotions/%s/%s", postID, userID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)
	expected := fmt.Sprintf(`{"emotion":"%s","post_id":"%s","user_id":"%s"}`, emotion, postID, userID)
	assert.JSONEq(t, expected, w.Body.String())
}