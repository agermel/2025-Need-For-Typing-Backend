package service

import (
	"io"
	"strconv"

	"type/dao"
	"type/models"
)

type AssetServiceInterface interface {
	GetAsset(assetID string) (*models.Asset, error)
	UploadAsset(assetID, qiniuPath, name string, reader io.Reader) (string, error)
	GetAllAssets() ([]models.Asset, error)
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

// UploadAsset 处理素材上传：上传文件并保存记录到数据库
func (s *AssetService) UploadAsset(assetID, qiniuPath, name string, reader io.Reader) (string, error) {
	// 调用上传函数，将本地文件上传到七牛云，返回文件 URL
	fileURL, err := UploadToQiniu(qiniuPath, reader)
	if err != nil {
		return "", err
	}

	// 将 assetID 转换为整数
	intID, err := strconv.Atoi(assetID)
	if err != nil {
		return "", err
	}

	asset := models.Asset{
		ID:      intID,
		Name:    name,
		File_id: fileURL,
	}

	if err := s.assetDAO.CreateAsset(&asset); err != nil {
		return "", err
	}

	return fileURL, nil
}

// GetAllAssets 获取所有素材信息
func (s *AssetService) GetAllAssets() ([]models.Asset, error) {
	return s.assetDAO.GetAllAssets()
}
