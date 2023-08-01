package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"simple-demo/common/config"
	"simple-demo/common/model"
	"simple-demo/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish save upload file to public directory
func Publish(c *gin.Context) {
	userID := c.Keys["auth_id"].(uint)
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d-%s", userID, filename)
	saveFile := filepath.Join(config.AppCfg.VideoPath, finalName)
	logrus.Println("save file: ", saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
	var video = model.Video{
		UserId:  userID,
		Title:   title,
		PlayUrl: fmt.Sprintf("http://%s/videos/%s", c.Request.Host, finalName),
	}
	logrus.Println("video url: ", video.PlayUrl)
	service.NewVideo().Publish(&video)
}

// PublishList all users have same publish video list
// isFavorate 还没有处理
func PublishList(c *gin.Context) {
	targetID, _ := strconv.Atoi(c.Query("user_id"))
	val, found := c.Keys["auth_id"]
	var userID uint
	if found {
		userID = val.(uint)
	} else {
		userID = 0
	}
	videos, err := service.NewVideo().GetPublishList(userID)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var targetUser model.User
	err = service.NewUser().Info(userID, uint(targetID), &targetUser)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var ctlVideos []Video
	for _, video := range videos {
		ctlVideos = append(ctlVideos, FromVideoModel(&video))
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: ctlVideos,
	})
}
