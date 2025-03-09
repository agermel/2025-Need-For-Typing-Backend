package controllers

import (
	"type/api/response"
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
// @Success 200 {object} response.Response{data:response.AssetResponse} "通信成功（通过code来判断具体情况）"
// @Router /api/assets [get]
func (sc *AssetController) GetAsset(c *gin.Context) {
	assetID := c.Query("asset_id")
	if assetID == "" {
		response.FailWithMessage("missing asset_id", c)
		return
	}

	asset, err := sc.assetService.GetAsset(assetID)
	if err != nil {
		response.FailWithMessage("asset not found", c)
		return
	}

	response.OKWithDetailed("Fetch asset successfully",
		response.AssetResponse{
			Asset_id: asset.File_id,
		}, c)
}

// GetAllAssets godoc
// @Summary 获取所有素材
// @Description 查询所有素材的信息
// @Tags 素材
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data:response.AssetsResponse{Asset_ids:[]models.Asset}} "通信成功（通过code来判断具体情况）"
// @Router /api/assets [get]
func (sc *AssetController) GetAllAssets(c *gin.Context) {
	assets, err := sc.assetService.GetAllAssets()
	if err != nil {
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OKWithDetailed("查询成功",
		response.AssetsResponse{
			Asset_ids: assets,
		}, c)
}

// UpdateList godoc
// @Summary 更新素材列表
// @Description 更新素材列表信息
// @Tags 素材
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /api/assets/update [post]
func (ac *AssetController) UpdateList(c *gin.Context) {
	if err := ac.assetService.SaveList(); err != nil {
		response.FailWithMessage("更新素材列表失败", c)
		return
	}
	response.OKWithMessage("update list successfully", c)
}
