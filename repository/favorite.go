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
	err = f.Model(&model.Favorite{}).Where(&model.Favorite{UserId: userId}).Count(&res).Error
	return
}
func (f *favoriteRepository) VideoFavoriteCount(video_id uint) (res int64, err error) {
	err = f.Model(&model.Favorite{}).Where(&model.Favorite{VideoId: video_id}).Count(&res).Error
	return
}

func (f *favoriteRepository) UserFavoriteList(userId uint) (res []model.Video, err error) {
	var favorites []model.Favorite
	err = f.Model(&model.Favorite{}).Where(&model.Favorite{UserId: userId}).Preload("Video").Preload("Video.Author").Find(&favorites).Error
	for i := range favorites {
		res = append(res, favorites[i].Video)
	}
	return
}

func (f *favoriteRepository) IsFavorite(userId uint, videoId uint) (res bool, err error) {
	var cnt int64
	err = f.Model(&model.Favorite{}).Where(&model.Favorite{UserId: userId, VideoId: videoId}).Count(&cnt).Error
	if err != nil {
		return
	}
	res = cnt > 0
	return
}

func (f *favoriteRepository) Create(userId uint, videoId uint) (err error) {
	return f.DB.Create(&model.Favorite{UserId: userId, VideoId: videoId}).Error
}

func (f *favoriteRepository) Delete(userId uint, videoId uint) (err error) {
	return f.DB.Delete(&model.Favorite{UserId: userId, VideoId: videoId}).Error
}
