package service

import (
	"errors"
	"time"

	"type/dao"
	"type/models"
	"type/utils"
)

// UserServiceInterface 定义用户服务的行为
type UserServiceInterface interface {
	RegisterUser(user *models.User) error
	LoginUser(username, password string) (string, error)
}

// UserService 实现 UserServiceInterface
type UserService struct {
	userDAO dao.UserDAOInterface
}

// 业务逻辑层
// NewUserService 创建 UserService
func NewUserService(userDAO dao.UserDAOInterface) UserServiceInterface {
	return &UserService{
		userDAO: userDAO,
	}
}

// RegisterUser 处理用户注册逻辑
func (service *UserService) RegisterUser(user *models.User) error {
	// 检查用户名是否已存在
	existingUser, _ := service.userDAO.GetUserByUsername(user.Username)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// 检查邮箱是否已注册
	existingUser, _ = service.userDAO.GetUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	user.EmailVerified = false

	// 创建用户
	if err := service.userDAO.CreateUser(user); err != nil {
		return err
	}

	// 生成验证码并发送邮件
	code := utils.GenerateRandomVerifyCode()
	SaveCode(user.Email, code, 30*time.Minute)

	if err := utils.SendMail(user.Email, "Typing_hero verification code", code); err != nil {
		return errors.New("failed to send verification email")
	}

	return nil
}

// LoginUser 处理用户登录逻辑
func (s *UserService) LoginUser(username, password string) (string, error) {
	user, err := s.userDAO.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user.Password != password {
		return "", errors.New("invalid password")
	}

	if !user.EmailVerified {
		return "", errors.New("email not verified")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(int(user.ID), user.Username)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}
