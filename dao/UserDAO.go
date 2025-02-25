package dao

import (
	"errors"
	"time"

	"type/database"
	"type/models"
)

// UserDAOInterface 定义 UserDAO 的行为
type UserDAOInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	RequestPasswordReset(email string, token string, expiresAt time.Time) error
	VerifyResetToken(email, token string) error
	ResetPassword(email, newPassword string) error
	VerifyEmail(email string) error
}

// 就是给 UserDAO 的方法的载体
type UserDAO struct{}

// NewUserDAO 创建 DAO 实例
func NewUserDAO() UserDAOInterface {
	return &UserDAO{}
}

// CreateUser 创建新用户
func (dao *UserDAO) CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// GetUserByUsername 根据用户名查找用户
func (dao *UserDAO) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱查找用户
func (dao *UserDAO) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RequestPasswordReset 生成密码重置 Token, 并存入数据库
func (dao *UserDAO) RequestPasswordReset(email string, token string, expiresAt time.Time) error {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.ResetToken = token
	user.TokenExpiresAt = expiresAt
	return database.DB.Save(user).Error
}

// VerifyResetToken 验证 Token 是否正确
func (dao *UserDAO) VerifyResetToken(email, token string) error {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证 Token 是否匹配
	if user.ResetToken != token {
		return errors.New("无效 Token")
	}

	if time.Now().After(user.TokenExpiresAt) {
		return errors.New("token 已过期")
	}

	return nil
}

// ResetPassword 修改用户密码
func (dao *UserDAO) ResetPassword(email, newPassword string) error {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.Password = newPassword
	user.ResetToken = ""
	user.TokenExpiresAt = time.Time{}

	return database.DB.Save(user).Error
}

// 修改邮箱是否验证为真
func (dao *UserDAO) VerifyEmail(email string) error {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	user.EmailVerified = true

	return database.DB.Save(user).Error
}
