package service

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"
)

type favoriteService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	favoriteInstance *favoriteService
	favoriteOnce     sync.Once
)

// NewService: construction function, injected by user repository
func NewFavorite() model.FavoriteService {
	userOnce.Do(func() {
		favoriteInstance = &favoriteService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return favoriteInstance
}

func (f *favoriteService) FavoriteAction(isFavorite bool, userId uint, videoId uint) error {
	if isFavorite {
		return f.Favorite(context.Background()).Create(userId, videoId)
	} else {
		return f.Favorite(context.Background()).Delete(userId, videoId)
	}
}

func (f *favoriteService) FavoriteList(currentId uint, targetId uint) (videos []model.Video, err error) {
	f.tximpl.Transaction(context.Background(), func(txctx context.Context) (err error) {
		videos, err = f.Favorite(txctx).UserFavoriteList(targetId)
		if err != nil {
			return
		}
		for i := range videos {
			err = f.Video(txctx).FillExtraData(currentId, &videos[i])
			if err != nil {
				return
			}
		}
		return
	})
	return
}
