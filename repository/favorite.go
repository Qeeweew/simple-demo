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

func (f *favoriteRepository) UserFavoriteCount(userId uint) (res int64, err error) {
	err = f.Model(&model.Favorite{UserId: userId}).Count(&res).Error
	return
}
func (f *favoriteRepository) VideoFavoriteCount(video_id uint) (res int64, err error) {
	err = f.Model(&model.Favorite{VideoId: video_id}).Count(&res).Error
	return
}

func (f *favoriteRepository) UserFavoriteList(userId uint) (res []model.Video, err error) {
	err = f.Model(&model.Favorite{UserId: userId}).Preload("Video").Preload("Video.Author").Find(&res).Error
	return
}

func (f *favoriteRepository) IsFavorite(userId uint, videoId uint) (res bool, err error) {
	var cnt int64
	err = f.Model(&model.Favorite{UserId: userId, VideoId: videoId}).Count(&cnt).Error
	if err != nil {
		return
	}
	res = cnt > 0
	return
}
