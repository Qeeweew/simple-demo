package service

import (
	"context"
	"errors"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"
)

type userService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	userInstance *userService
	userOnce     sync.Once
)

// NewService: construction function, injected by user repository
func NewUser() model.UserService {
	userOnce.Do(func() {
		userInstance = &userService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return userInstance
}

func (u *userService) Login(user *model.User) error {
	password := user.Password
	err := u.User(context.Background()).FindByName(user.Name, user, 0)
	if err != nil {
		return err
	}
	if user.Password != password {
		return errors.New("wrong password")
	}
	return nil
}

func (u *userService) Register(user *model.User) error {
	return u.User(context.Background()).Save(user)
}

func (u *userService) Info(userId uint, targetId uint, user *model.User) (err error) {
	err = u.User(context.Background()).FindById(targetId, user, 3)
	if err != nil {
		return
	}
	for _, fan := range user.Fans {
		if fan.Id == userId {
			user.IsFollow = true
		}
	}
	user.FollowCount = len(user.Follows)
	user.FanCount = len(user.Fans)
	return
}
