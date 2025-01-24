package controllers

import (
	"net/http"

	"type/database"
	"type/models"

	"github.com/gin-gonic/gin"
)

// @Summary      上传总分
// @Description  提交用户总分信息保存到数据库
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        totalScore  body  models.TotalScore  true  "总分信息"
// @Success      200         {object}  map[string]int  "返回分数 ID"
// @Failure      400         {object}  map[string]string  "请求体格式错误"
// @Failure      500         {object}  map[string]string  "服务器错误"
// @Router       /api/score [post]
func UploadTotalScore(c *gin.Context) {
	var totalScore models.TotalScore
	if err := c.ShouldBindJSON(&totalScore); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&totalScore).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score_id": totalScore.ID})
}

// @Summary      获取一首歌的所有总分信息
// @Description  根据歌曲 ID 获取该歌曲的所有分数信息
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        song_id  query  string  true  "歌曲 ID"
// @Success      200      {array}  map[string]interface{}  "返回分数列表"
// @Failure      400      {object}  map[string]string  "请求参数错误"
// @Failure      500      {object}  map[string]string  "服务器错误"
// @Router       /api/scores [get]
func GetAllTotalScores(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
		return
	}

	var scores []models.TotalScore
	if err := database.DB.
		Joins("JOIN game ON game.id = total_scores.game_id").
		Where("games.song_id = ?", songID).
		Order("total_scores.score DESC").
		Preload("User").      // 预加载User信息
		Preload("Game.song"). // 预加载关联的Song信息
		Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, score := range scores {
		result = append(result, gin.H{
			"user_id":    score.UserID,
			"username":   score.User.Username,
			"score":      score.Score,
			"song_title": score.Game.Song.Title,
			"time":       score.Game.Time,
		})
	}
	c.JSON(http.StatusOK, result)
}

// @Summary      获取用户所有最佳成绩
// @Description  根据用户 ID 获取其所有游戏中的最佳分数
// @Tags         Score
// @Accept       json
// @Produce      json
// @Param        user_id  query  string  true  "用户 ID"
// @Success      200      {object}  map[string]interface{}  "返回用户最佳成绩"
// @Failure      400      {object}  map[string]string  "缺少用户 ID"
// @Failure      404      {object}  map[string]string  "未找到用户"
// @Router       /api/user_scores [get]
func GetUserALLScores(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	var user models.User
	if err := database.DB.
		Preload("Games.Score").
		Preload("Games.Song").
		Where("id = ?", userID).
		Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 收集每场游戏的最好成绩
	var bestScores []gin.H
	for _, game := range user.Games {
		if len(game.Score) > 0 {
			// 找到该场游戏中的最高分
			bestScore := game.Score[0]
			for _, score := range game.Score {
				if score.Score > bestScore.Score {
					bestScore = score
				}
			}
			// 添加到返回结果中
			bestScores = append(bestScores, gin.H{
				"game_id":    game.ID,
				"song_title": game.Song.Title,
				"score":      bestScore.Score,
				"time":       game.Time,
			})
		}
	}

	// 返回用户最好成绩
	c.JSON(http.StatusOK, gin.H{
		"user_id":     user.ID,
		"username":    user.Username,
		"best_scores": bestScores,
	})
}
