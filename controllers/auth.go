package controllers

import (
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
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// 生成验证码
	code := utils.GenerateRandomVerifyCode()

	// 保存验证码
	service.SaveCode(email, code, 30*time.Minute)

	if err := utils.SendMail(email, "This is your verifing code not junk !!!", code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
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
func VerifyCode(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if service.VerifyCode(request.Email, request.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "Verification successful"})
		service.DeleteCode(request.Email)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired code"})
	}
}
