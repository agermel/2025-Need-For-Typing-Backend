package controllers

import (
	"strconv"
	"time"
	"type/api/response"

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
// @Param totalScore body request.ScoreRequest true "上传的总分信息"
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /api/score [post]
func (sc *ScoreController) UploadTotalScore(c *gin.Context) {
	var totalScore models.TotalScore
	userID, exists := c.Get("userID")
	if !exists {
		response.FailWithMessage("Token Error", c)
	}

	if err := c.ShouldBindJSON(&totalScore); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	id, ok := userID.(int)
	if !ok {
		response.FailWithMessage("userID 断言失败", c)
	}
	totalScore.UserID = id

	totalScore.CreatedAT = time.Now()

	if err := sc.scoreService.UploadTotalScore(&totalScore); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithMessage("upload total score successfully", c)
}

// GetAllTotalScores godoc
// @Summary 获取歌曲所有总分信息
// @Description 根据传入的 song_id 查询该歌曲所有的总分记录
// @Tags 分数
// @Accept json
// @Produce json
// @Param song_id query string true "歌曲ID"
// @Success 200 {object} response.Response{data=[]response.ScoreReponse} "通信成功（通过code来判断具体情况）"
// @Router /api/scores [get]
func (sc *ScoreController) GetAllTotalScores(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		response.FailWithMessage("song_id is required", c)
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		response.FailWithMessage("songID类型转换失败", c)
		return
	}

	result, err := sc.scoreService.GetAllTotalScores(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithDetailed("fetch all total scores successfully", result, c)
}

// GetUserAllScores godoc
// @Summary 获取用户所有最佳成绩
// @Description 根据传入的 user_id 查询该用户所有的最佳分数记录
// @Tags 分数
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.ScoreReponse} "通信成功（通过code来判断具体情况）"
// @Router /api/user_scores [get]
func (sc *ScoreController) GetUserAllScores(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		response.FailWithMessage("missing user_id", c)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		response.FailWithMessage("userID类型转换失败", c)
		return
	}

	result, err := sc.scoreService.GetUserAllBestScores(id)
	if err != nil {
		response.FailWithMessage("user not found", c)
		return
	}

	response.OKWithDetailed("fetch user all scores per song successfully", result, c)
}
