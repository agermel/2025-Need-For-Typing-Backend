package controllers

import (
	"net/http"

	"type/service"

	"github.com/gin-gonic/gin"
)

// SongController 处理歌曲相关的 HTTP 请求
type SongController struct {
	songService service.SongServiceInterface
}

// NewSongController 创建 SongController 实例，并注入 SongService
func NewSongController(service service.SongServiceInterface) *SongController {
	return &SongController{
		songService: service,
	}
}

// GetSong godoc
// @Summary 获取歌曲信息
// @Description 根据传入的 song_id 查询歌曲信息并返回文件ID
// @Tags 歌曲
// @Accept json
// @Produce json
// @Param song_id query string true "歌曲ID"
// @Success 200 {object} map[string]string "返回歌曲的 file_id"
// @Failure 400 {object} map[string]string "缺少 song_id 参数"
// @Failure 404 {object} map[string]string "歌曲未找到"
// @Router /song [get]
func (sc *SongController) GetSong(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing song_id"})
		return
	}

	song, err := sc.songService.GetSong(songID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": song.File_id})
}

// UpdateList godoc
// @Summary 更新歌曲列表
// @Description 更新歌曲列表信息
// @Tags 歌曲
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "更新成功"
// @Failure 500 {object} map[string]string "更新歌曲列表失败"
// @Router /songs/update [post]
func (sc *SongController) UpdateList(c *gin.Context) {
	if err := sc.songService.SaveList(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新歌曲列表失败"})
		return

	}
}

// GetAllSongs godoc
// @Summary 获取所有歌曲
// @Description 查询所有歌曲信息
// @Tags 歌曲
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "返回所有歌曲的列表"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /songs [get]
func (sc *SongController) GetAllSongs(c *gin.Context) {
	songs, err := sc.songService.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"songs": songs})
}
