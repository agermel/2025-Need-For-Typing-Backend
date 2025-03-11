package service

import (
	"type/dao"
	"type/models"
	"type/utils"
)

type SongServiceInterface interface {
	GetSong(songID int) (*models.Song, error)
	GetAllSongs() ([]models.Song, error)
	SaveList() error
}

// SongService 封装与歌曲相关的业务逻辑
type SongService struct {
	songDAO dao.SongDAOInterface
}

func NewSongService(dao dao.SongDAOInterface) SongServiceInterface {
	return &SongService{songDAO: dao}
}

// GetSong 获取歌曲信息
func (s *SongService) GetSong(songID int) (*models.Song, error) {
	return s.songDAO.GetSongByID(songID)
}

// 保存音乐文件队列
func (s *SongService) SaveList() error {
	songs, err := utils.GetList[models.Song]("music")
	if err != nil {
		return err
	}

	err = s.songDAO.SaveList(songs)
	if err != nil {
		return err
	}

	return nil
}

// GetAllSongs 获取所有歌曲信息
func (s *SongService) GetAllSongs() ([]models.Song, error) {
	return s.songDAO.GetAllSongs()
}
