package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type commentRepository struct {
	*gorm.DB
}

func NewCommentRepository(db *gorm.DB) model.CommentRepository {
	return &commentRepository{
		db,
	}
}

func (c *commentRepository) GetVideoCommentList(videoID uint) (res []model.Comment, err error) {
	err = c.Preload("User").Find(&res, "video_id = ?", videoID).Error
	return
}

func (c *commentRepository) GetVideoCommentCount(videoID uint) (res int64, err error) {
	err = c.Where("video_id = ?", videoID).Count(&res).Error
	return
}
