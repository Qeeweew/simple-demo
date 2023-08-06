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
	if u == nil {
		panic("???")
	}
	err := u.User(context.Background()).FindByName(user.Name, user)
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

func (u *userService) UserInfo(currentId uint, targetId uint) (user model.User, err error) {
	err = u.tximpl.Transaction(context.Background(), func(txctx context.Context) (err error) {
		err = u.User(txctx).FindById(targetId, &user)
		if err != nil {
			return
		}
		err = u.User(txctx).FillExtraData(currentId, &user, true)
		if err != nil {
			return
		}
		return
	})
	return
}
