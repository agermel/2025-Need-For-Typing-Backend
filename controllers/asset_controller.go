package controllers

import (
	"net/http"

	"type/service"

	"github.com/gin-gonic/gin"
)

// AssetController 处理素材相关的 HTTP 请求
type AssetController struct {
	assetService service.AssetServiceInterface
}

// NewAssetController 创建 AssetController，并注入 AssetService
func NewAssetController(assetService service.AssetServiceInterface) *AssetController {
	return &AssetController{
		assetService: assetService,
	}
}

// GetAsset godoc
// @Summary 获取单个素材
// @Description 根据传入的 asset_id 查询素材的 file_id
// @Tags 素材
// @Accept json
// @Produce json
// @Param asset_id query string true "素材ID"
// @Success 200 {object} map[string]string "返回素材的 file_id"
// @Failure 400 {object} map[string]string "缺少 asset_id"
// @Failure 404 {object} map[string]string "未找到素材"
// @Router /api/assets [get]
func (sc *AssetController) GetAsset(c *gin.Context) {
	assetID := c.Query("asset_id")
	if assetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing asset_id"})
		return
	}

	asset, err := sc.assetService.GetAsset(assetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": asset.File_id})
}

// GetAllAssets godoc
// @Summary 获取所有素材
// @Description 查询所有素材的信息
// @Tags 素材
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "返回所有素材信息"
// @Failure 500 {object} map[string]string "查询失败"
// @Router /api/assets [get]
func (sc *AssetController) GetAllAssets(c *gin.Context) {
	assets, err := sc.assetService.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"assets": assets})
}

// 更新列表
func (ac *AssetController) UpdateList(c *gin.Context) {
	if err := ac.assetService.SaveList(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新素材列表失败"})
		return
	}
}
