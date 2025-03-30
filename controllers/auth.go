package controllers

import (
	"context"
	"time"

	"type/api/response"
	"type/service"
	"type/utils"

	"github.com/gin-gonic/gin"
)

// @Summary      发送邮箱验证码
// @Description  生成并向用户邮箱发送验证码
// @Tags         Verification
// @Param        email  query   string  true  "邮箱地址"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router       /user/send_code [get]
func SendVerificationCode(c *gin.Context) {
	// 创建一个带定时关闭的子上下文
	ctx, cancel := context.WithTimeout(c, 1*time.Minute)
	defer cancel()

	email := c.Query("email")
	if email == "" {
		response.FailWithMessage("Email is required", c)
		return
	}

	// 生成验证码
	code := utils.GenerateRandomVerifyCode()

	// 保存验证码
	err := service.SaveCode(ctx, email, code, 30*time.Minute)
	if err != nil {
		response.FailWithMessage("Failed to save verification code"+err.Error(), c)
		return
	}

	if err := utils.SendMail(email, "Need for Typing account verification", code, "controllers/封面.png"); err != nil {
		response.FailWithMessage("Failed to send email"+err.Error(), c)
		return
	}

	response.OKWithMessage("Verification code sent", c)
}

// @Summary      验证验证码
// @Description  验证用户提交的验证码是否有效
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Param        request  body  request.VerifyRequest  true  "邮箱和验证码"
// @Success 200 {object} response.Response{data=response.VerificationResponse} "通信成功（通过code来判断具体情况）"
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
		response.FailWithMessage("Invalid request", c)
		return
	}

	user, err := uc.userService.VerifyCode(ctx, request.Email, request.Code)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithDetailed("Verification successful", response.VerificationResponse{
		Username: user.Username,
	}, c)

}
