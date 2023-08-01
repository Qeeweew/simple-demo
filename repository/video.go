package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type videoRepository struct {
	*gorm.DB
}

func NewVideoRepository(db *gorm.DB) model.VideoRepository {
	return &videoRepository{
		db,
	}
}

func (v *videoRepository) Save(video *model.Video) error {
	return v.Create(video).Error
}

func (v *videoRepository) FindListByUserID(userID uint, videos *[]model.Video) error {
	return v.Preload("Author").Find(videos, "user_id = ?", userID).Error
}

func (v *videoRepository) FeedList(limit uint, videos *[]model.Video) error {
	return v.Preload("Author").Limit(int(limit)).Order("created_at DESC").Find(videos).Error
}
