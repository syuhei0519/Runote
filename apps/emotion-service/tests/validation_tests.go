
package tests

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/syuhei0519/Runote/apps/emotion-service/mysql"
    "github.com/syuhei0519/Runote/apps/emotion-service/models"
    "github.com/syuhei0519/Runote/apps/emotion-service/testutils"
)

func TestRegisterEmotion_InvalidInput(t *testing.T) {
    testutils.Setup(t)
    router := testutils.SetupRouter()

    body := `{}`
    req, _ := http.NewRequest("POST", "/emotions", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    var res map[string]string
    json.Unmarshal(w.Body.Bytes(), &res)
    assert.Equal(t, "不正な入力です", res["error"])
}

func TestUpdatePostEmotion_InvalidInput(t *testing.T) {
    testutils.Setup(t)
    router := testutils.SetupRouter()

    body := `{"emotion_id": "invalid", "intensity": 3}`
    req, _ := http.NewRequest("PUT", "/post-emotions/1/1", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    var res map[string]string
    json.Unmarshal(w.Body.Bytes(), &res)
    assert.Equal(t, "不正な入力です", res["error"])
}

func TestUpdateEmotionName_EmptyName(t *testing.T) {
    testutils.Setup(t)

    emotion := models.Emotion{Name: "仮", IsPreset: false}
    mysql.DB.Create(&emotion)

    router := testutils.SetupRouter()

    body := `{"name": ""}`
    url := fmt.Sprintf("/emotions/%d", emotion.ID)
    req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    var res map[string]string
    json.Unmarshal(w.Body.Bytes(), &res)
    assert.Equal(t, "不正な入力です", res["error"])
}

func TestDeletePostEmotion_InvalidIDs(t *testing.T) {
    testutils.Setup(t)
    router := testutils.SetupRouter()

    req, _ := http.NewRequest("DELETE", "/post-emotions/9999/9999", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
    var res map[string]string
    json.Unmarshal(w.Body.Bytes(), &res)
    assert.Contains(t, res["error"], "not found")
}