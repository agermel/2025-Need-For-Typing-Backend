package models

type TotalScore struct {
	ID     uint `json:"id" gorm:"primaryKey;autoIncrement"` // 主键
	UserID uint `json:"user_id"`
	Score  int  `json:"total_score"`             // 分值
	GameID int  `json:"game_id" gorm:"not null"` // 外键，指向 Game 表的 ID
	User   User
	Game   Game
}
