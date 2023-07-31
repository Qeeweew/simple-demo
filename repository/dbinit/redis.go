package dbinit

// func RedisInit() {
// 	Redis := redis.NewClient(&redis.Options{
// 		Addr:     config.RedisCfg.Host,
// 		Password: config.RedisCfg.Password,
// 		DB:       0,
// 	})
// 	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
// 		logrus.Panic("connect redis failed: %v", err)
// 	}
// 	logrus.Info("Connect redis succeeded")
// }
