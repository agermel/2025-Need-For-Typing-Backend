package dao

import (
	"type/database"
	"type/models"
)

// AssetDAOInterface 定义 AssetDAO 的行为
type AssetDAOInterface interface {
	GetAssetByID(assetID string) (*models.Asset, error)
	CreateAsset(asset *models.Asset) error
	GetAllAssets() ([]models.Asset, error)
}

// AssetDAO 封装与素材有关的数据库操作
type AssetDAO struct{}

// NewAssertDAO 创建 AssetDAO 实例
func NewAssetDAO() AssetDAOInterface {
	return &AssetDAO{}
}

// GetAssetByID 根据素材 ID 查询素材信息
func (dao *AssetDAO) GetAssetByID(assetID string) (*models.Asset, error) {
	var asset models.Asset
	if err := database.DB.Where("id = ?", assetID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

// CreateAsset 将歌曲信息保存到数据库
func (dao *AssetDAO) CreateAsset(asset *models.Asset) error {
	return database.DB.Create(asset).Error
}

// GetAllAssets 查询所有素材
func (dao *AssetDAO) GetAllAssets() ([]models.Asset, error) {
	var assets []models.Asset
	if err := database.DB.Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}
