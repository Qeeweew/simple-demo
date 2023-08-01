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

func (f *favoriteRepository) GetUserFavoriteCount(userID uint) (res int64, err error) {
	err = f.Where("user_id = ?", userID).Count(&res).Error
	return
}

func (f *favoriteRepository) GetUserFavoriteList(userID uint) (res []model.Video, err error) {
	err = f.Preload("Video").Preload("Video.Author").Where("user_id = ?", userID).Find(&res).Error
	return
}
