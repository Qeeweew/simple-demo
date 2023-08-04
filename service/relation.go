package service

import (
	"context"
	"errors"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"simple-demo/utils"
	"sync"
)

type relationService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	ErrFollowSelf    = errors.New("can not follow yourself")
	relationOnce     sync.Once
	relationInstance *relationService
)

func NewRelation() model.RelationService {
	relationOnce.Do(func() {
		relationInstance = &relationService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return relationInstance
}

const (
	doFollow = iota + 1
	unFollow
)

// FollowAction
// 关注用户/取消关注
// 解析token，拿到当前用户的id
// 用户不能关注自己
// /*
func (r *relationService) FollowAction(token string, toUserId uint, actionType int) error {
	// 解析token
	userId, err := utils.ParseToken(token)
	if err != nil {
		log.Logger.Error("token parse error")
		return err
	}

	return r.tximpl.Transaction(context.Background(), func(txctx context.Context) error {
		// 用户不能关注自己
		if userId == toUserId {
			log.Logger.Error("can not follow yourself")
			return ErrFollowSelf
		}

		// 查询关注关系
		isFollow, err := r.Relation(txctx).CheckFollowRelationship(userId, toUserId)
		if err != nil {
			log.Logger.Error("check follow relationship error")
			return err
		}

		// 关注
		if actionType == doFollow && !isFollow {
			if err := r.Relation(txctx).Follow(userId, toUserId); err != nil {
				log.Logger.Error("follow error")
				return err
			}
		}

		// 取消关注
		if actionType == unFollow && isFollow {
			if err := r.Relation(txctx).UnFollow(userId, toUserId); err != nil {
				log.Logger.Error("unfollow error")
				return err
			}
		}
		return nil
	})

}

func (r *relationService) FollowList(token string, userId uint) ([]*model.User, error) {
	// 解析token
	userId, err := utils.ParseToken(token)
	if err != nil {
		log.Logger.Error("token parse error")
		return nil, err
	}

	// 查询关注列表
	users, err := r.Relation(context.Background()).FollowList(userId)
	if err != nil {
		log.Logger.Error("follow list error")
		return nil, err
	}
	return users, nil
}

func (r *relationService) FanList(token string, userId uint) ([]*model.User, error) {
	// 解析token
	userId, err := utils.ParseToken(token)
	if err != nil {
		log.Logger.Error("token parse error")
		return nil, err
	}

	// 查询粉丝列表
	users, err := r.Relation(context.Background()).FanList(userId)
	if err != nil {
		log.Logger.Error("fan list error")
		return nil, err
	}
	return users, nil
}
