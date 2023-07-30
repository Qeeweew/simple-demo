package service

import (
	"errors"
	"simple-demo/common/model"
	"sync"
)

type userService struct {
	repository model.UserRepository
}

var (
	userInstance *userService
	userOnce     sync.Once
)

// NewService: construction function, injected by user repository
func NewUserService(r model.UserRepository) model.UserService {
	userOnce.Do(func() {
		userInstance = &userService{
			repository: r,
		}
	})
	return userInstance
}

func (u *userService) Login(user *model.User) error {
	password := user.Password
	err := u.repository.FindByName(user.Name, user, 0)
	if err != nil {
		return err
	}
	if user.Password != password {
		return errors.New("wrong password")
	}
	return nil
}

func (u *userService) Register(user *model.User) error {
	return u.repository.Save(user)
}

func (u *userService) Info(userID uint, targetID uint, user *model.User) (err error) {
	err = u.repository.FindByID(targetID, user, 3)
	if err != nil {
		return
	}
	for _, fan := range user.Fans {
		if fan.ID == userID {
			user.IsFollow = true
		}
	}
	user.FollowCount = len(user.Follows)
	user.FanCount = len(user.Fans)
	return
}
