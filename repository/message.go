package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type messageRepository struct {
	*gorm.DB
}

func NewMessageRepository(db *gorm.DB) model.MessageRepository {
	return &messageRepository{
		db,
	}
}

func (m *messageRepository) Create(message *model.Message) (err error) {
	return m.DB.Create(message).Error
}

func (m *messageRepository) MessageList(preMsgTime int64, userId uint, friendId uint) (messages []model.Message, err error) {
	err = m.DB.Where("create_time > ?", preMsgTime).Where(
		m.DB.Where("from_user_id = ? AND to_user_id = ?", userId, friendId).
			Or("from_user_id = ? AND to_user_id = ?", friendId, userId)).
		Find(&messages).Error
	if err != nil {
		return
	}
	return
}
