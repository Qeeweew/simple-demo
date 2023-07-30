package controller

import (
	"github.com/gin-gonic/gin"
	"simple-demo/common/db"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/common/result"
	"simple-demo/repository"
	"simple-demo/service"
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

func getRelationService() model.RelationService {
	return service.NewRelationService(repository.NewRelationRepository(db.MySQL))
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

	if req.Token == "" || req.ToUserId <= 0 || req.ActionType > 2 || req.ActionType < 1 {
		log.Logger.Error("operation illegal")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	err = getRelationService().FollowAction(req.Token, req.ToUserId, req.ActionType)
	if err != nil {
		log.Logger.Error(err.Error())
		if req.ActionType == 1 {
			result.Error(c, result.Status{
				StatusCode: result.FollowErrorStatus.StatusCode,
				StatusMsg:  result.FollowErrorStatus.StatusMsg,
			})
			return
		} else {
			result.Error(c, result.Status{
				StatusCode: result.UnFollowErrorStatus.StatusCode,
				StatusMsg:  result.UnFollowErrorStatus.StatusMsg,
			})
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

	if req.UserId <= 0 || req.Token == "" {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	// 获取关注列表
	followList, err := getRelationService().FollowList(req.Token, req.UserId)
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

	if req.UserId <= 0 || req.Token == "" {
		log.Logger.Error("check params error")
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	// 获取粉丝列表
	followerList, err := getRelationService().FanList(req.Token, req.UserId)
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

// TODO: FriendList all users have same friend list
//func FriendList(c *gin.Context) {
//	c.JSON(http.StatusOK, UserListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		UserList: []User{DemoUser},
//	})
//}
