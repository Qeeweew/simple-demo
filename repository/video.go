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

func (v *videoRepository) FindListByUserID(userID uint, videos *[]model.Video, preload uint) error {
	var db = v.DB
	if preload&1 != 0 {
		db = db.Preload("Comments")
	}
	if preload&2 != 0 {
		db = db.Preload("Favors")
	}
	err := v.Find(videos, "user_id = ?", userID).Error
	if err != nil {
		return err
	}
	for i := range *videos {
		(*videos)[i].FavoriteCount = len((*videos)[i].Favors)
		(*videos)[i].CommentCount = len((*videos)[i].Comments)
	}
	return nil
}

func (v *videoRepository) FeedList(limit uint, videos *[]model.Video) error {
	err := v.Preload("Comments").Preload("Favors").Limit(int(limit)).Order("created_at DESC").Find(videos).Error
	if err != nil {
		return err
	}
	for i := range *videos {
		(*videos)[i].FavoriteCount = len((*videos)[i].Favors)
		(*videos)[i].CommentCount = len((*videos)[i].Comments)
	}
	return nil
}
