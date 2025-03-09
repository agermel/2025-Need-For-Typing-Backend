package controllers

import (
	"net/http"
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

// 待修改

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
