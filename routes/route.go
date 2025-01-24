package routes

import (
	"type/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/register", controllers.Register)
			userRoutes.POST("/login", controllers.Login)
			userRoutes.GET("/send_code", controllers.SendVerificationCode) // 发送验证码
			userRoutes.POST("/verify_code", controllers.VerifyCode)        // 验证验证码
		}
		api.GET("/song", controllers.GetSong) // 需要在路由上加入歌曲ID参数
		api.POST("/song", controllers.UploadSong)
		api.GET("/asset", controllers.GetAsset)
		api.POST("/asset", controllers.UploadAsset)
		api.POST("/score", controllers.UploadTotalScore)
		api.GET("/scores", controllers.GetAllTotalScores)
		api.GET("/user_scores", controllers.GetUserALLScores)
	}

	// 注册 WebSocket 路由，缺少房间，玩家参数
	r.GET("/", func(c *gin.Context) {
		controllers.ServeHTML(c.Writer, c.Request)
	})
	// 进入不同房间
	r.GET("/ws", func(c *gin.Context) {
		controllers.HandleWebSocket(c.Writer, c.Request)
	})
}
