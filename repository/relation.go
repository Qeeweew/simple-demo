package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type relationRepository struct {
	DB *gorm.DB
}

func (r *relationRepository) FanList(userId uint) ([]*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Fans").First(&user, userId).Error; err != nil {
		return nil, err
	}
	return user.Fans, nil
}

func (r *relationRepository) FollowList(userId uint) ([]*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Follows").First(&user, userId).Error; err != nil {
		return nil, err
	}
	return user.Follows, nil
}

func (r *relationRepository) UnFollow(userId uint, toUserId uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 检查 user_id 和 to_user_id 是否都存在于数据库中
		var user, toUser model.User
		if err := r.DB.First(&user, userId).Error; err != nil {
			return err
		}
		if err := r.DB.First(&toUser, toUserId).Error; err != nil {
			return err
		}
		// 使用 Association 方法取消关注关系
		if err := r.DB.Model(&user).Association("Follows").Delete(&toUser); err != nil {
			return err
		}
		return nil
	})
}

func (r *relationRepository) Follow(userId uint, toUserId uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// 检查 user_id 和 to_user_id 是否都存在于数据库中
		var user, toUser model.User
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		}
		if err := tx.First(&toUser, toUserId).Error; err != nil {
			return err
		}
		// 使用 Association 方法创建关注关系
		if err := tx.Model(&user).Association("Follows").Append(&toUser); err != nil {
			return err
		}
		return nil
	})
}

func (r *relationRepository) CheckFollowRelationship(userId uint, toUserId uint) (res bool, err error) {
	var user model.User
	if err = r.DB.First(&user).Error; err != nil {
		return
	}
	var cnt int64
	err = r.DB.Table("follow").Where("user_id = ? AND follow_id = ?", userId, toUserId).Count(&cnt).Error
	if err != nil {
		return
	}
	res = cnt > 0
	return
}

func NewRelationRepository(db *gorm.DB) model.RelationRepository {
	return &relationRepository{
		DB: db,
	}
}
