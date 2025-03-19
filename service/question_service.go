package service

import (
	"fmt"
	"type/dao"
	"type/models"
)

type QuestionService struct {
	QuestionDao dao.QuestionDAOInterface
}

type QuestionServiceInterface interface {
	SaveQuestion(question *models.Question) error
	GetContent(title string) (string, error)
	GetTitleList() ([]string, error)
}

func NewQuestionService(questionDAO dao.QuestionDAOInterface) QuestionServiceInterface {
	return &QuestionService{
		QuestionDao: questionDAO,
	}
}

func (s *QuestionService) SaveQuestion(question *models.Question) error {
	return s.QuestionDao.SaveQuestionWithTitle(question)
}

func (s *QuestionService) GetContent(title string) (string, error) {
	question, err := s.QuestionDao.GetContentWithTitle(title)
	if err != nil {
		return "", err
	}

	fmt.Println(question.Context)

	return question.Context, nil
}

func (s *QuestionService) GetTitleList() ([]string, error) {
	list, err := s.QuestionDao.GetTitleList()
	if err != nil {
		return nil, err
	}

	return list, nil
}
