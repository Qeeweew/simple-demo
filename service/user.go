package service

import (
	"context"
	"errors"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	userInstance *userService
	userOnce     sync.Once
	ErrUserExist = errors.New("user already exists")
	ErrPassword  = errors.New("password error")
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
	err := u.User(context.Background()).FindByName(user.Name, user)
	if err != nil {
		return err
	}
	// 密码校验
	if user.Password != password {
		log.Logger.Error("password error", zap.Any("user", user))
		return ErrPassword
	}
	return nil
}

func (u *userService) Register(user *model.User) error {
	// 检查用户是否注册过
	return u.tximpl.Transaction(context.Background(), func(txctx context.Context) error {
		err := u.User(txctx).FindByName(user.Name, user)

		// 数据库错误
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Logger.Error("mysql error", zap.Error(err))
			return err
		}

		// 用户未注册过
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.User(txctx).Save(user)
		}

		// 用户已注册
		return ErrUserExist

	})
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
