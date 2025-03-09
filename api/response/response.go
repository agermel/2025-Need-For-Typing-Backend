package response

import (
	"net/http"
	"time"
	"type/models"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"string"`
	Data interface{} `json:"data"`
}

type TokenResponse struct {
	User      string    `json:"user"`
	Token     string    `json:"token"`
	ExpiresAT time.Time `json:"expiresAT"`
}

type VerificationResponse struct {
	Username string `json:"username"`
}

type AssetResponse struct {
	Asset_id string `json:"asset_id"`
}

type AssetsResponse struct {
	Asset_ids []models.Asset `json:"asset_ids"`
}

type SongResponse struct {
	Song_id string `json:"song_id"`
}

type SongsResponse struct {
	Song_ids []models.Song `json:"song_ids"`
}

type UploadToken struct {
	UploadToken string `json:"upload_token"`
}

const (
	ERROR   = 0
	SUCCESS = 1
)

func Result(code int, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		msg,
		data,
	})
}

func OKWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, message, map[string]interface{}{}, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, message, map[string]interface{}{}, c)
}

func OKWithDetailed(message string, data interface{}, c *gin.Context) {
	Result(SUCCESS, message, data, c)
}
