package controllers

import (
	"net/http"

	"type/models"
	"type/service"

	"github.com/gin-gonic/gin"
)

// ScoreController 负责处理分数相关的 HTTP 请求
type ScoreController struct {
	scoreService service.ScoreServiceInterface
}

// NewScoreController 创建 ScoreController 实例，并注入 ScoreService
func NewScoreController(service service.ScoreServiceInterface) *ScoreController {
	return &ScoreController{
		scoreService: service,
	}
}

// UploadTotalScore godoc
// @Summary 上传总分
// @Description 接收 JSON 格式的总分数据并上传到服务器
// @Tags 分数
// @Accept json
// @Produce json
// @Param totalScore body models.TotalScore true "上传的总分信息"
// @Success 200 {object} map[string]interface{} "返回上传成功的分数ID"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /score [post]
func (sc *ScoreController) UploadTotalScore(c *gin.Context) {
	var totalScore models.TotalScore
	if err := c.ShouldBindJSON(&totalScore); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sc.scoreService.UploadTotalScore(&totalScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"score_id": totalScore.ID})
}

// GetAllTotalScores godoc
// @Summary 获取歌曲所有总分信息
// @Description 根据传入的 song_id 查询该歌曲所有的总分记录
// @Tags 分数
// @Accept json
// @Produce json
// @Param song_id query string true "歌曲ID"
// @Success 200 {object} map[string]interface{} "返回总分信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /scores [get]
func (sc *ScoreController) GetAllTotalScores(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "song_id is required"})
		return
	}

	result, err := sc.scoreService.GetAllTotalScores(songID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserAllScores godoc
// @Summary 获取用户所有最佳成绩
// @Description 根据传入的 user_id 查询该用户所有的最佳分数记录
// @Tags 分数
// @Accept json
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} map[string]interface{} "返回用户最佳分数记录"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "用户未找到"
// @Router /user_scores [get]
func (sc *ScoreController) GetUserAllScores(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	result, err := sc.scoreService.GetUserAllBestScores(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
