package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		db,
	}
}

func (u *userRepository) Save(user *model.User) error {
	return u.Create(user).Error
}

func (u *userRepository) FindById(userID uint, user *model.User) error {
	var db = u.DB
	return db.First(user, "id = ?", userID).Error
}

func (u *userRepository) FindByName(username string, user *model.User) error {
	var db = u.DB
	return db.First(user, "name = ?", username).Error
}

// TODO: 添加repository其他接口完成各个信息的查询并填充。（也可以直接在这里写）
func (u *userRepository) FillExtraData(currentUserId uint, targetUser *model.User) error {
	targetUser.FollowCount = 0
	targetUser.FanCount = 0
	targetUser.IsFollow = true
	targetUser.TotalFavorited = 0
	targetUser.WorkCount = 0
	targetUser.FanCount = 0
	return nil
}
