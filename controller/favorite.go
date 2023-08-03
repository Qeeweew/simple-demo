package controller

import (
	"net/http"
	"simple-demo/common/log"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	type Req struct {
		VideoId    int64 `form:"video_id"`
		ActionType int32 `form:"action_type"`
	}
	var UserId = c.Keys["auth_id"].(uint)
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	if req.ActionType != 1 && req.ActionType != 2 {
		log.Logger.Error("action type error")
		result.Error(c, result.QueryParamErrorStatus)
	}
	if err := service.NewFavorite().FavoriteAction(req.ActionType == 1, UserId, uint(req.VideoId)); err != nil {
		result.Error(c, result.FavoriteErrorStatus)
	}
	result.SuccessBasic(c)
	return
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
