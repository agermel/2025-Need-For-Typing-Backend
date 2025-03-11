package dao

import (
	"type/database"
	"type/models"
)

type GameDAOInterface interface {
	Creategame(game *models.Game) error
}

type GameDAO struct{}

func NewGameDAO() GameDAOInterface {
	return &GameDAO{}
}

func (dao *GameDAO) Creategame(game *models.Game) error {
	err := database.DB.Create(&game).Error
	if err != nil {
		return err
	}
	return nil
}
