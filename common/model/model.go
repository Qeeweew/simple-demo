package model

import (
	"context"
	"time"
)

/*
message User {
  required int64 id = 1; // 用户id
  required string name = 2; // 用户名称
  optional int64 follow_count = 3; // 关注总数
  optional int64 follower_count = 4; // 粉丝总数
  required bool is_follow = 5; // true-已关注，false-未关注
  optional string avatar = 6; //用户头像
  optional string background_image = 7; //用户个人页顶部大图
  optional string signature = 8; //个人简介
  optional int64 total_favorited = 9; //获赞数量
  optional int64 work_count = 10; //作品数量
  optional int64 favorite_count = 11; //点赞数量
}
*/

// User 用户表
type User struct {
	Id       uint   `gorm:"primarykey" json:"id,omitempty"`
	Name     string `gorm:"not null;unique;index" json:"name,omitempty"`
	Password string `gorm:"not null" json:"-"`

	Follows []*User `gorm:"many2many:follows;" json:"-"`
	Fans    []*User `gorm:"many2many:follows;joinForeignKey:follow_id;" json:"-"`

	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`

	// 不直接储存，需要查询得到
	FollowCount    int64 `gorm:"-:all" json:"follow_count,omitempty"`
	FanCount       int64 `gorm:"-:all" json:"follower_count,omitempty"`
	IsFollow       bool  `gorm:"-:all" json:"is_follow,omitempty"`
	TotalFavorited int64 `gorm:"-:all" json:"total_favorited,omitempty"`
	WorkCount      int64 `gorm:"-:all" json:"work_count,omitempty"`
	FavoriteCount  int64 `gorm:"-:all" json:"favorite_count,omitempty"`
}

/*
  required int64 id = 1; // 视频唯一标识
  required User author = 2; // 视频作者信息
  required string play_url = 3; // 视频播放地址
  required string cover_url = 4; // 视频封面地址
  required int64 favorite_count = 5; // 视频的点赞总数
  required int64 comment_count = 6; // 视频的评论总数
  required bool is_favorite = 7; // true-已点赞，false-未点赞
  required string title = 8; // 视频标题
*/
// Video 视频表
type Video struct {
	Id        uint      `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	UpdatedAt time.Time `gorm:"not null;index" json:"-"`
	AuthorId  uint      `gorm:"not null;index" json:"-"`
	Author    User      `gorm:"foreignKey:AuthorId" json:"author"`
	Title     string    `gorm:"not null;index" json:"title,omitempty"`
	PlayUrl   string    `gorm:"not null" json:"play_url,omitempty"`
	CoverUrl  string    `gorm:"not null" json:"cover_url,omitempty"`

	// 不直接储存，需要后续查询得到
	FavoriteCount int64 `gorm:"-:all" json:"favorite_count,omitempty"`
	CommentCount  int64 `gorm:"-:all" json:"comment_count,omitempty"`
	IsFavorite    bool  `gorm:"-:all" json:"is_favorite,omitempty"`
}

/*
message Comment {
  required int64 id = 1; // 视频评论id
  required User user =2; // 评论用户信息
  required string content = 3; // 评论内容
  required string create_date = 4; // 评论发布日期，格式 mm-dd
}
*/
// Comment 评论表
type Comment struct {
	Id         uint      `gorm:"primarykey" json:"id,omitempty"`
	Content    string    `gorm:"not null" json:"content,omitempty"`
	UserId     uint      `gorm:"not null;index" json:"-"`
	VideoId    uint      `gorm:"not null;index" json:"-"`
	CreatedAt  time.Time `gorm:"not null" json:"-"`
	CreateDate string    `gorm:"-:all" json:"create_date"`
	User       User      `gorm:"foreignKey:UserId"`
}

// Favorite 点赞表
type Favorite struct {
	UserId  uint  `gorm:"primaryKey;autoIncrement:false"`
	VideoId uint  `gorm:"primaryKey;autoIncrement:false"`
	Video   Video `gorm:"foreignKey:VideoId"`
}

