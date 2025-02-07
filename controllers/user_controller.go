package controllers

import (
	"net/http"

	"type/models"
	"type/service"

	"github.com/gin-gonic/gin"
)

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
