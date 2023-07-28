package model

import "gorm.io/gorm"

// User 用户表
type User struct {
	gorm.Model
	Name     string
	Password string
}

// Video 视频表
type Video struct {
	gorm.Model
	PlayUrl  string
	CoverUrl string
	UserID   uint
	User     User
}

// Comment 评论表
type Comment struct {
	gorm.Model
	Content string
	UserID  uint
	User    User
	VideoID uint
	Video   Video
}

// Like 点赞表
type Like struct {
	UserID  uint `gorm:"primaryKey;autoIncrement:false"`
	VideoID uint `gorm:"primaryKey;autoIncrement:false"`
}

// Follow 关注表
type Follow struct {
	UserID   uint `gorm:"primaryKey;autoIncrement:false"`
	FollowID uint `gorm:"primaryKey;autoIncrement:false"`
}
