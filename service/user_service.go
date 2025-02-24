package service

import (
	"errors"
	"time"

	"type/dao"
	"type/models"
	"type/utils"

	"github.com/google/uuid"
)

// UserServiceInterface 定义用户服务的行为
type UserServiceInterface interface {
	RegisterUser(user *models.User) error
	LoginUser(username, password string) (string, error)
	RequestPasswordReset(email string) error
	VerifyResetToken(email, token string) error
	ResetPassword(email, newPassword string) error
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

// RequestPasswordReset 请求密码重置
func (us *UserService) RequestPasswordReset(email string) error {
	// 查找用户是否存在
	user, err := us.userDAO.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 生成一个新的密码重置 Token
	token := uuid.New().String()

	// 设置Token的过期时间
	expiresAt := time.Now().Add(time.Hour)

	// 调用 DAO 层，保存 Token 和 过期时间
	err = us.userDAO.RequestPasswordReset(email, token, expiresAt)
	if err != nil {
		return errors.New("保存密码重置信息失败")
	}

	resetLink := "https://thusdaykfcv50.top/reset_password?token=" + token + "?email=" + user.Email
	err = utils.SendMail(user.Email, "邮箱验证", "点击修改密码"+resetLink)
	if err != nil {
		return errors.New("发送重置邮件失败")
	}

	return nil
}

func (us *UserService) VerifyResetToken(email, token string) error {
	// 查找用户
	user, err := us.userDAO.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证 token 是否匹配
	err = us.userDAO.VerifyResetToken(user.Email, token)
	if err != nil {
		return errors.New("无效的重置 token")
	}

	return nil
}

func (us *UserService) ResetPassword(email, newPassword string) error {
	// 查找用户
	user, err := us.userDAO.GetUserByEmail(email)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 更新密码
	if err := us.userDAO.ResetPassword(user.Email, newPassword); err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}
