package middlewares

import (
	"type/api/response"
	"type/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" {
			response.FailWithMessage("请求头中auth为空", c)
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(authHeader)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("userID", claims.UserID)
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
		c.Next()
	}
}
