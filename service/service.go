package service

import (
	"simple-demo/common/db"
	"simple-demo/common/model"
	"simple-demo/repository"

	"gorm.io/gorm"
)

func GetVideo() model.VideoService {
	return NewVideoService(repository.NewVideoRepository(db.MySQL))
}

func GetUser() model.UserService {
	return NewUserService(repository.NewUserRepository(db.MySQL))
}

func GetDB() *gorm.DB {
	return db.MySQL
}
