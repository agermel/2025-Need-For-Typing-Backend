package routes

import (
	"type/controllers"
	"type/middlewares"
	"type/utils"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	scoreController *controllers.ScoreController,
	userController *controllers.UserController,
	songController *controllers.SongController,
	assetController *controllers.AssetController,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.GET("/send_code", controllers.SendVerificationCode)    // 发送验证码
		userRoutes.POST("/verify_code", controllers.VerifyCode)           // 验证验证码
		userRoutes.GET("/forget_password", userController.ForgetPassword) // 发送忘记密码请求
		userRoutes.GET("/reset_password", userController.ResetPassword)
	}

	api := r.Group("/api")
	{
		api.GET("/song", songController.GetSong) // 需要在路由上加入歌曲ID参数
		api.POST("/song", songController.UploadSong)
		api.GET("/all_songs", songController.GetAllSongs)
		api.GET("/asset", assetController.GetAsset)
		api.POST("/asset", assetController.UploadAsset)
		api.GET("/all_assets", assetController.GetAllAssets)

		api.POST("/score", scoreController.UploadTotalScore)
		api.GET("/scores", scoreController.GetAllTotalScores)
		api.GET("/user_scores", scoreController.GetUserAllScores)
		api.GET("/essay", controllers.GetGeneratedEssay)
		api.POST("/uploadSong", songController.UploadSong)
	}

	// 加载测试用HTML
	r.GET("/WebSocket", func(c *gin.Context) {
		utils.ServeHTML(c.Writer, c.Request, "template\\client.html")
	})
	r.GET("/essay", func(c *gin.Context) {
		utils.ServeHTML(c.Writer, c.Request, "template\\llm.html")
	})
	r.GET("/upload_song", func(c *gin.Context) {
		utils.ServeHTML(c.Writer, c.Request, "template\\upload.html")
	})

	// 进入不同房间
	r.GET("/ws", func(c *gin.Context) {
		controllers.HandleWebSocket(c.Writer, c.Request)
	})

	return r
}
