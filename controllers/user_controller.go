package controllers

import (
	"time"

	"type/api/response"
	"type/models"
	"type/service"

	"github.com/gin-gonic/gin"
)

type VerifyToken struct {
	Token string `json:"token"`
}

type UserController struct {
	userService service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface) *UserController {
	return &UserController{
		userService: service,
	}
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口，接收用户信息创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body request.RegisterRequest true "用户注册信息"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := uc.userService.RegisterUser(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithMessage("user created successfully", c)
}

// Login godoc
// @Summary 用户登录接口
// @Description 用户登录，验证用户名和密码，生成 JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body request.LoginRequest true "用户登录信息"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	token, err := uc.userService.LoginUser(user.Username, user.Password)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithDetailed("login successfully",
		response.TokenResponse{
			User:      user.Username,
			Token:     token,
			ExpiresAT: time.Now().Add(12 * time.Hour),
		}, c)
}

// ForgetPassword godoc
// @Summary 请求密码重置
// @Description 用户通过邮箱请求密码重置，会收到包含重置链接的邮件
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param email query string true "用户注册邮箱"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /user/forget_password [get]
func (uc *UserController) ForgetPassword(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.FailWithMessage("Email need", c)
		return
	}

	if uc.userService == nil {
		response.FailWithMessage("Internal Server Error", c)
		return
	}
	err := uc.userService.RequestPasswordReset(email)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithMessage("Password reset link sent", c)
}

// ResetPassword godoc
// @Summary 重置密码
// @Description 用户通过邮件中提供的 token 进行密码重置
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param token query string true "密码重置令牌"
// @Param email query string true "用户注册邮箱"
// @Param new_password query string true "新密码"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /user/reset_password [post]
func (uc *UserController) ResetPassword(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")
	newPassword := c.Query("new_password")

	if token == "" || email == "" || newPassword == "" {
		response.FailWithMessage("需要token、email和newPassword", c)
		return
	}

	// 验证重置 token 是否有用
	if err := uc.userService.VerifyResetToken(email, token); err != nil {
		response.FailWithMessage("无效的 token 或过期", c)
		return
	}

	err := uc.userService.ResetPassword(email, newPassword)
	if err != nil {
		response.FailWithMessage("密码重置失败", c)
		return
	}

	response.OKWithMessage("密码重置成功", c)
}
