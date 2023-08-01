package service

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"
)

type videoService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	videoInstance *videoService
	videoOnce     sync.Once
)

func NewVideo() model.VideoService {
	videoOnce.Do(func() {
		videoInstance = &videoService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return videoInstance
}

func (v *videoService) Publish(video *model.Video) error {
	return v.Video(context.Background()).Save(video)
}
func (v *videoService) GetPublishList(userId uint) (videos []model.Video, err error) {
	err = v.Video(context.Background()).FindListByUserId(userId, &videos)
	return
}

func (v *videoService) GetFeedList(limit uint) (videos []model.Video, err error) {
	err = v.Video(context.Background()).FeedList(limit, &videos)
	return
}
