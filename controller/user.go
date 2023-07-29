package controller

import (
	"net/http"
	"simple-demo/common/db"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/service"
	"simple-demo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
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

func getService() model.UserService {
	return service.NewUserService(repository.NewUserRepository(db.MySQL))
}

func Register(c *gin.Context) {
	var user model.User
	user.Name = c.Query("username")
	user.Password = c.Query("password")
	err := getService().Register(&user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.ID),
			Token:    utils.CreateToken(user.ID),
		})
	}
}

func Login(c *gin.Context) {
	var user model.User
	user.Name = c.Query("username")
	user.Password = c.Query("password")
	err := getService().Login(&user)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.ID),
			Token:    utils.CreateToken(user.ID),
		})
	}
}

func UserInfo(c *gin.Context) {
	targetID, _ := strconv.Atoi(c.Query("user_id"))
	userID, _ := c.Keys["auth_id"].(uint)
	var targetUser model.User
	isFollow, err := getService().Info(userID, uint(targetID), &targetUser)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				Id:            int64(targetUser.ID),
				Name:          targetUser.Name,
				FollowCount:   int64(len(targetUser.Follows)),
				FollowerCount: int64(len(targetUser.Fans)),
				IsFollow:      isFollow,
			},
		})
	}
}
