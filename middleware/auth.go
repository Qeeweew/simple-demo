package middleware

import (
	"path"
	"simple-demo/common/result"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
)

var skipPaths = []string{
	"/douyin/user/register/",
	"/douyin/feed/",
	"/douyin/user/login/",
}

// JWTAuthMiddleware auth_id 记录token对应的用户也就是发起请求的用户的id
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要鉴权的
		for _, path := range skipPaths {
			if path == c.FullPath() {
				c.Next()
				return
			}
		}
		if matches, _ := path.Match("/videos/*", c.FullPath()); matches {
			return
		}
		type TokenReq struct {
			Token string `form:"token"`
		}
		var req TokenReq

		if err := c.ShouldBind(&req); err != nil {
			result.Error(c, result.TokenErrorStatus)
		}
		id, err := utils.ParseToken(req.Token)
		if err != nil {
			result.Error(c, result.TokenErrorStatus)
			return
		}
		c.Set("auth_id", id)
	}
}
