package models

import "time"

type User struct {
	ID             uint      `json:"user_id" gorm:"primaryKey"`           // id
	Username       string    `json:"username" gorm:"not null"`            // 用户名
	Password       string    `json:"password" gorm:"size:255;not null"`   // 密码
	Email          string    `json:"email" gorm:"unique;not null"`        // 邮箱
	EmailVerified  bool      `json:"email_verified" gorm:"default:false"` // 邮箱是否验证
	Games          []Game    `gorm:"foreignKey:UserID"`                   //
	ResetToken     string    `json:"reset_token" gorm:"size:255"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
}
