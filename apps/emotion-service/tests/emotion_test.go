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
)

func TestGetEmotionList_ReturnsRegisteredData(t *testing.T) {
	testutils.SetupTestDB()

	emotions := []models.Emotion{
		{Name: "楽しい", IsPreset: true},
		{Name: "疲れた", IsPreset: false},
	}
	for _, e := range emotions {
		mysql.DB.Create(&e)
	}

	router := testutils.SetupRouter()
	req, _ := http.NewRequest("GET", "/emotions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res []models.Emotion
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "楽しい", res[0].Name)
	assert.Equal(t, "疲れた", res[1].Name)
}

func TestRegisterEmotion_CreatesEmotion(t *testing.T) {
	testutils.SetupTestDB()

	router := testutils.SetupRouter()
	body := `{"name": "エモい"}`
	req, _ := http.NewRequest("POST", "/emotions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res models.Emotion
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "エモい", res.Name)
	assert.False(t, res.IsPreset)
}

func TestUpdateEmotionName_UpdatesSuccessfully(t *testing.T) {
	testutils.SetupTestDB()

	// 感情登録
	emotion := models.Emotion{Name: "仮", IsPreset: false}
	mysql.DB.Create(&emotion)

	// 正しくIDを取得できているか確認
	assert.NotZero(t, emotion.ID)

	router := testutils.SetupRouter()
	body := `{"name": "変更後の名前"}`
	url := fmt.Sprintf("/emotions/%d", emotion.ID)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res models.Emotion
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "変更後の名前", res.Name)
}

func TestDeleteEmotionByID_DeletesUnusedEmotion(t *testing.T) {
	testutils.SetupTestDB()

	emotion := models.Emotion{Name: "消していい", IsPreset: false}
	mysql.DB.Create(&emotion)

	router := testutils.SetupRouter()
	url := fmt.Sprintf("/emotions/%d", emotion.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	mysql.DB.Model(&models.Emotion{}).Where("id = ?", emotion.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteEmotionByID_FailsIfUsed(t *testing.T) {
	testutils.SetupTestDB()

	// 投稿を作成（外部キー制約を回避）
	post := models.Post{}
	mysql.DB.Create(&post)

	// 感情作成
	emotion := models.Emotion{Name: "使用中", IsPreset: false}
	mysql.DB.Create(&emotion)

	postEmotion := models.PostEmotion{
		PostID:    post.ID,
		EmotionID: emotion.ID,
		Intensity: 3,
	}
	mysql.DB.Create(&postEmotion)

	router := testutils.SetupRouter()
	url := fmt.Sprintf("/emotions/%d", emotion.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}