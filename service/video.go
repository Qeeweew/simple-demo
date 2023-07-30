package service

import (
	"simple-demo/common/model"
	"sync"
)

type videoService struct {
	repository model.VideoRepository
}

var (
	videoInstance *videoService
	videoOnce     sync.Once
)

func NewVideoService(r model.VideoRepository) model.VideoService {
	videoOnce.Do(func() {
		videoInstance = &videoService{
			repository: r,
		}
	})
	return videoInstance
}

func (v *videoService) Publish(video *model.Video) error {
	return v.repository.Save(video)
}
func (v *videoService) GetPublishList(userID uint) (videos []model.Video, err error) {
	err = v.repository.FindListByUserID(userID, &videos, 3)
	return
}

func (v *videoService) GetFeedList(limit uint) (videos []model.Video, err error) {
	err = v.repository.FeedList(limit, &videos)
	return
}
