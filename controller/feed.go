package controller

import (
	"net/http"
	"simple-demo/common/result"
	"simple-demo/service"
	"simple-demo/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type FeedRequest struct {
	Token      string `form:"token"`
	LatestTime int64  `form:"latest_time"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var req FeedRequest
	c.ShouldBind(&req)
	var userId uint = 0
	if req.Token != "" {
		id, err := utils.ParseToken(req.Token)
		if err != nil {
			result.Error(c, result.TokenErrorStatus)
			return
		}
		userId = id
	}
	videos, err := service.NewVideo().GetFeedList(userId, 30)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
