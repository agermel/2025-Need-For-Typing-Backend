package controllers

import (
	"net/http"

	"type/database"
	"type/models"

	"github.com/gin-gonic/gin"
)

func GetSong(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing song_id"})
		return
	}

	var song models.Song
	if err := database.DB.Where("id = ?", songID).First(&song).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": song.File_id})
}

// 使用URL参数song_id获取歌曲

func UploadSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&song).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": song.ID})
}
