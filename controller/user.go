package controller

import (
	"net/http"
	"simple-demo/common/db"
	"simple-demo/common/model"
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

	var user model.User

	if err := db.MySQL.Where("name = ?", username).First(user).Error; err != nil {
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
	token := c.Query("token")
	_, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 2, StatusMsg: "Invalid token"},
		})
		return
	}
	askID, _ := strconv.Atoi(c.Query("user_id"))
	var user model.User

	if err := db.MySQL.Where("id = ?", askID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				Id:            int64(user.ID),
				Name:          user.Name,
				FollowCount:   int64(len(user.Follows)),
				FollowerCount: int64(len(user.Fans)),
				IsFollow:      true,
			},
		})
	}
}
