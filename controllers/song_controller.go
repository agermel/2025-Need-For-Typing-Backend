package controllers

import (
	"net/http"
	"strconv"

	"type/api/response"
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
// @Success 200 {object} response.Response{data=response.SongResponse} "通信成功（通过code来判断具体情况）"
// @Router /api/song [get]
func (sc *SongController) GetSong(c *gin.Context) {
	songID := c.Query("song_id")
	if songID == "" {
		response.FailWithMessage("missing song_id", c)
		return
	}

	id, err := strconv.Atoi(songID)
	if err != nil {
		response.FailWithMessage("songID类型转换失败", c)
		return
	}

	song, err := sc.songService.GetSong(id)
	if err != nil {
		response.FailWithMessage("song not found", c)
		return
	}

	response.OKWithDetailed("Fetch song successfully",
		response.SongResponse{
			Song_id: song.File_id,
		}, c)
}

// UpdateList godoc
// @Summary 更新歌曲列表
// @Description 更新歌曲列表信息
// @Tags 歌曲
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "通信成功（通过code来判断具体情况）"
// @Router /api/update_list_song [get]
func (sc *SongController) UpdateList(c *gin.Context) {
	if err := sc.songService.SaveList(); err != nil {
		response.FailWithMessage("更新歌曲列表失败", c)
		return
	}
	response.OKWithMessage("update list successfully", c)
}

// GetAllSongs godoc
// @Summary 获取所有歌曲
// @Description 查询所有歌曲信息
// @Tags 歌曲
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=response.SongsResponse{Song_ids=[]models.Song}} "通信成功（通过code来判断具体情况）"
// @Router /api/all_songs [get]
func (sc *SongController) GetAllSongs(c *gin.Context) {
	songs, err := sc.songService.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	response.OKWithDetailed("查询成功",
		response.SongsResponse{
			Song_ids: songs,
		}, c)
}
