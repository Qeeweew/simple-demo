package controller

import (
	"net/http"
	"simple-demo/common/log"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
)

type RelationActionReq struct {
	Token      string `form:"token"`
	ToUserId   uint   `form:"to_user_id"`
	ActionType int    `form:"action_type"`
}

type RelationListReq struct {
	UserId uint   `form:"user_id"`
	Token  string `form:"token"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	var req RelationActionReq

	// 参数校验
	err := c.ShouldBind(&req)
	if err != nil {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	if req.ActionType > 2 || req.ActionType < 1 {
		log.Logger.Error("operation illegal")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	var currentId = c.Keys["auth_id"].(uint)
	err = service.NewRelation().FollowAction(currentId, req.ToUserId, req.ActionType)
	if err != nil {
		log.Logger.Error(err.Error())
		if req.ActionType == 1 {
			result.Error(c, result.FollowErrorStatus)
			return
		} else {
			result.Error(c, result.UnFollowErrorStatus)
			return
		}
	}

	result.Success(c, result.R{})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	var req RelationListReq

	// 参数校验
	err := c.ShouldBind(&req)
	if err != nil {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	var currentId = c.Keys["auth_id"].(uint)

	// 获取关注列表
	followList, err := service.NewRelation().FollowList(currentId, req.UserId)
	if err != nil {
		log.Logger.Error(err.Error())
		result.Error(c, result.Status{
			StatusCode: result.FollowListErrorStatus.StatusCode,
			StatusMsg:  result.FollowListErrorStatus.StatusMsg,
		})
		return
	}

	result.Success(c, result.R{
		"user_list": followList,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	var req RelationListReq

	// 参数校验
	err := c.ShouldBind(&req)
	if err != nil {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	var currentId = c.Keys["auth_id"].(uint)

	// 获取粉丝列表
	followerList, err := service.NewRelation().FanList(currentId, req.UserId)
	if err != nil {
		log.Logger.Error(err.Error())
		result.Error(c, result.Status{
			StatusCode: result.FollowerListErrorStatus.StatusCode,
			StatusMsg:  result.FollowerListErrorStatus.StatusMsg,
		})
		return
	}

	result.Success(c, result.R{
		"user_list": followerList,
	})
}

func FriendList(c *gin.Context) {
	var req RelationListReq
	// 参数校验
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	userList, err := service.NewRelation().FriendList(req.UserId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: err.Error()})
		return
	}
	result.Success(c, result.R{
		"user_list": userList,
	})
}
