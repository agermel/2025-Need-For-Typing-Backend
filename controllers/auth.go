package controllers

import (
	"context"
	"net/http"
	"time"

	"type/service"
	"type/utils"

	"github.com/gin-gonic/gin"
)

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// @Summary      发送邮箱验证码
// @Description  生成并向用户邮箱发送验证码
// @Tags         Verification
// @Param        email  query   string  true  "邮箱地址"
// @Success      200    {object}  map[string]string  "验证码发送成功"
// @Failure      400    {object}  map[string]string  "邮箱为空"
// @Failure      500    {object}  map[string]string  "发送失败"
// @Router       /user/send_code [get]
func SendVerificationCode(c *gin.Context) {
	// 创建一个带定时关闭的子上下文
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()

	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// 生成验证码
	code := utils.GenerateRandomVerifyCode()

	// 保存验证码
	err := service.SaveCode(ctx, email, code, 30*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to save verification code" + err.Error()})
	}

	if err := utils.SendMail(email, "This is your verifing code not junk !!!", code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

// @Summary      验证验证码
// @Description  验证用户提交的验证码是否有效
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Param        request  body  VerifyRequest  true  "邮箱和验证码"
// @Success      200      {object}  map[string]string  "验证成功"
// @Failure      400      {object}  map[string]string  "请求体无效"
// @Failure      401      {object}  map[string]string  "验证码无效或过期"
// @Router       /user/verify_code [post]
func (uc *UserController) VerifyCode(c *gin.Context) {
	// 创建一个带定时关闭的子上下文
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()

	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := uc.userService.VerifyCode(ctx, request.Email, request.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Verification successful", "user": user.Username})
}
