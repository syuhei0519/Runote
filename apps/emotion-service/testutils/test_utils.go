package testutils

import (
    "os"
    "github.com/syuhei0519/Runote/apps/emotion-service/mysql"
	"github.com/gin-gonic/gin"
	"github.com/syuhei0519/Runote/apps/emotion-service/handlers"
)

// func Setup(t *testing.T) {
//     _, b, _, _ := runtime.Caller(0)
//     basePath := filepath.Dir(b)
//     envPath := filepath.Join(basePath, "..", "test.env")

//     os.Setenv("ENV", "test")
//     if err := godotenv.Load(envPath); err != nil {
//         t.Fatalf("❌ test.env が読み込めません: %v", err)
//     }

//     redis.InitRedis()
//     redis.FlushAll()
// }

func SetupTestDB() {
	// test用の環境変数（DATABASE_URL）をセット
	os.Setenv("DATABASE_URL", "emotion:emotionpass@tcp(emotion-db:3306)/emotiondb?charset=utf8mb4&parseTime=True&loc=Local")

	// DB初期化（mysql/client.goのInit関数を使う）
	mysql.InitMySQL()

	// テスト用に全削除（安全確認してから）
	mysql.DB.Exec("DELETE FROM post_emotions")
	mysql.DB.Exec("DELETE FROM emotions")
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Emotionマスタ関連
	r.GET("/emotions", handlers.GetEmotionList)
	r.GET("/emotions/unused", handlers.GetUnusedEmotions)
	r.POST("/emotions", handlers.RegisterEmotion)
	r.PUT("/emotions/:id", handlers.UpdateEmotionName)
	r.DELETE("/emotions/:id", handlers.DeleteEmotionByID)

	// 投稿に紐づく感情関連
	r.GET("/post-emotions/:post_id/:user_id", handlers.GetEmotion)
	r.PUT("/post-emotions/:post_id/:user_id", handlers.UpdateEmotion)
	r.DELETE("/post-emotions/:post_id/:user_id", handlers.DeleteEmotion)

	return r
}
