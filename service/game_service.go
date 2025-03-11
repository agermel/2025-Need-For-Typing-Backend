package service

import (
	"time"
	"type/api/request"
	"type/dao"
	"type/models"
)

type GameServiceInterface interface {
}

type GameService struct {
	gameDAO dao.GameDAOInterface
}

func NewGameService(gameDAO dao.GameDAOInterface) GameServiceInterface {
	return &GameService{
		gameDAO: gameDAO,
	}
}

func (s *GameService) NewGame(gameRequest request.CreateGameRequest) error {
	game := models.Game{
		UserID: gameRequest.UserID,
		SongID: gameRequest.SongID,
		Time:   time.Now(),
		Score:  gameRequest.Score,
	}

	err := s.gameDAO.Creategame(&game)
	if err != nil {
		return err
	}

	return nil
}
