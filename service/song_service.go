package service

import (
	"io"
	"strconv"

	"type/dao"
	"type/models"
)

type SongServiceInterface interface {
	GetSong(songID string) (*models.Song, error)
	UploadSong(songID, qiniuPath, title string, reader io.Reader) (string, error)
	GetAllSongs() ([]models.Song, error)
}

// SongService 封装与歌曲相关的业务逻辑
type SongService struct {
	songDAO dao.SongDAOInterface
}

func NewSongService(dao dao.SongDAOInterface) SongServiceInterface {
	return &SongService{songDAO: dao}
}

// GetSong 获取歌曲信息
func (s *SongService) GetSong(songID string) (*models.Song, error) {
	return s.songDAO.GetSongByID(songID)
}

// UploadSong 处理歌曲上传：上传文件并保存记录到数据库
func (s *SongService) UploadSong(songID, qiniuPath, title string, reader io.Reader) (string, error) {
	// 调用上传函数，将本地文件上传到七牛云，返回文件 URL
	fileURL, err := UploadToQiniu(qiniuPath, reader)
	if err != nil {
		return "", err
	}

	// 将 songID 转换为整数
	intID, err := strconv.Atoi(songID)
	if err != nil {
		return "", err
	}

	song := models.Song{
		ID:      intID,
		Title:   title,
		File_id: fileURL,
	}

	if err := s.songDAO.CreateSong(&song); err != nil {
		return "", err
	}

	return fileURL, nil
}

// GetAllSongs 获取所有歌曲信息
func (s *SongService) GetAllSongs() ([]models.Song, error) {
	return s.songDAO.GetAllSongs()
}
