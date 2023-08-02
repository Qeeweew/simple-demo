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

func (c *commentRepository) VideoCommentList(videoID uint) (res []model.Comment, err error) {
	err = c.Model(&model.Comment{VideoId: videoID}).Preload("User").Error
	return
}

func (c *commentRepository) VideoCommentCount(videoID uint) (res int64, err error) {
	err = c.Model(&model.Comment{VideoId: videoID}).Count(&res).Error
	return
}
