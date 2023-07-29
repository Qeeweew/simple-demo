package controller

import (
	"net/http"
	"simple-demo/common/db"
	"simple-demo/common/model"
	"simple-demo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var user = model.User{
		Name:     username,
		Password: password,
	}
	if err := db.MySQL.Where("name = ?", username).First(user).Error; err != nil {
		//找不到
		db.MySQL.Create(&user)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   int64(user.ID),
			Token:    utils.CreateToken(user.ID),
		})
	} else {
		//找到
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	logrus.Println(username, password)

	var user model.User

	if err := db.MySQL.Where("name = ?", username).First(&user).Error; err != nil {
		//找不到
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		//找到
		if user.Password == password {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   int64(user.ID),
				Token:    utils.CreateToken(user.ID),
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 2, StatusMsg: "Wrong password"},
			})
		}
	}
}

func UserInfo(c *gin.Context) {
	targetID, _ := strconv.Atoi(c.Query("user_id"))
	userID, _ := c.Keys["auth_id"]

	var targetUser model.User
	if err := db.MySQL.Where("id = ?", targetID).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		// 发起请求的user
		var user model.User
		db.MySQL.First(&user, "ID = ?", userID)
		var isFollow = db.MySQL.Model(&user).Where("follows.user_id = ? AND user.id = ?", targetID, userID).Association("Follows").Count() > 0
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
