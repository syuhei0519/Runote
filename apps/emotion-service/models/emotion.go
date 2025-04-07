package models

import "time"

type Emotion struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex"` // ✅ 修正ポイント
	IsPreset  bool      `gorm:"default:false"`
	CreatedAt time.Time
	Posts     []PostEmotion
}

type Post struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	Emotions  []PostEmotion
}

type PostEmotion struct {
	PostID    uint `gorm:"primaryKey"`
	EmotionID uint `gorm:"primaryKey"`
	Intensity int
	Custom    string

	Post    Post
	Emotion Emotion
}