package controller

import (
	"simple-demo/common/log"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FavoriteAction no practical effect, just check if token is valid
type FavoriteActionRequest struct {
	VideoId    uint  `form:"video_id" binding:"required"`
	ActionType int32 `form:"action_type" binding:"required"`
}

type FavoriteListRequest struct {
	UserId uint `form:"user_id" binding:"required"`
}

func FavoriteAction(c *gin.Context) {
	var UserId = c.Keys["auth_id"].(uint)
	var req FavoriteActionRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	if req.ActionType != 1 && req.ActionType != 2 {
		log.Logger.Error("action type error")
		result.Error(c, result.QueryParamErrorStatus)
	}
	if err := service.NewFavorite().FavoriteAction(req.ActionType == 1, UserId, req.VideoId); err != nil {
		result.Error(c, result.FavoriteErrorStatus)
	}
	result.SuccessBasic(c)
	return
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	var UserId = c.Keys["auth_id"].(uint)
	var req FavoriteListRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	videos, err := service.NewFavorite().FavoriteList(UserId, req.UserId)
	if err != nil {
		result.Error(c, result.FavoriteErrorStatus)
	}
	result.Success(c, result.R{
		"video_list": videos,
	})
}
