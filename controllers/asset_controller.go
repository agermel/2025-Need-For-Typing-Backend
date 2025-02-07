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

// UploadAsset godoc
// @Summary 上传素材
// @Description 上传素材文件到七牛云，需提供 asset_id 与文件内容
// @Tags 素材
// @Accept multipart/form-data
// @Produce json
// @Param asset_id formData string true "素材ID"
// @Param file formData file true "上传的文件"
// @Success 200 {object} map[string]string "返回素材ID和文件URL"
// @Failure 400 {object} map[string]string "缺少 asset_id 或文件上传失败"
// @Failure 500 {object} map[string]string "上传到七牛云失败"
// @Router /api/assets [post]
func (sc *AssetController) UploadAsset(c *gin.Context) {
	// 获取前端提交的 asset_id
	assetID := c.PostForm("asset_id")
	if assetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 asset_id"})
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	defer file.Close()

	qiniuPath := "asset/" + assetID

	fileURL, err := sc.assetService.UploadAsset(assetID, qiniuPath, header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传到七牛云失败", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": assetID, "url": fileURL})
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
