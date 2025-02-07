package dao

import (
	"type/database"
	"type/models"
)

// UserDAOInterface 定义 UserDAO 的行为
type UserDAOInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
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
