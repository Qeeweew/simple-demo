package model

import (
	"context"
	"time"
)

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
	Extra *UserExtra `gorm:"-:all"`
}
type UserExtra struct {
	FollowCount    int64 `json:"follow_count,omitempty"`
	FanCount       int64 `json:"follower_count,omitempty"`
	IsFollow       bool  `json:"is_follow,omitempty"`
	TotalFavorited int64 `json:"total_favorited,omitempty"`
	WorkCount      int64 `json:"work_count,omitempty"`
	FavoriteCount  int64 `json:"favorite_count,omitempty"`
}

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
	Extra *VideoExtra `gorm:"-:all"`
}
type VideoExtra struct {
	FavoriteCount int64 `json:"favorite_count,omitempty"`
	CommentCount  int64 `json:"comment_count,omitempty"`
	IsFavorite    bool  `json:"is_favorite,omitempty"`
}

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
	FavoriteList(currentId uint, targetId uint) (videos []Video, err error)
}

type MessageService interface {
}

type MessageRepository interface {
}
