package controllers

import (
	"type/api/response"
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
// @Success 200 {object} response.Response{data=response.UploadToken} "通信成功（通过code来判断具体情况）"
// @Router /api/get_upload_token [get]
func (uc *UserController) GetToken(c *gin.Context) {
	uploadToken, err := utils.GetToken()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OKWithDetailed("fetch upload token successfully",
		response.UploadToken{UploadToken: uploadToken}, c)
}
