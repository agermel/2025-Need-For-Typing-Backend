package controllers

import (
	"type/api/response"
	"type/models"
	"type/service"

	"github.com/gin-gonic/gin"
)

// QuestionController handles question-related operations
type QuestionController struct {
	QuestionService service.QuestionServiceInterface
}

// NewQuestionController creates a new instance of QuestionController
func NewQuestionController(service service.QuestionServiceInterface) *QuestionController {
	return &QuestionController{
		QuestionService: service,
	}
}

// SaveQuestion godoc
// @Summary 保存问题
// @Description 保存一个包含标题和内容的问题
// @Tags 问题
// @Accept json
// @Produce json
// @Param question body models.Question true "问题数据"
// @Success 200 {object} response.Response
// @Router /api/save_question [post]
func (qc *QuestionController) SaveQuestion(c *gin.Context) {
	var question models.Question

	err := c.ShouldBindJSON(&question)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 手动检查字段是否为空
	if question.Title == "" || question.Context == "" {
		response.FailWithMessage("Title and context cannot be empty", c)
		return
	}

	err = qc.QuestionService.SaveQuestion(&question)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithMessage("成功保存", c)
}

// GetContent godoc
// @Summary 根据标题获取内容
// @Description 根据问题的标题获取其内容
// @Tags 问题
// @Accept json
// @Produce json
// @Param title query string true "问题标题"
// @Success 200 {object} response.Response
// @Router /api/get_context [get]
func (qc *QuestionController) GetContent(c *gin.Context) {
	title := c.Query("title")

	content, err := qc.QuestionService.GetContent(title)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithDetailed("fetched content success", content, c)
}

// GetTitleList godoc
// @Summary 获取问题标题列表
// @Description 获取所有问题的标题列表
// @Tags 问题
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/get_list [get]
func (qc *QuestionController) GetTitleList(c *gin.Context) {
	list, err := qc.QuestionService.GetTitleList()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OKWithDetailed("fetched list success", list, c)
}
