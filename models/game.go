package models

import (
	"time"
)

// 一场游戏的信息
type Game struct {
	ID      int          `json:"id" gorm:"primaryKey;autoIncrement"` // 主键，自动递增
	UserID  int          `json:"user_id"`                            // 外键，指向 User 表的 ID
	SongID  int          `json:"song_id" gorm:"not null"`            // 关联 Song 的 ID
	Song    Song         `json:"song" gorm:"foreignKey:SongID"`      // 关联 Song
	Score   []TotalScore `json:"score" gorm:"foreignKey:GameID"`     // 关联 TotalScore
	Time    time.Time    `json:"time"`
	Players []User       `json:"players" gorm:"many2many:game_users;foreignKey:ID;joinForeignKey:GameID;References:ID;JoinReferences:UserID"`
}
