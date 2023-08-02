package controller

import (
	"net/http"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/common/result"
	"simple-demo/service"
	"simple-demo/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	type Req = struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	var user = model.User{Name: req.Username, Password: req.Password}
	err := service.NewUser().Register(&user)
	if err != nil {
		log.Logger.Error("Register error",
			zap.String("username", user.Name), zap.String("password", user.Password),
			zap.String("err", err.Error()))
		result.Error(c, result.ServerErrorStatus)
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func Login(c *gin.Context) {
	type Req = struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	var user = model.User{Name: req.Username, Password: req.Password}
	err := service.NewUser().Login(&user)
	if err != nil {
		log.Logger.Error("Login error", zap.String("err", err.Error()))
		result.Error(c, result.LoginErrorStatus)
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func UserInfo(c *gin.Context) {
	type Req = struct {
		UserId int `form:"user_id"`
	}
	var req Req
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	targetId := req.UserId
	userId, _ := c.Keys["auth_id"].(uint)
	targetUser, err := service.NewUser().UserInfo(userId, uint(targetId))
	if err != nil {
		result.Error(c, result.ServerErrorStatus)
		log.Logger.Error("UserInfo error", zap.String("err", err.Error()))
	} else {
		log.Logger.Info("???")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     targetUser,
		})
	}
}
