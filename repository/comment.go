package repository

import (
	"errors"
	"fmt"
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

func (c *commentRepository) VideoCommentList(videoId uint) (res []model.Comment, err error) {
	err = c.Order("created_at DESC").Where("video_id = ?", videoId).Preload("User").Find(&res).Error
	return
}

func (c *commentRepository) VideoCommentCount(videoId uint) (res int64, err error) {
	err = c.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&res).Error
	return
}

func (c *commentRepository) Create(comment *model.Comment) error {
	err := c.DB.Create(comment).Error
	if err != nil {
		return err
	}
	_, mon, day := comment.CreatedAt.UTC().Date()
	comment.CreateDate = fmt.Sprintf("%02d:%02d", mon, day)
	return nil
}

func (c *commentRepository) Delete(comment *model.Comment) error {
	return c.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.First(&comment).Error; err != nil {
				return errors.New("Record Not Found")
			}
			return tx.Delete(comment).Error
		})
}
