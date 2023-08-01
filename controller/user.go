package controller

import (
	"net/http"
	"simple-demo/common/model"
	"simple-demo/service"
	"simple-demo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:          1,
		Name:        "zhanglei",
		FollowCount: 10,
		FanCount:    5,
		IsFollow:    true,
	},
}

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
	var user model.User
	user.Name = c.Query("username")
	user.Password = c.Query("password")
	err := service.NewUser().Register(&user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func Login(c *gin.Context) {
	var user model.User
	user.Name = c.Query("username")
	user.Password = c.Query("password")
	err := service.NewUser().Login(&user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.Id),
			Token:    utils.CreateToken(user.Id),
		})
	}
}

func UserInfo(c *gin.Context) {
	targetId, _ := strconv.Atoi(c.Query("user_id"))
	userId, _ := c.Keys["auth_id"].(uint)
	var targetUser model.User
	err := service.NewUser().Info(userId, uint(targetId), &targetUser)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     targetUser,
		})
	}
}
