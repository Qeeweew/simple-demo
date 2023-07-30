package controller

import (
	"net/http"
	"simple-demo/common/model"
	"simple-demo/service"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
// 点赞还没搞
func Feed(c *gin.Context) {
	videos, err := service.GetVideo().GetFeedList(30)
	userID := uint(0) //c.Keys["auth_id"].(uint)
	if err != nil {
		RegisterError(err, c)
		return
	}
	var ctlVideos []Video
	for i := range videos {
		var author model.User
		err = service.GetUser().Info(userID, videos[i].UserId, &author)
		if err != nil {
			RegisterError(err, c)
			return
		}
		ctlVideos = append(ctlVideos, FromVideoModel(&videos[i], &author))
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: ctlVideos,
		NextTime:  time.Now().Unix(),
	})
}
