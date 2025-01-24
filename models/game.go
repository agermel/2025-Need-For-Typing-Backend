package models

import (
	"time"
)

// 一场游戏的信息
type Game struct {
	ID      int          `json:"id" gorm:"primaryKey;autoIncrement"` // 主键，自动递增
	UserID  int          `json:"user_id"`                            // 外键，指向 User 表的 ID
	Song    Song         `gorm:"foreignKey:ID"`                      // 外键，指向 Song 表
	Score   []TotalScore `gorm:"foreignKey:GameID"`                  // 外键，指向 TotalScore 表
	Time    time.Time    `json:"time" `
	Players []User       `gorm:"many2many:game_users;foreignKey:ID;joinForeignKey:GameID;References:ID;JoinReferences:UserID"`
}
