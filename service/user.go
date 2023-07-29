package service

import (
	"errors"
	"simple-demo/common/model"
	"sync"
)

type userService struct {
	userRepository model.UserRepository
}

var (
	instance *userService
	once     sync.Once
)

// NewUserService: construction function, injected by user repository
func NewUserService(r model.UserRepository) model.UserService {
	once.Do(func() {
		instance = &userService{
			userRepository: r,
		}
	})
	return instance
}

func (u *userService) Login(user *model.User) error {
	password := user.Password
	err := u.userRepository.FindByName(user.Name, user, 0)
	if err != nil {
		return err
	}
	if user.Password != password {
		return errors.New("wrong password")
	}
	return nil
}

func (u *userService) Register(user *model.User) error {
	return u.userRepository.Save(user)
}

func (u *userService) Info(userID uint, targetID uint, user *model.User) (isFollow bool, err error) {
	err = u.userRepository.FindByID(targetID, user, 3)
	if err != nil {
		return
	}
	for _, fan := range user.Fans {
		if fan.ID == userID {
			isFollow = true
		}
	}
	return
}
