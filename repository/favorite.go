package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type favoriteRepository struct {
	*gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) model.FavoriteRepository {
	return &favoriteRepository{
		db,
	}
}

func (f *favoriteRepository) GetUserFavoriteCount(userId uint) (res int64, err error) {
	err = f.Where("user_id = ?", userId).Count(&res).Error
	return
}
func (f *favoriteRepository) GetVideoFavoriteCount(video_id uint) (res int64, err error) {
	err = f.Where("video_id = ?", video_id).Count(&res).Error
	return
}

func (f *favoriteRepository) GetUserFavoriteList(userId uint) (res []model.Video, err error) {
	err = f.Preload("Video").Preload("Video.Author").Where("user_id = ?", userId).Find(&res).Error
	return
}
