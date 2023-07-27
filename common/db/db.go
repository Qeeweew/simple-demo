package db

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var MySQL *gorm.DB
var Redis *redis.Client

func Init() {
	MySQLInit()

	// TODO: 等之后上缓存的时候再放开
	// RedisInit()
}
