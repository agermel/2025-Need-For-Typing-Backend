package dao

import (
	"type/database"
	"type/models"
)

type ScoreDAOInterface interface {
	CreateTotalScore(ts *models.TotalScore) error
	GetTotalScoresBySongID(songID int) (*[]models.TotalScore, error)
	GetScoresWithUserID(userID int) (*[]models.TotalScore, error)
}

type ScoreDAO struct{}

func NewScoreDAO() ScoreDAOInterface {
	return &ScoreDAO{}
}

func (dao *ScoreDAO) CreateTotalScore(ts *models.TotalScore) error {
	return database.DB.Create(ts).Error
}

func (dao *ScoreDAO) GetTotalScoresBySongID(songID int) (*[]models.TotalScore, error) {
	var scores []models.TotalScore

	err := database.DB.
		Where("song_id = ?", songID).
		Order("score DESC").
		Find(&scores).
		Error

	if err != nil {
		return nil, err
	}

	return &scores, nil
}

func (dao *ScoreDAO) GetScoresWithUserID(userID int) (*[]models.TotalScore, error) {
	var scores []models.TotalScore

	err := database.DB.
		Where("user_id = ?", userID).
		Find(scores).
		Error

	if err != nil {
		return nil, err
	}

	return &scores, nil
}
