package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"simple-demo/common/config"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/common/oss"
	"simple-demo/common/result"
	"simple-demo/service"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish save upload file to Aliyun oss
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

	// 获取文件
	file, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	// 判断是否为视频
	checkFile, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, checkFile); err != nil {
		log.Logger.Error("copy file error")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if filetype.IsVideo(buf.Bytes()) == false {
		log.Logger.Error("file is not video")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	checkFile.Close()

	// 存储到oss
	log.Logger.Info("start to upload video to oss, file type: ")
	ok, err := oss.UploadVideoToOss(config.AliyunCfg.BucketName, data.Filename, file)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 获取url 存储到数据库
	videoUrl, imgUrl, err := oss.GetOssVideoUrlAndImgUrl(config.AliyunCfg.BucketName, data.Filename)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	var video = model.Video{
		AuthorId: userId,
		Title:    title,
		PlayUrl:  videoUrl,
		CoverUrl: imgUrl,
	}
	if err := service.NewVideo().Publish(&video); err != nil {
		result.Error(c, result.ServerErrorStatus)
		log.Logger.Error("Saving video Failed", zap.String("err", err.Error()))
	}
	log.Logger.Info("Publish Succeed", zap.String("url", video.PlayUrl))
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  videoUrl + " uploaded successfully",
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
