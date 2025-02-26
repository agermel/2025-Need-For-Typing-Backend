package controllers

import (
	"net/http"
	"type/utils"

	"github.com/gin-gonic/gin"
)

// GetToken godoc
// @Summary 获取上传令牌
// @Description 验证用户令牌并返回上传令牌
// @Tags 上传
// @Accept json
// @Produce json
// @Param request body VerifyToken true "用户验证令牌"
// @Success 200 {object} map[string]string "返回上传令牌"
// @Failure 400 {object} map[string]string "请求错误或验证失败"
// @Router /upload/token [post]
func (uc *UserController) GetToken(c *gin.Context) {
	var verifyToken VerifyToken
	if err := c.ShouldBindJSON(&verifyToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := uc.userService.VerifyToken(verifyToken.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadToken, err := utils.GetToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uploadToken": uploadToken})
}
