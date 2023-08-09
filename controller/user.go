package controller

import (
	"errors"
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
type LoginRequest = struct {
	Username string `form:"username" binding:"required,min=1,max=32"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

type UserInfoRequest = struct {
	UserId int `form:"user_id" binding:"required"`
}

func Register(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}

	var user = model.User{Name: req.Username, Password: req.Password}
	err := service.NewUser().Register(&user)
	if err != nil {
		if errors.Is(err, service.ErrUserExist) {
			// 用户名已存在
			result.Error(c, result.UsernameExitErrorStatus)
			return
		} else {
			// 服务器错误
			log.Logger.Error("Register error",
				zap.String("username", user.Name), zap.String("password", user.Password),
				zap.String("err", err.Error()))
			result.Error(c, result.ServerErrorStatus)
		}
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error("check params error", zap.String("err", err.Error()))
		result.Error(c, result.QueryParamErrorStatus)
		return
	}
	var user = model.User{Name: req.Username, Password: req.Password}
	err := service.NewUser().Login(&user)
	if err != nil {
		if errors.Is(err, service.ErrPassword) {
			log.Logger.Info("password error", zap.Any("user", user))
			result.Error(c, result.PasswordErrorStatus)
			return
		} else {
			log.Logger.Error("Login error", zap.String("err", err.Error()))
			result.Error(c, result.LoginErrorStatus)
		}
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func UserInfo(c *gin.Context) {
	var req UserInfoRequest
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
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     targetUser,
		})
	}
}
