package controllers

import (
	"net/http"
	"time"

	"type/database"
	"type/models"
	"type/service"
	"type/utils"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// @Summary 用户注册接口
// @Description 用户注册，创建新用户并发送验证邮件
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param   user  body      models.User  true  "用户信息"
// @Success 201   {object}  map[string]interface{}  "返回注册成功消息"
// @Failure 400   {object}  map[string]interface{}  "请求参数错误"
// @Failure 409   {object}  map[string]interface{}  "用户名或邮箱已存在"
// @Failure 500   {object}  map[string]interface{}  "服务器内部错误"
// @Router /api/user/register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// 检查邮箱是否已注册
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user.EmailVerified = false

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	code := utils.GenerateRandomVerifyCode()
	service.SaveCode(user.Email, code, 30*time.Minute)
	if err := utils.SendMail(user.Email, "Typing_hero verification code", code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

// Login 用户登录
// @Summary 用户登录接口
// @Description 用户登录，验证用户名和密码，生成 JWT Token
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param   user  body      models.User  true  "用户登录信息"
// @Success 200   {object}  map[string]interface{}  "返回登录成功消息及 JWT Token"
// @Failure 400   {object}  map[string]interface{}  "请求参数错误"
// @Failure 401   {object}  map[string]interface{}  "用户名或密码错误"
// @Failure 403   {object}  map[string]interface{}  "邮箱未验证"
// @Failure 404   {object}  map[string]interface{}  "用户不存在"
// @Failure 500   {object}  map[string]interface{}  "服务器内部错误"
// @Router /api/user/login [post]
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在，并且把该用户写入existingUser
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if user.Password != existingUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	if !user.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(int(existingUser.ID), existingUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}

// 用前端来删除jwt来进行登出