/*
message Message {
  required int64 id = 1; // 消息id
  required int64 to_user_id = 2; // 该消息接收者的id
  required int64 from_user_id =3; // 该消息发送者的id
  required string content = 4; // 消息内容
  optional string create_time = 5; // 消息创建时间
}
*/
// Message 聊天消息表
type Message struct {
	Id         uint      `gorm:"primarykey" json:"id,omitempty"`
	FromId     uint      `gorm:"index" json:"from_user_id,omitempty"`
	ToUserId   uint      `gorm:"index" json:"to_user_id,omitempty"`
	Content    string    `gorm:"not null" json:"content,omitempty"`
	CreatedAt  time.Time `gorm:"not null" json:"-"`
	CreateDate string    `gorm:"-:all" json:"create_date"`
}

// // Friend 好友表
// type Friend struct {
// 	gorm.Model
// 	UserId   int64 `gorm:"index"`
// 	FriendId int64 `gorm:"index"`
// }

// 提供访问Repository的接口
type ServiceBase interface {
	User(ctx context.Context) UserRepository
	Video(ctx context.Context) VideoRepository
	Relation(ctx context.Context) RelationRepository
	Comment(ctx context.Context) CommentRepository
	Favorite(ctx context.Context) FavoriteRepository
}

// Service启动Transaction的接口
type ITransaction interface {
	Transaction(ctx context.Context, fn func(txctx context.Context) error) error
}

// UserService : represent the user's services
type UserService interface {
	ServiceBase
	Login(user *User) error
	Register(user *User) error
	UserInfo(curentId uint, targetId uint) (User, error)
}

// UserRepository : represent the user's repository contract
type UserRepository interface {
	Save(user *User) error
	FindById(userId uint, user *User) error
	FindByName(name string, user *User) error
	FillExtraData(currentUserId uint, targetUser *User) error
}

type VideoService interface {
	ServiceBase
	Publish(video *Video) error
	GetPublishList(userId uint, targetId uint) ([]Video, error)
	GetFeedList(userId uint) ([]Video, error)
}

type VideoRepository interface {
	Save(*Video) error
	FindListByUserId(uint, *[]Video) error
	FeedList(uint, *[]Video) error
	FillExtraData(userId uint, video *Video) error
}

type RelationService interface {
	ServiceBase
	FollowAction(token string, toUserId uint, actionType int) error
	FollowList(token string, userId uint) ([]*User, error)
	FanList(token string, userId uint) ([]*User, error)
	// TODO: 用户好友列表
	// ...
}

type RelationRepository interface {
	CheckFollowRelationship(userId uint, toUserId uint) (bool, error)
	Follow(userId uint, toUserId uint) error
	UnFollow(userId uint, toUserId uint) error
	FollowList(userId uint) ([]*User, error)
	FanList(userId uint) ([]*User, error)
}

type CommentRepository interface {
	VideoCommentList(videoId uint) (res []Comment, err error)
	VideoCommentCount(videoId uint) (res int64, err error)
	Create(*Comment) error
	Delete(*Comment) error
}

type CommentService interface {
	CommentAction(isComment bool, comment *Comment) error
	CommentList(userId uint, videoId uint) ([]Comment, error)
}

type FavoriteRepository interface {
	VideoFavoriteCount(videoId uint) (res int64, err error)
	UserFavoriteCount(userId uint) (res int64, err error)
	UserFavoriteList(userId uint) (res []Video, err error)
	IsFavorite(userId uint, videoId uint) (res bool, err error)
	Create(userId uint, videoId uint) error
	Delete(userId uint, videoId uint) error
}

type FavoriteService interface {
	FavoriteAction(isFavorite bool, userId uint, videoId uint) error
}

type MessageService interface {
}

type MessageRepository interface {
}
