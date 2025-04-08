package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syuhei0519/Runote/apps/emotion-service/models"
	"github.com/syuhei0519/Runote/apps/emotion-service/mysql"
	"github.com/syuhei0519/Runote/apps/emotion-service/testutils"
	redisPkg "github.com/syuhei0519/Runote/apps/emotion-service/redis"
	redislib "github.com/redis/go-redis/v9"
)

func TestUpdatePostEmotion_UpsertsData(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	// 必要なデータ作成
	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "元気", IsPreset: true}
	mysql.DB.Create(&emotion)

	router := testutils.SetupRouter()
	body := fmt.Sprintf(`{"emotion_id": %d, "intensity": 4}`, emotion.ID)
	url := fmt.Sprintf("/post-emotions/%d/1", post.ID) // user_id = 1
	req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.PostEmotion
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, emotion.ID, res.EmotionID)
	assert.Equal(t, 4, res.Intensity)
}

func TestGetPostEmotion_ReturnsDataFromDB(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "ワクワク", IsPreset: true}
	mysql.DB.Create(&emotion)

	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		UserID:    1,
		EmotionID: emotion.ID,
		Intensity: 5,
	}
	mysql.DB.Create(&postEmotion)

	router := testutils.SetupRouter()
	url := fmt.Sprintf("/post-emotions/%d/1", post.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.PostEmotion
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, 5, res.Intensity)
	assert.Equal(t, emotion.ID, res.EmotionID)
}

func TestDeletePostEmotion_DeletesRecord(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "モヤモヤ", IsPreset: false}
	mysql.DB.Create(&emotion)

	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		UserID:    1,
		EmotionID: emotion.ID,
		Intensity: 2,
	}
	mysql.DB.Create(&postEmotion)

	router := testutils.SetupRouter()
	url := fmt.Sprintf("/post-emotions/%d/1", post.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	mysql.DB.Model(&models.PostEmotion{}).
		Where("post_id = ? AND user_id = ?", post.ID, 1).
		Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestGetPostEmotion_ReturnsDataFromDBAndCache(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "ワクワク", IsPreset: true}
	mysql.DB.Create(&emotion)

	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		UserID:    1,
		EmotionID: emotion.ID,
		Intensity: 5,
	}
	mysql.DB.Create(&postEmotion)

	router := testutils.SetupRouter()

	url := fmt.Sprintf("/post-emotions/%d/1", post.ID)
	req, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.PostEmotion
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, 5, res.Intensity)

	// ✅ Redisキャッシュ確認
	key := fmt.Sprintf("emotion:%d:%d", post.ID, 1)
	val, err := redisPkg.Client.Get(redisPkg.Ctx, key).Result()
	assert.NoError(t, err)

	var cached models.PostEmotion
	err = json.Unmarshal([]byte(val), &cached)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, cached.PostID)
	assert.Equal(t, emotion.ID, cached.EmotionID)
	assert.Equal(t, 5, cached.Intensity)
}

func TestUpdatePostEmotion_UpdatesDBAndCache(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	// 投稿と感情を作成
	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "楽しい", IsPreset: true}
	mysql.DB.Create(&emotion)

	// 既存のPostEmotionを作成
	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		UserID:    1,
		EmotionID: emotion.ID,
		Intensity: 2,
	}
	mysql.DB.Create(&postEmotion)

	router := testutils.SetupRouter()

	// リクエストデータ（強度の更新）
	body := fmt.Sprintf(`{"emotion_id": %d, "intensity": 4}`, emotion.ID)
	url := fmt.Sprintf("/post-emotions/%d/1", post.ID)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// ステータス確認
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスの確認
	var updated models.PostEmotion
	err := json.Unmarshal(w.Body.Bytes(), &updated)
	assert.NoError(t, err)
	assert.Equal(t, 4, updated.Intensity)

	// ✅ Redisキャッシュの確認
	key := fmt.Sprintf("emotion:%d:%d", post.ID, 1)
	val, err := redisPkg.Client.Get(redisPkg.Ctx, key).Result()
	assert.NoError(t, err)

	var cached models.PostEmotion
	err = json.Unmarshal([]byte(val), &cached)
	assert.NoError(t, err)
	assert.Equal(t, 4, cached.Intensity)
	assert.Equal(t, post.ID, cached.PostID)
	assert.Equal(t, emotion.ID, cached.EmotionID)
}

func TestDeletePostEmotion_DeletesDBAndCache(t *testing.T) {
	testutils.SetupTestDB()
	testutils.Setup(t)

	// 投稿と感情を作成
	post := models.Post{}
	mysql.DB.Create(&post)

	emotion := models.Emotion{Name: "切ない", IsPreset: true}
	mysql.DB.Create(&emotion)

	// PostEmotionを作成
	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		UserID:    1,
		EmotionID: emotion.ID,
		Intensity: 5,
	}
	mysql.DB.Create(&postEmotion)

	// Redisにもキャッシュ登録
	key := fmt.Sprintf("emotion:%d:%d", post.ID, 1)
	cacheData, _ := json.Marshal(postEmotion)
	err := redisPkg.Client.Set(redisPkg.Ctx, key, cacheData, 0).Err()
	assert.NoError(t, err)

	// 確実にキャッシュされていることを確認
	val, err := redisPkg.Client.Get(redisPkg.Ctx, key).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, val)

	router := testutils.SetupRouter()

	// DELETE 実行
	url := fmt.Sprintf("/post-emotions/%d/1", post.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ✅ Redisからキーが削除されていることを確認
	_, err = redisPkg.Client.Get(redisPkg.Ctx, key).Result()
	assert.ErrorIs(t, err, redislib.Nil) 
}