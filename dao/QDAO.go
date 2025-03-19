package dao

import (
	"type/database"
	"type/models"
)

type QuestionDAOInterface interface {
	SaveQuestionWithTitle(*models.Question) error
	GetContentWithTitle(Title string) (*models.Question, error)
	GetTitleList() ([]string, error)
}

// 没有字段，但是有接口定义的方法
type QDAO struct{}

// 返回符合接口的实例
func NewQDAO() QuestionDAOInterface {
	return &QDAO{}
}

func (q *QDAO) SaveQuestionWithTitle(qs *models.Question) error {
	return database.DB.Create(qs).Error
}

func (q *QDAO) GetContentWithTitle(Title string) (*models.Question, error) {
	var question models.Question

	err := database.DB.Where("title = ?", Title).Find(&question).Error
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (q *QDAO) GetTitleList() ([]string, error) {
	var list []string

	err := database.DB.Model(&models.Question{}).Pluck("title", &list).Error
	if err != nil {
		return nil, err
	}

	return list, nil

}
