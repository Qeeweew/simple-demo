package repository

import (
	"simple-demo/repository/dbinit"

	"gorm.io/gorm"
)

var MySQL *gorm.DB

// var Redis *redis.Client

func Init() {
	MySQL = dbinit.MySQLInit()

	// TODO: 等之后上缓存的时候再放开
	// Re
}
