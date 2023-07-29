package middleware

import (
	"net/http"
	"simple-demo/controller"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
)

var skipPaths = []string{
	"/douyin/user/register",
	"/douyin/feed/",
	"/douyin/user/login/",
}

// auth_id 记录token对应用户的id

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要鉴权的
		for _, path := range skipPaths {
			if path == c.FullPath() {
				c.Next()
				return
			}
		}

		var defaultResponse = controller.Response{
			StatusCode: 1,
			StatusMsg:  "Invalid Token",
		}
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusOK, defaultResponse)
			return
		}
		id, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, defaultResponse)
			return
		}
		c.Set("auth_id", id)
	}
}
