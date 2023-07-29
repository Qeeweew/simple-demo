package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) Save(user *model.User) error {
	return u.DB.Create(user).Error
}

func (u *userRepository) FindByID(userID uint, user *model.User) error {
	return u.DB.First(user, "id = ?", userID).Error
}
func (u *userRepository) FindByName(username string, user *model.User) error {
	return u.DB.First(user, "name = ?", username).Error
}
