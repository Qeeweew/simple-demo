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

func (v *videoRepository) FindListByUserId(userID uint, videos *[]model.Video) error {
	return v.Preload("Author").Find(videos, "user_id = ?", userID).Error
}

func (v *videoRepository) FeedList(limit uint, videos *[]model.Video) error {
	return v.Preload("Author").Limit(int(limit)).Order("created_at DESC").Find(videos).Error
}

// TODO: Fill `comment_count` `isFavorate``

func (v *videoRepository) FillExtraData(userId uint, video *model.Video) (err error) {
	return v.Transaction(func(tx *gorm.DB) (err error) {
		err = NewUserRepository(tx).FillExtraData(userId, &video.Author)
		if err != nil {
			return
		}
		video.FavoriteCount, err = NewFavoriteRepository(tx).GetVideoFavoriteCount(video.Id)
		if err != nil {
			return
		}

		return
	})
}
