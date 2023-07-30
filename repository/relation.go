package repository

import (
	"gorm.io/gorm"
	"simple-demo/common/model"
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
}

func (r *relationRepository) Follow(userId uint, toUserId uint) error {
	// 检查 user_id 和 to_user_id 是否都存在于数据库中
	var user, toUser model.User
	if err := r.DB.First(&user, userId).Error; err != nil {
		return err
	}
	if err := r.DB.First(&toUser, toUserId).Error; err != nil {
		return err
	}

	// 使用 Association 方法创建关注关系
	if err := r.DB.Model(&user).Association("Follows").Append(&toUser); err != nil {
		return err
	}

	return nil
}

func (r *relationRepository) CheckFollowRelationship(userId uint, toUserId uint) (bool, error) {
	var user model.User
	if err := r.DB.Preload("Follows").First(&user, userId).Error; err != nil {
		return false, err
	}

	for _, follow := range user.Follows {
		if follow.ID == toUserId {
			return true, nil
		}
	}
	return false, nil
}

func NewRelationRepository(db *gorm.DB) model.RelationRepository {
	return &relationRepository{
		DB: db,
	}
}
