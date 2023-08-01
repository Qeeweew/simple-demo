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
