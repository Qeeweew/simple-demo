package middleware

import (
	"net/http"
	"simple-demo/controller"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var skipPaths = []string{
	"/douyin/user/register",
	"/douyin/feed/",
	"/douyin/user/login/",
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要鉴权的
		for _, path := range skipPaths {
			if path == c.FullPath() {
				c.Next()
				return
			}
		}

		var invalidTokenResponse = controller.Response{
			StatusCode: 1,
			StatusMsg:  "Invalid Token",
		}
		tokenString := c.Request.Header.Get("token")
		if tokenString == "" {
			c.JSON(http.StatusOK, invalidTokenResponse)
			return
		}
		var claims jwt.MapClaims
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return utils.SecretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusOK, invalidTokenResponse)
			return
		}
		aud, _ := claims.GetAudience()
		c.Set("user_id", aud[0])
	}
}
