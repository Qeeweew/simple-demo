package controller

import (
	"net/http"
	"simple-demo/common/model"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

func FromVideoModel(video *model.Video) Video {
	return Video{
		Id:            int64(video.ID),
		Author:        FromUserModel(&video.Author),
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: int64(video.FavoriteCount),
		CommentCount:  int64(video.CommentCount),
		IsFavorite:    video.IsFavorite,
	}
}

// need to initialize isFollow
func FromUserModel(user *model.User) User {
	return User{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   int64(user.FollowCount),
		FollowerCount: int64(user.FanCount),
		IsFollow:      user.IsFollow,
	}
}

func RegisterError(err error, c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{StatusCode: 1, StatusMsg: err.Error()}})
}
