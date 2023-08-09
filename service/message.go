package service

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"
)

type messageService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	messageInstance *messageService
	messageOnce     sync.Once
)

func NewMessage() model.MessageService {
	messageOnce.Do(func() {
		messageInstance = &messageService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return messageInstance
}

func (m *messageService) SendMessage(userId uint, toUserId uint, content string) error {
	err := m.tximpl.Transaction(
		context.Background(),
		func(txctx context.Context) (err error) {
			// TODO: Check friend relation
			err = m.Message(txctx).Create(&model.Message{
				FromUserId: userId,
				ToUserId:   toUserId,
				Content:    content,
			})
			return
		})
	if err == nil {
		m.RedisClient().HSet(context.Background(), genChatKey(userId, toUserId),
			"message", content, "sender", userId)
	}
	return err
}

func (m *messageService) ChatHistory(preMsgTime int64, userId uint, toUserId uint) ([]model.Message, error) {
	return m.Message(context.Background()).MessageList(preMsgTime, userId, toUserId)
}
