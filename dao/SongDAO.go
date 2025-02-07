package dao

import (
	"type/database"
	"type/models"
)

// SongDAOInterface 定义数据访问接口
type SongDAOInterface interface {
	GetSongByID(songID string) (*models.Song, error)
	CreateSong(song *models.Song) error
	GetAllSongs() ([]models.Song, error)
}

// SongDAO 封装与歌曲相关的数据库操作
type SongDAO struct{}

// NewSongDAO 创建 SongDAO 实例
func NewSongDAO() SongDAOInterface {
	return &SongDAO{}
}

// GetSongByID 根据歌曲 ID 查询歌曲信息
func (dao *SongDAO) GetSongByID(songID string) (*models.Song, error) {
	var song models.Song
	if err := database.DB.Where("id = ?", songID).First(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// CreateSong 将歌曲信息保存到数据库
func (dao *SongDAO) CreateSong(song *models.Song) error {
	return database.DB.Create(song).Error
}

// GetAllSongs 查询所有歌曲
func (dao *SongDAO) GetAllSongs() ([]models.Song, error) {
	var songs []models.Song
	if err := database.DB.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}
