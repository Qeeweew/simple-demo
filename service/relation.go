package service

import (
	"context"
	"errors"
	"simple-demo/common/log"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"

	"github.com/redis/go-redis/v9"
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

// 这里好友列表 = 粉丝列表 + 关注列表
// 这样保证 `好友关系` 是对称的
func (r *relationService) FriendList(userId uint) (friendUsers []model.FriendUser, err error) {
	// 查询好友列表
	var users []model.User
	err = r.tximpl.Transaction(context.Background(), func(txctx context.Context) (err error) {
		users, err = r.Relation(txctx).FriendList(userId)
		if err != nil {
			log.Logger.Error("friend list error")
			return
		}
		for i := range users {
			if err = r.User(txctx).FillExtraData(userId, &users[i], false); err != nil {
				return
			}
		}
		return
	})
	if err != nil {
		return
	}
	ctx := context.Background()

	// redis 查询好友间最后一条消息
	pipe := r.RedisClient().Pipeline()
	var vals []*redis.MapStringStringCmd
	for _, u := range users {
		vals = append(vals, pipe.HGetAll(ctx, genChatKey(userId, u.Id)))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	for i := range users {
		var msg MessageLatest
		var msgType int64
		if vals[i].Err() == nil {
			vals[i].Scan(&msg)
			if msg.Sender == userId {
				msgType = 1
			} else {
				msgType = 0
			}
		}

		friendUsers = append(friendUsers, model.FriendUser{
			User:    users[i],
			Message: msg.Message,
			MsgType: msgType,
		})
	}
	return
}
