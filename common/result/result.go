package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 200 OK
	SuccessStatus = newStatus(200, "success")

	// 400 BAD
	QueryParamErrorStatus   = newStatus(400, "请求的参数错误")
	LoginErrorStatus        = newStatus(400, "登录发生错误")
	RegisterErrorStatus     = newStatus(400, "注册发生错误")
	UsernameExitErrorStatus = newStatus(400, "用户名已存在")
	TokenErrorStatus        = newStatus(400, "token 错误")
	InfoErrorStatus         = newStatus(400, "无法获取该用户信息")
	FileErrorStatus         = newStatus(400, "文件上传失败")
	PublishErrorStatus      = newStatus(400, "发布时出现错误")
	FeedErrorStatus         = newStatus(400, "获取视频流出错")
	EmptyErrorStatus        = newStatus(400, "用户名或密码为空") // should be useless
	FollowErrorStatus       = newStatus(400, "关注失败")
	UnFollowErrorStatus     = newStatus(400, "取消关注失败")
	FavoriteErrorStatus     = newStatus(400, "点赞失败")
	FollowListErrorStatus   = newStatus(400, "获取关注列表时发生了错误")
	FollowerListErrorStatus = newStatus(400, "获取粉丝列表时发生了错误")
	MissingTokenErrorStatus = newStatus(400, "缺失Token")
	PasswordErrorStatus     = newStatus(400, "密码错误")

	// 401 WITHOUT PERMISSION
	NoLoginErrorStatus = newStatus(401, "用户未登录")

	// 403 ILLEGAL OPERATION
	PermissionErrorStatus = newStatus(403, "操作非法")

	// 404 NOT FOUND
	CommentNotExitErrorStatus = newStatus(404, "评论不存在")
	VideoNotExitErrorStatus   = newStatus(404, "视频不存在")

	// 500 INTERNAL ERROR
	ServerErrorStatus = newStatus(500, "服务器内部错误")
)

type Status struct {
	StatusCode int64
	StatusMsg  string
}

func (s Status) Code() int64 {
	return s.StatusCode
}

func (s Status) Mag() string {
	return s.StatusMsg
}

func newStatus(code int64, msg string) Status {
	return Status{code, msg}
}

type R map[string]interface{}

func Success(c *gin.Context, r R) {
	//value := reflect.ValueOf(data)
	h := gin.H{
		"status_code": 0,
		"status_msg":  "success",
	}
	for s, v := range r {
		h[s] = v
	}
	c.JSON(http.StatusOK, h)
}
func SuccessBasic(c *gin.Context) {
	h := gin.H{
		"status_code": 0,
		"status_msg":  "success",
	}
	c.JSON(http.StatusOK, h)
}

func Error(c *gin.Context, status Status) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": status.StatusCode,
		"status_msg":  status.StatusMsg,
	})
}
