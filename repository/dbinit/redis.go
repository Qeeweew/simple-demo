package dbinit

import (
	"context"
	"simple-demo/common/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func RedisInit() *redis.Client {
	Redis := redis.NewClient(&redis.Options{
		Addr:     config.RedisCfg.Host,
		Password: config.RedisCfg.Password,
		DB:       0,
	})
	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		logrus.Panicf("connect redis failed: %v", err)
	}
	logrus.Info("Connect redis succeeded")
	return Redis
}
