package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"simple-demo/common/config"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/common/result"
	"simple-demo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish save upload file to public directory
func Publish(c *gin.Context) {
	userId := c.Keys["auth_id"].(uint)
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
	finalName := fmt.Sprintf("%d-%s", userId, filename)
	saveFile := filepath.Join(config.AppCfg.VideoPath, finalName)
	// 暂时放这里了
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		result.Error(c, result.ServerErrorStatus)
		log.Logger.Error("Saving video Failed", zap.String("err", err.Error()))
		return
	}
	log.Logger.Info("Saving video Succeed", zap.String("File", finalName))
	var video = model.Video{
		AuthorId: userId,
		Title:    title,
		PlayUrl:  fmt.Sprintf("http://%s/videos/%s", c.Request.Host, finalName),
		// TODO: 生成视频封面
	}
	if err := service.NewVideo().Publish(&video); err != nil {
		result.Error(c, result.ServerErrorStatus)
		log.Logger.Error("Saving video Failed", zap.String("err", err.Error()))
	}
	log.Logger.Info("Publish Succeed", zap.String("url", video.PlayUrl))
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	type Req struct {
		UserId uint `form:"user_id"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	targetId := req.UserId
	val, found := c.Keys["auth_id"]
	var userId uint
	if found {
		userId = val.(uint)
	} else {
		userId = 0
	}
	videos, err := service.NewVideo().GetPublishList(userId, targetId)
	if err != nil {
		log.Logger.Error("PublishList error", zap.String("err", err.Error()))
		result.Error(c, result.ServerErrorStatus)
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
	})
}
