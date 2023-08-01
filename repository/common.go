package repository

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository/dbcore"
)

type tableVisitor struct{}

func NewTableVistor() *tableVisitor {
	return &tableVisitor{}
}

func (*tableVisitor) User(ctx context.Context) model.UserRepository {
	return NewUserRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Video(ctx context.Context) model.VideoRepository {
	return NewVideoRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Relation(ctx context.Context) model.RelationRepository {
	return NewRelationRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Favorite(ctx context.Context) model.FavoriteRepository {
	return NewFavoriteRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Comment(ctx context.Context) model.CommentRepository {
	return NewCommentRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Message(ctx context.Context) model.MessageRepository {
	return NewMessageRepository(dbcore.GetDB(ctx))
}
