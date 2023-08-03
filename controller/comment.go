package controller

import (
	"fmt"
	"net/http"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

type CommentActionRequest struct {
	Token       string `form:"token" binding:"required"`
	VideoId     uint   `form:"video_id" binding:"required"`
	ActionType  int    `form:"action_type" binding:"required"`
	CommentText string `form:"comment_text"`
	CommentId   uint   `form:"comment_id"`
}

type CommentListRequest struct {
	Token   string `form:"token" binding:"required"`
	VideoId uint   `form:"video_id" binding:"required"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	var req CommentActionRequest
	if err := c.ShouldBind(&req); err != nil {
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	log.Logger.Sugar().Infof("%v", req)
	// TODO: move some code to Service layer
	switch req.ActionType {
	case 1:
		var comment = model.Comment{UserId: c.Keys["auth_id"].(uint), VideoId: req.VideoId}
		err := service.NewComment().CommentAction(true, &comment)
		if err != nil {
			result.Error(c, result.ServerErrorStatus)
			log.Logger.Sugar().Error(err)
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Create Comment successfully"},
			Comment:  comment,
		})
		log.Logger.Info("Comment create Success", zap.String("Req", fmt.Sprint(req)), zap.String("Comment", fmt.Sprint(comment)))
	case 2:
		err := service.NewComment().CommentAction(false, &model.Comment{Id: req.CommentId})
		if err != nil {
			result.Error(c, result.CommentNotExitErrorStatus)
			log.Logger.Sugar().Error(err)
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Delete Comment successfully"},
		})
		log.Logger.Info("Comment create Success", zap.String("Comment Id", fmt.Sprint(req.CommentId)))
	default:
		result.Error(c, result.QueryParamErrorStatus)
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	var user_id = c.Keys["auth_id"].(uint)
	var req CommentListRequest
	if err := c.ShouldBind(&req); err != nil {
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	comments, err := service.NewComment().CommentList(user_id, req.VideoId)
	if err != nil {
		result.Error(c, result.ServerErrorStatus)
		log.Logger.Sugar().Errorf("%v", err)
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
