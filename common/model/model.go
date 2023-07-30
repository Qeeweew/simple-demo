package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	gorm.Model
	Name     string `gorm:"not null;unique;index"`
	Password string `gorm:"not null"`
	// Avatar          string `gorm:"not null"`
	// BackgroundImage string `gorm:"not null"`
	// Signature       string

	Follows []*User `gorm:"many2many:follows;"`
	Fans    []*User `gorm:"many2many:follows;joinForeignKey:follow_id;"`
}

// Video 视频表
type Video struct {
	gorm.Model
	UserId   uint   `gorm:"not null;index"`
	Title    string `gorm:"not null;index"`
	PlayUrl  string `gorm:"not null"`
	CoverUrl string `gorm:"not null"`

	// has many
	Comments []Comment
	Favors   []Favor
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

// UserService : represent the user's services
type UserService interface {
	Login(user *User) error
	Register(user *User) error
	Info(userID uint, targetID uint, user *User) (bool, error)
}

// UserRepository : represent the user's repository contract
type UserRepository interface {
	Save(user *User) error
	FindByID(userID uint, user *User, preload uint) error
	FindByName(name string, user *User, preload uint) error
}

type RelationService interface {
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
