package controller

import (
	"net/http"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")

	if _, err := utils.ParseToken(token); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
