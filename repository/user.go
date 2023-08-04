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
	return u.First(user, "id = ?", userID).Error
}

func (u *userRepository) FindByName(username string, user *model.User) error {
	return u.First(user, "name = ?", username).Error
}

func (u *userRepository) FillExtraData(currentUserId uint, targetUser *model.User) error {
	u.Transaction(func(tx *gorm.DB) (err error) {
		targetUser.FollowCount = tx.Where(&model.User{Id: targetUser.Id}).Association("Follows").Count()
		targetUser.FanCount = tx.Where(&model.User{Id: targetUser.Id}).Association("Fans").Count()
		targetUser.IsFollow, err = NewRelationRepository(tx).CheckFollowRelationship(currentUserId, targetUser.Id)
		if err != nil {
			return
		}
		targetUser.FavoriteCount, err = NewFavoriteRepository(tx).UserFavoriteCount(targetUser.Id)
		if err != nil {
			return
		}
		err = tx.Model(&model.Video{}).Where(&model.Video{AuthorId: targetUser.Id}).Count(&targetUser.WorkCount).Error
		if err != nil {
			return
		}
		// TODO: Replace with ORM operation
		err = tx.Raw("SELECT COUNT(*) FROM user INNER JOIN video ON user.id = video.author_id INNER JOIN favorite ON video.id = favorite.video_id WHERE user.id = ?", targetUser.Id).
			Scan(&targetUser.TotalFavorited).Error
		return
	})
	return nil
}
