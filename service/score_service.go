package service

import (
	"type/dao"
	"type/models"
)

type ScoreServiceInterface interface {
	UploadTotalScore(ts *models.TotalScore) error
	GetAllTotalScores(songID string) ([]map[string]interface{}, error)
	GetUserAllBestScores(userID string) (map[string]interface{}, error)
}

// ScoreService 封装与分数相关的业务逻辑
type ScoreService struct {
	scoreDAO dao.ScoreDAOInterface
}

func NewScoreService(dao dao.ScoreDAOInterface) ScoreServiceInterface {

	return &ScoreService{scoreDAO: dao}
}

// UploadTotalScore 处理上传总分的业务逻辑
func (s *ScoreService) UploadTotalScore(totalScore *models.TotalScore) error {
	return s.scoreDAO.CreateTotalScore(totalScore)
}

// GetAllTotalScores 根据歌曲 ID 获取所有总分信息，并整理成结果数据
func (s *ScoreService) GetAllTotalScores(songID string) ([]map[string]interface{}, error) {
	scores, err := s.scoreDAO.GetTotalScoresBySongID(songID)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, score := range scores {
		result = append(result, map[string]interface{}{
			"user_id":    score.UserID,
			"username":   score.User.Username,
			"score":      score.Score,
			"song_title": score.Game.Song.Title,
			"time":       score.Game.Time,
		})
	}
	return result, nil
}

// GetUserAllBestScores 根据用户 ID 获取用户所有游戏中的最佳成绩
func (s *ScoreService) GetUserAllBestScores(userID string) (map[string]interface{}, error) {
	user, err := s.scoreDAO.GetUserWithGames(userID)
	if err != nil {
		return nil, err
	}

	var bestScores []map[string]interface{}
	// 遍历用户的每场游戏，提取每场游戏中最高的分数
	for _, game := range user.Games {
		if len(game.Score) > 0 {
			bestScore := game.Score[0]
			for _, score := range game.Score {
				if score.Score > bestScore.Score {
					bestScore = score
				}
			}
			bestScores = append(bestScores, map[string]interface{}{
				"game_id":    game.ID,
				"song_title": game.Song.Title,
				"score":      bestScore.Score,
				"time":       game.Time,
			})
		}
	}

	result := map[string]interface{}{
		"user_id":     user.ID,
		"username":    user.Username,
		"best_scores": bestScores,
	}
	return result, nil
}
