package service

import (
	"time"
	"type/api/response"
	"type/dao"
	"type/models"
)

type ScoreServiceInterface interface {
	UploadTotalScore(ts *models.TotalScore) error
	GetAllTotalScores(songID int) ([]response.ScoreReponse, error)
	GetUserAllBestScores(userID int) ([]response.ScoreReponse, error)
}

// ScoreService 封装与分数相关的业务逻辑
type ScoreService struct {
	scoreDAO dao.ScoreDAOInterface
	userDAO  dao.UserDAOInterface
	songDAO  dao.SongDAOInterface
}

func NewScoreService(scoreDAO dao.ScoreDAOInterface,
	userDAO dao.UserDAOInterface,
	songDAO dao.SongDAOInterface) ScoreServiceInterface {
	return &ScoreService{scoreDAO: scoreDAO,
		userDAO: userDAO,
		songDAO: songDAO}
}

// UploadTotalScore 处理上传总分的业务逻辑
func (s *ScoreService) UploadTotalScore(totalScore *models.TotalScore) error {
	return s.scoreDAO.CreateTotalScore(totalScore)
}

// GetAllTotalScores 根据歌曲 ID 获取所有总分信息，并整理成结果数据
func (s *ScoreService) GetAllTotalScores(songid int) ([]response.ScoreReponse, error) {
	scores, err := s.scoreDAO.GetTotalScoresBySongID(songid)
	if err != nil {
		return []response.ScoreReponse{}, err
	}

	song, err := s.songDAO.GetSongByID(songid)
	if err != nil {
		return []response.ScoreReponse{}, err
	}

	var result []response.ScoreReponse
	for _, score := range *scores {

		user, err := s.userDAO.GetUserByID(score.UserID)
		if err != nil {
			return []response.ScoreReponse{}, err
		}

		result = append(result, response.ScoreReponse{
			UserID:   score.UserID,
			Username: user.Username,
			Score:    score.Score,
			Title:    song.Title,
			Time:     score.CreatedAT,
		})
	}
	return result, nil
}

// GetUserAllBestScores 根据用户 ID 获取用户所有已游玩过的歌曲中的最佳成绩
func (s *ScoreService) GetUserAllBestScores(userID int) ([]response.ScoreReponse, error) {
	var bestScoresMap = make(map[int]int)
	var bestScoresMaptime = make(map[int]time.Time)
	var bestScores []response.ScoreReponse

	// 获取该用户的所有成绩
	scores, err := s.scoreDAO.GetScoresWithUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, score := range *scores {
		// 如果当前 songID 还没有记录，或者当前分数比已存的高，则更新
		if maxScore, exists := bestScoresMap[score.SongID]; !exists || score.Score > maxScore {
			bestScoresMap[score.SongID] = score.Score
			bestScoresMaptime[score.SongID] = score.CreatedAT
		}
	}

	user, err := s.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 转换为结构体
	for songID, maxScore := range bestScoresMap {
		song, err := s.songDAO.GetSongByID(songID)
		if err != nil {
			return nil, err
		}

		bestScores = append(bestScores, response.ScoreReponse{
			UserID:   userID,
			Username: user.Username,
			Score:    maxScore,
			Title:    song.Title,
			Time:     bestScoresMaptime[songID],
		})
	}

	return bestScores, nil
}
