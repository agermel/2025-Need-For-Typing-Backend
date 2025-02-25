package service

import (
	"type/dao"
	"type/models"
	"type/utils"
)

type AssetServiceInterface interface {
	GetAsset(assetID string) (*models.Asset, error)
	GetAllAssets() ([]models.Asset, error)
	SaveList() error
}

// AssetService 封装与素材有关的业务逻辑
type AssetService struct {
	assetDAO dao.AssetDAOInterface
}

// NewAssetService 创建 AssetService 实例
func NewAssetService(assetDAO dao.AssetDAOInterface) AssetServiceInterface {
	return &AssetService{
		assetDAO: assetDAO,
	}
}

// GetAsset 获取素材信息
func (s *AssetService) GetAsset(assetID string) (*models.Asset, error) {
	return s.assetDAO.GetAssetByID(assetID)
}

// GetAllAssets 获取所有素材信息
func (s *AssetService) GetAllAssets() ([]models.Asset, error) {
	return s.assetDAO.GetAllAssets()
}

// 保存素材文件队列
func (s *AssetService) SaveList() error {
	assets, err := utils.GetList[models.Asset]("asset")
	if err != nil {
		return err
	}

	err = s.assetDAO.SaveList(assets)
	if err != nil {
		return err
	}

	return nil
}
