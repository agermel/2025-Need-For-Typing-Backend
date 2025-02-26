package controllers

import (
	"log"
	"net/http"

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
// @Param user body models.User true "用户注册信息"
// @Success 201 {object} map[string]string "用户创建成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userService.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

// Login godoc
// @Summary 用户登录接口
// @Description 用户登录，验证用户名和密码，生成 JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.User true "用户登录信息"
// @Success 200 {object} map[string]interface{} "返回登录成功消息及 JWT Token"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "用户名或密码错误"
// @Failure 403 {object} map[string]interface{} "邮箱未验证"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := uc.userService.LoginUser(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}

// 用前端来删除jwt来进行登出

// ForgetPassword godoc
// @Summary 忘记密码
// @Description 向后端发起忘记密码请求，通过邮箱发送重置密码的链接
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param email query string true "用户邮箱地址"
// @Success 200 {object} map[string]string "密码重置链接已发送"
// @Failure 400 {object} map[string]string "缺少邮箱参数"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /user/forget_password [get]
func (uc *UserController) ForgetPassword(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要邮箱提供"})
		return
	}

	if uc.userService == nil {
		log.Println("userService is nil")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	err := uc.userService.RequestPasswordReset(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
}

// ResetPassword godoc
// @Summary 重置密码
// @Description 提供重置密码的 token、邮箱和新密码，完成密码重置
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param token query string true "重置密码的 token"
// @Param email query string true "用户邮箱"
// @Param new_password query string true "新的密码"
// @Success 200 {object} map[string]string "密码重置成功"
// @Failure 400 {object} map[string]string "缺少参数或无效的 token"
// @Failure 500 {object} map[string]string "密码重置失败"
// @Router /user/reset_password [get]
func (uc *UserController) ResetPassword(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")
	newPassword := c.Query("new_password")

	if token == "" || email == "" || newPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需要token、email和newPassword"})
		return
	}

	// 验证重置 token 是否有用
	if err := uc.userService.VerifyResetToken(email, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 token 或过期"})
		return
	}

	err := uc.userService.ResetPassword(email, newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码重置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

// 验证玩家token有效性
func (uc *UserController) VerifyToken(c *gin.Context) {
	var verifyToken VerifyToken
	if err := c.ShouldBindJSON(&verifyToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.VerifyToken(verifyToken.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
