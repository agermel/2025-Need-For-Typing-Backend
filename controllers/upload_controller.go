package controllers

import (
	"net/http"
	"type/utils"

	"github.com/gin-gonic/gin"
)

func (uc *UserController) GetToken(c *gin.Context) {
	var verifyToken VerifyToken
	if err := c.ShouldBindJSON(&verifyToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := uc.userService.VerifyToken(verifyToken.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadToken, err := utils.GetToken()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uploadToken": uploadToken})
}
