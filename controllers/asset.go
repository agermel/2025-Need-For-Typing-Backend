package controllers

import (
	"net/http"

	"type/database"
	"type/models"

	"github.com/gin-gonic/gin"
)

// @Summary      获取资产信息
// @Description  根据 asset_id 查询资产信息
// @Tags         Asset
// @Param        asset_id  query   string  true  "资产 ID"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string  "缺少 asset_id"
// @Failure      404       {object}  map[string]string  "资产未找到"
// @Router       /api/asset [get]
func GetAsset(c *gin.Context) {
	AssetID := c.Query("asset_id")
	if AssetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing song_id"})
		return
	}

	var asset models.Asset
	if err := database.DB.Where("id = ?", AssetID).First(&asset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": asset.File_id})
}

// @Summary      上传资产信息
// @Description  提交 JSON 格式的资产信息保存到数据库
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Param        asset  body  models.Asset  true  "资产信息"
// @Success      200    {object}  map[string]int  "返回资产 ID"
// @Failure      400    {object}  map[string]string  "请求体格式错误"
// @Failure      500    {object}  map[string]string  "服务器错误"
// @Router       /api/asset [post]
func UploadAsset(c *gin.Context) {
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&asset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": asset.ID})
}
