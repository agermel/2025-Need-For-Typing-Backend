package dao

import (
	"type/database"
	"type/models"
)

type ScoreDAOInterface interface {
	CreateTotalScore(ts *models.TotalScore) error
	GetTotalScoresBySongID(songID string) ([]models.TotalScore, error)
	GetUserWithGames(userID string) (*models.User, error)
}

type ScoreDAO struct{}

func NewScoreDAO() ScoreDAOInterface {
	return &ScoreDAO{}
}

func (dao *ScoreDAO) CreateTotalScore(ts *models.TotalScore) error {
	return database.DB.Create(ts).Error
}

func (dao *ScoreDAO) GetTotalScoresBySongID(songID string) ([]models.TotalScore, error) {
	var scores []models.TotalScore
	err := database.DB.
		Joins("JOIN games ON games.id = total_scores.game_id").
		Where("games.song_id = ?", songID).
		Order("total_scores.score DESC").
		Preload("User").
		Preload("Game.Song").
		Find(&scores).Error
	if err != nil {
		return nil, err
	}
	return scores, nil
}

func (dao *ScoreDAO) GetUserWithGames(userID string) (*models.User, error) {
	var user models.User
	err := database.DB.
		Preload("Games.Score").
		Preload("Games.Song").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
