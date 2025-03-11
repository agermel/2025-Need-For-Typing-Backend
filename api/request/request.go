package request

import (
	"type/models"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ScoreRequest struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	Score  int  `json:"total_score"`
	GameID int  `json:"game_id"`
}

type CreateGameRequest struct {
	UserID int                 `json:"user_id"` // 外键，指向 User 表的 ID
	SongID int                 `json:"song_id"`
	Score  []models.TotalScore `json:"score"` // 关联 TotalScore
}
