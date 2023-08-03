package service

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository"
	"simple-demo/repository/dbcore"
	"sync"
)

type commentService struct {
	model.ServiceBase
	tximpl model.ITransaction
}

var (
	commentInstance *commentService
	commentOnce     sync.Once
)

// NewService: construction function, injected by user repository
func NewComment() model.CommentService {
	userOnce.Do(func() {
		commentInstance = &commentService{
			repository.NewTableVistor(),
			dbcore.NewTxImpl(),
		}
	})
	return commentInstance
}

func (c *commentService) CommentAction(isComment bool, comment *model.Comment) error {
	if isComment {
		return c.Comment(context.Background()).Create(comment)
	} else {
		return c.Comment(context.Background()).Delete(comment)
	}
}

func (c *commentService) CommentList(userId uint, videoId uint) (comments []model.Comment, err error) {
	err = c.tximpl.Transaction(
		context.Background(),
		func(txctx context.Context) (err error) {
			comments, err = c.Comment(context.Background()).VideoCommentList(videoId)
			if err != nil {
				return
			}
			for i := range comments {
				c.User(context.Background()).FillExtraData(userId, &comments[i].User)
			}
			return nil
		})
	return
}
