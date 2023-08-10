package repository

import (
	"simple-demo/common/model"
	"time"

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
	err = m.DB.Where("created_at > ?", preMsgTime).Where(
		m.DB.Where("from_user_id = ? AND to_user_id = ?", userId, friendId).
			Or("from_user_id = ? AND to_user_id = ?", friendId, userId)).
		Find(&messages).Error
	if err != nil {
		return
	}
	for i := range messages {
		messages[i].CreateDate = time.Unix(messages[i].CreatedAt, 0).Format(time.DateTime)
	}
	return
}
