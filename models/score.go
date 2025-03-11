package models

import "time"

type TotalScore struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id"`
	Score     int       `json:"total_score"`
	SongID    int       `json:"song_id"`
	CreatedAT time.Time `json:"created_at"`
}
