package main

import (
	"simple-demo/common/config"
	"simple-demo/controller"
	"simple-demo/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// 用于存放视频
	r.Static("/videos", config.AppCfg.VideoPath)

	apiRouter := r.Group("/douyin")

	// 测试
	apiRouter.Any("/ping", controller.Ping)

	// 视频流
	apiRouter.GET("/feed/", controller.Feed)

	// 用户相关
	userGroup := apiRouter.Group("/user")
	{
		// 获取用户登录信息
		userGroup.GET("/", middleware.JWTAuthMiddleware(), controller.UserInfo)

		// 新用户注册
		userGroup.POST("/register/", controller.Register)

		// 用户登录
		userGroup.POST("/login/", controller.Login)
	}

	// 视频发布相关
	publishGroup := apiRouter.Group("/publish", middleware.JWTAuthMiddleware())
	{
		// 用户上传视频
		publishGroup.POST("/action/", controller.Publish)

		// 直接列出用户投稿过的所有视频
		publishGroup.GET("/list/", controller.PublishList)
	}

	// 点赞相关
	favoriteGroup := apiRouter.Group("/favorite", middleware.JWTAuthMiddleware())
	{
		// 点赞操作
		favoriteGroup.POST("/action/", controller.FavoriteAction)

		// 获取点赞列表
		favoriteGroup.GET("/list/", controller.FavoriteList)
	}

	// 评论相关
	commentGroup := apiRouter.Group("/comment", middleware.JWTAuthMiddleware())
	{
		// 评论操作
		commentGroup.POST("/action/", controller.CommentAction)

		// 获取评论列表
		commentGroup.GET("/list/", controller.CommentList)
	}

	// 消息相关
	messageGroup := apiRouter.Group("/message", middleware.JWTAuthMiddleware())
	{
		// 聊天记录
		messageGroup.GET("/chat/", controller.MessageChat)

		// 消息操作
		messageGroup.POST("/action/", controller.MessageAction)
	}

	// 用户关系相关
	relationGroup := apiRouter.Group("/relation", middleware.JWTAuthMiddleware())
	{
		// 对指定用户 关注 取关
		relationGroup.POST("/action/", controller.RelationAction)

		// 获取用户关注列表
		relationGroup.GET("/follow/list/", controller.FollowList)

		// 获取用户粉丝列表
		relationGroup.GET("/follower/list/", controller.FollowerList)

		// TODO: 获取用户好友列表
		relationGroup.GET("/friend/list/", controller.FriendList)
	}
}
