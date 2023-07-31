package service

import (
	"simple-demo/common/model"
	"simple-demo/repository"
)

func GetVideo() model.VideoService {
	return NewVideoService(repository.NewVideoRepository(repository.MySQL))
}

func GetUser() model.UserService {
	return NewUserService(repository.NewUserRepository(repository.MySQL))
}

func GetRelation() model.RelationService {
	return NewRelationService(repository.NewRelationRepository(repository.MySQL))
}
