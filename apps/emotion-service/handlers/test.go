func TestCleanupHandler(db *gorm.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("NODE_ENV") != "test" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		// DBデータ削除
		db.Exec("DELETE FROM emotions")
		db.Exec("DELETE FROM post_emotions")

		// Redis削除
		redis.FlushDB(c) // DB0のみ

		c.JSON(http.StatusNoContent, nil)
	}
}
