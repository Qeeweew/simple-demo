package middleware

import (
	"net/http"
	"simple-demo/controller"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
)

var skipPaths = []string{
	"/douyin/user/register/",
	"/douyin/feed/",
	"/douyin/user/login/",
}

// auth_id 记录token对应的用户也就是发起请求的用户的id
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要鉴权的
		for _, path := range skipPaths {
			if path == c.FullPath() {
				c.Next()
				return
			}
		}

		tokenString := c.Query("token")
		if c.FullPath() == "/douyin/publish/action/" {
			tokenString = c.PostForm("token")
		}
		if tokenString == "" {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "missing token",
			})
			return
		}
		id, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  "invalid token",
			})
			return
		}
		c.Set("auth_id", id)
	}
}
