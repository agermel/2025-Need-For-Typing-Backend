package controllers

import (
	"type/api/request"
	"type/api/response"
	"type/service"

	"github.com/gin-gonic/gin"
)

type GameController struct {
	gameService service.GameServiceInterface
}

func NewGameController(service service.GameServiceInterface) *GameController {
	return &GameController{
		gameService: service,
	}
}

func (gc *GameController) CreateGame(c *gin.Context) {
	var gameRequest request.CreateGameRequest
	if err := c.ShouldBindJSON(&gameRequest); err != nil {
		response.FailWithMessage("fail to create game "+err.Error(), c)
		return
	}
	response.OKWithMessage("Game Created", c)
}
