package model

import (
	"context"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	Name     string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`

	Follows []*User `gorm:"many2many:follows;"`
	Fans    []*User `gorm:"many2many:follows;joinForeignKey:follow_id;"`

	// gorm会忽略，用于保存结果、转换到Controlle.User
	FollowCount     int    `gorm:"-:all"`
	FanCount        int    `gorm:"-:all"`
	Avatar          string `gorm:"-:all"`
	BackgroundImage string `gorm:"-:all"`
	Signature       string `gorm:"-:all"`
	IsFollow        bool   `gorm:"-:all"`
}

// Video 视频表
type Video struct {
	gorm.Model
	UserId   uint   `gorm:"not null;index"`
	Title    string `gorm:"not null;index"`
	PlayUrl  string `gorm:"not null"`
	CoverUrl string // `gorm:"not null"`

	// has many
	Comments []Comment
	Favors   []Favor

	// gorm会忽略，用于保存结果、转换到Controlle.Video
	FavoriteCount int  `gorm:"-:all"`
	CommentCount  int  `gorm:"-:all"`
	IsFavorite    bool `gorm:"-:all"`
}

// Comment 评论表
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint

	VideoID uint
}

// Favor 点赞表
type Favor struct {
	gorm.Model
	UserID uint

	VideoID uint
}

// Message 聊天消息表
type Message struct {
	gorm.Model
	FromId   int64  `gorm:"index"`
	ToUserId int64  `gorm:"index"`
	Content  string `gorm:"not null"`
}

// Friend 好友表
type Friend struct {
	gorm.Model
	UserId   int64 `gorm:"index"`
	FriendId int64 `gorm:"index"`
}

type ServiceBase interface {
	User(ctx context.Context) UserRepository
	Video(ctx context.Context) VideoRepository
	Relation(ctx context.Context) RelationRepository
}

type ITransaction interface {
	Transaction(ctx context.Context, fn func(txctx context.Context) error) error
}

// UserService : represent the user's services
type UserService interface {
	ServiceBase
	Login(user *User) error
	Register(user *User) error
	Info(userID uint, targetID uint, user *User) error
}

// UserRepository : represent the user's repository contract
type UserRepository interface {
	Save(user *User) error
	FindByID(userID uint, user *User, preload uint) error
	FindByName(name string, user *User, preload uint) error
}

type VideoService interface {
	ServiceBase
	Publish(video *Video) error
	GetPublishList(userID uint) ([]Video, error)
	GetFeedList(userID uint) ([]Video, error)
}

type VideoRepository interface {
	Save(*Video) error
	FindListByUserID(uint, *[]Video, uint) error
	FeedList(uint, *[]Video) error
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
