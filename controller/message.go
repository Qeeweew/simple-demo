package controller

import (
	"net/http"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

type MessageActionRequest struct {
	ToUserId   uint   `form:"to_user_id"`
	ActionType int32  `form:"action_type"`
	Content    string `form:"content"`
}

type MessageChatRequst struct {
	ToUserId   uint  `form:"to_user_id"`
	PreMsgTime int64 `form:"pre_msg_time"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	var req MessageActionRequest
	if err := c.ShouldBind(&req); err != nil {
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	userId := c.Keys["auth_id"].(uint)
	err := service.NewMessage().SendMessage(userId, req.ToUserId, req.Content)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	result.SuccessBasic(c)
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	var req MessageChatRequst
	if err := c.ShouldBind(&req); err != nil {
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	userId := c.Keys["auth_id"].(uint)
	messageList, err := service.NewMessage().ChatHistory(req.PreMsgTime, userId, req.ToUserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messageList})
}
