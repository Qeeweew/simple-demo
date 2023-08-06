package service

import (
	"context"
	"errors"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
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
func (r *relationService) FollowAction(currentId uint, targetId uint, actionType int) error {
	return r.tximpl.Transaction(context.Background(), func(txctx context.Context) error {
		// 用户不能关注自己
		if currentId == targetId {
			log.Logger.Error("can not follow yourself")
			return ErrFollowSelf
		}

		// 查询关注关系
		isFollow, err := r.Relation(txctx).CheckFollowRelationship(currentId, targetId)
		if err != nil {
			log.Logger.Error("check follow relationship error")
			return err
		}

		// 关注
		if actionType == doFollow && !isFollow {
			if err := r.Relation(txctx).Follow(currentId, targetId); err != nil {
				log.Logger.Error("follow error")
				return err
			}
		}

		// 取消关注
		if actionType == unFollow && isFollow {
			if err := r.Relation(txctx).UnFollow(currentId, targetId); err != nil {
				log.Logger.Error("unfollow error")
				return err
			}
		}
		return nil
	})

}

func (r *relationService) FollowList(currentId uint, userId uint) (users []*model.User, err error) {
	// 查询关注列表
	err = r.tximpl.Transaction(context.Background(), func(txctx context.Context) (err error) {
		users, err = r.Relation(txctx).FollowList(userId)
		if err != nil {
			log.Logger.Error("follow list error")
			return
		}
		for i := range users {
			if err = r.User(txctx).FillExtraData(currentId, users[i], false); err != nil {
				return
			}
		}
		return
	})
	return
}

func (r *relationService) FanList(currentId uint, userId uint) (users []*model.User, err error) {
	// 查询粉丝列表
	err = r.tximpl.Transaction(context.Background(), func(txctx context.Context) (err error) {
		users, err = r.Relation(txctx).FanList(userId)
		if err != nil {
			log.Logger.Error("follow list error")
			return
		}
		for i := range users {
			if err = r.User(txctx).FillExtraData(currentId, users[i], false); err != nil {
				return
			}
		}
		return
	})
	return
}
